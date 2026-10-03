package services

import (
	"context"
	"time"

	"github.com/JoaoNaif/Go-bid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewProductService(pool *pgxpool.Pool) ProductService {
	return ProductService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (ps *ProductService) CreateProduct(
	ctx context.Context,
	sellerID uuid.UUID,
	productName string,
	description string,
	basePrice float64,
	auctionEnd time.Time,
) (uuid.UUID, error) {
	args := pgstore.CreateProductParams{
		SellerID:    sellerID,
		ProductName: productName,
		Description: description,
		Baseprice:   basePrice,
		AuctionEnd:  auctionEnd,
	}

	id, err := ps.queries.CreateProduct(ctx, args)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}
