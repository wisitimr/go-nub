package response

import "github.com/xuri/excelize/v2"

type ExcelFile struct {
	File *excelize.File `json:"file"`
	Name string         `json:"name"`
}
