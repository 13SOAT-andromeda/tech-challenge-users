package http

import (
	"tech-challenge-users/internal/adapter/http/handlers"
	"tech-challenge-users/internal/adapter/http/middlewares"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine                 *gin.Engine
	userHandler            *handlers.UserHandler
	internalUserHandler    *handlers.InternalUserHandler
	customerHandler        *handlers.CustomerHandler
	vehicleHandler         *handlers.VehicleHandler
	customerVehicleHandler *handlers.CustomerVehicleHandler
	companyHandler         *handlers.CompanyHandler
}

func NewRouter(
	userHandler *handlers.UserHandler,
	internalUserHandler *handlers.InternalUserHandler,
	customerHandler *handlers.CustomerHandler,
	vehicleHandler *handlers.VehicleHandler,
	customerVehicleHandler *handlers.CustomerVehicleHandler,
	companyHandler *handlers.CompanyHandler,
) *Router {
	engine := gin.New()
	engine.Use(gin.Recovery())

	return &Router{
		engine:                 engine,
		userHandler:            userHandler,
		internalUserHandler:    internalUserHandler,
		customerHandler:        customerHandler,
		vehicleHandler:         vehicleHandler,
		customerVehicleHandler: customerVehicleHandler,
		companyHandler:         companyHandler,
	}
}

func (r *Router) Setup() {
	r.engine.GET("/health", handlers.Health)

	internal := r.engine.Group("/internal")
	r.setupInternalRoutes(internal)

	users := r.engine.Group("/users")
	users.Use(middlewares.AuthRequired())
	r.setupUsersRoutes(users)
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}

func (r *Router) setupInternalRoutes(g *gin.RouterGroup) {
	g.GET("/users/by-document", r.internalUserHandler.GetByDocument)
}

func (r *Router) setupUsersRoutes(g *gin.RouterGroup) {
	adminOrAttendant := middlewares.RoleRequired("administrator", "attendant")
	adminAttendantMechanic := middlewares.RoleRequired("administrator", "attendant", "mechanic")

	// Static sub-paths BEFORE dynamic :id

	customers := g.Group("/customers")
	customers.Use(adminOrAttendant)
	customers.POST("", r.customerHandler.Create)
	customers.GET("", r.customerHandler.List)
	customers.GET("/:id", r.customerHandler.GetByID)
	customers.PUT("/:id", r.customerHandler.Update)
	customers.DELETE("/:id", r.customerHandler.Delete)
	customers.GET("/:id/vehicles", r.customerVehicleHandler.List)
	customers.POST("/:id/vehicles/:vehicleId", r.customerVehicleHandler.Add)
	customers.DELETE("/:id/vehicles/:vehicleId", r.customerVehicleHandler.Remove)

	vehicles := g.Group("/vehicles")
	vehicles.Use(adminAttendantMechanic)
	vehicles.POST("", r.vehicleHandler.Create)
	vehicles.GET("", r.vehicleHandler.List)
	vehicles.GET("/:id", r.vehicleHandler.GetByID)
	vehicles.PUT("/:id", r.vehicleHandler.Update)
	vehicles.DELETE("/:id", r.vehicleHandler.Delete)

	adminOnly := middlewares.RoleRequired("administrator")
	companies := g.Group("/companies")
	companies.Use(adminOnly)
	companies.POST("", r.companyHandler.Create)
	companies.GET("/:id", r.companyHandler.GetByID)
	companies.PUT("/:id", r.companyHandler.Update)
	companies.DELETE("/:id", r.companyHandler.Delete)

	// Dynamic :id routes
	g.POST("", adminOrAttendant, r.userHandler.Create)
	g.GET("", adminOrAttendant, r.userHandler.List)
	g.GET("/:id", adminOrAttendant, r.userHandler.GetByID)
	g.PUT("/:id", adminOrAttendant, r.userHandler.Update)
	g.DELETE("/:id", adminOrAttendant, r.userHandler.Delete)
}
