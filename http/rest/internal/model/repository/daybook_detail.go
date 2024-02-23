package repository

import (
	"context"
	mDaybookDetail "findigitalservice/http/rest/internal/model/daybook_detail"
)

type DaybookDetailRepository interface {
	FindAll(ctx context.Context, query map[string][]string) ([]mDaybookDetail.DaybookDetail, error)
	FindById(ctx context.Context, id string) (mDaybookDetail.DaybookDetail, error)
	Create(ctx context.Context, payload mDaybookDetail.DaybookDetail) (mDaybookDetail.DaybookDetail, error)
	Update(ctx context.Context, payload mDaybookDetail.DaybookDetail) (mDaybookDetail.DaybookDetail, error)
}
