package rest

import (
	_handler "findigitalservice/http/rest/internal/handler"
	mCollection "findigitalservice/http/rest/internal/model/collection"
	mHandler "findigitalservice/http/rest/internal/model/handler"
	_repo "findigitalservice/http/rest/internal/repository"
	_service "findigitalservice/http/rest/internal/service"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handlers struct {
	authToken            *jwtauth.JWTAuth
	userHandler          mHandler.UserHandler
	accountHandler       mHandler.AccountHandler
	supplierHandler      mHandler.SupplierHandler
	customerHandler      mHandler.CustomerHandler
	productHandler       mHandler.ProductHandler
	companyHandler       mHandler.CompanyHandler
	daybookHandler       mHandler.DaybookHandler
	daybookDetailHandler mHandler.DaybookDetailHandler
	documentHandler      mHandler.DocumentHandler
	roleHandler          mHandler.RoleHandler
	materialHandler      mHandler.MaterialHandler
}

func Register(db *mongo.Database, logger *logrus.Logger) *chi.Mux {
	r := chi.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
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
	r.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Route("/auth", func(r chi.Router) {
				r.Post("/register", h.userHandler.Create)
				r.Post("/login", h.userHandler.Login)
			})
		})
		r.Group(func(r chi.Router) {
			// Seek, verify and validate JWT tokens
			r.Use(jwtauth.Verifier(h.authToken))

			// Handle valid / invalid tokens. In this example, we use
			// the provided authenticator middleware, but you can write your
			// own very easily, look at the Authenticator method in jwtauth.go
			// and tweak it, its not scary.
			r.Use(jwtauth.Authenticator(h.authToken))
			r.Route("/user", func(r chi.Router) {
				r.Get("/", h.userHandler.FindAll)
				r.Get("/{id}", h.userHandler.FindById)
				r.Get("/profile", h.userHandler.FindUserProfile)
				r.Get("/company", h.userHandler.FindUserCompany)
				r.Put("/{id}", h.userHandler.Update)
				r.Get("/count", h.userHandler.Count)
			})
			r.Route("/account", func(r chi.Router) {
				r.Get("/", h.accountHandler.FindAll)
				r.Get("/{id}", h.accountHandler.FindById)
				r.Post("/", h.accountHandler.Create)
				r.Put("/{id}", h.accountHandler.Update)
				r.Get("/count", h.accountHandler.Count)
				r.Delete("/{id}", h.accountHandler.Delete)
			})
			r.Route("/supplier", func(r chi.Router) {
				r.Get("/", h.supplierHandler.FindAll)
				r.Get("/{id}", h.supplierHandler.FindById)
				r.Post("/", h.supplierHandler.Create)
				r.Put("/{id}", h.supplierHandler.Update)
				r.Get("/count", h.supplierHandler.Count)
				r.Delete("/{id}", h.supplierHandler.Delete)
			})
			r.Route("/customer", func(r chi.Router) {
				r.Get("/", h.customerHandler.FindAll)
				r.Get("/{id}", h.customerHandler.FindById)
				r.Post("/", h.customerHandler.Create)
				r.Put("/{id}", h.customerHandler.Update)
				r.Get("/count", h.customerHandler.Count)
				r.Delete("/{id}", h.customerHandler.Delete)
			})
			r.Route("/product", func(r chi.Router) {
				r.Get("/", h.productHandler.FindAll)
				r.Get("/{id}", h.productHandler.FindById)
				r.Post("/", h.productHandler.Create)
				r.Put("/{id}", h.productHandler.Update)
				r.Get("/count", h.productHandler.Count)
				r.Delete("/{id}", h.productHandler.Delete)
			})
			r.Route("/material", func(r chi.Router) {
				r.Get("/", h.materialHandler.FindAll)
				r.Get("/{id}", h.materialHandler.FindById)
				r.Post("/", h.materialHandler.Create)
				r.Put("/{id}", h.materialHandler.Update)
				r.Get("/count", h.materialHandler.Count)
				r.Delete("/{id}", h.materialHandler.Delete)
			})
			r.Route("/company", func(r chi.Router) {
				r.Get("/", h.companyHandler.FindAll)
				r.Get("/{id}", h.companyHandler.FindById)
				r.Post("/", h.companyHandler.Create)
				r.Put("/{id}", h.companyHandler.Update)
			})
			r.Route("/daybook", func(r chi.Router) {
				r.Get("/", h.daybookHandler.FindAll)
				r.Get("/{id}", h.daybookHandler.FindById)
				r.Post("/", h.daybookHandler.Create)
				r.Put("/{id}", h.daybookHandler.Update)
				r.Get("/count", h.daybookHandler.Count)
				r.Route("/generate", func(r chi.Router) {
					r.Get("/excel/{id}", h.daybookHandler.GenerateExcel)
				})
			})
			r.Route("/daybook/detail", func(r chi.Router) {
				r.Get("/", h.daybookDetailHandler.FindAll)
				r.Get("/{id}", h.daybookDetailHandler.FindById)
				r.Post("/", h.daybookDetailHandler.Create)
				r.Put("/{id}", h.daybookDetailHandler.Update)
			})
			r.Route("/document", func(r chi.Router) {
				r.Get("/", h.documentHandler.FindAll)
				r.Get("/{id}", h.documentHandler.FindById)
				r.Post("/", h.documentHandler.Create)
				r.Put("/{id}", h.documentHandler.Update)
			})
			r.Route("/role", func(r chi.Router) {
				r.Get("/", h.roleHandler.FindAll)
				r.Get("/{id}", h.roleHandler.FindById)
				r.Post("/", h.roleHandler.Create)
				r.Put("/{id}", h.roleHandler.Update)
			})
		})
	})

	return r
}

