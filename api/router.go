package api

import (
	"github.com/gin-gonic/gin"

	"github.com/JasurbekUz/api-gateway/api/handlers/v1"
	"github.com/JasurbekUz/api-gateway/config"
	"github.com/JasurbekUz/api-gateway/pkg/logger"
	"github.com/JasurbekUz/api-gateway/services"
)

// OPTIONS
type Option struct {
	Conf config.Config
	Logger logger.Logger
	ServiceManager services.IServiceManager
}

// NEW
func New(option Option) gin.Engine {
	router := gin.Default()

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger: option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg: option.Conf,
	})

	api := router.Group("/v1")

	api.POST("/todo", handlerV1.CreateTodo)
	api.GET("/todo/:id", handlerV1.GetTodo)
	//api.GET("/todos", handlerV1.GetTodos)
	//api.GET("/todos/:time", handlerV1.GetTodosByDeadline)
	//api.PUT("/todo/:id", handlerV1.UpdateTodo)
	//api.DELETE("/todo/:id", handlerV1.DeleteTodo)

	return router
}
