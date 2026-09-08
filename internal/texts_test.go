//go:build integration

package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/skvdmt/skvdmt-back/internal/entities"
	"github.com/skvdmt/skvdmt-back/internal/model"
	_ "github.com/skvdmt/skvdmt-back/testing_init"
)

// TestGetTexts Тестирование получения текстов.
func TestGetTexts(t *testing.T) {
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

	textTests := []struct {
		name     string
		expected string
	}{
		{
			name:     "main",
			expected: "Dmitry Skidanov",
		},
		{
			name:     "tech",
			expected: "Main technologies",
		},
		{
			name:     "exam",
			expected: "Examples",
		},
		{
			name:     "favo",
			expected: "Favorites",
		},
		{
			name:     "soft",
			expected: "Software for development",
		},
		{
			name:     "libs",
			expected: "Golang external libraries i like",
		},
		{
			name:     "prof",
			expected: "full stack engineer.",
		},
		{
			name:     "abou",
			expected: fmt.Sprintf("Dmitry Skidanov — full stack engineer %d", time.Now().Year()),
		},
		{
			name:     "lock",
			expected: "Russian Federation, Moscow",
		},
	}

	for _, test := range textTests {
		t.Run(test.name, func(t *testing.T) {
			// Запрос
			resp, err := http.Get(fmt.Sprintf(
				"http://localhost:%d/text/%s",
				model.Config.Server.Port,
				test.name,
			))
			if err != nil {
				log.Fatal(err)
			}
			// Проверка статуса
			if http.StatusOK != resp.StatusCode {
				t.Errorf(
					"error: status code not equal\n expected: %d\n actual: %d\n",
					http.StatusOK, resp.StatusCode)
			}
			// Получение тела ответа
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			var txt entities.Text
			if err := json.Unmarshal(buf.Bytes(), &txt); err != nil {
				log.Fatal(err)
			}
			// Проверка тела ответа.
			if test.expected != txt.Text {
				t.Errorf(
					"\n\nerror: text not equal\n expected: %s\n actual: %s\n\n",
					test.expected, txt.Text)
			}
		})
	}
}
