package service

import (
	"context"
	mDocument "findigitalservice/http/rest/internal/model/document"
)

type DocumentService interface {
	FindAll(ctx context.Context, query map[string][]string) ([]mDocument.Document, error)
	FindById(ctx context.Context, id string) (mDocument.Document, error)
	Create(ctx context.Context, payload mDocument.Document) (mDocument.Document, error)
	Update(ctx context.Context, id string, payload mDocument.Document) (mDocument.Document, error)
}
