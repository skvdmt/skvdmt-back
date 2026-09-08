//go:build integration

package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/skvdmt/skvdmt-back/internal/entities"
	"github.com/skvdmt/skvdmt-back/internal/model"
	_ "github.com/skvdmt/skvdmt-back/testing_init"
)

// TestGetTechnologies Тестирование получения технологий.
func TestGetTechnologies(t *testing.T) {
	p, err := NewPrepareApp()
	if err != nil {
		log.Fatal(err)
	}
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := p.Stop(); err != nil {
			log.Fatal(err)
		}
	}()

	// Запрос
	resp, err := http.Get(fmt.Sprintf(
		"http://localhost:%d/technologies",
		model.Config.Server.Port,
	))
	if err != nil {
		log.Fatal(err)
	}
	t.Run("technologies", func(t *testing.T) {
		// Проверка статуса
		if http.StatusOK != resp.StatusCode {
			t.Errorf(
				"error: status code not equal\n expected: %d\n actual: %d\n",
				http.StatusOK, resp.StatusCode)
		}
	})

	// Получение тела ответа
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var actuals []*entities.Technology
	if err := json.Unmarshal(buf.Bytes(), &actuals); err != nil {
		log.Fatal(err)
	}

	expecteds := []string{
		"Go", "Postgres", "Docker", "REST API", "gRPC",
		"Git", "CI/CD", "JavaScript", "Vue",
	}
	for _, expected := range expecteds {
		t.Run(fmt.Sprintf("%s %s", "technology", expected), func(t *testing.T) {
			f := false
			for _, actual := range actuals {
				// Проверка тела ответа.
				if expected == actual.Title {
					f = true
					break
				}
			}
			if !f {
				t.Errorf("\n\nerror: technologi %s not found\n\n", expected)
			}
		})
	}
}
