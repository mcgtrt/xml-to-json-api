package producer

import (
	"context"
	"encoding/xml"
	"errors"
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

		var info types.NewListInformation
		if err := xml.NewDecoder(res.Body).Decode(&info); err != nil {
			logrus.Errorf("Producer parsing error: %s", err)
		}

		var (
			ctx             = context.Background()
			list            = info.ArticleList.Values
			newArticleCount = 0
		)
		for _, item := range list {
			var (
				_, err    = p.store.GetArticleByNewsArticleID(ctx, item.NewsArticleID)
				noDocCond = errors.Is(err, mongo.ErrNoDocuments)
			)

			if err != nil && !noDocCond {
				logrus.Errorf("Producer store error: %s", err)
			} else if err != nil && noDocCond {
				_, err := p.store.InsertArticle(ctx, &item)
				if err != nil {
					logrus.Errorf("Producer store insert error: %s", err)
				} else {
					newArticleCount++
				}
			}
		}

		logrus.WithFields(logrus.Fields{
			"from":             p.pollURL,
			"newArticlesCount": newArticleCount,
		}).Info("Producer data poll")

		// This should be set to time.Minute on production but for demonstration purposes
		// is set to seconds to prove it's working
		time.Sleep(time.Second * time.Duration(p.frequencyMins))
	}
}
