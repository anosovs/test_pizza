package getOrderByID

import (
	"encoding/json"
	"net/http"
	"test_pizza/internal/models"

	"github.com/go-chi/chi"
	"golang.org/x/exp/slog"
)

type Response struct {
	models.Order
}

type GetOrder interface{
	GetOrderById(order_id string) (models.Order, error)
}

func New(log *slog.Logger,  i GetOrder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request)  {
		order_id := chi.URLParam(r, "order_id")
		order, err := i.GetOrderById(order_id)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		log.Debug("order displayed", slog.String("id", order.Order_id))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(order)

	}
}