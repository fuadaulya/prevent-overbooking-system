package delivery

import (
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(uc Usecase) *Handler {
	return &Handler{
		usecase: uc,
	}
}

// NewRouter initializes the HTTP router and routes the requests.
func NewRouter(uc Usecase) *httprouter.Router {
	router := httprouter.New()

	handler := NewHandler(uc)

	// Endpoint
	router.POST("/cart/:cartID/item", handler.AddItemToCart)
	router.POST("/order/checkout", handler.CheckoutOrder)

	return router
}
