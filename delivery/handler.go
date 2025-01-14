package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) AddItemToCart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse cartID from URL
	cartID, err := strconv.Atoi(ps.ByName("cartID"))
	if err != nil {
		http.Error(w, "invalid cartID", http.StatusBadRequest)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		http.Error(w, "missing user ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Call usecase
	err = h.usecase.AddItemToCart(r.Context(), cartID, userID, req.ProductID, req.Quantity)
	if err != nil {
		if err.Error() == "invalid quantity or insufficient stock" {
			http.Error(w, "Not enough stock available", http.StatusBadRequest)
			return
		}

		http.Error(w, "failed to add item to cart: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Item successfully added to cart"))
}

func (h *Handler) CheckoutOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		http.Error(w, "missing user ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	// Retrieve the cartID from the request body or query
	var req struct {
		CartID int `json:"cart_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate Cart and User Ownership
	if !h.usecase.ValidateCartOwnership(ctx, req.CartID, userID) {
		http.Error(w, "cart does not belong to the user", http.StatusUnauthorized)
		return
	}

	// Call usecase
	err = h.usecase.Checkout(ctx, req.CartID, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to checkout: %+v", err), http.StatusInternalServerError)
		return
	}

	// Successful Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "checkout successful"})
}
