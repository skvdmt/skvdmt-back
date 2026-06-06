package usecase

import (
	"context"

	"github.com/skvdmt/skvdmt-back/internal/entities"
)

// Repository Интерфейс репозиторного слоя.
type Repository interface {
	// Запуск.
	Start(ctx context.Context) error
	// Остановка.
	Stop(ctx context.Context) error
	// Репозиторий текста.
	Text(ctx context.Context, name string) (*entities.Text, error)
	// Репозиторий технологий.
	Technologies(ctx context.Context) ([]*entities.Technology, error)
	// Репозиторий примеров.
	Examples(ctx context.Context) ([]*entities.Example, error)
	// Репозиторий программ.
	Software(ctx context.Context) ([]*entities.Software, error)
	// Репозиторий библиотек.
	Libs(ctx context.Context) ([]*entities.Lib, error)
	// Репозиторий ссылок.
	Links(ctx context.Context) ([]*entities.Link, error)
}
