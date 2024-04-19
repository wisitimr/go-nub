package rest

import (
	"net/http"
	_handler "nub/internal/handler"
	mCollection "nub/internal/model/collection"
	mHandler "nub/internal/model/handler"
	mRepo "nub/internal/model/repository"
	mService "nub/internal/model/service"
	_repo "nub/internal/repository"
	_service "nub/internal/service"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(db *mongo.Database, logger *logrus.Logger) *chi.Mux {
	r := chi.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(c.Handler)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	h := initRouter(db, logger)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	r.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Route("/auth", func(r chi.Router) {
				r.Post("/register", h.User.Create)
				r.Post("/login", h.User.Login)
			})
		})
		r.Group(func(r chi.Router) {
			// Seek, verify and validate JWT tokens
			r.Use(jwtauth.Verifier(h.AuthToken))

			// Handle valid / invalid tokens. In this example, we use
			// the provided authenticator middleware, but you can write your
			// own very easily, look at the Authenticator method in jwtauth.go
			// and tweak it, its not scary.
			r.Use(jwtauth.Authenticator(h.AuthToken))
			r.Route("/user", func(r chi.Router) {
				r.Get("/", h.User.FindAll)
				r.Get("/{id}", h.User.FindById)
				r.Get("/profile", h.User.FindUserProfile)
				r.Get("/company", h.User.FindUserCompany)
				r.Put("/{id}", h.User.Update)
				r.Get("/count", h.User.Count)
			})
			r.Route("/account", func(r chi.Router) {
				r.Get("/", h.Account.FindAll)
				r.Get("/{id}", h.Account.FindById)
				r.Post("/", h.Account.Create)
				r.Put("/{id}", h.Account.Update)
				r.Get("/count", h.Account.Count)
				r.Delete("/{id}", h.Account.Delete)
			})
			r.Route("/accountType", func(r chi.Router) {
				r.Get("/", h.AccountType.FindAll)
				r.Get("/{id}", h.AccountType.FindById)
				r.Post("/", h.AccountType.Create)
				r.Put("/{id}", h.AccountType.Update)
				r.Get("/count", h.AccountType.Count)
				r.Delete("/{id}", h.AccountType.Delete)
			})
			r.Route("/forward", func(r chi.Router) {
				r.Get("/", h.ForwardAccount.FindAll)
				r.Get("/{id}", h.ForwardAccount.FindById)
				r.Post("/", h.ForwardAccount.Create)
				r.Put("/{id}", h.ForwardAccount.Update)
				r.Get("/count", h.ForwardAccount.Count)
				r.Delete("/{id}", h.ForwardAccount.Delete)
			})
			r.Route("/supplier", func(r chi.Router) {
				r.Get("/", h.Supplier.FindAll)
				r.Get("/{id}", h.Supplier.FindById)
				r.Post("/", h.Supplier.Create)
				r.Put("/{id}", h.Supplier.Update)
				r.Get("/count", h.Supplier.Count)
				r.Delete("/{id}", h.Supplier.Delete)
			})
			r.Route("/customer", func(r chi.Router) {
				r.Get("/", h.Customer.FindAll)
				r.Get("/{id}", h.Customer.FindById)
				r.Post("/", h.Customer.Create)
				r.Put("/{id}", h.Customer.Update)
				r.Get("/count", h.Customer.Count)
				r.Delete("/{id}", h.Customer.Delete)
			})
			r.Route("/product", func(r chi.Router) {
				r.Get("/", h.Product.FindAll)
				r.Get("/{id}", h.Product.FindById)
				r.Post("/", h.Product.Create)
				r.Put("/{id}", h.Product.Update)
				r.Get("/count", h.Product.Count)
				r.Delete("/{id}", h.Product.Delete)
			})
			r.Route("/material", func(r chi.Router) {
				r.Get("/", h.Material.FindAll)
				r.Get("/{id}", h.Material.FindById)
				r.Post("/", h.Material.Create)
				r.Put("/{id}", h.Material.Update)
				r.Get("/count", h.Material.Count)
				r.Delete("/{id}", h.Material.Delete)
			})
			r.Route("/company", func(r chi.Router) {
				r.Get("/", h.Company.FindAll)
				r.Get("/{id}", h.Company.FindById)
				r.Post("/", h.Company.Create)
				r.Put("/{id}", h.Company.Update)
			})
			r.Route("/daybook", func(r chi.Router) {
				r.Get("/", h.Daybook.FindAll)
				r.Get("/{id}", h.Daybook.FindById)
				r.Post("/", h.Daybook.Create)
				r.Put("/{id}", h.Daybook.Update)
				r.Get("/count", h.Daybook.Count)
			})
			r.Route("/report", func(r chi.Router) {
				r.Get("/account/ledger/{company}/{year}", h.Daybook.FindLedgerAccount)
				r.Get("/account/balance/{company}/{year}", h.Daybook.FindAccountBalance)
				r.Route("/generate", func(r chi.Router) {
					r.Get("/excel/{id}", h.Daybook.GenerateExcel)
					r.Get("/financial/{company}/{year}", h.Daybook.GenerateFinancialStatement)
				})
			})
			r.Route("/daybook/detail", func(r chi.Router) {
				r.Get("/", h.DaybookDetail.FindAll)
				r.Get("/{id}", h.DaybookDetail.FindById)
				r.Post("/", h.DaybookDetail.Create)
				r.Put("/{id}", h.DaybookDetail.Update)
			})
			r.Route("/document", func(r chi.Router) {
				r.Get("/", h.Document.FindAll)
				r.Get("/{id}", h.Document.FindById)
				r.Post("/", h.Document.Create)
				r.Put("/{id}", h.Document.Update)
			})
			r.Route("/payment/method", func(r chi.Router) {
				r.Get("/", h.PaymentMethod.FindAll)
				r.Get("/{id}", h.PaymentMethod.FindById)
				r.Post("/", h.PaymentMethod.Create)
				r.Put("/{id}", h.PaymentMethod.Update)
			})
			r.Route("/role", func(r chi.Router) {
				r.Get("/", h.Role.FindAll)
				r.Get("/{id}", h.Role.FindById)
				r.Post("/", h.Role.Create)
				r.Put("/{id}", h.Role.Update)
			})
		})
	})

	return r
}

