package store

import (
	"context"

	"github.com/mcgtrt/xml-to-json-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Ideally we could use more interface methods for CRUD operations
// but including only those required to complete the task.
type ArticleStorer interface {
	GetArticleByID(context.Context, string) (*types.Article, error)
	GetArticles(context.Context, bson.M) ([]*types.Article, error)
	InsertArticle(context.Context, *types.Article) (*types.Article, error)
}

type MongoArticleStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoArticleStore(client *mongo.Client, dbname string) ArticleStorer {
	return &MongoArticleStore{
		client: client,
		coll:   client.Database(dbname).Collection("articles"),
	}
}

func (s *MongoArticleStore) GetArticleByID(ctx context.Context, id string) (*types.Article, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var (
		filter  = bson.M{"_id": oid}
		res     = s.coll.FindOne(ctx, filter)
		article *types.Article
	)

	if err := res.Decode(&article); err != nil {
		return nil, err
	}

	return article, nil
}

func (s *MongoArticleStore) GetArticles(ctx context.Context, filter bson.M) ([]*types.Article, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var articles []*types.Article
	if err := cur.All(ctx, &articles); err != nil {
		return nil, err
	}

	return articles, nil
}

func (s *MongoArticleStore) InsertArticle(ctx context.Context, article *types.Article) (*types.Article, error) {
	res, err := s.coll.InsertOne(ctx, article)
	if err != nil {
		return nil, err
	}

	article.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return article, nil
}
