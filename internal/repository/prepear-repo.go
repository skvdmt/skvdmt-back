package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/skvdmt/skvdmt-back/internal/model"
)

// PrepearRepo Подготовка репозитория для юнит тестов.
type PrepearRepo struct {
	app *App
	wg  *sync.WaitGroup
}

// NewPrepearRepo Конструктор.
func NewPrepearRepo() (*PrepearRepo, error) {
	// Создание логгера.
	if err := model.LoadLogger(); err != nil {
		return nil, err
	}
	// Загрузка конфигурации.
	if err := model.LoadConfig(); err != nil {
		return nil, err
	}
	// Создаем глобальный канал ошибок.
	model.Errors = make(chan error)
	r, err := NewApp(context.Background())
	if err != nil {
		return nil, err
	}
	return &PrepearRepo{
		app: r,
		wg:  &sync.WaitGroup{},
	}, nil
}

// Start Запуск.
func (p *PrepearRepo) Start() error {
	if err := p.app.Start(context.Background()); err != nil {
		return err
	}
	return p.HealthCheckLock()
}

func (p *PrepearRepo) HealthCheckLock() error {
	s := time.Now()
	for {
		if p.app.Ready() {
			return nil
		}
		time.Sleep(time.Millisecond * 20)
		if time.Since(s) > time.Second*3 {
			return fmt.Errorf("health check repository timeout")
		}
	}
}

func (p *PrepearRepo) Stop() error {
	return p.app.Stop(context.Background())
}

// getExampleIdByTitle Получить id примера по названию.
func (p *PrepearRepo) getExampleIdByTitle(title string) (*uuid.UUID, error) {
	examples, err := p.app.Examples(context.Background())
	if err != nil {
		return nil, err
	}
	for _, example := range examples {
		if example.Title == title {
			return &example.Id, nil
		}
	}
	return nil, fmt.Errorf("example %s not found", title)
}
