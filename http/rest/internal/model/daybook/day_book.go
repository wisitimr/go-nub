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
	DaybookDetails    []DaybookDetails    `bson:"daybookDetails" json:"daybookDetails"`
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
	DaybookDetails    []DaybookDetails   `bson:"daybookDetails" json:"daybookDetails"`
	DebitTotalAmount  float64            `bson:"debitTotalAmount" json:"debitTotalAmount"`
	CreditTotalAmount float64            `bson:"creditTotalAmount" json:"creditTotalAmount"`
	CreatedBy         primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt         time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedBy         primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt         time.Time          `bson:"updatedAt" json:"updatedAt"`
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

type DaybookDetails struct {
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

type Account struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Code        string             `bson:"code" json:"code"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Type        string             `bson:"type" json:"type"`
	Company     primitive.ObjectID `bson:"company" json:"company"`
}

type Document struct {
	Id   primitive.ObjectID `bson:"_id" json:"id"`
	Code string             `bson:"code" json:"code"`
	Name string             `bson:"name" json:"name"`
}
