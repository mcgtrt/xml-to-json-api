package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mcgtrt/xml-to-json-api/api"
	"github.com/mcgtrt/xml-to-json-api/store"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburi  = "mongodb://localhost:27017"
	dbname = "xmlToJsonApi"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dburi))
	if err != nil {
		panic(err)
	}

	var (
		articleStore   = store.NewMongoArticleStore(client, dbname)
		articleHandler = api.NewArticleHandler(articleStore)

		listenAddr = flag.String("listenAddr", ":3000", "api gateway http listen address")
		r          = mux.NewRouter()
	)
	flag.Parse()

	{
		r.HandleFunc("/article/{id}", makeHandlerFunc(articleHandler.HandleGetArticle)).Methods("GET")
		r.HandleFunc("/article", makeHandlerFunc(articleHandler.HandleGetArticles)).Methods("GET")
		r.HandleFunc("/article", makeHandlerFunc(articleHandler.HandlePostArticle)).Methods("POST")

		http.Handle("/", r)
	}

	logrus.Infof("Starting HTTP server at port %s", *listenAddr)
	http.ListenAndServe(*listenAddr, nil)
}

func makeHandlerFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			api.WriteJSON(w, http.StatusBadRequest, map[string]string{"er: ": err.Error()})
		}
	}
}
