package post

import (
	"context"
	"example/es_golang/internal/pkg/storage"
	"github.com/google/uuid"
	"time"
)

type service struct {
	storage storage.PostStorer
}

func (s service) create(ctx context.Context, req createRequest) (createResponse, error) {
	id := uuid.New().String()
	cr := time.Now().UTC()

	doc := storage.Post{
		ID:        id,
		Title:     req.Title,
		Text:      req.Text,
		Tags:      req.Tags,
		CreatedAt: &cr,
	}

	if err := s.storage.Insert(ctx, doc); err != nil {
		return createResponse{}, err
	}

	return createResponse{ID: id}, nil
}

func (s service) update(ctx context.Context, req updateRequest) error {
	doc := storage.Post{
		ID:    req.ID,
		Title: req.Title,
		Text:  req.Text,
		Tags:  req.Tags,
	}

	if err := s.storage.Update(ctx, doc); err != nil {
		return err
	}

	return nil
}
