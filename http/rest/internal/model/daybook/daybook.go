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
	Code        string        `json:"code"`
	Name        string        `json:"name"`
	MonthDetail []MonthDetail `json:"monthDetail"`
}

type MonthDetail struct {
	Month         string          `json:"month"`
	AccountDetail []AccountDetail `json:"accountDetail"`
}

type AccountDetail struct {
	Date     int     `json:"date"`
	Detail   string  `json:"description"`
	Number   string  `json:"number"`
	AmountDr float64 `json:"amountDr"`
	AmountCr float64 `json:"amountCr"`
}
