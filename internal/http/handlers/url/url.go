package url

import (
	"URLShortener/internal/http/response"
	"URLShortener/internal/service"
	"encoding/json"
	"errors"
	"net/http"
)

type Service interface {
	SaveURL(originalURL string) (string, error)
	GetAllURLs() (map[string]string, error)
}

type CreateURLRequest struct {
	URL string `json:"url"`
}

type CreateURLResponse struct {
	Alias string `json:"alias"`
}

type AllURLsResponse struct {
	URLs map[string]string `json:"urls"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// AllURLs GoDoc
//
//	@Summary		Get all URLs
//	@Description	Returns all saved URLs.
//	@Tags			URL
//	@Produce		json
//	@Success		200	{object}	AllURLsResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/urls [get]
func (h *Handler) AllURLs(w http.ResponseWriter, _ *http.Request) {
	urls, err := h.service.GetAllURLs()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{
			Error: response.ErrorInternalServer,
		})
		return
	}

	response.WriteJSON(w, http.StatusOK, AllURLsResponse{
		URLs: urls,
	})
}

// CreateURL GoDoc
//
//	@Summary		Create short URL
//	@Description	Creates an alias for the specified URL.
//	@Tags			URL
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateURLRequest	true	"Original URL"
//	@Success		201		{object}	CreateURLResponse
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/url [post]
func (h *Handler) CreateURL(w http.ResponseWriter, r *http.Request) {
	var request CreateURLRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse{
			Error: response.ErrorInvalidRequestBody,
		})
		return
	}

	var alias string
	alias, err = h.service.SaveURL(request.URL)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) {
			response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{
			Error: response.ErrorInternalServer,
		})
		return
	}

	response.WriteJSON(w, http.StatusCreated, CreateURLResponse{
		Alias: alias,
	})
}
