package handler

import "github.com/go-chi/jwtauth/v5"

type Handler struct {
	AuthToken      *jwtauth.JWTAuth
	User           UserHandler
	Account        AccountHandler
	ForwardAccount ForwardAccountHandler
	Supplier       SupplierHandler
	Customer       CustomerHandler
	Document       DocumentHandler
	PaymentMethod  PaymentMethodHandler
	Product        ProductHandler
	Company        CompanyHandler
	Daybook        DaybookHandler
	DaybookDetail  DaybookDetailHandler
	Role           RoleHandler
	Material       MaterialHandler
}
