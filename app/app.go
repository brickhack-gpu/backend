package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"gpu/model"
	"gpu/routes"
    "gpu/scan"
)

type App struct {
	Router        *mux.Router
	DB            *bun.DB
	JwtSecret     string
	StripeSecret  string
	StripeWebhook string
	GCPComputeKey string
	DEV           bool
}

func MakeTables(db *bun.DB) error {
	ctx := context.Background()
	_, err := db.NewCreateTable().Model((*model.User)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateTable().Model((*model.Deposit)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateTable().Model((*model.Product)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateTable().Model((*model.Purchase)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateTable().Model((*model.Notification)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateTable().Model((*model.Template)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateTable().Model((*model.ServerConfig)(nil)).IfNotExists().Exec(ctx)
	return err
}

func (a *App) Initialize(user, password, dbname, jwtSecret, stripeSecret, stripeWebhook, gcpComputeKey string, dev bool) {
	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", user, password, dbname)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connectionString)))

	a.DB = bun.NewDB(sqldb, pgdialect.New())
	err := MakeTables(a.DB)
	if err != nil {
		log.Fatal(err)
	}

    go scan.ScanBalance(a.DB)

	a.Router = mux.NewRouter()
	a.JwtSecret = jwtSecret
	a.StripeSecret = stripeSecret
	a.StripeWebhook = stripeWebhook
	a.GCPComputeKey = gcpComputeKey
	a.DEV = dev

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func (a *App) initializeRoutes() {
	cor := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://test.local", "http://localhost:5173"},
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		Debug: true,
		AllowedHeaders: []string{
			"*",
		},
	}).Handler

	router := routes.NewRouter(a.DB, a.JwtSecret, a.StripeSecret, a.StripeWebhook, a.GCPComputeKey, a.DEV)

	a.Router.Handle("/register", cor(http.HandlerFunc(router.Register))).Methods("OPTIONS", "POST")
	a.Router.Handle("/login", cor(http.HandlerFunc(router.Login))).Methods("OPTIONS", "POST")
	a.Router.Handle("/refresh", cor(http.HandlerFunc(router.RefreshToken))).Methods("OPTIONS", "POST")
	a.Router.Handle("/profile", cor(router.AuthMiddleware(http.HandlerFunc(router.Profile)))).Methods("OPTIONS", "GET")
	a.Router.Handle("/transactions", cor(router.AuthMiddleware(http.HandlerFunc(router.Transactions)))).Methods("OPTIONS", "GET")
	a.Router.Handle("/products", cor(router.AuthMiddleware(http.HandlerFunc(router.Products)))).Methods("OPTIONS", "GET")
	a.Router.Handle("/templates", cor(router.AuthMiddleware(http.HandlerFunc(router.Templates)))).Methods("OPTIONS", "GET")
	a.Router.Handle("/search", cor(router.AuthMiddleware(http.HandlerFunc(router.Search)))).Methods("OPTIONS", "GET")
	a.Router.Handle("/data", cor(router.AuthMiddleware(http.HandlerFunc(router.Data)))).Methods("OPTIONS", "GET")

	a.Router.Handle("/servers/start", cor(router.AuthMiddleware(http.HandlerFunc(router.SpinServer)))).Methods("OPTIONS", "POST")
	a.Router.Handle("/servers/destroy", cor(router.AuthMiddleware(http.HandlerFunc(router.KillServer)))).Methods("OPTIONS", "POST")
}
