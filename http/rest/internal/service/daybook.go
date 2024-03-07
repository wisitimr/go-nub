package service

import (
	"context"
	"findigitalservice/http/rest/internal/auth"
	mDaybook "findigitalservice/http/rest/internal/model/daybook"
	mDaybookDetail "findigitalservice/http/rest/internal/model/daybook_detail"
	mForwardAccount "findigitalservice/http/rest/internal/model/forward_account"
	mRepo "findigitalservice/http/rest/internal/model/repository"
	mRes "findigitalservice/http/rest/internal/model/response"
	mService "findigitalservice/http/rest/internal/model/service"
	mUser "findigitalservice/http/rest/internal/model/user"
	"findigitalservice/http/rest/internal/util"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type daybookService struct {
	Daybook        mRepo.DaybookRepository
	PaymentMethod  mRepo.PaymentMethodRepository
	Account        mRepo.AccountRepository
	ForwardAccount mRepo.ForwardAccountRepository
	logger         *logrus.Logger
}

func InitDaybookService(repo mRepo.Repository, logger *logrus.Logger) mService.DaybookService {
	return &daybookService{
		Daybook:        repo.Daybook,
		PaymentMethod:  repo.PaymentMethod,
		Account:        repo.Account,
		ForwardAccount: repo.ForwardAccount,
		logger:         logger,
	}
}

func (s daybookService) Count(ctx context.Context, query map[string][]string) (mRes.CountDto, error) {
	count, err := s.Daybook.Count(ctx, query)
	if err != nil {
		return mRes.CountDto{Count: 0}, err
	}
	return mRes.CountDto{Count: count}, nil
}

func (s daybookService) FindAll(ctx context.Context, query map[string][]string) ([]mDaybook.DaybookList, error) {
	res, err := s.Daybook.FindAll(ctx, query)
	if err != nil {
		return []mDaybook.DaybookList{}, err
	}
	return res, nil
}

func (s daybookService) FindById(ctx context.Context, id string) (mDaybook.DaybookResponse, error) {
	res, err := s.Daybook.FindById(ctx, id)
	if err != nil {
		return mDaybook.DaybookResponse{}, err
	}
	return res, nil
}

func (s daybookService) GenerateExcel(ctx context.Context, id string) (*excelize.File, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	res, err := s.Daybook.FindByIdForExcel(ctx, id)
	if err != nil {
		return nil, err
	}
	xlsx, err := excelize.OpenFile(fmt.Sprintf("config/templates/daybook/%s.xlsx", res.Company.Id.Hex()))

	if err != nil {
		return nil, err
	}
	fm := "_-* #,##0.00_-;-* #,##0.00_-;_-* \"-\"??_-;_-@_-"
	sheetName := "Sheet1"
	// if err := xlsx.AddPicture(sheetName, "A2", "658e542c6aebff64cf245e43.png", nil); err != nil {
	// 	return nil, err
	// }
	xlsx.SetCellValue(sheetName, "B2", res.Company.Name)
	xlsx.SetCellValue(sheetName, "B3", res.Company.Address)
	xlsx.SetCellValue(sheetName, "A5", res.Document.Name)
	if res.Supplier != nil {
		xlsx.SetCellValue(sheetName, "A7", res.Supplier.Code)
		xlsx.SetCellValue(sheetName, "B8", res.Supplier.Name)
	}
	if res.Customer != nil {
		xlsx.SetCellValue(sheetName, "A7", res.Customer.Code)
		xlsx.SetCellValue(sheetName, "B8", res.Customer.Name)
	}
	xlsx.SetCellValue(sheetName, "G6", res.Number)
	xlsx.SetCellValue(sheetName, "G7", res.TransactionDate.Format("02/01/2006"))
	xlsx.SetCellValue(sheetName, "G8", res.Invoice)
	priceStyle, err := xlsx.NewStyle(&excelize.Style{
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
		return nil, err
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
		return nil, err
	}
	row := 10
	for _, detail := range res.DaybookDetails {
		err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), textStyle)
		if err != nil {
			return nil, err
		}
		err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), textStyle)
		if err != nil {
			return nil, err
		}
		xlsx.SetCellValue(sheetName, fmt.Sprintf("A%d", row), detail.Account.Code)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("B%d", row), detail.Account.Name)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("E%d", row), detail.Name)
		switch detail.Type {
		case "DR":
			xlsx.SetCellValue(sheetName, fmt.Sprintf("G%d", row), detail.Amount)
		case "CR":
			xlsx.SetCellValue(sheetName, fmt.Sprintf("H%d", row), detail.Amount)
		}
		err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), priceStyle)
		if err != nil {
			return nil, err
		}
		err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), priceStyle)
		if err != nil {
			return nil, err
		}
		row++
	}
	defaultTableRecord := 20
	if row < defaultTableRecord {
		length := defaultTableRecord - row
		for i := 0; i < length; i++ {
			err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), textStyle)
			if err != nil {
				return nil, err
			}
			err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), textStyle)
			if err != nil {
				return nil, err
			}
			err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), priceStyle)
			if err != nil {
				return nil, err
			}
			err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), priceStyle)
			if err != nil {
				return nil, err
			}
			if i < length-1 {
				row++
			}
		}
	}
	row++
	bahtUnitColumn := fmt.Sprintf("A%d", row)
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
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, bahtUnitColumn, bahtUnitColumn, style)
	if err != nil {
		return nil, err
	}
	totalTextColumn := fmt.Sprintf("F%d", row)
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
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, totalTextColumn, totalTextColumn, style)
	if err != nil {
		return nil, err
	}
	sumDrColumn := fmt.Sprintf("G%d", row)
	xlsx.SetCellFormula(sheetName, sumDrColumn, fmt.Sprintf("SUM(G10:G%d)", row-1))
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
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, sumDrColumn, sumDrColumn, style)
	if err != nil {
		return nil, err
	}
	sumCrColumn := fmt.Sprintf("H%d", row)
	xlsx.SetCellFormula(sheetName, sumCrColumn, fmt.Sprintf("SUM(H10:H%d)", row-1))
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
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, sumCrColumn, sumCrColumn, style)
	if err != nil {
		return nil, err
	}
	bahtTextColumn := fmt.Sprintf("B%d", row)
	err = xlsx.MergeCell(sheetName, bahtTextColumn, fmt.Sprintf("E%d", row))
	if err != nil {
		return nil, err
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
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, bahtTextColumn, bahtTextColumn, style)
	if err != nil {
		return nil, err
	}
	xlsx.SetCellFormula(sheetName, bahtTextColumn, fmt.Sprintf("BAHTTEXT(%s)", sumDrColumn))
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), style)
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), style)
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), style)
	if err != nil {
		return nil, err
	}
	row++
	allPaymentMethod, _ := s.PaymentMethod.FindAll(ctx, make(map[string][]string))
	var methodColumn = [...]string{"A", "B", "C", "D"}
	for i, method := range allPaymentMethod {
		styleCol := methodColumn[i]
		err = xlsx.AddFormControl(sheetName, excelize.FormControl{
			Cell: fmt.Sprintf("%s%d", styleCol, row),
			Type: excelize.FormControlCheckBox,
			Paragraph: []excelize.RichTextRun{
				{
					Font: &excelize.Font{
						Family: "TH Sarabun New",
						Size:   26,
						Color:  "000000",
					},
					Text: method.Name,
				},
			},
			Checked: method.Id == res.PaymentMethod,
		})
		if err != nil {
			return nil, err
		}
		if styleCol == "A" {
			style, err = xlsx.NewStyle(&excelize.Style{
				Border: []excelize.Border{
					{Type: "left", Color: "000000", Style: 1},
				},
			})
			if err != nil {
				return nil, err
			}
			err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("%s%d", styleCol, row), fmt.Sprintf("%s%d", styleCol, row), style)
			if err != nil {
				return nil, err
			}
		}
		if styleCol == "D" {
			style, err = xlsx.NewStyle(&excelize.Style{
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
				},
			})
			if err != nil {
				return nil, err
			}
			err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("%s%d", styleCol, row), fmt.Sprintf("%s%d", styleCol, row), style)
			if err != nil {
				return nil, err
			}
		}
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), style)
	if err != nil {
		return nil, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), style)
	if err != nil {
		return nil, err
	}
	row++
	bankTextColumn := fmt.Sprintf("A%d", row)
	err = xlsx.MergeCell(sheetName, bankTextColumn, fmt.Sprintf("D%d", row))
	if err != nil {
		return nil, err
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
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, bankTextColumn, bankTextColumn, style)
	if err != nil {
		return nil, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), style)
	if err != nil {
		return nil, err
	}
	creatorTextColumn := fmt.Sprintf("E%d", row)
	err = xlsx.MergeCell(sheetName, creatorTextColumn, fmt.Sprintf("F%d", row))
	if err != nil {
		return nil, err
	}
	xlsx.SetCellValue(sheetName, creatorTextColumn, fmt.Sprintf(".......%s %s.......ผู้จัดทำ", user.FirstName, user.LastName))
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
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, creatorTextColumn, creatorTextColumn, style)
	if err != nil {
		return nil, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), style)
	if err != nil {
		return nil, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), style)
	if err != nil {
		return nil, err
	}
	///
	row++
	checkNumberTextColumn := fmt.Sprintf("A%d", row)
	err = xlsx.MergeCell(sheetName, checkNumberTextColumn, fmt.Sprintf("D%d", row))
	if err != nil {
		return nil, err
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
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, checkNumberTextColumn, checkNumberTextColumn, style)
	if err != nil {
		return nil, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), style)
	if err != nil {
		return nil, err
	}
	bookKeeperTextColumn := fmt.Sprintf("E%d", row)
	err = xlsx.MergeCell(sheetName, bookKeeperTextColumn, fmt.Sprintf("F%d", row))
	if err != nil {
		return nil, err
	}
	xlsx.SetCellValue(sheetName, bookKeeperTextColumn, fmt.Sprintf(".......%s %s.......ผู้บันทึกบัญชี", user.FirstName, user.LastName))
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
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, bookKeeperTextColumn, bookKeeperTextColumn, style)
	if err != nil {
		return nil, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), style)
	if err != nil {
		return nil, err
	}
	dotAreaColumn := fmt.Sprintf("G%d", row)
	err = xlsx.MergeCell(sheetName, dotAreaColumn, fmt.Sprintf("H%d", row))
	if err != nil {
		return nil, err
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
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, dotAreaColumn, dotAreaColumn, style)
	if err != nil {
		return nil, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), style)
	if err != nil {
		return nil, err
	}
	///
	row++
	datedTextColumn := fmt.Sprintf("A%d", row)
	err = xlsx.MergeCell(sheetName, datedTextColumn, fmt.Sprintf("D%d", row))
	if err != nil {
		return nil, err
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
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, datedTextColumn, datedTextColumn, style)
	if err != nil {
		return nil, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), style)
	if err != nil {
		return nil, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), style)
	if err != nil {
		return nil, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), style)
	if err != nil {
		return nil, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), style)
	if err != nil {
		return nil, err
	}
	checkerTextColumn := fmt.Sprintf("E%d", row)
	err = xlsx.MergeCell(sheetName, checkerTextColumn, fmt.Sprintf("F%d", row))
	if err != nil {
		return nil, err
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
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, checkerTextColumn, checkerTextColumn, style)
	if err != nil {
		return nil, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), style)
	if err != nil {
		return nil, err
	}
	approverTextColumn := fmt.Sprintf("G%d", row)
	err = xlsx.MergeCell(sheetName, approverTextColumn, fmt.Sprintf("H%d", row))
	if err != nil {
		return nil, err
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
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, approverTextColumn, approverTextColumn, style)
	if err != nil {
		return nil, err
	}
	style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), style)
	if err != nil {
		return nil, err
	}
	row = row + 3
	endTextColumn := fmt.Sprintf("A%d", row)
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
		return nil, err
	}
	err = xlsx.SetCellStyle(sheetName, endTextColumn, endTextColumn, style)
	if err != nil {
		return nil, err
	}

	return xlsx, nil
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
	inv.PaymentMethod = payload.PaymentMethod
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
		err = s.Daybook.BulkCreateDaybookDetail(ctx, docs)
		if err != nil {
			return payload, err
		}
	}

	_, err = s.Daybook.Create(ctx, inv)
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
	res, err := s.Daybook.Update(ctx, payload)
	if err != nil {
		return mDaybook.Daybook{}, err
	}
	return res, nil
}

