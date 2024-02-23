package service

import (
	"context"
	"fmt"
	"saved/http/rest/internal/auth"
	mDaybook "saved/http/rest/internal/model/daybook"
	mDaybookDetail "saved/http/rest/internal/model/daybook_detail"
	mRepo "saved/http/rest/internal/model/repository"
	mRes "saved/http/rest/internal/model/response"
	mService "saved/http/rest/internal/model/service"
	mUser "saved/http/rest/internal/model/user"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type daybookService struct {
	daybookRepo       mRepo.DaybookRepository
	daybookDetailRepo mRepo.DaybookDetailRepository
	logger            *logrus.Logger
}

func InitDaybookService(daybookRepo mRepo.DaybookRepository, daybookDetailRepo mRepo.DaybookDetailRepository, logger *logrus.Logger) mService.DaybookService {
	return &daybookService{
		daybookRepo:       daybookRepo,
		daybookDetailRepo: daybookDetailRepo,
		logger:            logger,
	}
}

func (s daybookService) Count(ctx context.Context, query map[string][]string) (mRes.CountDto, error) {
	count, err := s.daybookRepo.Count(ctx, query)
	if err != nil {
		return mRes.CountDto{Count: 0}, err
	}
	return mRes.CountDto{Count: count}, nil
}

func (s daybookService) FindAll(ctx context.Context, query map[string][]string) ([]mDaybook.DaybookList, error) {
	res, err := s.daybookRepo.FindAll(ctx, query)
	if err != nil {
		return []mDaybook.DaybookList{}, err
	}
	return res, nil
}

func (s daybookService) FindById(ctx context.Context, id string) (mDaybook.DaybookResponse, error) {
	res, err := s.daybookRepo.FindById(ctx, id)
	if err != nil {
		return mDaybook.DaybookResponse{}, err
	}
	return res, nil
}

func (s daybookService) GenerateExcel(ctx context.Context, id string) (mRes.ExcelFile, error) {
	res, err := s.daybookRepo.FindByIdForExcel(ctx, id)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	xlsx, err := excelize.OpenFile(fmt.Sprintf("config/templates/daybook/%s.xlsx", res.Company.Id.Hex()))

	if err != nil {
		return mRes.ExcelFile{}, err
	}
	fm := "_-* #,##0.00_-;-* #,##0.00_-;_-* \"-\"??_-;_-@_-"
	sheetName := "Sheet1"
	// if err := xlsx.AddPicture(sheetName, "A2", "658e542c6aebff64cf245e43.png", nil); err != nil {
	// 	return mRes.ExcelFile{}, err
	// }
	xlsx.SetCellValue(sheetName, "B2", res.Company.Name)
	xlsx.SetCellValue(sheetName, "B3", res.Company.Address)
	xlsx.SetCellValue(sheetName, "A5", res.Document.Name)
	var fileName string
	if res.Supplier != nil {
		fileName = fmt.Sprintf("%s-%s.xlsx", res.Number, res.Supplier.Name)
		xlsx.SetCellValue(sheetName, "A7", res.Supplier.Code)
		xlsx.SetCellValue(sheetName, "B8", res.Supplier.Name)
	}
	if res.Customer != nil {
		fileName = fmt.Sprintf("%s-%s.xlsx", res.Number, res.Customer.Name)
		xlsx.SetCellValue(sheetName, "A7", res.Customer.Code)
		xlsx.SetCellValue(sheetName, "B8", res.Customer.Name)
	}
	xlsx.SetCellValue(sheetName, "G6", res.Number)
	xlsx.SetCellValue(sheetName, "G7", res.TransactionDate.Format("02/01/2006"))
	xlsx.SetCellValue(sheetName, "G8", res.Invoice)
	numberStyle, err := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "TH Sarabun New",
			Size:   24,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Indent:     1,
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &fm,
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	textStyle, err := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "TH Sarabun New",
			Size:   24,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Indent:     1,
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	cell := 10
	for _, detail := range res.DaybookDetails {
		err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("A%d", cell), fmt.Sprintf("A%d", cell), textStyle)
		if err != nil {
			return mRes.ExcelFile{}, err
		}
		err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("D%d", cell), fmt.Sprintf("D%d", cell), textStyle)
		if err != nil {
			return mRes.ExcelFile{}, err
		}
		xlsx.SetCellValue(sheetName, fmt.Sprintf("A%d", cell), detail.Account.Code)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("B%d", cell), detail.Account.Name)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("E%d", cell), detail.Name)
		switch detail.Type {
		case "DR":
			xlsx.SetCellValue(sheetName, fmt.Sprintf("G%d", cell), detail.Amount)
		case "CR":
			xlsx.SetCellValue(sheetName, fmt.Sprintf("H%d", cell), detail.Amount)
		}
		err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("G%d", cell), fmt.Sprintf("G%d", cell), numberStyle)
		if err != nil {
			return mRes.ExcelFile{}, err
		}
		err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("H%d", cell), fmt.Sprintf("H%d", cell), numberStyle)
		if err != nil {
			return mRes.ExcelFile{}, err
		}
		cell++
	}
	defaultTableRecord := 20
	if cell < defaultTableRecord {
		length := defaultTableRecord - cell
		for i := 0; i < length; i++ {
			err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("A%d", cell), fmt.Sprintf("A%d", cell), textStyle)
			if err != nil {
				return mRes.ExcelFile{}, err
			}
			err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("D%d", cell), fmt.Sprintf("D%d", cell), textStyle)
			if err != nil {
				return mRes.ExcelFile{}, err
			}
			err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("G%d", cell), fmt.Sprintf("G%d", cell), numberStyle)
			if err != nil {
				return mRes.ExcelFile{}, err
			}
			err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("H%d", cell), fmt.Sprintf("H%d", cell), numberStyle)
			if err != nil {
				return mRes.ExcelFile{}, err
			}
			if i < length-1 {
				cell++
			}
		}
	}
	cell++
	bahtUnitColumn := fmt.Sprintf("A%d", cell)
	xlsx.SetCellValue(sheetName, bahtUnitColumn, "บาท:")
	style, err := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "TH Sarabun New",
			Size:   26,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Indent:     1,
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, bahtUnitColumn, bahtUnitColumn, style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	totalTextColumn := fmt.Sprintf("F%d", cell)
	xlsx.SetCellValue(sheetName, totalTextColumn, "รวม")
	style, err = xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "TH Sarabun New",
			Size:   26,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Indent:     1,
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, totalTextColumn, totalTextColumn, style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	sumDrColumn := fmt.Sprintf("G%d", cell)
	xlsx.SetCellFormula(sheetName, sumDrColumn, fmt.Sprintf("SUM(G10:G%d)", cell-1))
	style, err = xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "TH Sarabun New",
			Size:   26,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Indent:     1,
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		CustomNumFmt: &fm,
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, sumDrColumn, sumDrColumn, style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	sumCrColumn := fmt.Sprintf("H%d", cell)
	xlsx.SetCellFormula(sheetName, sumCrColumn, fmt.Sprintf("SUM(H10:H%d)", cell-1))
	style, err = xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "TH Sarabun New",
			Size:   26,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Indent:     1,
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		CustomNumFmt: &fm,
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, sumCrColumn, sumCrColumn, style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	bahtTextColumn := fmt.Sprintf("B%d", cell)
	err = xlsx.MergeCell(sheetName, bahtTextColumn, fmt.Sprintf("E%d", cell))
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "TH Sarabun New",
			Size:   26,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Indent:     1,
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, bahtTextColumn, bahtTextColumn, style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	xlsx.SetCellFormula(sheetName, bahtTextColumn, fmt.Sprintf("BAHTTEXT(%s)", sumDrColumn))
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("C%d", cell), fmt.Sprintf("C%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("D%d", cell), fmt.Sprintf("D%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("E%d", cell), fmt.Sprintf("E%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	cell++
	err = xlsx.AddFormControl(sheetName, excelize.FormControl{
		Cell: fmt.Sprintf("A%d", cell),
		Type: excelize.FormControlCheckBox,
		Paragraph: []excelize.RichTextRun{
			{
				Font: &excelize.Font{
					Family: "TH Sarabun New",
					Size:   26,
					Color:  "000000",
				},
				Text: "เงินสด",
			},
		},
		Checked: true,
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("A%d", cell), fmt.Sprintf("A%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.AddFormControl(sheetName, excelize.FormControl{
		Cell: fmt.Sprintf("B%d", cell),
		Type: excelize.FormControlCheckBox,
		Paragraph: []excelize.RichTextRun{
			{
				Font: &excelize.Font{
					Family: "TH Sarabun New",
					Size:   26,
					Color:  "000000",
				},
				Text: "เช็ค",
			},
		},
		Checked: false,
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.AddFormControl(sheetName, excelize.FormControl{
		Cell: fmt.Sprintf("C%d", cell),
		Type: excelize.FormControlCheckBox,
		Paragraph: []excelize.RichTextRun{
			{
				Font: &excelize.Font{
					Family: "TH Sarabun New",
					Size:   26,
					Color:  "000000",
				},
				Text: "โอนMB",
			},
		},
		Checked: false,
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.AddFormControl(sheetName, excelize.FormControl{
		Cell: fmt.Sprintf("D%d", cell),
		Type: excelize.FormControlCheckBox,
		Paragraph: []excelize.RichTextRun{
			{
				Font: &excelize.Font{
					Family: "TH Sarabun New",
					Size:   26,
					Color:  "000000",
				},
				Text: "หักบัญชี",
			},
		},
		Checked: false,
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("D%d", cell), fmt.Sprintf("D%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("F%d", cell), fmt.Sprintf("F%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("H%d", cell), fmt.Sprintf("H%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	cell++
	bankTextColumn := fmt.Sprintf("A%d", cell)
	err = xlsx.MergeCell(sheetName, bankTextColumn, fmt.Sprintf("D%d", cell))
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	xlsx.SetCellValue(sheetName, bankTextColumn, "ธนาคาร………………………...…………………….")
	style, err = xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "TH Sarabun New",
			Size:   26,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Indent:     1,
			Vertical:   "bottom",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, bankTextColumn, bankTextColumn, style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("D%d", cell), fmt.Sprintf("D%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	creatorTextColumn := fmt.Sprintf("E%d", cell)
	err = xlsx.MergeCell(sheetName, creatorTextColumn, fmt.Sprintf("F%d", cell))
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	xlsx.SetCellValue(sheetName, creatorTextColumn, fmt.Sprintf(".......%s.......ผู้จัดทำ", res.Company.Contact))
	style, err = xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "TH Sarabun New",
			Size:   26,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Indent:     1,
			Vertical:   "bottom",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, creatorTextColumn, creatorTextColumn, style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("F%d", cell), fmt.Sprintf("F%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("H%d", cell), fmt.Sprintf("H%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	///
	cell++
	checkNumberTextColumn := fmt.Sprintf("A%d", cell)
	err = xlsx.MergeCell(sheetName, checkNumberTextColumn, fmt.Sprintf("D%d", cell))
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	xlsx.SetCellValue(sheetName, checkNumberTextColumn, "เช็คเลขที่…….……………………………..…….….")
	style, err = xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "TH Sarabun New",
			Size:   26,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Indent:     1,
			Vertical:   "bottom",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, checkNumberTextColumn, checkNumberTextColumn, style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("D%d", cell), fmt.Sprintf("D%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	bookKeeperTextColumn := fmt.Sprintf("E%d", cell)
	err = xlsx.MergeCell(sheetName, bookKeeperTextColumn, fmt.Sprintf("F%d", cell))
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	xlsx.SetCellValue(sheetName, bookKeeperTextColumn, fmt.Sprintf(".......%s.......ผู้บันทึกบัญชี", res.Company.Contact))
	style, err = xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "TH Sarabun New",
			Size:   26,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Indent:     1,
			Vertical:   "bottom",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, bookKeeperTextColumn, bookKeeperTextColumn, style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("F%d", cell), fmt.Sprintf("F%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	dotAreaColumn := fmt.Sprintf("G%d", cell)
	err = xlsx.MergeCell(sheetName, dotAreaColumn, fmt.Sprintf("H%d", cell))
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	xlsx.SetCellValue(sheetName, dotAreaColumn, "…………………………….")
	style, err = xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "TH Sarabun New",
			Size:   26,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Indent:     1,
			Vertical:   "bottom",
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, dotAreaColumn, dotAreaColumn, style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("H%d", cell), fmt.Sprintf("H%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	///
	cell++
	datedTextColumn := fmt.Sprintf("A%d", cell)
	err = xlsx.MergeCell(sheetName, datedTextColumn, fmt.Sprintf("D%d", cell))
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	xlsx.SetCellValue(sheetName, datedTextColumn, "ลงวันที่……………………………….…..….……….")
	style, err = xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "TH Sarabun New",
			Size:   26,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Indent:     1,
			Vertical:   "bottom",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, datedTextColumn, datedTextColumn, style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("B%d", cell), fmt.Sprintf("B%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("C%d", cell), fmt.Sprintf("C%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("D%d", cell), fmt.Sprintf("D%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("D%d", cell), fmt.Sprintf("D%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	checkerTextColumn := fmt.Sprintf("E%d", cell)
	err = xlsx.MergeCell(sheetName, checkerTextColumn, fmt.Sprintf("F%d", cell))
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	xlsx.SetCellValue(sheetName, checkerTextColumn, "......………………………...…………….ผู้ตรวจ")
	style, err = xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "TH Sarabun New",
			Size:   26,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Indent:     1,
			Vertical:   "bottom",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, checkerTextColumn, checkerTextColumn, style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("F%d", cell), fmt.Sprintf("F%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	approverTextColumn := fmt.Sprintf("G%d", cell)
	err = xlsx.MergeCell(sheetName, approverTextColumn, fmt.Sprintf("H%d", cell))
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	xlsx.SetCellValue(sheetName, approverTextColumn, "ผู้อนุมัติ")
	style, err = xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "TH Sarabun New",
			Size:   26,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Indent:     1,
			Vertical:   "bottom",
		},
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, approverTextColumn, approverTextColumn, style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("H%d", cell), fmt.Sprintf("H%d", cell), style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	cell = cell + 3
	endTextColumn := fmt.Sprintf("A%d", cell)
	xlsx.SetCellValue(sheetName, endTextColumn, "ผู้รับเงิน……………………………………………………………….วันที่…………………………………………………………………")
	style, err = xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "TH Sarabun New",
			Size:   26,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Indent:     1,
			Vertical:   "right",
		},
	})
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	err = xlsx.SetCellStyle(sheetName, endTextColumn, endTextColumn, style)
	if err != nil {
		return mRes.ExcelFile{}, err
	}
	f := mRes.ExcelFile{}

	f.File = xlsx
	f.Name = fileName

	return f, nil
}

func (s daybookService) Create(ctx context.Context, payload mDaybook.DaybookPayload) (mDaybook.DaybookPayload, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	// new daybook
	newId := primitive.NewObjectID()
	newDt := time.Now()
	inv := mDaybook.Daybook{}
	inv.Id = newId
	inv.Number = payload.Number
	inv.Invoice = payload.Invoice
	inv.Document = payload.Document
	inv.TransactionDate = payload.TransactionDate
	inv.Company = payload.Company
	inv.Supplier = payload.Supplier
	inv.Customer = payload.Customer
	inv.CreatedBy = user.Id
	inv.CreatedAt = newDt
	inv.UpdatedBy = user.Id
	inv.UpdatedAt = newDt

	// new daybook detail
	payload.Id = inv.Id
	payload.CreatedBy = inv.CreatedBy
	payload.CreatedAt = inv.CreatedAt
	payload.UpdatedBy = inv.UpdatedBy
	payload.UpdatedAt = inv.UpdatedAt
	if len(payload.DaybookDetails) > 0 {
		docs := make([]interface{}, len(payload.DaybookDetails))
		for i, doc := range payload.DaybookDetails {
			newId = primitive.NewObjectID()
			newDt = time.Now()
			docs[i] = mDaybookDetail.DaybookDetail{
				Id:        newId,
				Name:      doc.Name,
				Type:      doc.Type,
				Amount:    doc.Amount,
				Account:   doc.Account,
				Daybook:   payload.Id,
				CreatedBy: user.Id,
				CreatedAt: newDt,
				UpdatedBy: user.Id,
				UpdatedAt: newDt,
			}
			inv.DaybookDetails = append(inv.DaybookDetails, newId)
			payload.DaybookDetails[i].Id = newId
			payload.DaybookDetails[i].CreatedBy = user.Id
			payload.DaybookDetails[i].CreatedAt = newDt
			payload.DaybookDetails[i].UpdatedBy = user.Id
			payload.DaybookDetails[i].UpdatedAt = newDt
		}
		err = s.daybookRepo.BulkCreateDaybookDetail(ctx, docs)
		if err != nil {
			return payload, err
		}
	}

	_, err = s.daybookRepo.Create(ctx, inv)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

func (s daybookService) Update(ctx context.Context, id string, payload mDaybook.Daybook) (mDaybook.Daybook, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mDaybook.Daybook{}, err
	}
	payload.Id = doc
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.daybookRepo.Update(ctx, payload)
	if err != nil {
		return mDaybook.Daybook{}, err
	}
	return res, nil
}
