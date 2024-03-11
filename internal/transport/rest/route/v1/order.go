package v1

import (
	"applicationDesignTest/internal/transport/rest/controllers"
	"github.com/go-chi/chi/v5"
)

func RegisterOrderRoutes(r chi.Router, ctr *controllers.OrderController) {
	r.Route("/v1/orders", func(r chi.Router) {
		r.Post("/", ctr.CreateOrder)
	})
}
