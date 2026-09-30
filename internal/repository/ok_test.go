//go:build unit2

package repository

import (
	"log"
	"testing"

	_ "github.com/skvdmt/skvdmt-back/testing_init"
)

// TestOK Unit2 тест OK.
func TestOK(t *testing.T) {
	t.Run("ok test", func(t *testing.T) {
		p, err := NewPrepareRepo()
		if err != nil {
			log.Fatal(err)
		}
		// if err := p.Start(); err != nil {
		// 	log.Fatal(err)
		// }
		// defer func() {
		// 	if err := p.Stop(); err != nil {
		// 		log.Fatal(err)
		// 	}
		// }()
	})
}
