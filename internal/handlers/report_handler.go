package handlers

import (
	"net/http"
	"taskmanager/internal/service"
	"taskmanager/internal/utils"
)

type ReportHandler struct {
	ReportService *service.ReportService
}

func NewReportHandler() *ReportHandler {
	return &ReportHandler{ReportService: service.NewReportService()}
}

func (h *ReportHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")

		return
	}

	ctx := r.Context()

	if err := h.ReportService.GenerateReport(ctx); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())

		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
