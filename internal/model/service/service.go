package service

type Service struct {
	User           UserService
	Account        AccountService
	AccountType    AccountTypeService
	ForwardAccount ForwardAccountService
	Supplier       SupplierService
	Customer       CustomerService
	Document       DocumentService
	PaymentMethod  PaymentMethodService
	Product        ProductService
	Company        CompanyService
	Daybook        DaybookService
	DaybookDetail  DaybookDetailService
	Role           RoleService
	Material       MaterialService
}
