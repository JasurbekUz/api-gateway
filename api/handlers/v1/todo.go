package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/JasurbekUz/api-gateway/genproto"
	l "github.com/JasurbekUz/api-gateway/pkg/logger"
	"github.com/JasurbekUz/api-gateway/pkg/utils"
)

// CreateTodo ...
// @Summary CreateTodo
// @Description This API for creating a new todo
// @Tags todo
// @Accept  json
// @Produce  json
// @Param Todo request body models.CreateTodo true "TodoCreateRequest"
// @Success 200 {object} models.GetTodo
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/todo/ [post]
func (h *handlerV1) CreateTodo(c *gin.Context) {
	var (
		body        pb.Todo
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.TodoService().Create(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetTodo ...
// @Summary GetTodo
// @Description This API for getting todo detail
// @Tags todo
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Success 200 {object} models.GetTodo
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/todo/{id} [get]
//
// GetTodo gets todo by id
// route /v1/todo/{id} [get]
func (h *handlerV1) GetTodo(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.TodoService().Get(
		ctx, &pb.ByIdReq{
			Id: guid,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetTodos ...
// @Summary GetTodos
// @Description This API for getting list of todos
// @Tags todo
// @Accept  json
// @Produce  json
// @Param page query string false "Page"
// @Param limit query string false "Limit"
// @Success 200 {object} models.GetTodos
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/todos [get]
// ListTodos returns list of todos
// route /v1/todos/ [get]
func (h *handlerV1) GetTodos(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	params, errStr := utils.ParseQueryParams(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		h.log.Error("failed to parse query params json" + errStr[0])
		return
	}

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.TodoService().List(
		ctx, &pb.ListReq{
			Limit: params.Limit,
			Page:  params.Page,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list todos", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetTodos ...
// @Summary GetTodosByDeadline
// @Description This API for getting list of todos
// @Tags todo
// @Accept  json
// @Produce  json
// @Param time path string true "time"
// @Param page query string false "Page"
// @Param limit query string false "Limit"
// @Success 200 {object} models.GetTodos
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/todos/{time} [get]
// ListTodos by deadline returns list of todos
// route /v1/todos/ [get]
func (h *handlerV1) GetTodosByDeadline(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	guid := c.Param("time")
	params, errStr := utils.ParseQueryParams(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		h.log.Error("failed to parse query params json" + errStr[0])
		return
	}

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.TodoService().ListOverdue(
		ctx, &pb.ListTime{
			ListPage: &pb.ListReq{
				Page:  params.Page,
				Limit: params.Limit,
			},
			ToTime: guid,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list todos", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateTodo ...
// @Summary UpdateTodo
// @Description This API for updating todo
// @Tags todo
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Param User request body models.CreateTodo true "userUpdateRequest"
// @Success 200 {object} models.GetTodo
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/todo/{id} [put]
// UpdateTodo updates todo by id
// route /v1/todos/{id} [put]
func (h *handlerV1) UpdateTodo(c *gin.Context) {
	var (
		body        pb.Todo
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.Id = c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.TodoService().Update(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update todo", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteTodo ...
// @Summary DeleteTodo
// @Description This API for deleting todo
// @Tags todo
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Success 200
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/todo/{id} [delete]
// DeleteTodo deletes todo by id
// route /v1/todo/{id} [delete]
func (h *handlerV1) DeleteTodo(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.TodoService().Delete(
		ctx, &pb.ByIdReq{
			Id: guid,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete todo", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
