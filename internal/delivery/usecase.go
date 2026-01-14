package delivery

import (
	"context"

	"github.com/skvdmt/skvdmt-back/internal/entities"
)

// Usecase Интерфейс сервисного слоя.
type Usecase interface {
	// Остановка.
	Stop(ctx context.Context) error
	// Сервис текстов.
	Text(ctx context.Context, name string) (*entities.Text, error)
	// Сервис технологий.
	Technologies(ctx context.Context) (*[]entities.Technology, error)
	// Сервис примеров.
	Examples(ctx context.Context) (*[]entities.Example, error)
	// Сервис программ.
	Software(ctx context.Context) (*[]entities.Software, error)
	// Сервис библиотек.
	Libs(ctx context.Context) (*[]entities.Lib, error)
	// Сервис ссылкок.
	Links(ctx context.Context) (*[]entities.Link, error)
}
