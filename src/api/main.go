package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	handler "github.com/marcosmariano/go-clean-arch-api/src/api/handlers"
	"github.com/marcosmariano/go-clean-arch-api/src/api/middleware"
	"github.com/marcosmariano/go-clean-arch-api/src/config"
	"github.com/marcosmariano/go-clean-arch-api/src/infraestructure/repository"
	"github.com/marcosmariano/go-clean-arch-api/src/pkg/metric"
	"github.com/marcosmariano/go-clean-arch-api/src/usecase/book"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func main() {

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_DATABASE)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	bookRepo := repository.NewBookMySQL(db)
	bookService := book.NewService(bookRepo)

	metricService, err := metric.NewPrometheusService()
	if err != nil {
		log.Fatal(err.Error())
	}
	r := mux.NewRouter()
	//handlers
	n := negroni.New(
		negroni.HandlerFunc(middleware.Cors),
		negroni.HandlerFunc(middleware.Metrics(metricService)),
		negroni.NewLogger(),
	)
	//book
	handler.MakeBookHandlers(r, *n, bookService)

	http.Handle("/", r)
	http.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + strconv.Itoa(config.API_PORT),
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
