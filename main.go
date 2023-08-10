package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mcgtrt/xml-to-json-api/api"
	"github.com/mcgtrt/xml-to-json-api/producer"
	"github.com/mcgtrt/xml-to-json-api/store"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dburi          string
	dbname         string
	httpListenAddr string
	fetchuri       = "https://www.htafc.com/api/incrowd/getnewlistinformation?count=50"
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

		producer = producer.NewProducer(articleStore, fetchuri, 5)
		r        = mux.NewRouter()
	)

	{
		r.HandleFunc("/article/{id}", makeHandlerFunc(articleHandler.HandleGetArticle)).Methods("GET")
		r.HandleFunc("/article", makeHandlerFunc(articleHandler.HandleGetArticles)).Methods("GET")
		r.HandleFunc("/article", makeHandlerFunc(articleHandler.HandlePostArticle)).Methods("POST")

		http.Handle("/", r)
	}

	go func() {
		producer.Start()
	}()

	logrus.Infof("Starting HTTP server at port %s", httpListenAddr)
	http.ListenAndServe(httpListenAddr, nil)
}

func makeHandlerFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			apiErr, ok := err.(api.Error)
			if ok {
				api.WriteJSON(w, apiErr.Status, apiErr)
				return
			}
			apiErr = api.ErrInternalServerError()
			api.WriteJSON(w, apiErr.Status, apiErr)
		}
	}
}

func init() {
	dburi = os.Getenv("MONGO_DB_URI")
	dbname = os.Getenv("MONGO_DB_NAME")
	httpListenAddr = os.Getenv("HTTP_LISTEN_ADDR")
}
