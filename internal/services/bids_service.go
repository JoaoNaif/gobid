package services

import (
	"context"
	"errors"
	"time"

	"github.com/JoaoNaif/Go-bid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrAuctionEnded     = errors.New("auction has already ended")
	ErrBidOnOwnProduct  = errors.New("seller cannot bid on own product")
	ErrBidTooLow        = errors.New("bid amount is too low")
	ErrNoBidsForProduct = errors.New("no bids for this product")
)

type BidService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewBidService(pool *pgxpool.Pool) BidService {
	return BidService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (bs *BidService) CreateBid(ctx context.Context, productID, bidderID uuid.UUID, amount float64) (pgstore.Bid, error) {
	tx, err := bs.pool.Begin(ctx)
	if err != nil {
		return pgstore.Bid{}, err
	}
	defer tx.Rollback(ctx)

	qtx := bs.queries.WithTx(tx)

	// Trava a linha do produto para que lances simultâneos sejam processados em ordem.
	product, err := qtx.GetProductByIdForUpdate(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, ErrProductNotFound
		}
		return pgstore.Bid{}, err
	}

	if product.IsSold || time.Now().After(product.AuctionEnd) {
		return pgstore.Bid{}, ErrAuctionEnded
	}

	if product.SellerID == bidderID {
		return pgstore.Bid{}, ErrBidOnOwnProduct
	}

	highestBid, err := qtx.GetHighestBidByProductId(ctx, productID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return pgstore.Bid{}, err
	}

	// Sem lances: precisa ser >= preço base. Com lances: precisa superar o maior.
	if errors.Is(err, pgx.ErrNoRows) {
		if amount < product.Baseprice {
			return pgstore.Bid{}, ErrBidTooLow
		}
	} else if amount <= highestBid.BidAmount {
		return pgstore.Bid{}, ErrBidTooLow
	}

	bid, err := qtx.CreateBid(ctx, pgstore.CreateBidParams{
		ProductID: productID,
		BidderID:  bidderID,
		BidAmount: amount,
	})
	if err != nil {
		return pgstore.Bid{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return pgstore.Bid{}, err
	}

	return bid, nil
}

func (bs *BidService) GetBidsByProductId(ctx context.Context, productID uuid.UUID) ([]pgstore.Bid, error) {
	if _, err := bs.queries.GetProductById(ctx, productID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	bids, err := bs.queries.GetBidsByProductId(ctx, productID)
	if err != nil {
		return nil, err
	}

	return bids, nil
}

func (bs *BidService) GetHighestBidByProductId(ctx context.Context, productID uuid.UUID) (pgstore.Bid, error) {
	if _, err := bs.queries.GetProductById(ctx, productID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, ErrProductNotFound
		}
		return pgstore.Bid{}, err
	}

	bid, err := bs.queries.GetHighestBidByProductId(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, ErrNoBidsForProduct
		}
		return pgstore.Bid{}, err
	}

	return bid, nil
}
