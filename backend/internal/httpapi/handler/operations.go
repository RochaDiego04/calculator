package handler

import (
	"net/http"

	"github.com/RochaDiego04/calculator/backend/internal/calc"
	"github.com/RochaDiego04/calculator/backend/internal/httpapi/respond"
)

func (h *Handler) Operations(w http.ResponseWriter, r *http.Request) {
	respond.JSON(w, http.StatusOK, calc.Operations())
}
