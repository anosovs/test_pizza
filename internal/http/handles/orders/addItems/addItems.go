package addItemsOrder

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

type Request struct {
	Items []int `json:"items" validate:"required"`
}


type AddItems interface{
	AddItems(order_id string, items []int) error
}

func New(log *slog.Logger,  i AddItems) http.HandlerFunc {
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
		order_id := chi.URLParam(r, "order_id")
		err = i.AddItems(order_id, req.Items)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		log.Debug("Items updated", slog.String("order_id", order_id))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}