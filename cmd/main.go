package main

import (
	"net/http"
	"os"
	"test_pizza/internal/config"
	addItemsOrder "test_pizza/internal/http/handles/orders/addItems"
	createOrder "test_pizza/internal/http/handles/orders/create"
	"test_pizza/internal/http/handles/orders/finishOrder"
	"test_pizza/internal/http/handles/orders/getOrderByID"
	"test_pizza/internal/http/handles/orders/getOrdersByStatus"
	"test_pizza/internal/storage"
	postge "test_pizza/internal/storage/postgre"
	"test_pizza/internal/storage/ram"

	"github.com/go-chi/chi"
	"golang.org/x/exp/slog"
)



func main() {
	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}


	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("starting service.", slog.String("env", cfg.Env))


	var storage storage.Storage	
	if cfg.Env == config.EnvProd {
		storage, err = postge.New(cfg.DBDSN)
		if err != nil {
			panic(err)
		}
	}  else {
		storage = ram.New()
	}
	
	
	//TODO chi
	
	r := chi.NewRouter()
	r.Post("/orders", createOrder.New(log, storage))
	r.Post("/orders/{order_id}/items", addItemsOrder.New(log, storage))
	r.Get("/orders/{order_id}", getOrderByID.New(log, storage))
	r.Post("/orders/{order_id}/done", finishOrder.New(log, storage))
	r.Get("/orders", getOrdersByStatus.New(log, storage))

	//TODO run server
	
	srv := &http.Server{
		Addr: cfg.ServerHost,
		Handler: r,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	// TODO gracefulshutdown
	// RM

	

}