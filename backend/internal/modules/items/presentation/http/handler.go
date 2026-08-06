package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/ports"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/httpx"
)

const (
	ProductIdMuxPattern = "productId"
)

func GetItems(usecase ports.GetAllItemsUsecase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		items, err := usecase.GetAllItems(request.Context())
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		httpx.EncodeJSON(writer, items, http.StatusOK)
	}
}

func GetItem(usecase ports.GetItemUsecase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := uuid.Parse(chi.URLParam(request, ProductIdMuxPattern))
		if err != nil {
			http.Error(writer, "invalid item id", http.StatusBadRequest)
			return
		}

		item, err := usecase.GetItem(request.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrItemNotFound):
				http.Error(writer, "item not found", http.StatusNotFound)
				return
			default:
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		
		httpx.EncodeJSON(writer, item, http.StatusOK)
	}
}
