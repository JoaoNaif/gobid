package api

import (
	"github.com/JoaoNaif/Go-bid/internal/services"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

type Api struct {
	Router         *chi.Mux
	UserService    services.UserService
	ProductService services.ProductService
	BidService     services.BidService
	Sessions       *scs.SessionManager
}
