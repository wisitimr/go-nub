package daybook

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Daybook struct {
	Id              primitive.ObjectID   `bson:"_id" json:"id"`
	Number          string               `bson:"number" json:"number"`
	Invoice         string               `bson:"invoice" json:"invoice"`
	Document        primitive.ObjectID   `bson:"document" json:"document"`
	TransactionDate time.Time            `bson:"transactionDate" json:"transactionDate"`
	Company         primitive.ObjectID   `bson:"company" json:"company"`
	Customer        *primitive.ObjectID  `bson:"customer" json:"customer"`
	Supplier        *primitive.ObjectID  `bson:"supplier" json:"supplier"`
	PaymentMethod   *primitive.ObjectID  `bson:"paymentMethod" json:"paymentMethod"`
	DaybookDetails  []primitive.ObjectID `bson:"daybookDetails" json:"daybookDetails"`
	CreatedBy       primitive.ObjectID   `bson:"createdBy" json:"createdBy"`
	CreatedAt       time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedBy       primitive.ObjectID   `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt       time.Time            `bson:"updatedAt" json:"updatedAt"`
}

type DaybookPayload struct {
	Id              primitive.ObjectID     `bson:"_id" json:"id"`
	Number          string                 `bson:"number" json:"number"`
	Invoice         string                 `bson:"invoice" json:"invoice"`
	Document        primitive.ObjectID     `bson:"document" json:"document"`
	TransactionDate time.Time              `bson:"transactionDate" json:"transactionDate"`
	Company         primitive.ObjectID     `bson:"company" json:"company"`
	Customer        *primitive.ObjectID    `bson:"customer" json:"customer"`
	Supplier        *primitive.ObjectID    `bson:"supplier" json:"supplier"`
	PaymentMethod   *primitive.ObjectID    `bson:"paymentMethod" json:"paymentMethod"`
	DaybookDetails  []DaybookPayloadDetail `bson:"daybookDetails" json:"daybookDetails"`
	CreatedBy       primitive.ObjectID     `bson:"createdBy" json:"createdBy"`
	CreatedAt       time.Time              `bson:"createdAt" json:"createdAt"`
	UpdatedBy       primitive.ObjectID     `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt       time.Time              `bson:"updatedAt" json:"updatedAt"`
}

type DaybookList struct {
	Id              primitive.ObjectID   `bson:"_id" json:"id"`
	Number          string               `bson:"number" json:"number"`
	Invoice         string               `bson:"invoice" json:"invoice"`
	Document        Document             `bson:"document" json:"document"`
	TransactionDate time.Time            `bson:"transactionDate" json:"transactionDate"`
	Company         Company              `bson:"company" json:"company"`
	Customer        *Customer            `bson:"customer" json:"customer"`
	Supplier        *Supplier            `bson:"supplier" json:"supplier"`
	PaymentMethod   *primitive.ObjectID  `bson:"paymentMethod" json:"paymentMethod"`
	DaybookDetails  []primitive.ObjectID `bson:"daybookDetails" json:"daybookDetails"`
	CreatedBy       primitive.ObjectID   `bson:"createdBy" json:"createdBy"`
	CreatedAt       time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedBy       primitive.ObjectID   `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt       time.Time            `bson:"updatedAt" json:"updatedAt"`
}

type DaybookPayloadDetail struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Type      string             `bson:"type" json:"type"`
	Amount    float64            `bson:"amount" json:"amount"`
	Account   primitive.ObjectID `bson:"account" json:"account"`
	CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type DaybookResponse struct {
	Id                primitive.ObjectID  `bson:"_id" json:"id"`
	Number            string              `bson:"number" json:"number"`
	Invoice           string              `bson:"invoice" json:"invoice"`
	Document          primitive.ObjectID  `bson:"document" json:"document"`
	TransactionDate   time.Time           `bson:"transactionDate" json:"transactionDate"`
	Company           primitive.ObjectID  `bson:"company" json:"company"`
	Customer          *primitive.ObjectID `bson:"customer" json:"customer"`
	Supplier          *primitive.ObjectID `bson:"supplier" json:"supplier"`
	PaymentMethod     *primitive.ObjectID `bson:"paymentMethod" json:"paymentMethod"`
	DaybookDetails    []DaybookDetail     `bson:"daybookDetails" json:"daybookDetails"`
	DebitTotalAmount  float64             `bson:"debitTotalAmount" json:"debitTotalAmount"`
	CreditTotalAmount float64             `bson:"creditTotalAmount" json:"creditTotalAmount"`
	CreatedBy         primitive.ObjectID  `bson:"createdBy" json:"createdBy"`
	CreatedAt         time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedBy         primitive.ObjectID  `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt         time.Time           `bson:"updatedAt" json:"updatedAt"`
}

