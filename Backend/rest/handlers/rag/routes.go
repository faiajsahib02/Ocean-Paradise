package raghandler

import (
	"net/http"

	middleware "oasis/backend/rest/middlewares"
)

// RegisterRoutes registers all RAG routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middleware.Manager) {
	mux.HandleFunc("POST /api/rag/ingest", h.IngestPDF)
	mux.HandleFunc("GET /api/rag/health", h.HealthCheck)
	mux.HandleFunc("POST /api/rag/ask", h.AskConcierge)
}
