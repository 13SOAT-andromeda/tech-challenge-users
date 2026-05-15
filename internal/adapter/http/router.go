package http

import (
	"tech-challenge-users/internal/adapter/http/handlers"
	"tech-challenge-users/internal/adapter/http/middlewares"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter() *Router {
	engine := gin.New()
	engine.Use(gin.Recovery())

	return &Router{engine: engine}
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

// setupInternalRoutes wires /internal/* handlers (populated by module-users).
func (r *Router) setupInternalRoutes(_ *gin.RouterGroup) {}

// setupUsersRoutes wires /users/* handlers (populated by each domain module).
func (r *Router) setupUsersRoutes(_ *gin.RouterGroup) {}
