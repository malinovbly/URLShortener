package redirect

import (
	"URLShortener/internal/http/response"
	"URLShortener/internal/storage"
	"errors"
	"net/http"
)

type Service interface {
	GetURL(alias string) (string, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Redirect GoDoc
//
//	@Summary		Redirect to original URL
//	@Description	Redirects the client using the specified URL alias.
//	@Tags			URL
//	@Produce		json
//	@Param			alias	path	string	true	"URL alias"
//	@Success		302
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/{alias} [get]
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	alias := r.PathValue("alias")

	originalURL, err := h.service.GetURL(alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			response.WriteJSON(w, http.StatusNotFound, response.ErrorResponse{
				Error: "url not found",
			})
			return
		}

		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{
			Error: "internal server error",
		})
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}
