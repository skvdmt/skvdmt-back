//go:build unit

package repository

import (
	"context"
	"log"
	"testing"

	_ "github.com/skvdmt/skvdmt-back/testing_init"
)

// TestUpdates Unit тест проверки обновлений данных.
func TestUpdates(t *testing.T) {
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

	t.Run("update technologies", func(t *testing.T) {
		e, err := p.app.Technologies(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		if len(e) == 0 {
			t.Errorf("\n\nerror: technologies not found after update\n\n")
		}
	})

	t.Run("update examples", func(t *testing.T) {
		e, err := p.app.Examples(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		if len(e) == 0 {
			t.Errorf("\n\nerror: examples not found after update\n\n")
		}
	})

	t.Run("update software", func(t *testing.T) {
		e, err := p.app.Software(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		if len(e) == 0 {
			t.Errorf("\n\nerror: software not found after update\n\n")
		}
	})

	t.Run("update libs", func(t *testing.T) {
		e, err := p.app.Libs(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		if len(e) == 0 {
			t.Errorf("\n\nerror: libs not found after update\n\n")
		}
	})

	t.Run("update links", func(t *testing.T) {
		e, err := p.app.Links(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		if len(e) == 0 {
			t.Errorf("\n\nerror: libs not found after update\n\n")
		}
	})

}
