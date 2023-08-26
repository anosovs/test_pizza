package getOrdersByStatus

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test_pizza/internal/models"

	"golang.org/x/exp/slog"
)

type GetOrder interface{
	GetOrdersByStatus(done int)  ([]models.Order, error)
}

// GetOrdersByStatus godoc
// @Summary Get orders
// @Description Get orders
// @Tags orders
// @Accept  json
// @Produce  json
// @Param done query int false "Filter by status"
// @Success 200 {array} []models.Order
// @Failure      401  {string}  string
// @Failure      500  {string}  string
// @Router /orders [get]
func New(log *slog.Logger,  i GetOrder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request)  {
		status := r.URL.Query().Get("done")
		if !(status=="1" || status=="0") {
			status = "-1"
		}
		statusInt, err := strconv.Atoi(status)
		if err!=nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		orders, err := i.GetOrdersByStatus(statusInt)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		log.Debug("orders displayed by status", slog.String("status", strconv.Itoa(statusInt)))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(orders)

	}
}