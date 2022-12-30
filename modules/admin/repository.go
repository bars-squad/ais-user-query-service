package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bars-squad/ais-user-query-service/entity"
	"github.com/bars-squad/ais-user-query-service/exception"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	Save(ctx context.Context, admin entity.Admin) (err error)
}

type elasticsearchRequest interface {
	Do(ctx context.Context, transport esapi.Transport) (resp *esapi.Response, err error)
}

type RepositoryImpl struct {
	Logger  *logrus.Logger
	Client  *elasticsearch.Client
	Index   string
	DocType string
}

func NewRepository(logger *logrus.Logger, client *elasticsearch.Client) Repository {
	return &RepositoryImpl{
		Logger:  logger,
		Client:  client,
		Index:   "administrator",
		DocType: "_doc",
	}
}

func (r *RepositoryImpl) do(ctx context.Context, req elasticsearchRequest) (responseBodyBuff []byte, err error) {
	var resp *esapi.Response
	if resp, err = req.Do(ctx, r.Client); err != nil {
		r.Logger.Error(err)
		err = exception.ErrInternalServer
		return
	}

	defer resp.Body.Close()

	if resp.IsError() {
		if resp.StatusCode != http.StatusNotFound {
			r.Logger.Error(resp.String())
			err = exception.ErrInternalServer
			return
		}
		err = exception.ErrNotFound
		return
	}

	responseBodyBuff, _ = ioutil.ReadAll(resp.Body)
	return
}

func (r *RepositoryImpl) Save(ctx context.Context, admin entity.Admin) (err error) {

	tpAdminBuff := new(bytes.Buffer)
	json.NewEncoder(tpAdminBuff).Encode(&admin)

	req := esapi.IndexRequest{
		Index:        r.Index,
		DocumentType: r.DocType,
		DocumentID:   admin.ID,
		Body:         tpAdminBuff,
	}

	_, err = r.do(ctx, &req)

	return
}