func (s daybookService) GenerateFinancialStatement(ctx context.Context, company string, year string) (*excelize.File, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	financial, err := s.Daybook.GenerateFinancialStatement(ctx, company, year)
	if err != nil {
		return nil, err
	}
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return nil, err
	}
	var xlsx *excelize.File
	if len(financial) > 0 {
		mapFin := make(map[string]mDaybook.DaybookFinancialStatement)
		for _, v := range financial {
			mapFin[v.Code] = v
		}
		xlsx, err = excelize.OpenFile(fmt.Sprintf("config/templates/financial_statement/%s.xlsx", company))

		if err != nil {
			s.logger.Error(err)
		}
		fm := "_-* #,##0.00_-;-* #,##0.00_-;_-* \"-\"??_-;_-@_-"
		sheet := "TB12"
		// TB12
		query := make(map[string][]string)
		query["company"] = append(query["company"], company)
		account, err := s.Account.FindAll(ctx, query)
		if err != nil {
			return nil, err
		}
		if financial[0].Company.Name != "" {
			xlsx.SetCellValue(sheet, "A1", financial[0].Company.Name)
		}
		var fromDate string
		if len(financial[0].DaybookDetails) > 0 && !financial[0].DaybookDetails[0].Daybook.TransactionDate.IsZero() {
			year := financial[0].DaybookDetails[0].Daybook.TransactionDate.Format("2006")
			fromDate = fmt.Sprintf("From Date :  1 Jan %s To  31 December %s", year, year)
		}
		xlsx.SetCellValue(sheet, "A3", fromDate)
		row := 7
		groupCode := 1
		accountType := ""
		// s.logger.Warn(janDr)
		sumStyle, err := xlsx.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold:   true,
				Family: "TH Sarabun New",
				Size:   9,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "right",
				Vertical:   "bottom",
			},
			Border: []excelize.Border{
				{Type: "top", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "left", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
			// Fill: excelize.Fill{Type: "gradient", Color: []string{"E3F2FD"}, Shading: 1},
			Fill:         excelize.Fill{Type: "pattern", Color: []string{"E0ECF4"}, Pattern: 1},
			CustomNumFmt: &fm,
		})
		if err != nil {
			return nil, err
		}
		titleStyle, err := xlsx.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold:   true,
				Family: "TH Sarabun New",
				Size:   9,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "bottom",
			},
			Border: []excelize.Border{
				{Type: "top", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "left", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
			// Fill: excelize.Fill{Type: "gradient", Color: []string{"E3F2FD"}, Shading: 1},
			Fill:         excelize.Fill{Type: "pattern", Color: []string{"E0ECF4"}, Pattern: 1},
			CustomNumFmt: &fm,
		})
		if err != nil {
			return nil, err
		}
		priceStyle, err := xlsx.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Family: "TH Sarabun New",
				Size:   9,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "right",
				Vertical:   "bottom",
			},
			Border: []excelize.Border{
				{Type: "top", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "left", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
			CustomNumFmt: &fm,
		})
		if err != nil {
			return nil, err
		}
		blankStyle, err := xlsx.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Family: "TH Sarabun New",
				Size:   9,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "right",
				Vertical:   "bottom",
			},
			Border: []excelize.Border{
				{Type: "top", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "left", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
		})
		if err != nil {
			return nil, err
		}
		var totalAccountForwardDrSum []string
		var totalAccountForwardCrSum []string
		var janAccountDrSum []string
		var janAccountCrSum []string
		var febAccountDrSum []string
		var febAccountCrSum []string
		var marAccountDrSum []string
		var marAccountCrSum []string
		var aprAccountDrSum []string
		var aprAccountCrSum []string
		var mayAccountDrSum []string
		var mayAccountCrSum []string
		var junAccountDrSum []string
		var junAccountCrSum []string
		var julAccountDrSum []string
		var julAccountCrSum []string
		var augAccountDrSum []string
		var augAccountCrSum []string
		var sepAccountDrSum []string
		var sepAccountCrSum []string
		var octAccountDrSum []string
		var octAccountCrSum []string
		var novAccountDrSum []string
		var novAccountCrSum []string
		var decAccountDrSum []string
		var decAccountCrSum []string
		var totalAccountDrSum []string
		var totalAccountCrSum []string
		var totalAccountForwardDr []string
		var totalAccountForwardCr []string
		var janAccountDr []string
		var janAccountCr []string
		var febAccountDr []string
		var febAccountCr []string
		var marAccountDr []string
		var marAccountCr []string
		var aprAccountDr []string
		var aprAccountCr []string
		var mayAccountDr []string
		var mayAccountCr []string
		var junAccountDr []string
		var junAccountCr []string
		var julAccountDr []string
		var julAccountCr []string
		var augAccountDr []string
		var augAccountCr []string
		var sepAccountDr []string
		var sepAccountCr []string
		var octAccountDr []string
		var octAccountCr []string
		var novAccountDr []string
		var novAccountCr []string
		var decAccountDr []string
		var decAccountCr []string
		var totalAccountDr []string
		var totalAccountCr []string
		var totalAccountEnding []string
		var sumDr []string
		var sumCr []string
		var sumDrTmp float64
		var sumCrTmp float64
		var resultNetProfitLoss []string
		var resultDiffAssetsLiabilitiesOwnerEquity []string
		var resultDifference []string
		fwMap := make(map[string]string)
		var typeList []mDaybook.FinancialStatement
		for i := 0; i < len(account); i++ {
			acc := account[i]
			if accountType == "" {
				accountType = acc.Type
			}
			isTotal := false
			accountFirstNo, _ := strconv.Atoi(acc.Code[0:1])
			if accountFirstNo == groupCode || (accountFirstNo > 5 && accountFirstNo <= 9) {
				isTotal = false
				accCode := fmt.Sprintf("A%d", row)
				xlsx.SetCellValue(sheet, accCode, acc.Code)
				style, err := xlsx.NewStyle(&excelize.Style{
					Font: &excelize.Font{
						Family: "TH Sarabun New",
						Size:   9,
					},
					Alignment: &excelize.Alignment{
						Horizontal: "center",
						Vertical:   "bottom",
					},
					Border: []excelize.Border{
						{Type: "top", Color: "000000", Style: 1},
						{Type: "right", Color: "000000", Style: 1},
						{Type: "left", Color: "000000", Style: 1},
						{Type: "bottom", Color: "000000", Style: 1},
					},
				})
				if err != nil {
					return nil, err
				}
				err = xlsx.SetCellStyle(sheet, accCode, accCode, style)
				if err != nil {
					return nil, err
				}
				accName := fmt.Sprintf("B%d", row)
				xlsx.SetCellValue(sheet, accName, acc.Name)
				style, err = xlsx.NewStyle(&excelize.Style{
					Font: &excelize.Font{
						Family: "TH Sarabun New",
						Size:   9,
					},
					Alignment: &excelize.Alignment{
						Horizontal: "general",
						Vertical:   "bottom",
					},
					Border: []excelize.Border{
						{Type: "top", Color: "000000", Style: 1},
						{Type: "right", Color: "000000", Style: 1},
						{Type: "left", Color: "000000", Style: 1},
						{Type: "bottom", Color: "000000", Style: 1},
					},
				})
				if err != nil {
					return nil, err
				}
				err = xlsx.SetCellStyle(sheet, accName, accName, style)
				if err != nil {
					return nil, err
				}
				data := mapFin[acc.Code]
				var janDr float64
				var janCr float64
				var febDr float64
				var febCr float64
				var marDr float64
				var marCr float64
				var aprDr float64
				var aprCr float64
				var mayDr float64
				var mayCr float64
				var junDr float64
				var junCr float64
				var julDr float64
				var julCr float64
				var augDr float64
				var augCr float64
				var sepDr float64
				var sepCr float64
				var octDr float64
				var octCr float64
				var novDr float64
				var novCr float64
				var decDr float64
				var decCr float64
				var forwardingDr float64
				var forwardingCr float64
				if len(totalAccountForwardDr) == 0 {
					totalAccountForwardDr = append(totalAccountForwardDr, fmt.Sprintf("C%d", row))
				}
				if len(totalAccountForwardCr) == 0 {
					totalAccountForwardCr = append(totalAccountForwardCr, fmt.Sprintf("D%d", row))
				}
				if len(janAccountDr) == 0 {
					janAccountDr = append(janAccountDr, fmt.Sprintf("E%d", row))
				}
				sumDr = append(sumDr, fmt.Sprintf("E%d", row))
				if len(janAccountCr) == 0 {
					janAccountCr = append(janAccountCr, fmt.Sprintf("F%d", row))
				}
				sumCr = append(sumCr, fmt.Sprintf("F%d", row))
				if len(febAccountDr) == 0 {
					febAccountDr = append(febAccountDr, fmt.Sprintf("G%d", row))
				}
				sumDr = append(sumDr, fmt.Sprintf("G%d", row))
				if len(febAccountCr) == 0 {
					febAccountCr = append(febAccountCr, fmt.Sprintf("H%d", row))
				}
				sumCr = append(sumCr, fmt.Sprintf("H%d", row))
				if len(marAccountDr) == 0 {
					marAccountDr = append(marAccountDr, fmt.Sprintf("I%d", row))
				}
				sumDr = append(sumDr, fmt.Sprintf("I%d", row))
				if len(marAccountCr) == 0 {
					marAccountCr = append(marAccountCr, fmt.Sprintf("J%d", row))
				}
				sumCr = append(sumCr, fmt.Sprintf("J%d", row))
				if len(aprAccountDr) == 0 {
					aprAccountDr = append(aprAccountDr, fmt.Sprintf("K%d", row))
				}
				sumDr = append(sumDr, fmt.Sprintf("K%d", row))
				if len(aprAccountCr) == 0 {
					aprAccountCr = append(aprAccountCr, fmt.Sprintf("L%d", row))
				}
				sumCr = append(sumCr, fmt.Sprintf("L%d", row))
				if len(mayAccountDr) == 0 {
					mayAccountDr = append(mayAccountDr, fmt.Sprintf("M%d", row))
				}
				sumDr = append(sumDr, fmt.Sprintf("M%d", row))
				if len(mayAccountCr) == 0 {
					mayAccountCr = append(mayAccountCr, fmt.Sprintf("N%d", row))
				}
				sumCr = append(sumCr, fmt.Sprintf("N%d", row))
				if len(junAccountDr) == 0 {
					junAccountDr = append(junAccountDr, fmt.Sprintf("O%d", row))
				}
				sumDr = append(sumDr, fmt.Sprintf("O%d", row))
				if len(junAccountCr) == 0 {
					junAccountCr = append(junAccountCr, fmt.Sprintf("P%d", row))
				}
				sumCr = append(sumCr, fmt.Sprintf("P%d", row))
				if len(julAccountDr) == 0 {
					julAccountDr = append(julAccountDr, fmt.Sprintf("Q%d", row))
				}
				sumDr = append(sumDr, fmt.Sprintf("Q%d", row))
				if len(julAccountCr) == 0 {
					julAccountCr = append(julAccountCr, fmt.Sprintf("R%d", row))
				}
				sumCr = append(sumCr, fmt.Sprintf("R%d", row))
				if len(augAccountDr) == 0 {
					augAccountDr = append(augAccountDr, fmt.Sprintf("S%d", row))
				}
				sumDr = append(sumDr, fmt.Sprintf("S%d", row))
				if len(augAccountCr) == 0 {
					augAccountCr = append(augAccountCr, fmt.Sprintf("T%d", row))
				}
				sumCr = append(sumCr, fmt.Sprintf("T%d", row))
				if len(sepAccountDr) == 0 {
					sepAccountDr = append(sepAccountDr, fmt.Sprintf("U%d", row))
				}
				sumDr = append(sumDr, fmt.Sprintf("U%d", row))
				if len(sepAccountCr) == 0 {
					sepAccountCr = append(sepAccountCr, fmt.Sprintf("V%d", row))
				}
				sumCr = append(sumCr, fmt.Sprintf("V%d", row))
				if len(octAccountDr) == 0 {
					octAccountDr = append(octAccountDr, fmt.Sprintf("W%d", row))
				}
				sumDr = append(sumDr, fmt.Sprintf("W%d", row))
				if len(octAccountCr) == 0 {
					octAccountCr = append(octAccountCr, fmt.Sprintf("X%d", row))
				}
				sumCr = append(sumCr, fmt.Sprintf("X%d", row))
				if len(novAccountDr) == 0 {
					novAccountDr = append(novAccountDr, fmt.Sprintf("Y%d", row))
				}
				sumDr = append(sumDr, fmt.Sprintf("Y%d", row))
				if len(novAccountCr) == 0 {
					novAccountCr = append(novAccountCr, fmt.Sprintf("Z%d", row))
				}
				sumCr = append(sumCr, fmt.Sprintf("Z%d", row))
				if len(decAccountDr) == 0 {
					decAccountDr = append(decAccountDr, fmt.Sprintf("AA%d", row))
				}
				sumDr = append(sumDr, fmt.Sprintf("AA%d", row))
				if len(decAccountCr) == 0 {
					decAccountCr = append(decAccountCr, fmt.Sprintf("AB%d", row))
				}
				sumCr = append(sumCr, fmt.Sprintf("AB%d", row))
				if len(totalAccountDr) == 0 {
					totalAccountDr = append(totalAccountDr, fmt.Sprintf("AC%d", row))
				}
				if len(totalAccountCr) == 0 {
					totalAccountCr = append(totalAccountCr, fmt.Sprintf("AD%d", row))
				}
				if len(totalAccountEnding) == 0 {
					totalAccountEnding = append(totalAccountEnding, fmt.Sprintf("AE%d", row))
				}
				var rowType mDaybook.FinancialStatement
				rowType.Code = acc.Code
				rowType.Name = acc.Name
				mapMonth := make(map[string][]mDaybook.AccountDetail)
				for _, detail := range data.DaybookDetails {
					var rowAccount mDaybook.AccountDetail
					rowAccount.Number = detail.Daybook.Number
					rowAccount.Description = detail.Name
					rowAccount.Date = detail.Daybook.TransactionDate.Day()
					m := util.Month(detail.Daybook.TransactionDate)
					// 1113-01 = BBL-S/A-# 101-883310-1
					if detail.Type == "DR" && acc.Code == "1113-01" {
						rowAccount.AmountDr = detail.Amount
					} else if detail.Type == "CR" {
						rowAccount.AmountCr = detail.Amount
					}
					mapMonth[m] = append(mapMonth[m], rowAccount)
					if util.IsValidMonth(1, detail.Daybook.TransactionDate) {
						switch detail.Type {
						case "DR":
							janDr += detail.Amount
							sumDrTmp += detail.Amount
						case "CR":
							janCr += detail.Amount
							sumCrTmp += detail.Amount
						}
					} else if util.IsValidMonth(2, detail.Daybook.TransactionDate) {
						switch detail.Type {
						case "DR":
							febDr += detail.Amount
							sumDrTmp += detail.Amount
						case "CR":
							febCr += detail.Amount
							sumCrTmp += detail.Amount
						}
					} else if util.IsValidMonth(3, detail.Daybook.TransactionDate) {
						switch detail.Type {
						case "DR":
							marDr += detail.Amount
							sumDrTmp += detail.Amount
						case "CR":
							marCr += detail.Amount
							sumCrTmp += detail.Amount
						}
					} else if util.IsValidMonth(4, detail.Daybook.TransactionDate) {
						switch detail.Type {
						case "DR":
							aprDr += detail.Amount
							sumDrTmp += detail.Amount
						case "CR":
							aprCr += detail.Amount
							sumCrTmp += detail.Amount
						}
					} else if util.IsValidMonth(5, detail.Daybook.TransactionDate) {
						switch detail.Type {
						case "DR":
							mayDr += detail.Amount
							sumDrTmp += detail.Amount
						case "CR":
							mayCr += detail.Amount
							sumCrTmp += detail.Amount
						}
					} else if util.IsValidMonth(6, detail.Daybook.TransactionDate) {
						switch detail.Type {
						case "DR":
							junDr += detail.Amount
							sumDrTmp += detail.Amount
						case "CR":
							junCr += detail.Amount
							sumCrTmp += detail.Amount
						}
					} else if util.IsValidMonth(7, detail.Daybook.TransactionDate) {
						switch detail.Type {
						case "DR":
							julDr += detail.Amount
							sumDrTmp += detail.Amount
						case "CR":
							julCr += detail.Amount
							sumCrTmp += detail.Amount
						}
					} else if util.IsValidMonth(8, detail.Daybook.TransactionDate) {
						switch detail.Type {
						case "DR":
							augDr += detail.Amount
							sumDrTmp += detail.Amount
						case "CR":
							augCr += detail.Amount
							sumCrTmp += detail.Amount
						}
					} else if util.IsValidMonth(9, detail.Daybook.TransactionDate) {
						switch detail.Type {
						case "DR":
							sepDr += detail.Amount
							sumDrTmp += detail.Amount
						case "CR":
							sepCr += detail.Amount
							sumCrTmp += detail.Amount
						}
					} else if util.IsValidMonth(10, detail.Daybook.TransactionDate) {
						switch detail.Type {
						case "DR":
							octDr += detail.Amount
							sumDrTmp += detail.Amount
						case "CR":
							octCr += detail.Amount
							sumCrTmp += detail.Amount
						}
					} else if util.IsValidMonth(11, detail.Daybook.TransactionDate) {
						switch detail.Type {
						case "DR":
							novDr += detail.Amount
							sumDrTmp += detail.Amount
						case "CR":
							novCr += detail.Amount
							sumCrTmp += detail.Amount
						}
					} else if util.IsValidMonth(12, detail.Daybook.TransactionDate) {
						switch detail.Type {
						case "DR":
							decDr += detail.Amount
							sumDrTmp += detail.Amount
						case "CR":
							decCr += detail.Amount
							sumCrTmp += detail.Amount
						}
					}
				}
				for k, v := range mapMonth {
					var rowMonth mDaybook.MonthDetail
					rowMonth.Month = k
					rowMonth.AccountDetail = append(rowMonth.AccountDetail, v...)
					rowType.MonthDetail = append(rowType.MonthDetail, rowMonth)
				}
				typeList = append(typeList, rowType)
				// s.logger.Error(fmt.Sprintf("E%d", row))
				totalForwardingDr := fmt.Sprintf("C%d", row)
				if forwardingDr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, totalForwardingDr, forwardingDr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, totalForwardingDr, totalForwardingDr, style)
				if err != nil {
					return nil, err
				}
				totalForwardingCr := fmt.Sprintf("D%d", row)
				if forwardingCr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, totalForwardingCr, forwardingCr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, totalForwardingCr, totalForwardingCr, style)
				if err != nil {
					return nil, err
				}
				query := make(map[string][]string)
				query["account"] = append(query["account"], acc.Id.Hex())
				query["year"] = append(query["year"], strconv.Itoa(yearInt-1))
				query["company"] = append(query["company"], acc.Company.Hex())
				fwr, err := s.ForwardAccount.FindOne(ctx, query)
				if err == nil {
					switch fwr.Type {
					case "DR":
						forwardDrColumn := fmt.Sprintf("C%d", row)
						xlsx.SetCellValue(sheet, forwardDrColumn, fwr.Amount)
						err = xlsx.SetCellStyle(sheet, forwardDrColumn, forwardDrColumn, priceStyle)
						if err != nil {
							return nil, err
						}
						fwMap[acc.Code] = fmt.Sprintf("'%s'!%s", sheet, forwardDrColumn)
					case "CR":
						forwardCrColumn := fmt.Sprintf("D%d", row)
						xlsx.SetCellValue(sheet, forwardCrColumn, fwr.Amount)
						err = xlsx.SetCellStyle(sheet, forwardCrColumn, forwardCrColumn, priceStyle)
						if err != nil {
							return nil, err
						}
						fwMap[acc.Code] = fmt.Sprintf("'%s'!%s", sheet, forwardCrColumn)
					}
				}
				janTotalDr := fmt.Sprintf("E%d", row)
				if janDr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, janTotalDr, janDr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, janTotalDr, janTotalDr, style)
				if err != nil {
					return nil, err
				}
				janTotalCr := fmt.Sprintf("F%d", row)
				if janCr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, janTotalCr, janCr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, janTotalCr, janTotalCr, style)
				if err != nil {
					return nil, err
				}
				febTotalDr := fmt.Sprintf("G%d", row)
				if febDr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, febTotalDr, febDr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, febTotalDr, febTotalDr, style)
				if err != nil {
					return nil, err
				}
				febTotalCr := fmt.Sprintf("H%d", row)
				if febCr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, febTotalCr, febCr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, febTotalCr, febTotalCr, style)
				if err != nil {
					return nil, err
				}
				marTotalDr := fmt.Sprintf("I%d", row)
				if marDr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, marTotalDr, marDr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, marTotalDr, marTotalDr, style)
				if err != nil {
					return nil, err
				}
				marTotalCr := fmt.Sprintf("J%d", row)
				if marCr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, marTotalCr, marCr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, marTotalCr, marTotalCr, style)
				if err != nil {
					return nil, err
				}
				aprTotalDr := fmt.Sprintf("K%d", row)
				if aprDr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, aprTotalDr, aprDr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, aprTotalDr, aprTotalDr, style)
				if err != nil {
					return nil, err
				}
				aprTotalCr := fmt.Sprintf("L%d", row)
				if aprCr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, aprTotalCr, aprCr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, aprTotalCr, aprTotalCr, style)
				if err != nil {
					return nil, err
				}
				mayTotalDr := fmt.Sprintf("M%d", row)
				if mayDr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, mayTotalDr, mayDr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, mayTotalDr, mayTotalDr, style)
				if err != nil {
					return nil, err
				}
				mayTotalCr := fmt.Sprintf("N%d", row)
				if mayCr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, mayTotalCr, mayCr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, mayTotalCr, mayTotalCr, style)
				if err != nil {
					return nil, err
				}
				junTotalDr := fmt.Sprintf("O%d", row)
				if junDr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, junTotalDr, junDr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, junTotalDr, junTotalDr, style)
				if err != nil {
					return nil, err
				}
				junTotalCr := fmt.Sprintf("P%d", row)
				if junCr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, junTotalCr, junCr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, junTotalCr, junTotalCr, style)
				if err != nil {
					return nil, err
				}
				julTotalDr := fmt.Sprintf("Q%d", row)
				if julDr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, julTotalDr, julDr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, julTotalDr, julTotalDr, style)
				if err != nil {
					return nil, err
				}
				julTotalCr := fmt.Sprintf("R%d", row)
				if julCr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, julTotalCr, julCr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, julTotalCr, julTotalCr, style)
				if err != nil {
					return nil, err
				}
				augTotalDr := fmt.Sprintf("S%d", row)
				if augDr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, augTotalDr, augDr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, augTotalDr, augTotalDr, style)
				if err != nil {
					return nil, err
				}
				augTotalCr := fmt.Sprintf("T%d", row)
				if augCr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, augTotalCr, augCr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, augTotalCr, augTotalCr, style)
				if err != nil {
					return nil, err
				}
				sepTotalDr := fmt.Sprintf("U%d", row)
				if sepDr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, sepTotalDr, sepDr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, sepTotalDr, sepTotalDr, style)
				if err != nil {
					return nil, err
				}
				sepTotalCr := fmt.Sprintf("V%d", row)
				if sepCr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, sepTotalCr, sepCr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, sepTotalCr, sepTotalCr, style)
				if err != nil {
					return nil, err
				}
				octTotalDr := fmt.Sprintf("W%d", row)
				if octDr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, octTotalDr, octDr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, octTotalDr, octTotalDr, style)
				if err != nil {
					return nil, err
				}
				octTotalCr := fmt.Sprintf("X%d", row)
				if octCr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, octTotalCr, octCr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, octTotalCr, octTotalCr, style)
				if err != nil {
					return nil, err
				}
				novTotalDr := fmt.Sprintf("Y%d", row)
				if novDr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, novTotalDr, novDr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, novTotalDr, novTotalDr, style)
				if err != nil {
					return nil, err
				}
				novTotalCr := fmt.Sprintf("Z%d", row)
				if novCr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, novTotalCr, novCr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, novTotalCr, novTotalCr, style)
				if err != nil {
					return nil, err
				}
				decTotalDr := fmt.Sprintf("AA%d", row)
				if decDr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, decTotalDr, decDr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, decTotalDr, decTotalDr, style)
				if err != nil {
					return nil, err
				}
				decTotalCr := fmt.Sprintf("AB%d", row)
				if decCr == 0 {
					style = blankStyle
				} else {
					xlsx.SetCellValue(sheet, decTotalCr, decCr)
					style = priceStyle
				}
				err = xlsx.SetCellStyle(sheet, decTotalCr, decTotalCr, style)
				if err != nil {
					return nil, err
				}

				// Total DR
				allDr := fmt.Sprintf("AC%d", row)
				xlsx.SetCellFormula(sheet, allDr, fmt.Sprintf("SUM(%s)", strings.Join(sumDr, "+")))
				allCr := fmt.Sprintf("AD%d", row)
				xlsx.SetCellFormula(sheet, allCr, fmt.Sprintf("SUM(%s)", strings.Join(sumCr, "+")))
				resultEnding := fmt.Sprintf("AE%d", row)
				xlsx.SetCellFormula(sheet, resultEnding, fmt.Sprintf("%s-%s", fmt.Sprintf("AC%d", row), fmt.Sprintf("AD%d", row)))
				err = xlsx.SetCellStyle(sheet, allDr, resultEnding, priceStyle)
				if err != nil {
					return nil, err
				}
				twoDigit := fmt.Sprintf("%.2f", sumDrTmp-sumCrTmp)
				e, err := strconv.ParseFloat(twoDigit, 64)
				if err != nil {
					return nil, err
				}
				if e != 0 && accountFirstNo <= 3 {
					var forward mForwardAccount.ForwardAccount
					forward.Account = acc.Id
					switch accountFirstNo {
					case 1:
						if acc.Code == "1420-01" ||
							acc.Code == "1420-02" ||
							acc.Code == "1420-03" ||
							acc.Code == "1420-04" {
							forward.Type = "CR"
						} else {
							forward.Type = "DR"
						}
					case 2:
						forward.Type = "CR"
					case 3:
						forward.Type = "CR"
					}
					forward.Amount = e
					forward.Year = yearInt
					forward.Company = acc.Company
					forward.UpdatedBy = user.Id
					forward.UpdatedAt = time.Now()
					query := make(map[string][]string)
					query["account"] = append(query["account"], acc.Id.Hex())
					query["year"] = append(query["year"], year)
					query["company"] = append(query["company"], acc.Company.Hex())
					fw, err := s.ForwardAccount.FindOne(ctx, query)
					if err != nil {
						forward.Id = primitive.NewObjectID()
						forward.CreatedBy = user.Id
						forward.CreatedAt = time.Now()
						_, err = s.ForwardAccount.Create(ctx, forward)
						if err != nil {
							s.logger.Error(err)
						}
					} else {
						forward.Id = fw.Id
						_, err = s.ForwardAccount.Update(ctx, forward)
						if err != nil {
							s.logger.Error(err)
						}
					}
				}
				sumDr = []string{}
				sumCr = []string{}
				sumDrTmp = 0
				sumCrTmp = 0
			} else {
				isTotal = true
			}
			end := i+1 == len(account)
			if isTotal || end {
				if end {
					row++
				}
				err := xlsx.SetRowHeight(sheet, row, 21.75)
				total := fmt.Sprintf("A%d", row)
				if err != nil {
					return nil, err
				}
				err = xlsx.MergeCell(sheet, total, fmt.Sprintf("B%d", row))
				if err != nil {
					return nil, err
				}
				xlsx.SetCellValue(sheet, total, fmt.Sprintf("รวม%s", accountType))
				accountType = acc.Type
				err = xlsx.SetCellStyle(sheet, total, fmt.Sprintf("B%d", row), titleStyle)
				if err != nil {
					return nil, err
				}
				totalAccountForwardDrColumn := fmt.Sprintf("C%d", row)
				totalAccountForwardDr = append(totalAccountForwardDr, fmt.Sprintf("C%d", row-1))
				totalAccountForwardDrSum = append(totalAccountForwardDrSum, totalAccountForwardDrColumn)
				xlsx.SetCellFormula(sheet, totalAccountForwardDrColumn, fmt.Sprintf("SUM(%s)", strings.Join(totalAccountForwardDr, ":")))
				err = xlsx.SetCellStyle(sheet, totalAccountForwardDrColumn, totalAccountForwardDrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				totalAccountForwardCrColumn := fmt.Sprintf("D%d", row)
				totalAccountForwardCr = append(totalAccountForwardCr, fmt.Sprintf("D%d", row-1))
				totalAccountForwardCrSum = append(totalAccountForwardCrSum, totalAccountForwardCrColumn)
				xlsx.SetCellFormula(sheet, totalAccountForwardCrColumn, fmt.Sprintf("SUM(%s)", strings.Join(totalAccountForwardCr, ":")))
				err = xlsx.SetCellStyle(sheet, totalAccountForwardCrColumn, totalAccountForwardCrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				janAccountDrColumn := fmt.Sprintf("E%d", row)
				janAccountDr = append(janAccountDr, fmt.Sprintf("E%d", row-1))
				janAccountDrSum = append(janAccountDrSum, janAccountDrColumn)
				xlsx.SetCellFormula(sheet, janAccountDrColumn, fmt.Sprintf("SUM(%s)", strings.Join(janAccountDr, ":")))
				err = xlsx.SetCellStyle(sheet, janAccountDrColumn, janAccountDrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				janAccountCrColumn := fmt.Sprintf("F%d", row)
				janAccountCr = append(janAccountCr, fmt.Sprintf("F%d", row-1))
				janAccountCrSum = append(janAccountCrSum, janAccountCrColumn)
				xlsx.SetCellFormula(sheet, janAccountCrColumn, fmt.Sprintf("SUM(%s)", strings.Join(janAccountCr, ":")))
				err = xlsx.SetCellStyle(sheet, janAccountCrColumn, janAccountCrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				febAccountDrColumn := fmt.Sprintf("G%d", row)
				febAccountDr = append(febAccountDr, fmt.Sprintf("G%d", row-1))
				febAccountDrSum = append(febAccountDrSum, febAccountDrColumn)
				xlsx.SetCellFormula(sheet, febAccountDrColumn, fmt.Sprintf("SUM(%s)", strings.Join(febAccountDr, ":")))
				err = xlsx.SetCellStyle(sheet, febAccountDrColumn, febAccountDrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				febAccountCrColumn := fmt.Sprintf("H%d", row)
				febAccountCr = append(febAccountCr, fmt.Sprintf("H%d", row-1))
				febAccountCrSum = append(febAccountCrSum, febAccountCrColumn)
				xlsx.SetCellFormula(sheet, febAccountCrColumn, fmt.Sprintf("SUM(%s)", strings.Join(febAccountCr, ":")))
				err = xlsx.SetCellStyle(sheet, febAccountCrColumn, febAccountCrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				marAccountDrColumn := fmt.Sprintf("I%d", row)
				marAccountDr = append(marAccountDr, fmt.Sprintf("I%d", row-1))
				marAccountDrSum = append(marAccountDrSum, marAccountDrColumn)
				xlsx.SetCellFormula(sheet, marAccountDrColumn, fmt.Sprintf("SUM(%s)", strings.Join(marAccountDr, ":")))
				err = xlsx.SetCellStyle(sheet, marAccountDrColumn, marAccountDrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				marAccountCrColumn := fmt.Sprintf("J%d", row)
				marAccountCr = append(marAccountCr, fmt.Sprintf("J%d", row-1))
				marAccountCrSum = append(marAccountCrSum, marAccountCrColumn)
				xlsx.SetCellFormula(sheet, marAccountCrColumn, fmt.Sprintf("SUM(%s)", strings.Join(marAccountCr, ":")))
				err = xlsx.SetCellStyle(sheet, marAccountCrColumn, marAccountCrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				aprAccountDrColumn := fmt.Sprintf("K%d", row)
				aprAccountDr = append(aprAccountDr, fmt.Sprintf("K%d", row-1))
				aprAccountDrSum = append(aprAccountDrSum, aprAccountDrColumn)
				xlsx.SetCellFormula(sheet, aprAccountDrColumn, fmt.Sprintf("SUM(%s)", strings.Join(aprAccountDr, ":")))
				err = xlsx.SetCellStyle(sheet, aprAccountDrColumn, aprAccountDrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				aprAccountCrColumn := fmt.Sprintf("L%d", row)
				aprAccountCr = append(aprAccountCr, fmt.Sprintf("L%d", row-1))
				aprAccountCrSum = append(aprAccountCrSum, aprAccountCrColumn)
				xlsx.SetCellFormula(sheet, aprAccountCrColumn, fmt.Sprintf("SUM(%s)", strings.Join(aprAccountCr, ":")))
				err = xlsx.SetCellStyle(sheet, aprAccountCrColumn, aprAccountCrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				mayAccountDrColumn := fmt.Sprintf("M%d", row)
				mayAccountDr = append(mayAccountDr, fmt.Sprintf("M%d", row-1))
				mayAccountDrSum = append(mayAccountDrSum, mayAccountDrColumn)
				xlsx.SetCellFormula(sheet, mayAccountDrColumn, fmt.Sprintf("SUM(%s)", strings.Join(mayAccountDr, ":")))
				err = xlsx.SetCellStyle(sheet, mayAccountDrColumn, mayAccountDrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				mayAccountCrColumn := fmt.Sprintf("N%d", row)
				mayAccountCr = append(mayAccountCr, fmt.Sprintf("N%d", row-1))
				mayAccountCrSum = append(mayAccountCrSum, mayAccountCrColumn)
				xlsx.SetCellFormula(sheet, mayAccountCrColumn, fmt.Sprintf("SUM(%s)", strings.Join(mayAccountCr, ":")))
				err = xlsx.SetCellStyle(sheet, mayAccountCrColumn, mayAccountCrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				junAccountDrColumn := fmt.Sprintf("O%d", row)
				junAccountDr = append(junAccountDr, fmt.Sprintf("O%d", row-1))
				junAccountDrSum = append(junAccountDrSum, junAccountDrColumn)
				xlsx.SetCellFormula(sheet, junAccountDrColumn, fmt.Sprintf("SUM(%s)", strings.Join(junAccountDr, ":")))
				err = xlsx.SetCellStyle(sheet, junAccountDrColumn, junAccountDrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				junAccountCrColumn := fmt.Sprintf("P%d", row)
				junAccountCr = append(junAccountCr, fmt.Sprintf("P%d", row-1))
				junAccountCrSum = append(junAccountCrSum, junAccountCrColumn)
				xlsx.SetCellFormula(sheet, junAccountCrColumn, fmt.Sprintf("SUM(%s)", strings.Join(junAccountCr, ":")))
				err = xlsx.SetCellStyle(sheet, junAccountCrColumn, junAccountCrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				julAccountDrColumn := fmt.Sprintf("Q%d", row)
				julAccountDr = append(julAccountDr, fmt.Sprintf("Q%d", row-1))
				julAccountDrSum = append(julAccountDrSum, julAccountDrColumn)
				xlsx.SetCellFormula(sheet, julAccountDrColumn, fmt.Sprintf("SUM(%s)", strings.Join(julAccountDr, ":")))
				err = xlsx.SetCellStyle(sheet, julAccountDrColumn, julAccountDrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				julAccountCrColumn := fmt.Sprintf("R%d", row)
				julAccountCr = append(julAccountCr, fmt.Sprintf("R%d", row-1))
				julAccountCrSum = append(julAccountCrSum, julAccountCrColumn)
				xlsx.SetCellFormula(sheet, julAccountCrColumn, fmt.Sprintf("SUM(%s)", strings.Join(julAccountCr, ":")))
				err = xlsx.SetCellStyle(sheet, julAccountCrColumn, julAccountCrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				augAccountDrColumn := fmt.Sprintf("S%d", row)
				augAccountDr = append(augAccountDr, fmt.Sprintf("S%d", row-1))
				augAccountDrSum = append(augAccountDrSum, augAccountDrColumn)
				xlsx.SetCellFormula(sheet, augAccountDrColumn, fmt.Sprintf("SUM(%s)", strings.Join(augAccountDr, ":")))
				err = xlsx.SetCellStyle(sheet, augAccountDrColumn, augAccountDrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				augAccountCrColumn := fmt.Sprintf("T%d", row)
				augAccountCr = append(augAccountCr, fmt.Sprintf("T%d", row-1))
				augAccountCrSum = append(augAccountCrSum, augAccountCrColumn)
				xlsx.SetCellFormula(sheet, augAccountCrColumn, fmt.Sprintf("SUM(%s)", strings.Join(augAccountCr, ":")))
				err = xlsx.SetCellStyle(sheet, augAccountCrColumn, augAccountCrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				sepAccountDrColumn := fmt.Sprintf("U%d", row)
				sepAccountDr = append(sepAccountDr, fmt.Sprintf("U%d", row-1))
				sepAccountDrSum = append(sepAccountDrSum, sepAccountDrColumn)
				xlsx.SetCellFormula(sheet, sepAccountDrColumn, fmt.Sprintf("SUM(%s)", strings.Join(sepAccountDr, ":")))
				err = xlsx.SetCellStyle(sheet, sepAccountDrColumn, sepAccountDrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				sepAccountCrColumn := fmt.Sprintf("V%d", row)
				sepAccountCr = append(sepAccountCr, fmt.Sprintf("V%d", row-1))
				sepAccountCrSum = append(sepAccountCrSum, sepAccountCrColumn)
				xlsx.SetCellFormula(sheet, sepAccountCrColumn, fmt.Sprintf("SUM(%s)", strings.Join(sepAccountCr, ":")))
				err = xlsx.SetCellStyle(sheet, sepAccountCrColumn, sepAccountCrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				octAccountDrColumn := fmt.Sprintf("W%d", row)
				octAccountDr = append(octAccountDr, fmt.Sprintf("W%d", row-1))
				octAccountDrSum = append(octAccountDrSum, octAccountDrColumn)
				xlsx.SetCellFormula(sheet, octAccountDrColumn, fmt.Sprintf("SUM(%s)", strings.Join(octAccountDr, ":")))
				err = xlsx.SetCellStyle(sheet, octAccountDrColumn, octAccountDrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				octAccountCrColumn := fmt.Sprintf("X%d", row)
				octAccountCr = append(octAccountCr, fmt.Sprintf("X%d", row-1))
				octAccountCrSum = append(octAccountCrSum, octAccountCrColumn)
				xlsx.SetCellFormula(sheet, octAccountCrColumn, fmt.Sprintf("SUM(%s)", strings.Join(octAccountCr, ":")))
				err = xlsx.SetCellStyle(sheet, octAccountCrColumn, octAccountCrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				novAccountDrColumn := fmt.Sprintf("Y%d", row)
				novAccountDr = append(novAccountDr, fmt.Sprintf("Y%d", row-1))
				novAccountDrSum = append(novAccountDrSum, novAccountDrColumn)
				xlsx.SetCellFormula(sheet, novAccountDrColumn, fmt.Sprintf("SUM(%s)", strings.Join(novAccountDr, ":")))
				err = xlsx.SetCellStyle(sheet, novAccountDrColumn, novAccountDrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				novAccountCrColumn := fmt.Sprintf("Z%d", row)
				novAccountCr = append(novAccountCr, fmt.Sprintf("Z%d", row-1))
				novAccountCrSum = append(novAccountCrSum, novAccountCrColumn)
				xlsx.SetCellFormula(sheet, novAccountCrColumn, fmt.Sprintf("SUM(%s)", strings.Join(novAccountCr, ":")))
				err = xlsx.SetCellStyle(sheet, novAccountCrColumn, novAccountCrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				decAccountDrColumn := fmt.Sprintf("AA%d", row)
				decAccountDr = append(decAccountDr, fmt.Sprintf("AA%d", row-1))
				decAccountDrSum = append(decAccountDrSum, decAccountDrColumn)
				xlsx.SetCellFormula(sheet, decAccountDrColumn, fmt.Sprintf("SUM(%s)", strings.Join(decAccountDr, ":")))
				err = xlsx.SetCellStyle(sheet, decAccountDrColumn, decAccountDrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				decAccountCrColumn := fmt.Sprintf("AB%d", row)
				decAccountCr = append(decAccountCr, fmt.Sprintf("AB%d", row-1))
				decAccountCrSum = append(decAccountCrSum, decAccountCrColumn)
				xlsx.SetCellFormula(sheet, decAccountCrColumn, fmt.Sprintf("SUM(%s)", strings.Join(decAccountCr, ":")))
				err = xlsx.SetCellStyle(sheet, decAccountCrColumn, decAccountCrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				totalAccountDrColumn := fmt.Sprintf("AC%d", row)
				totalAccountDr = append(totalAccountDr, fmt.Sprintf("AC%d", row-1))
				totalAccountDrSum = append(totalAccountDrSum, totalAccountDrColumn)
				xlsx.SetCellFormula(sheet, totalAccountDrColumn, fmt.Sprintf("SUM(%s)", strings.Join(totalAccountDr, ":")))
				err = xlsx.SetCellStyle(sheet, totalAccountDrColumn, totalAccountDrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				totalAccountCrColumn := fmt.Sprintf("AD%d", row)
				totalAccountCr = append(totalAccountCr, fmt.Sprintf("AD%d", row-1))
				totalAccountCrSum = append(totalAccountCrSum, totalAccountCrColumn)
				xlsx.SetCellFormula(sheet, totalAccountCrColumn, fmt.Sprintf("SUM(%s)", strings.Join(totalAccountCr, ":")))
				err = xlsx.SetCellStyle(sheet, totalAccountCrColumn, totalAccountCrColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				totalAccountEndingColumn := fmt.Sprintf("AE%d", row)
				totalAccountEnding = append(totalAccountEnding, fmt.Sprintf("AE%d", row-1))
				xlsx.SetCellFormula(sheet, totalAccountEndingColumn, fmt.Sprintf("SUM(%s)", strings.Join(totalAccountEnding, ":")))
				err = xlsx.SetCellStyle(sheet, totalAccountEndingColumn, totalAccountEndingColumn, sumStyle)
				if err != nil {
					return nil, err
				}
				if groupCode == 4 || groupCode == 9 {
					// รวมรายได้ || รวมรายจ่าย
					resultNetProfitLoss = append(resultNetProfitLoss, totalAccountEndingColumn)
				}
				if groupCode == 1 || groupCode == 2 || groupCode == 3 {
					resultDiffAssetsLiabilitiesOwnerEquity = append(resultDiffAssetsLiabilitiesOwnerEquity, totalAccountEndingColumn)
				}
				groupCode++
				totalAccountForwardDr = []string{}
				totalAccountForwardCr = []string{}
				janAccountDr = []string{}
				janAccountCr = []string{}
				febAccountDr = []string{}
				febAccountCr = []string{}
				marAccountDr = []string{}
				marAccountCr = []string{}
				aprAccountDr = []string{}
				aprAccountCr = []string{}
				mayAccountDr = []string{}
				mayAccountCr = []string{}
				junAccountDr = []string{}
				junAccountCr = []string{}
				julAccountDr = []string{}
				julAccountCr = []string{}
				augAccountDr = []string{}
				augAccountCr = []string{}
				sepAccountDr = []string{}
				sepAccountCr = []string{}
				octAccountDr = []string{}
				octAccountCr = []string{}
				novAccountDr = []string{}
				novAccountCr = []string{}
				decAccountDr = []string{}
				decAccountCr = []string{}
				totalAccountDr = []string{}
				totalAccountCr = []string{}
				totalAccountEnding = []string{}
			}
			row++
		}
		// b, err := json.Marshal(typeList)
		// if err != nil {
		// 	s.logger.Error(err)
		// }
		// s.logger.Warn(string(b))
		// #################### sum result ####################
		err = xlsx.SetRowHeight(sheet, row, 21.75)
		blankColumn := fmt.Sprintf("A%d", row)
		if err != nil {
			return nil, err
		}
		err = xlsx.MergeCell(sheet, blankColumn, fmt.Sprintf("B%d", row))
		if err != nil {
			return nil, err
		}
		err = xlsx.SetCellStyle(sheet, blankColumn, fmt.Sprintf("B%d", row), titleStyle)
		if err != nil {
			return nil, err
		}
		totalAccountForwardDrSumColumn := fmt.Sprintf("C%d", row)
		xlsx.SetCellFormula(sheet, totalAccountForwardDrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(totalAccountForwardDrSum, ":")))
		err = xlsx.SetCellStyle(sheet, totalAccountForwardDrSumColumn, totalAccountForwardDrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		totalAccountForwardCrSumColumn := fmt.Sprintf("D%d", row)
		xlsx.SetCellFormula(sheet, totalAccountForwardCrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(totalAccountForwardCrSum, ":")))
		err = xlsx.SetCellStyle(sheet, totalAccountForwardCrSumColumn, totalAccountForwardCrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		janAccountDrSumColumn := fmt.Sprintf("E%d", row)
		xlsx.SetCellFormula(sheet, janAccountDrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(janAccountDrSum, ":")))
		err = xlsx.SetCellStyle(sheet, janAccountDrSumColumn, janAccountDrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		janAccountCrSumColumn := fmt.Sprintf("F%d", row)
		xlsx.SetCellFormula(sheet, janAccountCrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(janAccountCrSum, ":")))
		err = xlsx.SetCellStyle(sheet, janAccountCrSumColumn, janAccountCrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		febAccountDrSumColumn := fmt.Sprintf("G%d", row)
		xlsx.SetCellFormula(sheet, febAccountDrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(febAccountDrSum, ":")))
		err = xlsx.SetCellStyle(sheet, febAccountDrSumColumn, febAccountDrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		febAccountCrSumColumn := fmt.Sprintf("H%d", row)
		xlsx.SetCellFormula(sheet, febAccountCrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(febAccountCrSum, ":")))
		err = xlsx.SetCellStyle(sheet, febAccountCrSumColumn, febAccountCrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		marAccountDrSumColumn := fmt.Sprintf("I%d", row)
		xlsx.SetCellFormula(sheet, marAccountDrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(marAccountDrSum, ":")))
		err = xlsx.SetCellStyle(sheet, marAccountDrSumColumn, marAccountDrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		marAccountCrSumColumn := fmt.Sprintf("J%d", row)
		xlsx.SetCellFormula(sheet, marAccountCrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(marAccountCrSum, ":")))
		err = xlsx.SetCellStyle(sheet, marAccountCrSumColumn, marAccountCrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		aprAccountDrSumColumn := fmt.Sprintf("K%d", row)
		xlsx.SetCellFormula(sheet, aprAccountDrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(aprAccountDrSum, ":")))
		err = xlsx.SetCellStyle(sheet, aprAccountDrSumColumn, aprAccountDrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		aprAccountCrSumColumn := fmt.Sprintf("L%d", row)
		xlsx.SetCellFormula(sheet, aprAccountCrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(aprAccountCrSum, ":")))
		err = xlsx.SetCellStyle(sheet, aprAccountCrSumColumn, aprAccountCrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		mayAccountDrSumColumn := fmt.Sprintf("M%d", row)
		xlsx.SetCellFormula(sheet, mayAccountDrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(mayAccountDrSum, ":")))
		err = xlsx.SetCellStyle(sheet, mayAccountDrSumColumn, mayAccountDrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		mayAccountCrSumColumn := fmt.Sprintf("N%d", row)
		xlsx.SetCellFormula(sheet, mayAccountCrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(mayAccountCrSum, ":")))
		err = xlsx.SetCellStyle(sheet, mayAccountCrSumColumn, mayAccountCrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		junAccountDrSumColumn := fmt.Sprintf("O%d", row)
		xlsx.SetCellFormula(sheet, junAccountDrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(junAccountDrSum, ":")))
		err = xlsx.SetCellStyle(sheet, junAccountDrSumColumn, junAccountDrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		junAccountCrSumColumn := fmt.Sprintf("P%d", row)
		xlsx.SetCellFormula(sheet, junAccountCrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(junAccountCrSum, ":")))
		err = xlsx.SetCellStyle(sheet, junAccountCrSumColumn, junAccountCrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		julAccountDrSumColumn := fmt.Sprintf("Q%d", row)
		xlsx.SetCellFormula(sheet, julAccountDrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(julAccountDrSum, ":")))
		err = xlsx.SetCellStyle(sheet, julAccountDrSumColumn, julAccountDrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		julAccountCrSumColumn := fmt.Sprintf("R%d", row)
		xlsx.SetCellFormula(sheet, julAccountCrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(julAccountCrSum, ":")))
		err = xlsx.SetCellStyle(sheet, julAccountCrSumColumn, julAccountCrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		augAccountDrSumColumn := fmt.Sprintf("S%d", row)
		xlsx.SetCellFormula(sheet, augAccountDrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(augAccountDrSum, ":")))
		err = xlsx.SetCellStyle(sheet, augAccountDrSumColumn, augAccountDrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		augAccountCrSumColumn := fmt.Sprintf("T%d", row)
		xlsx.SetCellFormula(sheet, augAccountCrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(augAccountCrSum, ":")))
		err = xlsx.SetCellStyle(sheet, augAccountCrSumColumn, augAccountCrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		sepAccountDrSumColumn := fmt.Sprintf("U%d", row)
		xlsx.SetCellFormula(sheet, sepAccountDrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(sepAccountDrSum, ":")))
		err = xlsx.SetCellStyle(sheet, sepAccountDrSumColumn, sepAccountDrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		sepAccountCrSumColumn := fmt.Sprintf("V%d", row)
		xlsx.SetCellFormula(sheet, sepAccountCrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(sepAccountCrSum, ":")))
		err = xlsx.SetCellStyle(sheet, sepAccountCrSumColumn, sepAccountCrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		octAccountDrSumColumn := fmt.Sprintf("W%d", row)
		xlsx.SetCellFormula(sheet, octAccountDrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(octAccountDrSum, ":")))
		err = xlsx.SetCellStyle(sheet, octAccountDrSumColumn, octAccountDrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		octAccountCrSumColumn := fmt.Sprintf("X%d", row)
		xlsx.SetCellFormula(sheet, octAccountCrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(octAccountCrSum, ":")))
		err = xlsx.SetCellStyle(sheet, octAccountCrSumColumn, octAccountCrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		novAccountDrSumColumn := fmt.Sprintf("Y%d", row)
		xlsx.SetCellFormula(sheet, novAccountDrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(novAccountDrSum, ":")))
		err = xlsx.SetCellStyle(sheet, novAccountDrSumColumn, novAccountDrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		novAccountCrSumColumn := fmt.Sprintf("Z%d", row)
		xlsx.SetCellFormula(sheet, novAccountCrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(novAccountCrSum, ":")))
		err = xlsx.SetCellStyle(sheet, novAccountCrSumColumn, novAccountCrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		decAccountDrSumColumn := fmt.Sprintf("AA%d", row)
		xlsx.SetCellFormula(sheet, decAccountDrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(decAccountDrSum, ":")))
		err = xlsx.SetCellStyle(sheet, decAccountDrSumColumn, decAccountDrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		decAccountCrSumColumn := fmt.Sprintf("AB%d", row)
		xlsx.SetCellFormula(sheet, decAccountCrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(decAccountCrSum, ":")))
		err = xlsx.SetCellStyle(sheet, decAccountCrSumColumn, decAccountCrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		totalAccountDrSumColumn := fmt.Sprintf("AC%d", row)
		xlsx.SetCellFormula(sheet, totalAccountDrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(totalAccountDrSum, ":")))
		err = xlsx.SetCellStyle(sheet, totalAccountDrSumColumn, totalAccountDrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		totalAccountCrSumColumn := fmt.Sprintf("AD%d", row)
		xlsx.SetCellFormula(sheet, totalAccountCrSumColumn, fmt.Sprintf("SUM(%s)", strings.Join(totalAccountCrSum, ":")))
		err = xlsx.SetCellStyle(sheet, totalAccountCrSumColumn, totalAccountCrSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		endingSumColumn := fmt.Sprintf("AE%d", row)
		xlsx.SetCellFormula(sheet, endingSumColumn, fmt.Sprintf("%s-%s", fmt.Sprintf("AC%d", row), fmt.Sprintf("AD%d", row)))
		err = xlsx.SetCellStyle(sheet, endingSumColumn, endingSumColumn, sumStyle)
		if err != nil {
			return nil, err
		}

		// #################### กำไร (ขาดทุน) สุทธิ ####################
		row++
		err = xlsx.SetRowHeight(sheet, row, 21.75)
		netProfitLossColumn := fmt.Sprintf("A%d", row)
		if err != nil {
			return nil, err
		}
		err = xlsx.MergeCell(sheet, netProfitLossColumn, fmt.Sprintf("B%d", row))
		if err != nil {
			return nil, err
		}
		xlsx.SetCellValue(sheet, netProfitLossColumn, "กำไร (ขาดทุน) สุทธิ")
		err = xlsx.SetCellStyle(sheet, netProfitLossColumn, fmt.Sprintf("B%d", row), titleStyle)
		if err != nil {
			return nil, err
		}
		totalForwardDrSumNetProfitLossColumn := fmt.Sprintf("C%d", row)
		err = xlsx.SetCellStyle(sheet, totalForwardDrSumNetProfitLossColumn, totalForwardDrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		totalForwardCrSumNetProfitLossColumn := fmt.Sprintf("D%d", row)
		err = xlsx.SetCellStyle(sheet, totalForwardCrSumNetProfitLossColumn, totalForwardCrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		janDrSumNetProfitLossColumn := fmt.Sprintf("E%d", row)
		err = xlsx.SetCellStyle(sheet, janDrSumNetProfitLossColumn, janDrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		janCrSumNetProfitLossColumn := fmt.Sprintf("F%d", row)
		xlsx.SetCellFormula(sheet, janCrSumNetProfitLossColumn, fmt.Sprintf("%s-%s", fmt.Sprintf("E%d", row-1), fmt.Sprintf("F%d", row-1)))
		err = xlsx.SetCellStyle(sheet, janCrSumNetProfitLossColumn, janCrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		febDrSumNetProfitLossColumn := fmt.Sprintf("G%d", row)
		err = xlsx.SetCellStyle(sheet, febDrSumNetProfitLossColumn, febDrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		febCrSumNetProfitLossColumn := fmt.Sprintf("H%d", row)
		xlsx.SetCellFormula(sheet, febCrSumNetProfitLossColumn, fmt.Sprintf("%s-%s", fmt.Sprintf("G%d", row-1), fmt.Sprintf("H%d", row-1)))
		err = xlsx.SetCellStyle(sheet, febCrSumNetProfitLossColumn, febCrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		marDrSumNetProfitLossColumn := fmt.Sprintf("I%d", row)
		err = xlsx.SetCellStyle(sheet, marDrSumNetProfitLossColumn, marDrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		marCrSumNetProfitLossColumn := fmt.Sprintf("J%d", row)
		xlsx.SetCellFormula(sheet, marCrSumNetProfitLossColumn, fmt.Sprintf("%s-%s", fmt.Sprintf("I%d", row-1), fmt.Sprintf("J%d", row-1)))
		err = xlsx.SetCellStyle(sheet, marCrSumNetProfitLossColumn, marCrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		aprDrSumNetProfitLossColumn := fmt.Sprintf("K%d", row)
		err = xlsx.SetCellStyle(sheet, aprDrSumNetProfitLossColumn, aprDrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		aprCrSumNetProfitLossColumn := fmt.Sprintf("L%d", row)
		xlsx.SetCellFormula(sheet, aprCrSumNetProfitLossColumn, fmt.Sprintf("%s-%s", fmt.Sprintf("K%d", row-1), fmt.Sprintf("L%d", row-1)))
		err = xlsx.SetCellStyle(sheet, aprCrSumNetProfitLossColumn, aprCrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		mayDrSumNetProfitLossColumn := fmt.Sprintf("M%d", row)
		err = xlsx.SetCellStyle(sheet, mayDrSumNetProfitLossColumn, mayDrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		mayCrSumNetProfitLossColumn := fmt.Sprintf("N%d", row)
		xlsx.SetCellFormula(sheet, mayCrSumNetProfitLossColumn, fmt.Sprintf("%s-%s", fmt.Sprintf("M%d", row-1), fmt.Sprintf("N%d", row-1)))
		err = xlsx.SetCellStyle(sheet, mayCrSumNetProfitLossColumn, mayCrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		junDrSumNetProfitLossColumn := fmt.Sprintf("O%d", row)
		err = xlsx.SetCellStyle(sheet, junDrSumNetProfitLossColumn, junDrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		junCrSumNetProfitLossColumn := fmt.Sprintf("P%d", row)
		xlsx.SetCellFormula(sheet, junCrSumNetProfitLossColumn, fmt.Sprintf("%s-%s", fmt.Sprintf("O%d", row-1), fmt.Sprintf("P%d", row-1)))
		err = xlsx.SetCellStyle(sheet, junCrSumNetProfitLossColumn, junCrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		julDrSumNetProfitLossColumn := fmt.Sprintf("Q%d", row)
		err = xlsx.SetCellStyle(sheet, julDrSumNetProfitLossColumn, julDrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		julCrSumNetProfitLossColumn := fmt.Sprintf("R%d", row)
		xlsx.SetCellFormula(sheet, julCrSumNetProfitLossColumn, fmt.Sprintf("%s-%s", fmt.Sprintf("Q%d", row-1), fmt.Sprintf("R%d", row-1)))
		err = xlsx.SetCellStyle(sheet, julCrSumNetProfitLossColumn, julCrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		augDrSumNetProfitLossColumn := fmt.Sprintf("S%d", row)
		err = xlsx.SetCellStyle(sheet, augDrSumNetProfitLossColumn, augDrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		augCrSumNetProfitLossColumn := fmt.Sprintf("T%d", row)
		xlsx.SetCellFormula(sheet, augCrSumNetProfitLossColumn, fmt.Sprintf("%s-%s", fmt.Sprintf("S%d", row-1), fmt.Sprintf("T%d", row-1)))
		err = xlsx.SetCellStyle(sheet, augCrSumNetProfitLossColumn, augCrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		sepDrSumNetProfitLossColumn := fmt.Sprintf("U%d", row)
		err = xlsx.SetCellStyle(sheet, sepDrSumNetProfitLossColumn, sepDrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		sepCrSumNetProfitLossColumn := fmt.Sprintf("V%d", row)
		xlsx.SetCellFormula(sheet, sepCrSumNetProfitLossColumn, fmt.Sprintf("%s-%s", fmt.Sprintf("U%d", row-1), fmt.Sprintf("V%d", row-1)))
		err = xlsx.SetCellStyle(sheet, sepCrSumNetProfitLossColumn, sepCrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		octDrSumNetProfitLossColumn := fmt.Sprintf("W%d", row)
		err = xlsx.SetCellStyle(sheet, octDrSumNetProfitLossColumn, octDrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		octCrSumNetProfitLossColumn := fmt.Sprintf("X%d", row)
		xlsx.SetCellFormula(sheet, octCrSumNetProfitLossColumn, fmt.Sprintf("%s-%s", fmt.Sprintf("W%d", row-1), fmt.Sprintf("X%d", row-1)))
		err = xlsx.SetCellStyle(sheet, octCrSumNetProfitLossColumn, octCrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		novDrSumNetProfitLossColumn := fmt.Sprintf("Y%d", row)
		err = xlsx.SetCellStyle(sheet, novDrSumNetProfitLossColumn, novDrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		novCrSumNetProfitLossColumn := fmt.Sprintf("Z%d", row)
		xlsx.SetCellFormula(sheet, novCrSumNetProfitLossColumn, fmt.Sprintf("%s-%s", fmt.Sprintf("Y%d", row-1), fmt.Sprintf("Z%d", row-1)))
		err = xlsx.SetCellStyle(sheet, novCrSumNetProfitLossColumn, novCrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		decDrSumNetProfitLossColumn := fmt.Sprintf("AA%d", row)
		err = xlsx.SetCellStyle(sheet, decDrSumNetProfitLossColumn, decDrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		decCrSumNetProfitLossColumn := fmt.Sprintf("AB%d", row)
		xlsx.SetCellFormula(sheet, decCrSumNetProfitLossColumn, fmt.Sprintf("%s-%s", fmt.Sprintf("AA%d", row-1), fmt.Sprintf("AB%d", row-1)))
		err = xlsx.SetCellStyle(sheet, decCrSumNetProfitLossColumn, decCrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		totalDrSumNetProfitLossColumn := fmt.Sprintf("AC%d", row)
		err = xlsx.SetCellStyle(sheet, totalDrSumNetProfitLossColumn, totalDrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		totalCrSumNetProfitLossColumn := fmt.Sprintf("AD%d", row)
		xlsx.SetCellFormula(sheet, totalCrSumNetProfitLossColumn, fmt.Sprintf("%s-%s", fmt.Sprintf("AC%d", row-1), fmt.Sprintf("AD%d", row-1)))
		err = xlsx.SetCellStyle(sheet, totalCrSumNetProfitLossColumn, totalCrSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		endingSumNetProfitLossColumn := fmt.Sprintf("AE%d", row)
		xlsx.SetCellFormula(sheet, endingSumNetProfitLossColumn, fmt.Sprintf("-%s", strings.Join(resultNetProfitLoss, "-")))
		err = xlsx.SetCellStyle(sheet, endingSumNetProfitLossColumn, endingSumNetProfitLossColumn, sumStyle)
		if err != nil {
			return nil, err
		}
		resultDifference = append(resultDifference, endingSumNetProfitLossColumn)
		// #################### ผลต่างระหว่างสินทรัพย์กับหนี้สินและส่วนของเจ้าของ ####################
		row++
		titleDiffStyle, err := xlsx.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold:   true,
				Family: "TH Sarabun New",
				Size:   9,
				Color:  "FF0006",
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "bottom",
			},
			Border: []excelize.Border{
				{Type: "top", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "left", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"FFECDC"}, Pattern: 1},
		})
		if err != nil {
			return nil, err
		}
		titleRedStyle, err := xlsx.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold:   true,
				Family: "TH Sarabun New",
				Size:   9,
				Color:  "FF0006",
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "bottom",
			},
			Border: []excelize.Border{
				{Type: "top", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "left", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"FFECDC"}, Pattern: 1},
		})
		if err != nil {
			return nil, err
		}
		err = xlsx.SetRowHeight(sheet, row, 21.75)
		diffAssetsLiabilitiesOwnerEquityColumn := fmt.Sprintf("A%d", row)
		if err != nil {
			return nil, err
		}
		err = xlsx.MergeCell(sheet, diffAssetsLiabilitiesOwnerEquityColumn, fmt.Sprintf("B%d", row))
		if err != nil {
			return nil, err
		}
		xlsx.SetCellValue(sheet, diffAssetsLiabilitiesOwnerEquityColumn, "ผลต่างระหว่างสินทรัพย์กับหนี้สินและส่วนของเจ้าของ")
		err = xlsx.SetCellStyle(sheet, diffAssetsLiabilitiesOwnerEquityColumn, diffAssetsLiabilitiesOwnerEquityColumn, titleDiffStyle)
		if err != nil {
			return nil, err
		}
		err = xlsx.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("AE%d", row), titleRedStyle)
		if err != nil {
			return nil, err
		}
		xlsx.SetCellFormula(sheet, fmt.Sprintf("AE%d", row), strings.Join(resultDiffAssetsLiabilitiesOwnerEquity, "+"))
		resultDifference = append(resultDifference, fmt.Sprintf("AE%d", row))
		// #################### ผลต่าง ####################
		row++
		err = xlsx.SetRowHeight(sheet, row, 21.75)
		differenceColumn := fmt.Sprintf("A%d", row)
		if err != nil {
			return nil, err
		}
		err = xlsx.MergeCell(sheet, differenceColumn, fmt.Sprintf("B%d", row))
		if err != nil {
			return nil, err
		}
		xlsx.SetCellValue(sheet, differenceColumn, "ผลต่าง")
		err = xlsx.SetCellStyle(sheet, differenceColumn, differenceColumn, titleDiffStyle)
		if err != nil {
			return nil, err
		}
		err = xlsx.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("AE%d", row), titleRedStyle)
		if err != nil {
			return nil, err
		}
		xlsx.SetCellFormula(sheet, fmt.Sprintf("AE%d", row), strings.Join(resultDifference, "-"))

		// typeList
		sheet = "แยกประเภท"
		headAccountNameStyle, err := xlsx.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				// Bold:   true,
				Family: "TH Sarabun New",
				Size:   14,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "bottom",
			},
			// Border: []excelize.Border{
			// 	{Type: "top", Color: "000000", Style: 1},
			// 	{Type: "right", Color: "000000", Style: 1},
			// 	{Type: "left", Color: "000000", Style: 1},
			// 	{Type: "bottom", Color: "000000", Style: 1},
			// },
			// Fill: excelize.Fill{Type: "pattern", Color: []string{"F0F4DC"}, Pattern: 1},
		})
		if err != nil {
			return nil, err
		}
		headAccountCodeStyle, err := xlsx.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				// Bold:   true,
				Family: "TH Sarabun New",
				Size:   14,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "left",
				Vertical:   "bottom",
			},
			// Border: []excelize.Border{
			// 	{Type: "top", Color: "000000", Style: 1},
			// 	{Type: "right", Color: "000000", Style: 1},
			// 	{Type: "left", Color: "000000", Style: 1},
			// 	{Type: "bottom", Color: "000000", Style: 1},
			// },
			// Fill: excelize.Fill{Type: "pattern", Color: []string{"F0F4DC"}, Pattern: 1},
		})
		if err != nil {
			return nil, err
		}
		headAccountStyle, err := xlsx.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				// Bold:   true,
				Family: "TH Sarabun New",
				Size:   14,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				{Type: "top", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "left", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"F0F4DC"}, Pattern: 1},
		})
		if err != nil {
			return nil, err
		}
		subjectAccountStyle, err := xlsx.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				// Bold:   true,
				Family: "TH Sarabun New",
				Size:   14,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				// {Type: "top", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "left", Color: "000000", Style: 1},
				// {Type: "bottom", Color: "000000", Style: 1},
			},
		})
		if err != nil {
			return nil, err
		}
		detailAccountStyle, err := xlsx.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				// Bold:   true,
				Family: "TH Sarabun New",
				Size:   14,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "left",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				// {Type: "top", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "left", Color: "000000", Style: 1},
				// {Type: "bottom", Color: "000000", Style: 1},
			},
		})
		if err != nil {
			return nil, err
		}
		priceAccountStyle, err := xlsx.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				// Bold:   true,
				Family: "TH Sarabun New",
				Size:   14,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "right",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				// {Type: "top", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "left", Color: "000000", Style: 1},
				// {Type: "bottom", Color: "000000", Style: 1},
			},
			CustomNumFmt: &fm,
		})
		if err != nil {
			return nil, err
		}
		endingAccountStyle, err := xlsx.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				// Bold:   true,
				Family: "TH Sarabun New",
				Size:   14,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "left",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				{Type: "top", Color: "000000", Style: 1},
				// {Type: "right", Color: "000000", Style: 1},
				// {Type: "left", Color: "000000", Style: 1},
				// {Type: "bottom", Color: "000000", Style: 1},
			},
		})
		if err != nil {
			return nil, err
		}
		emptyEndingAccountStyle, err := xlsx.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				// Bold:   true,
				Family: "TH Sarabun New",
				Size:   14,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "left",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				{Type: "top", Color: "000000", Style: 1},
				// {Type: "right", Color: "000000", Style: 1},
				// {Type: "left", Color: "000000", Style: 1},
				// {Type: "bottom", Color: "000000", Style: 1},
			},
		})
		if err != nil {
			return nil, err
		}
		row = 1
		for i := 0; i < len(typeList); i++ {
			t := typeList[i]
			accNameColumn := fmt.Sprintf("A%d", row)
			err = xlsx.MergeCell(sheet, accNameColumn, fmt.Sprintf("F%d", row))
			if err != nil {
				return nil, err
			}
			xlsx.SetCellValue(sheet, accNameColumn, t.Name)
			err = xlsx.SetCellStyle(sheet, accNameColumn, accNameColumn, headAccountNameStyle)
			if err != nil {
				return nil, err
			}
			accCodeColumn := fmt.Sprintf("G%d", row)
			xlsx.SetCellValue(sheet, accCodeColumn, fmt.Sprintf("เลขที่ %s", t.Code))
			err = xlsx.SetCellStyle(sheet, accCodeColumn, accCodeColumn, headAccountCodeStyle)
			if err != nil {
				return nil, err
			}
			row++
			yearColumn := fmt.Sprintf("A%d", row)
			err = xlsx.MergeCell(sheet, yearColumn, fmt.Sprintf("B%d", row))
			if err != nil {
				return nil, err
			}
			xlsx.SetCellValue(sheet, yearColumn, strconv.Itoa(yearInt+543))
			err = xlsx.SetCellStyle(sheet, yearColumn, fmt.Sprintf("B%d", row), headAccountStyle)
			if err != nil {
				return nil, err
			}
			listColumn := fmt.Sprintf("C%d", row)
			xlsx.SetCellValue(sheet, listColumn, "รายการ")
			err = xlsx.MergeCell(sheet, listColumn, fmt.Sprintf("C%d", row+1))
			if err != nil {
				return nil, err
			}
			err = xlsx.SetCellStyle(sheet, listColumn, fmt.Sprintf("C%d", row+1), headAccountStyle)
			if err != nil {
				return nil, err
			}
			frontAccountColumn := fmt.Sprintf("D%d", row)
			xlsx.SetCellValue(sheet, frontAccountColumn, "หน้าบัญชี")
			err = xlsx.MergeCell(sheet, frontAccountColumn, fmt.Sprintf("D%d", row+1))
			if err != nil {
				return nil, err
			}
			err = xlsx.SetCellStyle(sheet, frontAccountColumn, fmt.Sprintf("D%d", row+1), headAccountStyle)
			if err != nil {
				return nil, err
			}
			drColumn := fmt.Sprintf("E%d", row)
			xlsx.SetCellValue(sheet, drColumn, "เดบิต")
			err = xlsx.SetCellStyle(sheet, drColumn, drColumn, headAccountStyle)
			if err != nil {
				return nil, err
			}
			crColumn := fmt.Sprintf("F%d", row)
			xlsx.SetCellValue(sheet, crColumn, "เครดิต")
			err = xlsx.SetCellStyle(sheet, crColumn, crColumn, headAccountStyle)
			if err != nil {
				return nil, err
			}
			totalAmountColumn := fmt.Sprintf("G%d", row)
			xlsx.SetCellValue(sheet, totalAmountColumn, "ยอดคงเหลือ")
			err = xlsx.SetCellStyle(sheet, totalAmountColumn, totalAmountColumn, headAccountStyle)
			if err != nil {
				return nil, err
			}
			row++
			monthColumn := fmt.Sprintf("A%d", row)
			xlsx.SetCellValue(sheet, monthColumn, "เดือน")
			err = xlsx.SetCellStyle(sheet, monthColumn, monthColumn, headAccountStyle)
			if err != nil {
				return nil, err
			}
			dateColumn := fmt.Sprintf("B%d", row)
			xlsx.SetCellValue(sheet, dateColumn, "วันที่")
			err = xlsx.SetCellStyle(sheet, dateColumn, dateColumn, headAccountStyle)
			if err != nil {
				return nil, err
			}
			bahtDrColumn := fmt.Sprintf("E%d", row)
			xlsx.SetCellValue(sheet, bahtDrColumn, "บาท")
			err = xlsx.SetCellStyle(sheet, bahtDrColumn, bahtDrColumn, headAccountStyle)
			if err != nil {
				return nil, err
			}
			bahtCrColumn := fmt.Sprintf("F%d", row)
			xlsx.SetCellValue(sheet, bahtCrColumn, "บาท")
			err = xlsx.SetCellStyle(sheet, bahtCrColumn, bahtCrColumn, headAccountStyle)
			if err != nil {
				return nil, err
			}
			bahtTotalAmountColumn := fmt.Sprintf("G%d", row)
			xlsx.SetCellValue(sheet, bahtTotalAmountColumn, "บาท")
			err = xlsx.SetCellStyle(sheet, bahtTotalAmountColumn, bahtTotalAmountColumn, headAccountStyle)
			if err != nil {
				return nil, err
			}
			row++
			subjectColumnStart := ""
			detailColumnStart := ""
			frontAccountColumnStart := ""
			amountColumnStart := ""
			startEndingColumn := ""
			endEndingColumn := ""
			if fwMap[t.Code] != "" {
				if subjectColumnStart == "" {
					subjectColumnStart = fmt.Sprintf("A%d", row)
				}
				eDateColumn := fmt.Sprintf("B%d", row)
				xlsx.SetCellValue(sheet, eDateColumn, 1)
				detailColumn := fmt.Sprintf("C%d", row)
				if detailColumnStart == "" {
					detailColumnStart = detailColumn
				}
				xlsx.SetCellValue(sheet, detailColumn, "ยอดยกมา")
				eFrontAccountColumn := fmt.Sprintf("D%d", row)
				if frontAccountColumnStart == "" {
					frontAccountColumnStart = eFrontAccountColumn
				}
				// xlsx.SetCellValue(sheet, eFrontAccountColumn, a.Number)
				eDrColumn := fmt.Sprintf("E%d", row)
				if amountColumnStart == "" {
					amountColumnStart = eDrColumn
				}
				eCrColumn := fmt.Sprintf("F%d", row)
				eTotalColumn := fmt.Sprintf("G%d", row)
				accountFirstNo, _ := strconv.Atoi(t.Code)
				if accountFirstNo > 1 || (t.Code == "1420-01" ||
					t.Code == "1420-02" ||
					t.Code == "1420-03" ||
					t.Code == "1420-04") {
					xlsx.SetCellFormula(sheet, eCrColumn, fwMap[t.Code])
				} else {
					xlsx.SetCellFormula(sheet, eDrColumn, fwMap[t.Code])
				}
				xlsx.SetCellFormula(sheet, eTotalColumn, fmt.Sprintf("%s-%s", fmt.Sprintf("E%d", row), fmt.Sprintf("F%d", row)))
				row++
			}
			if t.MonthDetail != nil {
				for i := 0; i < len(t.MonthDetail); i++ {
					m := t.MonthDetail[i]
					shortMonthColumn := fmt.Sprintf("A%d", row)
					if subjectColumnStart == "" {
						subjectColumnStart = shortMonthColumn
						xlsx.SetCellValue(sheet, shortMonthColumn, m.Month)
					} else {
						xlsx.SetCellValue(sheet, fmt.Sprintf("A%d", row-1), m.Month)
					}
					if m.AccountDetail != nil {
						for j := 0; j < len(m.AccountDetail); j++ {
							a := m.AccountDetail[j]
							eDateColumn := fmt.Sprintf("B%d", row)
							xlsx.SetCellValue(sheet, eDateColumn, a.Date)
							detailColumn := fmt.Sprintf("C%d", row)
							if detailColumnStart == "" {
								detailColumnStart = detailColumn
							}
							xlsx.SetCellValue(sheet, detailColumn, a.Description)
							eFrontAccountColumn := fmt.Sprintf("D%d", row)
							if frontAccountColumnStart == "" {
								frontAccountColumnStart = eFrontAccountColumn
							}
							xlsx.SetCellValue(sheet, eFrontAccountColumn, a.Number)
							eDrColumn := fmt.Sprintf("E%d", row)
							if amountColumnStart == "" {
								amountColumnStart = eDrColumn
							}
							xlsx.SetCellValue(sheet, eDrColumn, a.AmountDr)
							eCrColumn := fmt.Sprintf("F%d", row)
							xlsx.SetCellValue(sheet, eCrColumn, a.AmountCr)
							eTotalColumn := fmt.Sprintf("G%d", row)

							if j == 0 && fwMap[t.Code] == "" {
								xlsx.SetCellFormula(sheet, eTotalColumn, fmt.Sprintf("%s-%s", fmt.Sprintf("E%d", row), fmt.Sprintf("F%d", row)))
							} else {
								xlsx.SetCellFormula(sheet, eTotalColumn, fmt.Sprintf("%s+%s-%s", fmt.Sprintf("G%d", row-1), fmt.Sprintf("E%d", row), fmt.Sprintf("F%d", row)))
							}
							if len(m.AccountDetail) == j+1 {
								row++
							}
							err = xlsx.SetCellStyle(sheet, subjectColumnStart, fmt.Sprintf("B%d", row), subjectAccountStyle)
							if err != nil {
								return nil, err
							}
							err = xlsx.SetCellStyle(sheet, detailColumnStart, fmt.Sprintf("C%d", row), detailAccountStyle)
							if err != nil {
								return nil, err
							}
							err = xlsx.SetCellStyle(sheet, frontAccountColumnStart, fmt.Sprintf("D%d", row), subjectAccountStyle)
							if err != nil {
								return nil, err
							}
							err = xlsx.SetCellStyle(sheet, amountColumnStart, fmt.Sprintf("G%d", row), priceAccountStyle)
							if err != nil {
								return nil, err
							}
							row++
							startEndingColumn = fmt.Sprintf("A%d", row)
							endEndingColumn = fmt.Sprintf("G%d", row)
						}
					}
				}
				err = xlsx.SetCellStyle(sheet, startEndingColumn, endEndingColumn, endingAccountStyle)
				if err != nil {
					return nil, err
				}
			} else {
				if subjectColumnStart == "" {
					subjectColumnStart = fmt.Sprintf("A%d", row)
				} else {
					subjectColumnStart = fmt.Sprintf("A%d", row-1)
				}
				err = xlsx.SetCellStyle(sheet, subjectColumnStart, fmt.Sprintf("B%d", row+1), subjectAccountStyle)
				if err != nil {
					return nil, err
				}
				if detailColumnStart == "" {
					detailColumnStart = fmt.Sprintf("C%d", row)
				} else {
					detailColumnStart = fmt.Sprintf("C%d", row-1)
				}
				err = xlsx.SetCellStyle(sheet, detailColumnStart, fmt.Sprintf("C%d", row+1), detailAccountStyle)
				if err != nil {
					return nil, err
				}
				if frontAccountColumnStart == "" {
					frontAccountColumnStart = fmt.Sprintf("D%d", row)
				} else {
					frontAccountColumnStart = fmt.Sprintf("D%d", row-1)
				}
				err = xlsx.SetCellStyle(sheet, frontAccountColumnStart, fmt.Sprintf("D%d", row+1), subjectAccountStyle)
				if err != nil {
					return nil, err
				}
				if amountColumnStart == "" {
					amountColumnStart = fmt.Sprintf("E%d", row)
				} else {
					amountColumnStart = fmt.Sprintf("E%d", row-1)
				}
				err = xlsx.SetCellStyle(sheet, amountColumnStart, fmt.Sprintf("G%d", row+1), priceAccountStyle)
				if err != nil {
					return nil, err
				}
				err = xlsx.SetCellStyle(sheet, fmt.Sprintf("A%d", row+2), fmt.Sprintf("G%d", row+2), emptyEndingAccountStyle)
				if err != nil {
					return nil, err
				}
				row += 2
			}
			row++
		}
	}
	return xlsx, nil
}