type DaybookExpand struct {
	Id                primitive.ObjectID `bson:"_id" json:"id"`
	Number            string             `bson:"number" json:"number"`
	Invoice           string             `bson:"invoice" json:"invoice"`
	Document          Document           `bson:"document" json:"document"`
	TransactionDate   time.Time          `bson:"transactionDate" json:"transactionDate"`
	Company           Company            `bson:"company" json:"company"`
	Customer          *Customer          `bson:"customer" json:"customer"`
	Supplier          *Supplier          `bson:"supplier" json:"supplier"`
	PaymentMethod     primitive.ObjectID `bson:"paymentMethod" json:"paymentMethod"`
	DaybookDetails    []DaybookDetail    `bson:"daybookDetails" json:"daybookDetails"`
	DebitTotalAmount  float64            `bson:"debitTotalAmount" json:"debitTotalAmount"`
	CreditTotalAmount float64            `bson:"creditTotalAmount" json:"creditTotalAmount"`
	CreatedBy         primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt         time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedBy         primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt         time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type DaybookFinancialStatement struct {
	// Id             AccountFinState        `bson:"_id" json:"account"`
	Id             string                 `bson:"id" json:"id"`
	Code           string                 `bson:"code" json:"code"`
	Name           string                 `bson:"name" json:"name"`
	Description    string                 `bson:"description" json:"description"`
	Type           string                 `bson:"type" json:"type"`
	Company        Company                `bson:"company" json:"company"`
	DaybookDetails []AccountDaybookDetail `bson:"daybookDetails" json:"daybookDetails"`
	// Number            string             `bson:"number" json:"number"`
	// Invoice           string             `bson:"invoice" json:"invoice"`
	// Document          Document           `bson:"document" json:"document"`
	// TransactionDate   time.Time          `bson:"transactionDate" json:"transactionDate"`
	// Company           Company            `bson:"company" json:"company"`
	// Customer          *Customer          `bson:"customer" json:"customer"`
	// Supplier          *Supplier          `bson:"supplier" json:"supplier"`
	// PaymentMethod     primitive.ObjectID `bson:"paymentMethod" json:"paymentMethod"`
	// DaybookDetails    []DaybookDetail   `bson:"daybookDetails" json:"daybookDetails"`
	// DebitTotalAmount  float64            `bson:"debitTotalAmount" json:"debitTotalAmount"`
	// CreditTotalAmount float64            `bson:"creditTotalAmount" json:"creditTotalAmount"`
	// CreatedBy         primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	// CreatedAt         time.Time          `bson:"createdAt" json:"createdAt"`
	// UpdatedBy         primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
	// UpdatedAt         time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Customer struct {
	Id      primitive.ObjectID `bson:"_id" json:"id"`
	Code    string             `bson:"code" json:"code"`
	Name    string             `bson:"name" json:"name"`
	Address string             `bson:"address" json:"address"`
	Phone   string             `bson:"phone" json:"phone"`
	Contact string             `bson:"contact" json:"contact"`
}

type Supplier struct {
	Id      primitive.ObjectID `bson:"_id" json:"id"`
	Code    string             `bson:"code" json:"code"`
	Name    string             `bson:"name" json:"name"`
	Address string             `bson:"address" json:"address"`
	Phone   string             `bson:"phone" json:"phone"`
	Contact string             `bson:"contact" json:"contact"`
}

type Company struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Address     string             `bson:"address" json:"address"`
	Phone       string             `bson:"phone" json:"phone"`
	Contact     string             `bson:"contact" json:"contact"`
}

type OutDaybookDetails struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Type      string             `bson:"type" json:"type"`
	Amount    float64            `bson:"amount" json:"amount"`
	Account   Account            `bson:"account" json:"account"`
	CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type DaybookDetail struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Type      string             `bson:"type" json:"type"`
	Amount    float64            `bson:"amount" json:"amount"`
	Account   Account            `bson:"account" json:"account"`
	CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type AccountDaybookDetail struct {
	Id      primitive.ObjectID `bson:"_id" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Type    string             `bson:"type" json:"type"`
	Detail  string             `bson:"detail" json:"detail"`
	Amount  float64            `bson:"amount" json:"amount"`
	Daybook AccountDaybook     `bson:"daybook" json:"daybook"`
}

type Account struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Code        string             `bson:"code" json:"code"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Type        string             `bson:"type" json:"type"`
	Company     primitive.ObjectID `bson:"company" json:"company"`
}

