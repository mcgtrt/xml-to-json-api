package producer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mcgtrt/xml-to-json-api/store"
	"github.com/mcgtrt/xml-to-json-api/types"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Producer struct {
	store store.ArticleStorer

	pollURL       string
	frequencyMins int

	isRunning bool
}

func NewProducer(store store.ArticleStorer, url string, freq int) *Producer {
	return &Producer{
		store:         store,
		pollURL:       url,
		frequencyMins: freq,
		isRunning:     false,
	}
}

func (p *Producer) Start() {
	p.isRunning = true
	logrus.Info("Data producing started")
	p.loopProduce()
}

func (p *Producer) Stop() {
	p.isRunning = false
}

func (p *Producer) loopProduce() {
	for p.isRunning {
		res, err := http.Get(p.pollURL)
		if err != nil {
			logrus.Errorf("Producer error: %s", err)
		}

		var info *types.NewListInformation
		if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
			logrus.Errorf("Producer parsing error: %s", err)
		}

		var (
			ctx  = context.Background()
			list = info.ArticleList.Values
		)
		for _, item := range list {
			article, err := p.store.GetArticleByNewsArticleID(ctx, fmt.Sprintf("%d", item.NewsArticleID))
			if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
				logrus.Errorf("Producer store error: %s", &err)
			} else {
				insertedArticle, err := p.store.InsertArticle(ctx, article)
				if err != nil {
					logrus.Errorf("Producer store insert error: %s", err)
				} else {
					logrus.Infof("Successfully inserted new article with ID: %s", insertedArticle.ID)
				}
			}
		}

		// This should be set to time.Minute on production but for demonstration purposes
		// is set to seconds to prove it's working
		time.Sleep(time.Second * time.Duration(p.frequencyMins))
	}
}
