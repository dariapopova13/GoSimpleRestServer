package main

import (
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Route struct {
	Name        string
	Path        string
	Method      string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

const (
	address      = "127.0.0.1:8080"
	writeTimeout = 15 * time.Second
	readTimeout  = 15 * time.Second
)

func main() {
	routes := Routes{
		Route{
			Name:        "AllBooks",
			Path:        "/books",
			Method:      http.MethodGet,
			HandlerFunc: AllBooksHandler,
		},
		Route{
			Name:        "BookById",
			Path:        "/books/{id}",
			Method:      http.MethodGet,
			HandlerFunc: SingleBookHandler,
		},
		Route{
			Name:        "InsertBook",
			Path:        "/books",
			Method:      http.MethodPost,
			HandlerFunc: InsertBookHandler,
		},
		Route{
			Name:        "UpdateBook",
			Path:        "/books",
			Method:      http.MethodPut,
			HandlerFunc: UpdateBookHandler,
		},
		Route{
			Name:        "DeleteBook",
			Path:        "/books/{id}",
			Method:      http.MethodDelete,
			HandlerFunc: DeleteBookHandler,
		},
	}
	router := mux.NewRouter()
	// add all routes
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	server := initServer(router)
	logger.Fatal(server.ListenAndServe())
}

func initServer(r *mux.Router) (server *http.Server) {
	server = &http.Server{
		Handler:           r,
		Addr:              address,
		WriteTimeout:      writeTimeout,
		ReadHeaderTimeout: readTimeout,
	}
	return
}

func CheckError(e error, message string) {
	if e != nil {
		logger.Error(message, e)
		panic(e)
	}
}
