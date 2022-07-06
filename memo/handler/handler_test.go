package handler_test

import (
	"encoding/json"
	"memoapi/memo/handler"
	"memoapi/memo/mocks"
	"memoapi/model"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type memoHandlerTestSuite struct {
	suite.Suite
	echo *echo.Echo
}

func TestMemoHandlerTestSuite(t *testing.T) {
	suite.Run(t, &memoHandlerTestSuite{})
}

func (s *memoHandlerTestSuite) SetupSuite() {
	s.echo = echo.New()
}

func (s *memoHandlerTestSuite) TestMemoHandler_HandleCreateMemo() {
	reqBody := `{"title": "Hello, Go", "content": "Hello, World!"}`
	expected := model.Memo{
		ID:      5,
		Title:   "Hello, Go",
		Content: "Hello, World!",
	}
	req := httptest.NewRequest(
		http.MethodPost,
		"/memo",
		strings.NewReader(reqBody),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)

	// DI
	mockUsecase := &mocks.Usecase{}
	mockUsecase.On("CreateMemo", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		memo := args.Get(1).(*model.Memo)
		memo.ID = expected.ID
	})
	h := handler.NewMemoHTTPHandler(mockUsecase)

	s.Require().NoError(h.HandleCreateMemo(c))
	s.Equal(http.StatusCreated, rec.Code)
	respObj := struct {
		Message string `json:"message"`
		Data    struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"data"`
	}{}

	s.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &respObj))
	s.Equal("success", respObj.Message)
	s.Equal(expected.ID, respObj.Data.ID)
	s.Equal(expected.Title, respObj.Data.Title)
	s.Equal(expected.Content, respObj.Data.Content)
}

func (s *memoHandlerTestSuite) TestMemoHandler_HandleGetMemoByID() {
	expected := model.Memo{
		ID:      5,
		Title:   "Hello, Go",
		Content: "Hello, World!",
	}
	req := httptest.NewRequest(
		http.MethodGet,
		"/memo",
		nil,
	)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/memo/:memo_id")
	c.SetParamNames("memo_id")
	c.SetParamValues(strconv.Itoa(expected.ID))

	// DI
	mockUsecase := &mocks.Usecase{}
	mockUsecase.On("GetMemoByID", mock.Anything, int(expected.ID)).Return(&expected, nil)
	h := handler.NewMemoHTTPHandler(mockUsecase)

	s.Require().NoError(h.HandleGetMemoByID(c))
	s.Equal(http.StatusOK, rec.Code)
	respObj := struct {
		Message string `json:"message"`
		Data    struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"data"`
	}{}

	s.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &respObj))
	s.Equal("success", respObj.Message)
	s.Equal(expected.ID, respObj.Data.ID)
	s.Equal(expected.Title, respObj.Data.Title)
	s.Equal(expected.Content, respObj.Data.Content)
}

func (s *memoHandlerTestSuite) TestMemoHandler_HandleGetMemoByID_MemoDoesNotExist() {
	req := httptest.NewRequest(
		http.MethodGet,
		"/memo",
		nil,
	)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/memo/:memo_id")
	c.SetParamNames("memo_id")
	c.SetParamValues("1")

	// DI
	mockUsecase := &mocks.Usecase{}
	mockUsecase.On("GetMemoByID", mock.Anything, int(1)).Return(nil, nil)
	h := handler.NewMemoHTTPHandler(mockUsecase)

	s.Require().NoError(h.HandleGetMemoByID(c))
	s.Equal(http.StatusNotFound, rec.Code)
	respObj := struct {
		Message string `json:"message"`
	}{}

	s.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &respObj))
	s.NotEmpty(respObj.Message)
}

func (s *memoHandlerTestSuite) TestMemoHandler_HandleGetAllMemo() {
	expected := []*model.Memo{
		{ID: 1, Title: "first program", Content: "Hello, World!"},
		{ID: 2, Title: "last program", Content: "Goodbye, World!"},
	}
	req := httptest.NewRequest(
		http.MethodGet,
		"/memo",
		nil,
	)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)

	// DI
	mockUsecase := &mocks.Usecase{}
	mockUsecase.On("GetAllMemo", mock.Anything).Return(expected, nil)
	h := handler.NewMemoHTTPHandler(mockUsecase)

	s.Require().NoError(h.HandleGetAllMemo(c))
	s.Equal(http.StatusOK, rec.Code)
	respObj := struct {
		Message string `json:"message"`
		Data    []struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"data"`
	}{}

	s.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &respObj))
	s.Equal("success", respObj.Message)
	s.Len(respObj.Data, len(expected))
}
