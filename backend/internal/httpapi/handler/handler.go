package handler

type Handler struct {
	maxBodyBytes int64
}

func New(maxBodyBytes int64) *Handler {
	return &Handler{maxBodyBytes: maxBodyBytes}
}
