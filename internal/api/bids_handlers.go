package api

import (
	"errors"
	"net/http"

	"github.com/JoaoNaif/Go-bid/internal/services"
	"github.com/JoaoNaif/Go-bid/internal/use-cases/bid"
	"github.com/JoaoNaif/Go-bid/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *Api) handleCreateBid(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(chi.URLParam(r, "product_id"))
	if err != nil {
		_ = utils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "invalid product id",
		})
		return
	}

	data, problems, err := utils.DecodeValidJson[bid.CreateBidReq](r)
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

	createdBid, err := api.BidService.CreateBid(r.Context(), productID, userId, data.BidAmount)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrProductNotFound):
			_ = utils.EncodeJson(w, r, http.StatusNotFound, map[string]any{"error": err.Error()})
		case errors.Is(err, services.ErrAuctionEnded),
			errors.Is(err, services.ErrBidOnOwnProduct),
			errors.Is(err, services.ErrBidTooLow):
			_ = utils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{"error": err.Error()})
		default:
			_ = utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
				"error": "unexpected internal server error",
			})
		}
		return
	}

	_ = utils.EncodeJson(w, r, http.StatusCreated, createdBid)
}

func (api *Api) handleGetBidsByProductId(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(chi.URLParam(r, "product_id"))
	if err != nil {
		_ = utils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "invalid product id",
		})
		return
	}

	bids, err := api.BidService.GetBidsByProductId(r.Context(), productID)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			_ = utils.EncodeJson(w, r, http.StatusNotFound, map[string]any{"error": err.Error()})
			return
		}

		_ = utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	_ = utils.EncodeJson(w, r, http.StatusOK, bids)
}

func (api *Api) handleGetHighestBidByProductId(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(chi.URLParam(r, "product_id"))
	if err != nil {
		_ = utils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "invalid product id",
		})
		return
	}

	highestBid, err := api.BidService.GetHighestBidByProductId(r.Context(), productID)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) || errors.Is(err, services.ErrNoBidsForProduct) {
			_ = utils.EncodeJson(w, r, http.StatusNotFound, map[string]any{"error": err.Error()})
			return
		}

		_ = utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	_ = utils.EncodeJson(w, r, http.StatusOK, highestBid)
}