func initRouter(db *mongo.Database, logger *logrus.Logger) Handlers {
	// init collection
	collection := mCollection.Collection{
		User:          db.Collection("users"),
		Account:       db.Collection("accounts"),
		Supplier:      db.Collection("suppliers"),
		Customer:      db.Collection("customers"),
		Document:      db.Collection("documents"),
		Product:       db.Collection("products"),
		Company:       db.Collection("companies"),
		Daybook:       db.Collection("daybooks"),
		DaybookDetail: db.Collection("daybook_details"),
		Role:          db.Collection("roles"),
		Material:      db.Collection("materials"),
	}
	// init repository
	userRepo := _repo.InitUserRepository(collection, logger)
	accountRepo := _repo.InitAccountRepository(collection, logger)
	supplierRepo := _repo.InitSupplierRepository(collection, logger)
	customerRepo := _repo.InitCustomerRepository(collection, logger)
	productRepo := _repo.InitProductRepository(collection, logger)
	companyRepo := _repo.InitCompanyRepository(collection, logger)
	daybookRepo := _repo.InitDaybookRepository(collection, logger)
	daybookDetailRepo := _repo.InitDaybookDetailRepository(collection, logger)
	documentRepo := _repo.InitDocumentRepository(collection, logger)
	roleRepo := _repo.InitRoleRepository(collection, logger)
	materialRepo := _repo.InitMaterialRepository(collection, logger)

	// init service
	userService := _service.InitUserService(userRepo, logger)
	accountService := _service.InitAccountService(accountRepo, logger)
	supplierService := _service.InitSupplierService(supplierRepo, logger)
	customerService := _service.InitCustomerService(customerRepo, logger)
	productService := _service.InitProductService(productRepo, logger)
	companyService := _service.InitCompanyService(companyRepo, logger)
	daybookService := _service.InitDaybookService(daybookRepo, daybookDetailRepo, logger)
	daybookDetailService := _service.InitDaybookDetailService(daybookDetailRepo, daybookRepo, logger)
	documentService := _service.InitDocumentService(documentRepo, logger)
	roleService := _service.InitRoleService(roleRepo, logger)
	materialService := _service.InitMaterialService(materialRepo, logger)

	// init handler
	userHandler := _handler.InitUserHandler(userService, logger)
	accountHandler := _handler.InitAccountHandler(accountService, logger)
	supplierHandler := _handler.InitSupplierHandler(supplierService, logger)
	customerHandler := _handler.InitCustomerHandler(customerService, logger)
	productHandler := _handler.InitProductHandler(productService, logger)
	companyHandler := _handler.InitCompanyHandler(companyService, logger)
	daybookHandler := _handler.InitDaybookHandler(daybookService, logger)
	daybookDetailHandler := _handler.InitDaybookDetailHandler(daybookDetailService, logger)
	documentHandler := _handler.InitDocumentHandler(documentService, logger)
	roleHandler := _handler.InitRoleHandler(roleService, logger)
	materialHandler := _handler.InitMaterialHandler(materialService, logger)

	return Handlers{
		authToken:            jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil),
		userHandler:          userHandler,
		accountHandler:       accountHandler,
		supplierHandler:      supplierHandler,
		customerHandler:      customerHandler,
		productHandler:       productHandler,
		companyHandler:       companyHandler,
		daybookHandler:       daybookHandler,
		daybookDetailHandler: daybookDetailHandler,
		documentHandler:      documentHandler,
		roleHandler:          roleHandler,
		materialHandler:      materialHandler,
	}
}
