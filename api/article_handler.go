package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mcgtrt/xml-to-json-api/store"
	"github.com/mcgtrt/xml-to-json-api/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DEFAULT_PAGE  = 1
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

	article, err := h.store.GetArticleByID(context.Background(), id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, article)
}

func (h *ArticleHandler) HandleGetArticles(w http.ResponseWriter, r *http.Request) error {
	opts, errs := makeFindOptions(r)
	if len(errs) > 0 {
		return WriteJSON(w, http.StatusBadRequest, errs)
	}

	articles, err := h.store.GetArticles(context.Background(), opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNoDocuments()
		}
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

func makeFindOptions(r *http.Request) (*options.FindOptions, map[string]string) {
	var (
		errors      = make(map[string]string)
		pageString  = r.URL.Query().Get("page")
		limitString = r.URL.Query().Get("limit")
		pageInt     = DEFAULT_PAGE
		limitInt    = DEFAULT_LIMIT

		opts = &options.FindOptions{}
	)

	if pageString != "" {
		p, err := strconv.Atoi(pageString)
		if err != nil {
			errors["page"] = "invalid page parameter value"
		} else {
			pageInt = p
		}
	}

	if limitString != "" {
		l, err := strconv.Atoi(limitString)
		if err != nil {
			errors["limit"] = "invalid limit parameter value"
		} else {
			limitInt = l
		}
	}

	if len(errors) > 0 {
		return nil, errors
	}

	opts.SetSkip(int64((pageInt - 1) * limitInt))
	opts.SetLimit(int64(limitInt))
	return opts, nil
}
