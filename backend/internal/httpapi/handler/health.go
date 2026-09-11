package handler

import (
	"net/http"

	"github.com/RochaDiego04/calculator/backend/internal/httpapi/respond"
)

type healthResponse struct {
	Status string `json:"status"`
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	respond.JSON(w, http.StatusOK, healthResponse{Status: "ok"})
}
