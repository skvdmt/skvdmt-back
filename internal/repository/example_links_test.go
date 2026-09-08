//go:build unit

package repository

import (
	"context"
	"fmt"
	"log"
	"testing"

	_ "github.com/skvdmt/skvdmt-back/testing_init"
)

// TestExampleLinks Unit тест получения ссылок примера.
func TestExampleLinks(t *testing.T) {
	p, err := NewPrepearRepo()
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

	title := "Telegram bot"
	id, err := p.getExampleIdByTitle(title)
	if err != nil {
		log.Fatal(err)
	}

	expecteds := []string{
		"https://t.me/skidanovdima_msgs_bot",
		"https://msgs.skvdmt.ru/",
	}
	t.Run(fmt.Sprintf("example %s links", title), func(t *testing.T) {
		actuals, err := p.app.exampleLinks(context.Background(), *id)
		if err != nil {
			log.Fatal(err)
		}
		for _, expected := range expecteds {
			f := false
			for _, actual := range *actuals {
				if expected == actual.Url {
					f = true
					break
				}
			}
			if !f {
				t.Errorf("\n\nerror: example %s link %s not found\n\n", title, expected)
			}
		}
	})
}
