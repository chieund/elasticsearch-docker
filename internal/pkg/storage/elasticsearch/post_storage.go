package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"example/es_golang/internal/pkg/domain"
	"example/es_golang/internal/pkg/storage"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"time"
)

var _ storage.PostStorer = &PostStorage{}

type PostStorage struct {
	elastic ElasticSearch
	timeout time.Duration
}

func NewPostStorage(elastic ElasticSearch) (PostStorage, error) {
	return PostStorage{
		elastic: elastic,
		timeout: time.Second * 10,
	}, nil
}

func (p *PostStorage) Insert(ctx context.Context, post storage.Post) error {
	bdy, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("insert: marshall: %w", err)
	}

	req := esapi.CreateRequest{
		Index:      p.elastic.alias,
		DocumentID: post.ID,
		Body:       bytes.NewReader(bdy),
	}
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := req.Do(ctx, p.elastic.client)
	if err != nil {
		return fmt.Errorf("insert: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 409 {
		return domain.ErrConflict
	}
	if res.IsError() {
		return fmt.Errorf("insert: response: %s", res.String())
	}

	return nil
}
