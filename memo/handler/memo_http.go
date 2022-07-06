package handler

import (
	"context"
	"memoapi/memo"
	"memoapi/model"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// MemoHTTPHandler represents memo http handler.
type MemoHTTPHandler struct {
	usecase memo.Usecase
}

// NewMemoHTTPHandler returns new instance of MemoHTTPHandler.
func NewMemoHTTPHandler(usecase memo.Usecase) *MemoHTTPHandler {
	return &MemoHTTPHandler{usecase: usecase}
}

// HandleCreateMemo handles create a memo.
func (h *MemoHTTPHandler) HandleCreateMemo(c echo.Context) error {
	memo := model.Memo{}
	if err := c.Bind(&memo); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &Response{
			Message: err.Error(),
			Data:    nil,
		})
	}

	if err := h.usecase.CreateMemo(context.Background(), &memo); err != nil {
		statusCode := http.StatusBadRequest
		if _, ok := err.(validator.ValidationErrors); ok {
			statusCode = http.StatusBadRequest
		}
		return c.JSON(statusCode, &Response{
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, &Response{
		Message: "success",
		Data:    memo,
	})
}

// HandleGetMemoByID handles get memo by id.
func (h *MemoHTTPHandler) HandleGetMemoByID(c echo.Context) error {
	memoIDStr := c.Param("memo_id")
	id, err := strconv.Atoi(memoIDStr)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{
			Message: "memo is not found",
			Data:    nil,
		})
	}

	memo, err := h.usecase.GetMemoByID(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Message: err.Error(),
			Data:    nil,
		})
	}

	if memo == nil {
		return c.JSON(http.StatusNotFound, &Response{
			Message: "memo is not found",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, &Response{
		Message: "success",
		Data:    memo,
	})
}

// HandleGetAllMemo handles get all memos.
func (h *MemoHTTPHandler) HandleGetAllMemo(c echo.Context) error {
	res, err := h.usecase.GetAllMemo(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, &Response{
		Message: "success",
		Data:    res,
	})
}
