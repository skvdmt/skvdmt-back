//go:build unit

package repository

import (
	"context"
	"fmt"
	"log"
	"testing"

	_ "github.com/skvdmt/skvdmt-back/testing_init"
)

// TestExampleTechnologies Unit тест получения технологий примера.
func TestExampleTechnologies(t *testing.T) {
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

	title := "Chess game"
	id, err := p.getExampleIdByTitle(title)
	if err != nil {
		log.Fatal(err)
	}

	expecteds := []string{
		"Go", "Docker", "Git", "Vue",
	}
	t.Run(fmt.Sprintf("example %s technologies", title), func(t *testing.T) {
		actuals, err := p.app.exampleTechnologies(context.Background(), *id)
		if err != nil {
			log.Fatal(err)
		}
		for _, expected := range expecteds {
			f := false
			for _, actual := range *actuals {
				if expected == actual.Title {
					f = true
					break
				}
			}
			if !f {
				t.Errorf("\n\nerror: example %s technology %s not found\n\n", title, expected)
			}
		}
	})
}
