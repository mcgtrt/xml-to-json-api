package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mcgtrt/xml-to-json-api/store"
	"github.com/mcgtrt/xml-to-json-api/types"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	DEFAULT_PAGE  = 0
	DEFAULT_LIMIT = 50
)

type ArticleHandler struct {
	store store.ArticleStorer
}

func NewArticleHandler(store store.ArticleStorer) *ArticleHandler {
	return &ArticleHandler{
		store: store,
	}
}

func (h *ArticleHandler) HandleGetArticle(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	fmt.Printf("Received query id: %s\n", id)
	article, err := h.store.GetArticleByID(context.Background(), id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, article)
}

func (h *ArticleHandler) HandleGetArticles(w http.ResponseWriter, r *http.Request) error {
	var (
		ctx         = context.Background()
		pageString  = r.URL.Query().Get("page")
		limitString = r.URL.Query().Get("limit")
		pageInt     = DEFAULT_PAGE
		limitInt    = DEFAULT_LIMIT
	)

	if pageString != "" {
		p, err := strconv.Atoi(pageString)
		if err != nil {
			return fmt.Errorf("bad request")
		}
		pageInt = p
	}

	if limitString != "" {
		l, err := strconv.Atoi(limitString)
		if err != nil {
			return fmt.Errorf("bad request")
		}
		limitInt = l
	}

	articles, err := h.store.GetArticles(ctx, bson.M{
		"$skip":  pageInt * limitInt,
		"$limit": limitInt,
	})
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, articles)
}

func (h *ArticleHandler) HandlePostArticle(w http.ResponseWriter, r *http.Request) error {
	var article *types.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		return err
	}

	insertedArticle, err := h.store.InsertArticle(context.Background(), article)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, insertedArticle)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
