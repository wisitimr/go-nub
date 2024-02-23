package service

import (
	"context"
	mDaybook "findigitalservice/http/rest/internal/model/daybook"
	mRes "findigitalservice/http/rest/internal/model/response"
)

type DaybookService interface {
	Count(ctx context.Context, query map[string][]string) (mRes.CountDto, error)
	FindAll(ctx context.Context, query map[string][]string) ([]mDaybook.DaybookList, error)
	FindById(ctx context.Context, id string) (mDaybook.DaybookResponse, error)
	GenerateExcel(ctx context.Context, id string) (mRes.ExcelFile, error)
	Create(ctx context.Context, payload mDaybook.DaybookPayload) (mDaybook.DaybookPayload, error)
	Update(ctx context.Context, id string, payload mDaybook.Daybook) (mDaybook.Daybook, error)
}
