package raghandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"oasis/backend/util"
	"oasis/backend/rag"
)

type Handler struct {
	ragSvc *rag.Service
}

func NewHandler(ragSvc *rag.Service) *Handler {
	return &Handler{
		ragSvc: ragSvc,
	}
}

// IngestPDF handles PDF file uploads and ingestion
func (h *Handler) IngestPDF(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse the multipart form
	err := r.ParseMultipartForm(10 * 1024 * 1024) // 10MB max
	if err != nil {
		util.SendError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("pdf")
	if err != nil {
		util.SendError(w, http.StatusBadRequest, "No file uploaded")
		return
	}
	defer file.Close()

	// Save the file temporarily
	tempPath := fmt.Sprintf("./uploads/%s", header.Filename)
	
	// Create uploads directory if it doesn't exist
	os.MkdirAll("./uploads", 0755)
	
	tempFile, err := os.Create(tempPath)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, "Failed to write file")
		return
	}

	// Ingest the PDF
	err = h.ragSvc.IngestPDF(ctx, tempPath)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to ingest PDF: %v", err))
		return
	}

	// Clean up
	os.Remove(tempPath)

	util.SendData(w, http.StatusOK, map[string]string{
		"message": "PDF ingested successfully",
		"file":    header.Filename,
	})
}

// GetVersion returns API info
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	util.SendData(w, http.StatusOK, map[string]string{
		"status": "RAG service is running",
	})
}

// AskConcierge handles concierge questions via HTTP
func (h *Handler) AskConcierge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Question string `json:"question"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Question == "" {
		util.SendError(w, http.StatusBadRequest, "Question is required")
		return
	}

	answer, err := h.ragSvc.GetAnswer(ctx, req.Question)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get answer: %v", err))
		return
	}

	util.SendData(w, http.StatusOK, map[string]string{
		"type":   "concierge_response",
		"answer": answer,
	})
}
