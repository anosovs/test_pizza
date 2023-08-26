package main

import (
	"fmt"
	"net/http"
	"os"
	"test_pizza/internal/config"
	addItemsOrder "test_pizza/internal/http/handles/orders/addItems"
	createOrder "test_pizza/internal/http/handles/orders/create"
	"test_pizza/internal/http/handles/orders/finishOrder"
	"test_pizza/internal/http/handles/orders/getOrderByID"
	"test_pizza/internal/http/handles/orders/getOrdersByStatus"
	xapikey "test_pizza/internal/http/middleware"
	"test_pizza/internal/storage"
	postge "test_pizza/internal/storage/postgre"
	"test_pizza/internal/storage/ram"

	_ "test_pizza/docs"

	"github.com/go-chi/chi"
	"golang.org/x/exp/slog"
)

// @title Orders API
// @version 1.0
// @description This is a sample serice for managing orders
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name X-Auth-Key
// @description Very secret code, like qwerty123

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
	
	
	
	r := chi.NewRouter()
	r.Post("/orders", createOrder.New(log, storage))
	r.Post("/orders/{order_id}/items", addItemsOrder.New(log, storage))
	r.Get("/orders/{order_id}", getOrderByID.New(log, storage))
	r.With(xapikey.CheckApiKey).Post("/orders/{order_id}/done", finishOrder.New(log, storage))
	r.With(xapikey.CheckApiKey).Get("/orders", getOrdersByStatus.New(log, storage))
	// r.Get("/swagger/*", httpSwagger.Handler(
	// 	httpSwagger.URL("./docs/swagger.json"), 
	// ))
	
	srv := &http.Server{
		Addr: cfg.ServerHost,
		Handler: r,
	}
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err)
		log.Error("failed to start server")
	}

	// TODO gracefulshutdown	

}