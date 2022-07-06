//go:build integration
// +build integration

package it

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"memoapi/config"
	"memoapi/infra/db"
	"memoapi/infra/server"
	"memoapi/model"
	"net/http"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type e2eTestSuite struct {
	suite.Suite
	config      config.Config
	db          *gorm.DB
	dbMigration *migrate.Migrate
	server      *server.Server
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {
	var err error

	s.config = config.Load()
	s.db, err = db.NewPostgresDB(s.config.RdbConfig.ConnectionString())
	s.Require().NoError(err)

	s.dbMigration, err = migrate.New("file://../initdb/db-test", s.config.RdbConfig.ConnectionString())
	s.Require().NoError(err)
	if err := s.dbMigration.Up(); err != nil && err != migrate.ErrNoChange {
		s.Require().NoError(err)
	}

	serverReady := make(chan interface{}, 1)
	s.server = server.NewServer(s.config.AppPort(), s.config.AppHost(), s.db, serverReady)

	go s.server.Start()
	<-serverReady
}

func (s *e2eTestSuite) TearDownSuite() {
	s.server.Shutdown()
}

func (s *e2eTestSuite) SetupTest() {
	if err := s.dbMigration.Up(); err != nil && err != migrate.ErrNoChange {
		s.Require().NoError(err)
	}
}

func (s *e2eTestSuite) TearDownTest() {
	s.NoError(s.dbMigration.Down())
}

func (s *e2eTestSuite) TestCreateMemo() {
	reqBody := `{"title": "Hello, Go", "content": "Hello, World!"}`
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://localhost:%d/memo", s.config.AppPort()),
		strings.NewReader(reqBody),
	)
	s.Require().NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}
	resp, err := client.Do(req)
	s.Require().NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	s.Require().NoError(err)

	respObj := struct {
		Message string `json:"message"`
		Data    struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"data"`
	}{}
	s.Require().NoError(json.Unmarshal(byteBody, &respObj))

	s.Equal("success", respObj.Message)
	s.NotEmpty(respObj.Data.ID)
	s.Equal("Hello, Go", respObj.Data.Title)
	s.Equal("Hello, World!", respObj.Data.Content)
}

func (s *e2eTestSuite) TestGetMemoByID() {
	dataCount := 5
	memos := make([]model.Memo, dataCount)
	for i := 0; i < len(memos); i++ {
		memos[i].Title = uuid.NewString()
		memos[i].Content = uuid.NewString()
	}
	s.db.Create(memos)

	expected := &memos[dataCount/2]
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://localhost:%d/memo/%d", s.config.AppPort(), expected.ID),
		nil,
	)
	s.Require().NoError(err)

	client := http.Client{}
	resp, err := client.Do(req)
	s.Require().NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	s.Require().NoError(err)

	respObj := struct {
		Message string `json:"message"`
		Data    struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"data"`
	}{}

	s.Require().NoError(json.Unmarshal(byteBody, &respObj))
	s.Equal("success", respObj.Message)
	s.Equal(expected.ID, respObj.Data.ID)
	s.Equal(expected.Title, respObj.Data.Title)
	s.Equal(expected.Content, respObj.Data.Content)
}

func (s *e2eTestSuite) TestGetAllMemoByID() {
	dataCount := 5
	memos := make([]model.Memo, dataCount)
	for i := 0; i < len(memos); i++ {
		memos[i].Title = uuid.NewString()
		memos[i].Content = uuid.NewString()
	}
	s.db.Create(memos)

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://localhost:%d/memo", s.config.AppPort()),
		nil,
	)
	s.Require().NoError(err)

	client := http.Client{}
	resp, err := client.Do(req)
	s.Require().NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	s.Require().NoError(err)

	respObj := struct {
		Message string `json:"message"`
		Data    []struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"data"`
	}{}

	s.Require().NoError(json.Unmarshal(byteBody, &respObj))
	s.Equal("success", respObj.Message)
	s.Equal(dataCount, len(respObj.Data))
	for i, memo := range respObj.Data {
		s.Equal(memos[i].ID, memo.ID)
		s.Equal(memos[i].Title, memo.Title)
		s.Equal(memos[i].Content, memo.Content)
	}
}
