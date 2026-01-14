package usecase

import (
	"context"

	"github.com/skvdmt/skvdmt-back/internal/entities"
)

// Repository Интерфейс репозиторного слоя.
type Repository interface {
	// Остановка.
	Stop(ctx context.Context) error
	// Репозиторий текста.
	Text(c context.Context, name string) (*entities.Text, error)
	// Репозиторий технологий.
	Technologies(c context.Context) (*[]entities.Technology, error)
	// Репозиторий примеров.
	Examples(c context.Context) (*[]entities.Example, error)
	// Репозиторий программ.
	Software(c context.Context) (*[]entities.Software, error)
	// Репозиторий библиотек.
	Libs(c context.Context) (*[]entities.Lib, error)
	// Репозиторий ссылок.
	Links(c context.Context) (*[]entities.Link, error)
}
