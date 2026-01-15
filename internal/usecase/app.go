package usecase

import (
	"context"

	"github.com/skvdmt/skvdmt-back/internal/entities"
	"github.com/skvdmt/skvdmt-back/internal/model"
	"github.com/skvdmt/skvdmt-back/internal/repository"
)

// App Сервисный слой.
type App struct {
	// Репозиторный слой.
	repository Repository
}

// NewHome Конструктор.
func NewApp(ctx context.Context) (*App, error) {
	model.Logs.Info.Info("usecase layer creating")
	// Создание репозиторного слоя.
	rep, err := repository.NewApp(ctx)
	if err != nil {
		return nil, err
	}
	return &App{
		repository: rep,
	}, nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	// Остановка репозиторного слоя.
	if err := a.repository.Stop(ctx); err != nil {
		return err
	}
	model.Logs.Info.Info("usecase layer stopped")
	return nil
}

// Text Сервис текстов.
func (a *App) Text(ctx context.Context, name string) (*entities.Text, error) {
	return a.repository.Text(ctx, name)
}

// Technologies Сервис технологий.
func (a *App) Technologies(ctx context.Context) ([]*entities.Technology, error) {
	return a.repository.Technologies(ctx)
}

// Examples Сервис примеров.
func (a *App) Examples(ctx context.Context) ([]*entities.Example, error) {
	return a.repository.Examples(ctx)
}

// Software Сервис программ.
func (a *App) Software(ctx context.Context) ([]*entities.Software, error) {
	return a.repository.Software(ctx)
}

// Libs Сервис библиотек.
func (a *App) Libs(ctx context.Context) ([]*entities.Lib, error) {
	return a.repository.Libs(ctx)
}

// Links Сервис ссылок.
func (a *App) Links(ctx context.Context) ([]*entities.Link, error) {
	return a.repository.Links(ctx)
}
