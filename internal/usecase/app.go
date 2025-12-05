package usecase

import (
	"context"

	"github.com/skvdmt/skvdmt-back/internal/entities"
	"github.com/skvdmt/skvdmt-back/internal/repository"
)

// App usecase layer application
type App struct {
	repository Repository
}

// NewHome constructor usecase layer application
func NewApp() (*App, error) {
	rep, err := repository.NewApp()
	if err != nil {
		return nil, err
	}
	return &App{
		repository: rep,
	}, nil
}

// Text service homepage implementation
func (a *App) Text(c context.Context, name string) (*entities.Text, error) {
	return a.repository.Text(c, name)
}

// Technologies service homepage implementation
func (a *App) Technologies(c context.Context) (*[]entities.Technology, error) {
	return a.repository.Technologies(c)
}

// Examples service homepage implementation
func (a *App) Examples(c context.Context) (*[]entities.Example, error) {
	return a.repository.Examples(c)
}

// Software service homepage implementation
func (a *App) Software(c context.Context) (*[]entities.Software, error) {
	return a.repository.Software(c)
}

// Libs service homepage implementation
func (a *App) Libs(c context.Context) (*[]entities.Lib, error) {
	return a.repository.Libs(c)
}

// Links service homepage implementation
func (a *App) Links(c context.Context) (*[]entities.Link, error) {
	return a.repository.Links(c)
}