type AccountFinState struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Code        string             `bson:"code" json:"code"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Type        string             `bson:"type" json:"type"`
}

type AccountDaybook struct {
	Id              primitive.ObjectID  `bson:"_id" json:"id"`
	Number          string              `bson:"number" json:"number"`
	Invoice         string              `bson:"invoice" json:"invoice"`
	Document        primitive.ObjectID  `bson:"document" json:"document"`
	TransactionDate time.Time           `bson:"transactionDate" json:"transactionDate"`
	Company         primitive.ObjectID  `bson:"company" json:"company"`
	Customer        *primitive.ObjectID `bson:"customer" json:"customer"`
	Supplier        *primitive.ObjectID `bson:"supplier" json:"supplier"`
	PaymentMethod   *primitive.ObjectID `bson:"paymentMethod" json:"paymentMethod"`
}

type Document struct {
	Id   primitive.ObjectID `bson:"_id" json:"id"`
	Code string             `bson:"code" json:"code"`
	Name string             `bson:"name" json:"name"`
}

type FinancialStatement struct {
	Code          string          `json:"code"`
	Name          string          `json:"name"`
	AccountDetail []AccountDetail `json:"accountDetail"`
}

type AccountDetail struct {
	Month    string  `json:"month"`
	Date     int     `json:"date"`
	Detail   string  `json:"detail"`
	Number   string  `json:"number"`
	AmountDr float64 `json:"amountDr"`
	AmountCr float64 `json:"amountCr"`
}

type AccountBalance struct {
	AccountGroup string                `json:"accountGroup"`
	SumForwardDr float64               `json:"sumForwardDr"`
	SumForwardCr float64               `json:"sumForwardCr"`
	SumJanDr     float64               `json:"sumJanDr"`
	SumJanCr     float64               `json:"sumJanCr"`
	SumFebDr     float64               `json:"sumFebDr"`
	SumFebCr     float64               `json:"sumFebCr"`
	SumMarDr     float64               `json:"sumMarDr"`
	SumMarCr     float64               `json:"sumMarCr"`
	SumAprDr     float64               `json:"sumAprDr"`
	SumAprCr     float64               `json:"sumAprCr"`
	SumMayDr     float64               `json:"sumMayDr"`
	SumMayCr     float64               `json:"sumMayCr"`
	SumJunDr     float64               `json:"sumJunDr"`
	SumJunCr     float64               `json:"sumJunCr"`
	SumJulDr     float64               `json:"sumJulDr"`
	SumJulCr     float64               `json:"sumJulCr"`
	SumAugDr     float64               `json:"sumAugDr"`
	SumAugCr     float64               `json:"sumAugCr"`
	SumSepDr     float64               `json:"sumSepDr"`
	SumSepCr     float64               `json:"sumSepCr"`
	SumOctDr     float64               `json:"sumOctDr"`
	SumOctCr     float64               `json:"sumOctCr"`
	SumNovDr     float64               `json:"sumNovDr"`
	SumNovCr     float64               `json:"sumNovCr"`
	SumDecDr     float64               `json:"sumDecDr"`
	SumDecCr     float64               `json:"sumDecCr"`
	SumTotalDr   float64               `json:"sumTotalDr"`
	SumTotalCr   float64               `json:"sumTotalCr"`
	SumBalance   float64               `json:"sumBalance"`
	Child        []ChildAccountBalance `json:"child"`
}

type ChildAccountBalance struct {
	AccountCode string  `json:"accountCode"`
	AccountName string  `json:"accountName"`
	ForwardDr   float64 `json:"forwardDr"`
	ForwardCr   float64 `json:"forwardCr"`
	JanDr       float64 `json:"janDr"`
	JanCr       float64 `json:"janCr"`
	FebDr       float64 `json:"febDr"`
	FebCr       float64 `json:"febCr"`
	MarDr       float64 `json:"marDr"`
	MarCr       float64 `json:"marCr"`
	AprDr       float64 `json:"aprDr"`
	AprCr       float64 `json:"aprCr"`
	MayDr       float64 `json:"mayDr"`
	MayCr       float64 `json:"mayCr"`
	JunDr       float64 `json:"junDr"`
	JunCr       float64 `json:"junCr"`
	JulDr       float64 `json:"julDr"`
	JulCr       float64 `json:"julCr"`
	AugDr       float64 `json:"augDr"`
	AugCr       float64 `json:"augCr"`
	SepDr       float64 `json:"sepDr"`
	SepCr       float64 `json:"sepCr"`
	OctDr       float64 `json:"octDr"`
	OctCr       float64 `json:"octCr"`
	NovDr       float64 `json:"novDr"`
	NovCr       float64 `json:"novCr"`
	DecDr       float64 `json:"decDr"`
	DecCr       float64 `json:"decCr"`
	TotalDr     float64 `json:"totalDr"`
	TotalCr     float64 `json:"totalCr"`
	Balance     float64 `json:"balance"`
}
