package repository

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/skvdmt/skvdmt-back/internal/model"
)

// PrepareRepo Подготовка репозитория для юнит тестов.
type PrepareRepo struct {
	app *App
	wg  *sync.WaitGroup
}

// NewPrepareRepo Конструктор.
func NewPrepareRepo() (*PrepareRepo, error) {
	// Создание логгера.
	if err := model.CreateLogger(); err != nil {
		return nil, err
	}
	var err error
	// Загрузка конфигурации.
	model.Config, err = model.NewConfig()
	if err != nil {
		return nil, err
	}
	// Создаем глобальный канал ошибок.
	model.Errors = make(chan error)
	go func() {
		for {
			err := <-model.Errors
			model.Logs.Error.Error(err.Error())
			_ = model.Logs.Close()
			os.Exit(1)
		}
	}()
	r, err := NewApp(context.Background())
	if err != nil {
		return nil, err
	}
	return &PrepareRepo{
		app: r,
		wg:  &sync.WaitGroup{},
	}, nil
}

// Start Запуск.
func (p *PrepareRepo) Start() error {
	if err := p.app.Start(context.Background()); err != nil {
		return err
	}
	return p.HealthCheckLock()
}

func (p *PrepareRepo) HealthCheckLock() error {
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

func (p *PrepareRepo) Stop() error {
	return p.app.Stop(context.Background())
}

// getExampleIdByTitle Получить id примера по названию.
func (p *PrepareRepo) getExampleIdByTitle(title string) (*uuid.UUID, error) {
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
