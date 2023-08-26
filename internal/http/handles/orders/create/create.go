package createOrder

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"test_pizza/internal/models"

	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

type Request struct {
	Items []int `json:"items" validate:"required"`
}


type CreateOrder interface{
	CreateOrder(items []int) (models.Order, error)
}
// CreateOrder godoc
// @Summary Insert new order
// @Description Add new order 
// @Tags orders
// @Accept  json
// @Produce  json
// @Param items body []int true "Create order"
// @Success 200 {object} models.Order
// @Failure      400  {string}  string
// @Failure      500  {string}  string
// @Router /orders [post]
func New(log *slog.Logger,  i CreateOrder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request)  {
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		
		defer r.Body.Close()
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			http.Error(w, "request body is empty", 400)
			return 
		}
		if err != nil {
			log.Error("failed to decode request. error: " + err.Error())
			http.Error(w, "failed to decode request", 400)
			return
		}

		validate := validator.New()
		err = validate.Struct(req)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "items is required", 400)
			return
		}

		order, err := i.CreateOrder(req.Items)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		log.Debug("order added", slog.String("id", order.Order_id))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)

	}
}