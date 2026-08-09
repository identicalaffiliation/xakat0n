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
	ItemIdMuxPattern = "itemId"
)

func GetItems(usecase ports.GetAllItemsUsecase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		items, err := usecase.GetAllItems(request.Context())
		if err != nil {
			httpx.WriteError(writer, http.StatusInternalServerError, "internal_error", "internal server error")
			return
		}

		httpx.EncodeJSON(writer, items, http.StatusOK)
	}
}

func GetItem(usecase ports.GetItemUsecase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := uuid.Parse(chi.URLParam(request, ItemIdMuxPattern))
		if err != nil {
			httpx.WriteError(writer, http.StatusBadRequest, "bad_request", "invalid item id")
			return
		}

		item, err := usecase.GetItem(request.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrItemNotFound):
				httpx.WriteError(writer, http.StatusNotFound, "item_not_found", "item not found")
			default:
				httpx.WriteError(writer, http.StatusInternalServerError, "internal_error", "internal server error")
			}
			return
		}

		httpx.EncodeJSON(writer, item, http.StatusOK)
	}
}
