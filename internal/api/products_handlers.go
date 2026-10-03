package api

import (
	"net/http"

	"github.com/JoaoNaif/Go-bid/internal/use-cases/product"
	"github.com/JoaoNaif/Go-bid/internal/utils"
	"github.com/google/uuid"
)

func (api *Api) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := utils.DecodeValidJson[product.CreateProductReq](r)
	if err != nil {
		_ = utils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	userId, ok := api.Sessions.Get(r.Context(), "AuthenticateUserId").(uuid.UUID)
	if !ok {
		_ = utils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{
			"error": "unexpected error, try to login again",
		})
		return
	}

	id, err := api.ProductService.CreateProduct(r.Context(),
		userId,
		data.ProductName,
		data.Description,
		data.Baseprice,
		data.AuctionEnd,
	)
	if err != nil {
		_ = utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "failed to create product auction, try again later",
		})
		return
	}

	_ = utils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"message":    "product created with success",
		"product_id": id,
	})
}
