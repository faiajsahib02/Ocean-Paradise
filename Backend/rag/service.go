package rag

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/ledongthuc/pdf"
	"github.com/pgvector/pgvector-go"
	"github.com/sashabaranov/go-openai"
)

type Service struct {
	db     *sqlx.DB
	client *openai.Client
}

func NewService(db *sqlx.DB, apiKey string) *Service {
	return &Service{
		db:     db,
		client: openai.NewClient(apiKey),
	}
}

// 1. IngestPDF reads a PDF file, chunks it, and saves vectors
func (s *Service) IngestPDF(ctx context.Context, filePath string) error {
	// A. Read the PDF
	content, err := readPdf(filePath)
	if err != nil {
		return err
	}

	// B. Clean up the text: collapse excessive whitespace
	// PDF extraction often puts each word on its own line
	cleaned := strings.Join(strings.Fields(content), " ")

	// C. Chunk the text by bullet points (●) or double newlines
	chunks := strings.Split(cleaned, "●")
	// If no bullet points found, try splitting by sentences (~300 chars each)
	if len(chunks) <= 1 {
		chunks = chunkBySize(cleaned, 300)
	}

	for _, chunk := range chunks {
		chunk = strings.TrimSpace(chunk)
		if len(chunk) < 20 {
			continue // Skip empty or tiny lines
		}

		// C. Get Embedding from OpenAI
		resp, err := s.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
			Input: []string{chunk},
			Model: openai.AdaEmbeddingV2,
		})
		if err != nil {
			return fmt.Errorf("openai error: %w", err)
		}

		// D. Save to Database
		vector := pgvector.NewVector(resp.Data[0].Embedding)
		_, err = s.db.ExecContext(ctx,
			`INSERT INTO documents (content, embedding) VALUES ($1, $2)`,
			chunk, vector,
		)
		if err != nil {
			return fmt.Errorf("db error: %w", err)
		}
	}
	return nil
}

// 2. GetAnswer takes a user question, finds context, and asks LLM
func (s *Service) GetAnswer(ctx context.Context, userQuestion string) (string, error) {
	// A. Embed the User's Question
	resp, err := s.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{userQuestion},
		Model: openai.AdaEmbeddingV2,
	})
	if err != nil {
		return "", fmt.Errorf("embedding error: %w", err)
	}
	questionVector := pgvector.NewVector(resp.Data[0].Embedding)

	// B. Search Postgres for the 3 most relevant paragraphs
	// We use the '<=>' operator for Cosine Distance (Similarity)
	var chunks []string
	err = s.db.SelectContext(ctx, &chunks,
		`SELECT content FROM documents 
		 ORDER BY embedding <=> $1 LIMIT 3`,
		questionVector,
	)
	if err != nil {
		return "", fmt.Errorf("search error: %w", err)
	}

	// C. Construct the "Sandwich" Prompt
	contextBlock := strings.Join(chunks, "\n\n")
	systemPrompt := fmt.Sprintf(`
You are a helpful Concierge at Oasis Hotel.
Answer the guest's question using ONLY the context below.
If the answer is not in the context, politely say you don't know and offer to call the front desk.

CONTEXT:
%s
`, contextBlock)

	// D. Generate Answer with GPT-4 (or gpt-3.5-turbo)
	chatResp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4, // Or openai.GPT3Dot5Turbo
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userQuestion},
		},
	})
	if err != nil {
		return "", fmt.Errorf("llm error: %w", err)
	}

	return chatResp.Choices[0].Message.Content, nil
}

// Helper function to read PDF text
func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var textBuilder strings.Builder
	// Basic text extraction (for advanced PDFs, use a better library)
	totalPage := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		text, _ := p.GetPlainText(nil)
		textBuilder.WriteString(text)
		textBuilder.WriteString("\n\n") // Add separation
	}
	return textBuilder.String(), nil
}

// chunkBySize splits text into chunks of approximately maxSize characters
func chunkBySize(text string, maxSize int) []string {
	var chunks []string
	sentences := strings.Split(text, ". ")
	current := ""
	for _, s := range sentences {
		if len(current)+len(s) > maxSize && current != "" {
			chunks = append(chunks, strings.TrimSpace(current))
			current = ""
		}
		current += s + ". "
	}
	if strings.TrimSpace(current) != "" {
		chunks = append(chunks, strings.TrimSpace(current))
	}
	return chunks
}
