package http

import (
	"time"

	"tech-challenge-users/internal/adapter/http/handlers"
	"tech-challenge-users/internal/adapter/http/middlewares"

	gintrace "github.com/DataDog/dd-trace-go/contrib/gin-gonic/gin/v2"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Router struct {
	engine                 *gin.Engine
	userHandler            *handlers.UserHandler
	internalUserHandler    *handlers.InternalUserHandler
	customerHandler        *handlers.CustomerHandler
	vehicleHandler         *handlers.VehicleHandler
	customerVehicleHandler *handlers.CustomerVehicleHandler
	companyHandler         *handlers.CompanyHandler
	employeeHandler        *handlers.EmployeeHandler
}

func NewRouter(
	logger *zap.Logger,
	serviceName string,
	userHandler *handlers.UserHandler,
	internalUserHandler *handlers.InternalUserHandler,
	customerHandler *handlers.CustomerHandler,
	vehicleHandler *handlers.VehicleHandler,
	customerVehicleHandler *handlers.CustomerVehicleHandler,
	companyHandler *handlers.CompanyHandler,
	employeeHandler *handlers.EmployeeHandler,
) *Router {
	engine := gin.New()

	engine.Use(gintrace.Middleware(serviceName,
		gintrace.WithUseGinErrors(),
		gintrace.WithAnalytics(true),
		gintrace.WithIgnoreRequest(func(c *gin.Context) bool {
			return c.Request.URL.Path == "/health"
		}),
		gintrace.WithStatusCheck(func(statusCode int) bool {
			return statusCode >= 400
		}),
	))
	engine.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
			fields := []zapcore.Field{}
			span, ok := tracer.SpanFromContext(c.Request.Context())
			if ok {
				fields = append(fields, zap.String("trace_id", span.Context().TraceID()))
				fields = append(fields, zap.String("span_id", span.String()))
			}
			return fields
		}),
	}))
	engine.Use(ginzap.RecoveryWithZap(logger, true))

	return &Router{
		engine:                 engine,
		userHandler:            userHandler,
		internalUserHandler:    internalUserHandler,
		customerHandler:        customerHandler,
		vehicleHandler:         vehicleHandler,
		customerVehicleHandler: customerVehicleHandler,
		companyHandler:         companyHandler,
		employeeHandler:        employeeHandler,
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

	employees := g.Group("/employees")
	employees.Use(adminOrAttendant)
	employees.GET("/:id", r.employeeHandler.GetByID)

	// Dynamic :id routes
	g.POST("", adminOrAttendant, r.userHandler.Create)
	g.GET("", adminOrAttendant, r.userHandler.List)
	g.GET("/:id", adminOrAttendant, r.userHandler.GetByID)
	g.PUT("/:id", adminOrAttendant, r.userHandler.Update)
	g.DELETE("/:id", adminOrAttendant, r.userHandler.Delete)
}
