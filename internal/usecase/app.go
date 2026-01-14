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
func NewApp() (*App, error) {
	model.Logs.Info.Info("usecase layer creating")
	// Создание репозиторного слоя.
	rep, err := repository.NewApp()
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
func (a *App) Text(c context.Context, name string) (*entities.Text, error) {
	return a.repository.Text(c, name)
}

// Technologies Сервис технологий.
func (a *App) Technologies(c context.Context) (*[]entities.Technology, error) {
	return a.repository.Technologies(c)
}

// Examples Сервис примеров.
func (a *App) Examples(c context.Context) (*[]entities.Example, error) {
	return a.repository.Examples(c)
}

// Software Сервис программ.
func (a *App) Software(c context.Context) (*[]entities.Software, error) {
	return a.repository.Software(c)
}

// Libs Сервис библиотек.
func (a *App) Libs(c context.Context) (*[]entities.Lib, error) {
	return a.repository.Libs(c)
}

// Links Сервис ссылок.
func (a *App) Links(c context.Context) (*[]entities.Link, error) {
	return a.repository.Links(c)
}
