package internal

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"net/http/httptest"
// 	"net/url"
// 	"strings"
// 	"testing"

// 	"github.com/labstack/echo/v4"
// 	"github.com/skvdmt/skvdmt-back/internal/entities"
// 	"github.com/skvdmt/skvdmt-back/internal/pool"
// 	_ "github.com/skvdmt/skvdmt-back/testing_init"
// 	"github.com/stretchr/testify/assert"
// )

// type CreateSample struct {
// 	name    string
// 	in      *entities.Sample
// 	expCode int
// 	expBody string
// }

// const (
// 	titleSample        = "testing sample"
// 	updatedTitleSample = "updated sample title"
// )

// // TestCreateSample testing create new Sample
// func TestCreateSample(t *testing.T) {
// 	ss, err := NewSamples(WithoutParseFlags())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	tss := []CreateSample{
// 		{
// 			name: "create short title",
// 			in: &entities.Sample{
// 				Title: "жжж",
// 			},
// 			expCode: http.StatusBadRequest,
// 			expBody: fmt.Sprintf(`{"message":"%v"}`, ss.bundle.Errors[pool.ErrTitleLen]),
// 		},
// 		{
// 			name: "create long title",
// 			in: &entities.Sample{
// 				Title: "жжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжж",
// 			},
// 			expCode: http.StatusBadRequest,
// 			expBody: fmt.Sprintf(`{"message":"%v"}`, ss.bundle.Errors[pool.ErrTitleLen]),
// 		},
// 		{
// 			name: "create success ok",
// 			in: &entities.Sample{
// 				Title: titleSample,
// 			},
// 			expCode: http.StatusCreated,
// 			expBody: `{"id":4}`,
// 		},
// 	}

// 	for _, ts := range tss {
// 		t.Run(ts.name, func(t *testing.T) {
// 			in, err := json.Marshal(ts.in)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			sin := string(in)
// 			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(sin))
// 			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 			rec := httptest.NewRecorder()
// 			c := ss.bundle.Router.NewContext(req, rec)
// 			if assert.NoError(t, ss.handlers.Create(c)) {
// 				assert.Equal(t, ts.expCode, rec.Code)
// 				assert.Equal(t, ts.expBody, strings.TrimSpace(rec.Body.String()))
// 			}
// 		})
// 	}
// }

// type ReadSample struct {
// 	name    string
// 	id      int
// 	expCode int
// 	expBody string
// }

// // TestReadSample testing read Sample
// func TestReadSample(t *testing.T) {
// 	ss, err := NewSamples(WithoutParseFlags())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	tss := []ReadSample{
// 		{
// 			name:    "read undefined sample",
// 			id:      5,
// 			expCode: http.StatusNotFound,
// 			expBody: fmt.Sprintf(`{"message":"%v"}`, ss.bundle.Errors[pool.ErrSampleNotFound]),
// 		},
// 		{
// 			name:    "read sample",
// 			id:      4,
// 			expCode: http.StatusOK,
// 			expBody: fmt.Sprintf(`{"title":"%s"}`, titleSample),
// 		},
// 	}
// 	for _, ts := range tss {
// 		t.Run(ts.name, func(t *testing.T) {
// 			req := httptest.NewRequest(http.MethodGet, "/", nil)
// 			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 			rec := httptest.NewRecorder()
// 			c := ss.bundle.Router.NewContext(req, rec)
// 			c.SetPath("/:id")
// 			c.SetParamNames("id")
// 			c.SetParamValues(fmt.Sprint(ts.id))
// 			if assert.NoError(t, ss.handlers.Read(c)) {
// 				assert.Equal(t, ts.expCode, rec.Code)
// 				assert.Equal(t, ts.expBody, strings.TrimSpace(rec.Body.String()))
// 			}
// 		})
// 	}
// }

// type PutSample struct {
// 	name    string
// 	id      int
// 	in      *entities.Sample
// 	expCode int
// 	expBody string
// }

// // TestPutSample testing update Sample
// func TestPutSample(t *testing.T) {
// 	ss, err := NewSamples(WithoutParseFlags())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	tss := []PutSample{
// 		{
// 			name: "update undefined sample",
// 			id:   5,
// 			in: &entities.Sample{
// 				Title: updatedTitleSample,
// 			},
// 			expCode: http.StatusNotFound,
// 			expBody: fmt.Sprintf(`{"message":"%v"}`, ss.bundle.Errors[pool.ErrSampleNotFound]),
// 		},
// 		{
// 			name: "complete update sample",
// 			id:   4,
// 			in: &entities.Sample{
// 				Title: updatedTitleSample,
// 			},
// 			expCode: http.StatusOK,
// 			expBody: fmt.Sprintf(`{"title":"%s"}`, updatedTitleSample),
// 		},
// 	}
// 	for _, ts := range tss {
// 		t.Run(ts.name, func(t *testing.T) {
// 			q := make(url.Values)
// 			q.Set("title", ts.in.Title)
// 			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/?%s", q.Encode()), nil)
// 			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 			rec := httptest.NewRecorder()
// 			c := ss.bundle.Router.NewContext(req, rec)
// 			c.SetPath("/:id")
// 			c.SetParamNames("id")
// 			c.SetParamValues(fmt.Sprint(ts.id))
// 			if assert.NoError(t, ss.handlers.Update(c)) {
// 				assert.Equal(t, ts.expCode, rec.Code)
// 				assert.Equal(t, ts.expBody, strings.TrimSpace(rec.Body.String()))
// 			}
// 		})
// 	}
// }

// type DelSample struct {
// 	name    string
// 	id      int
// 	expCode int
// 	expBody string
// }

// // TestDeleteSample testing delete Sample
// func TestDeleteSample(t *testing.T) {
// 	ss, err := NewSamples(WithoutParseFlags())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	tss := []DelSample{
// 		{
// 			name:    "delete undefined sample",
// 			id:      5,
// 			expCode: http.StatusNotFound,
// 			expBody: fmt.Sprintf(`{"message":"%v"}`, ss.bundle.Errors[pool.ErrSampleNotFound]),
// 		},
// 		{
// 			name:    "complete delete sample",
// 			id:      4,
// 			expCode: http.StatusOK,
// 			expBody: `{}`,
// 		},
// 	}
// 	for _, ts := range tss {
// 		t.Run(ts.name, func(t *testing.T) {
// 			req := httptest.NewRequest(http.MethodPut, "/", nil)
// 			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 			rec := httptest.NewRecorder()
// 			c := ss.bundle.Router.NewContext(req, rec)
// 			c.SetPath("/:id")
// 			c.SetParamNames("id")
// 			c.SetParamValues(fmt.Sprint(ts.id))
// 			if assert.NoError(t, ss.handlers.Delete(c)) {
// 				assert.Equal(t, ts.expCode, rec.Code)
// 				assert.Equal(t, ts.expBody, strings.TrimSpace(rec.Body.String()))
// 			}
// 		})
// 	}
// }
