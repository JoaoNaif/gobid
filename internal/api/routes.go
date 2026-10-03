package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *Api) BindRoutes() {
	api.Router.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger, api.Sessions.LoadAndSave)

	// csrfMiddleware := csrf.Protect(
	//	 []byte(os.Getenv("GOBID_CSRF_KEY")),
	//	 csrf.Secure(false), // DEV ONLY
	//	 csrf.Path("/"),     // 1 único cookie _gorilla_csrf para toda a API
	// )

	// DEV ONLY: marca a request como HTTP puro para o gorilla/csrf (v1.7.3+)
	// não exigir header Origin/Referer, permitindo testar via Postman.
	// api.Router.Use(func(next http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		next.ServeHTTP(w, csrf.PlaintextHTTPRequest(r))
	// 	})
	// })

	// api.Router.Use(csrfMiddleware)

	api.Router.Route("/api", func(r chi.Router) {

		r.Route("/v1", func(r chi.Router) {
			//r.Get("/csrftoken", api.HandleGetCSRFtoken)
			r.Route("/users", func(r chi.Router) {
				r.Post("/signup", api.handleSignupUser)
				r.Post("/login", api.handleLoginUser)
				r.With(api.AuthMiddleware).Post("/logout", api.handleLogoutUser)
			})

			r.Route("/products", func(r chi.Router) {
				r.Use(api.AuthMiddleware)
				r.Post("/", api.handleCreateProduct)

				r.Route("/{product_id}/bids", func(r chi.Router) {
					r.Post("/", api.handleCreateBid)
					r.Get("/", api.handleGetBidsByProductId)
					r.Get("/highest", api.handleGetHighestBidByProductId)
				})
			})
		})
	})
}
