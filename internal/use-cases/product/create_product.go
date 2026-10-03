package product

import (
	"context"
	"time"

	"github.com/JoaoNaif/Go-bid/internal/validator"
)

// O seller_id não vem no body: é obtido da sessão do usuário autenticado.
type CreateProductReq struct {
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	Baseprice   float64   `json:"baseprice"`
	AuctionEnd  time.Time `json:"auction_end"`
}

const minAuctionDuration = 2 * time.Hour

func (req CreateProductReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.ProductName), "product_name", "This field cannot be empty.")
	eval.CheckField(validator.MaxChars(req.ProductName, 255), "product_name", "This field must have at most 255 chars.")

	eval.CheckField(validator.NotBlank(req.Description), "description", "This field cannot be empty.")

	eval.CheckField(req.Baseprice > 0, "baseprice", "This field must be a positive number.")

	eval.CheckField(req.AuctionEnd.After(time.Now().Add(minAuctionDuration)), "auction_end", "Auction end time must be at least 2 hours from now.")

	return eval
}
