//go:build unit

package repository

import (
	"context"
	"fmt"
	"log"
	"testing"

	_ "github.com/skvdmt/skvdmt-back/testing_init"
)

// TestExampleSources Unit тест получения исходников примера.
func TestExampleSources(t *testing.T) {
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

	title := "Authentication"
	id, err := p.getExampleIdByTitle(title)
	if err != nil {
		log.Fatal(err)
	}

	expecteds := []string{
		"https://github.com/skvdmt/secret-front",
		"https://github.com/skvdmt/secret-back",
		"https://github.com/skvdmt/auth-back",
		"https://github.com/skvdmt/valid-access-back",
		"https://github.com/skvdmt/jwt",
	}
	t.Run(fmt.Sprintf("example %s sources", title), func(t *testing.T) {
		actuals, err := p.app.exampleSources(context.Background(), *id)
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
				t.Errorf("\n\nerror: example %s source %s not found\n\n", title, expected)
			}
		}
	})
}
