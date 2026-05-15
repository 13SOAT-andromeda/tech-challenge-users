package http

import (
	"tech-challenge-users/internal/adapter/http/handlers"
	"tech-challenge-users/internal/adapter/http/middlewares"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine              *gin.Engine
	userHandler         *handlers.UserHandler
	internalUserHandler *handlers.InternalUserHandler
}

func NewRouter(userHandler *handlers.UserHandler, internalUserHandler *handlers.InternalUserHandler) *Router {
	engine := gin.New()
	engine.Use(gin.Recovery())

	return &Router{
		engine:              engine,
		userHandler:         userHandler,
		internalUserHandler: internalUserHandler,
	}
}

func (r *Router) Setup() {
	// Public routes
	r.engine.GET("/health", handlers.Health)

	// Internal routes — no auth, protected by NetworkPolicy
	internal := r.engine.Group("/internal")
	r.setupInternalRoutes(internal)

	// Public /users base group — all sub-resources under /users
	// IMPORTANT: static sub-paths (customers, vehicles, companies) must be
	// registered before dynamic param routes (:id) to ensure correct resolution.
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

	g.POST("", adminOrAttendant, r.userHandler.Create)
	g.GET("", adminOrAttendant, r.userHandler.List)
	g.GET("/:id", adminOrAttendant, r.userHandler.GetByID)
	g.PUT("/:id", adminOrAttendant, r.userHandler.Update)
	g.DELETE("/:id", adminOrAttendant, r.userHandler.Delete)
}
