package finishOrder

import (
	"net/http"

	"github.com/go-chi/chi"
	"golang.org/x/exp/slog"
)


type FinisherOrder interface{
	FinishOrder(order_id string) error
}

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