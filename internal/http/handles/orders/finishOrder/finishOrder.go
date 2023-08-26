package finishOrder

import (
	"net/http"

	"github.com/go-chi/chi"
	"golang.org/x/exp/slog"
)


type FinisherOrder interface{
	FinishOrder(order_id string) error
}

// FinishOrder godoc
// @Summary Finish order by ID
// @Description Finish order by ID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order_id   path string true "Order ID"
// @Success 200
// @Failure      500  {string}  string
// @Failure      401  {string}  string
// @Router /orders/{order_id}/done [get]
func New(log *slog.Logger,  i FinisherOrder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request)  {
		order_id := chi.URLParam(r, "order_id")
		err := i.FinishOrder(order_id)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		log.Debug("order done", slog.String("order_id", order_id))
		w.WriteHeader(http.StatusOK)
	}
}