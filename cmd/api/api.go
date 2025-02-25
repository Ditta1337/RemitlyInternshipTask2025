package main

import (
	"fmt"
	docsPkg "github.com/Ditta1337/RemitlyInternshipTask2025/docs"
	"github.com/Ditta1337/RemitlyInternshipTask2025/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

type config struct {
	addr       string
	apiURL     string
	db         dbConfig
	env        string
	apiVersion string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Timeout(60 * time.Second))

	baseRoute := "/" + app.config.apiVersion
	r.Route(baseRoute, func(r chi.Router) {

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))
		app.logger.Infof("docsURL: %v", docsURL)
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// docs
	docsPkg.SwaggerInfo.Version = version
	docsPkg.SwaggerInfo.BasePath = "/" + app.config.apiVersion

	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infof("server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}
