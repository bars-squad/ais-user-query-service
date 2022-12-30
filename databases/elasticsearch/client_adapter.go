package es

import (
	"context"

	"github.com/bars-squad/ais-user-query-service/config"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
)

type Client interface {
	Connect(ctx context.Context, logger *logrus.Logger) (client *elasticsearch.Client, err error)
	// Disconnect(ctx context.Context) (err error)
}

type ClientAdapter struct {
	Config *config.Config
}

func NewClientAdapter() Client {
	return &ClientAdapter{
		Config: &config.Config{},
	}
}

func (ca *ClientAdapter) Connect(ctx context.Context, logger *logrus.Logger) (client *elasticsearch.Client, err error) {
	es, err := elasticsearch.NewClient(ca.Config.Elasticsearch)
	if err != nil {
		logger.Fatal(err)
		return nil, err
	}

	res, err := es.Info()
	if err != nil {
		logger.WithField("errors", err.Error()).Fatal("Elasticsearch disconnected")
		return nil, err
	}

	logger.Info("Elasticsearch connected")
	defer res.Body.Close()
	return es, nil
}
