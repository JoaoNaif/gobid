package bid

import (
	"context"

	"github.com/JoaoNaif/Go-bid/internal/validator"
)

// O product_id vem da URL e o bidder_id da sessão do usuário autenticado.
type CreateBidReq struct {
	BidAmount float64 `json:"bid_amount"`
}

func (req CreateBidReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(req.BidAmount > 0, "bid_amount", "This field must be a positive number.")

	return eval
}
