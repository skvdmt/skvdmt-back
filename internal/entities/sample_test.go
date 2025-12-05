package entities

// import (
// 	"log"
// 	"testing"

// 	"github.com/labstack/echo/v4"
// 	"github.com/skvdmt/skvdmt-back/internal/pool"
// 	"github.com/skvdmt/skvdmt-back/internal/pool/logger"
// 	_ "github.com/skvdmt/skvdmt-back/testing_init"
// )

// type Validate struct {
// 	Name               string
// 	Title              string
// 	ExpectedErrMessage error
// }

// func TestValidate(t *testing.T) {
// 	l, err := logger.New()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	b, err := pool.NewBundle(l, false)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	s := NewSample(b)

// 	sls := []Validate{
// 		{
// 			Name:               "short title",
// 			Title:              "жжж",
// 			ExpectedErrMessage: b.Errors[pool.ErrTitleLen],
// 		},
// 		{
// 			Name:               "long title",
// 			Title:              "жжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжжж",
// 			ExpectedErrMessage: b.Errors[pool.ErrTitleLen],
// 		},
// 		{
// 			Name:               "ok title",
// 			Title:              "жжжж",
// 			ExpectedErrMessage: nil,
// 		},
// 	}
// 	for _, sl := range sls {
// 		t.Run(sl.Name, func(t *testing.T) {
// 			s.Title = sl.Title
// 			err = s.Validate()
// 			if err != nil {
// 				e, ok := err.(*echo.HTTPError)
// 				if !ok {
// 					log.Fatalf("can't convert error: %v to *echo.HTTPError", err)
// 				}
// 				if sl.ExpectedErrMessage != e.Message {
// 					t.Fatalf("Expected %v got %v", sl.ExpectedErrMessage, err)
// 				}
// 			} else if sl.ExpectedErrMessage != nil {
// 				t.Fatalf("Expected %v got %v", sl.ExpectedErrMessage, err)
// 			}
// 		})
// 	}
// }