func initRouter(db *mongo.Database, logger *logrus.Logger) mHandler.Handler {
	// init collection
	collection := mCollection.Collection{
		User:           db.Collection("users"),
		Account:        db.Collection("accounts"),
		AccountType:    db.Collection("accountTypes"),
		ForwardAccount: db.Collection("forward_accounts"),
		Supplier:       db.Collection("suppliers"),
		Customer:       db.Collection("customers"),
		Document:       db.Collection("documents"),
		PaymentMethod:  db.Collection("payment_methods"),
		Product:        db.Collection("products"),
		Company:        db.Collection("companies"),
		Daybook:        db.Collection("daybooks"),
		DaybookDetail:  db.Collection("daybook_details"),
		Role:           db.Collection("roles"),
		Material:       db.Collection("materials"),
	}
	// init repository
	repo := mRepo.Repository{
		User:           _repo.InitUserRepository(collection, logger),
		Account:        _repo.InitAccountRepository(collection, logger),
		AccountType:    _repo.InitAccountTypeRepository(collection, logger),
		ForwardAccount: _repo.InitForwardAccountRepository(collection, logger),
		Supplier:       _repo.InitSupplierRepository(collection, logger),
		Customer:       _repo.InitCustomerRepository(collection, logger),
		Product:        _repo.InitProductRepository(collection, logger),
		Company:        _repo.InitCompanyRepository(collection, logger),
		Daybook:        _repo.InitDaybookRepository(collection, logger),
		DaybookDetail:  _repo.InitDaybookDetailRepository(collection, logger),
		Document:       _repo.InitDocumentRepository(collection, logger),
		PaymentMethod:  _repo.InitPaymentMethodRepository(collection, logger),
		Role:           _repo.InitRoleRepository(collection, logger),
		Material:       _repo.InitMaterialRepository(collection, logger),
	}

	// init service
	service := mService.Service{
		User:           _service.InitUserService(repo, logger),
		Account:        _service.InitAccountService(repo, logger),
		AccountType:    _service.InitAccountTypeService(repo, logger),
		ForwardAccount: _service.InitForwardAccountService(repo, logger),
		Supplier:       _service.InitSupplierService(repo, logger),
		Customer:       _service.InitCustomerService(repo, logger),
		Product:        _service.InitProductService(repo, logger),
		Company:        _service.InitCompanyService(repo, logger),
		Daybook:        _service.InitDaybookService(repo, logger),
		DaybookDetail:  _service.InitDaybookDetailService(repo, logger),
		Document:       _service.InitDocumentService(repo, logger),
		PaymentMethod:  _service.InitPaymentMethodService(repo, logger),
		Role:           _service.InitRoleService(repo, logger),
		Material:       _service.InitMaterialService(repo, logger),
	}

	// init handler
	handler := mHandler.Handler{
		AuthToken:      jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil),
		User:           _handler.InitUserHandler(service, logger),
		Account:        _handler.InitAccountHandler(service, logger),
		AccountType:    _handler.InitAccountTypeHandler(service, logger),
		ForwardAccount: _handler.InitForwardAccountHandler(service, logger),
		Supplier:       _handler.InitSupplierHandler(service, logger),
		Customer:       _handler.InitCustomerHandler(service, logger),
		Product:        _handler.InitProductHandler(service, logger),
		Company:        _handler.InitCompanyHandler(service, logger),
		Daybook:        _handler.InitDaybookHandler(service, logger),
		DaybookDetail:  _handler.InitDaybookDetailHandler(service, logger),
		Document:       _handler.InitDocumentHandler(service, logger),
		PaymentMethod:  _handler.InitPaymentMethodHandler(service, logger),
		Role:           _handler.InitRoleHandler(service, logger),
		Material:       _handler.InitMaterialHandler(service, logger),
	}

	return handler
}
