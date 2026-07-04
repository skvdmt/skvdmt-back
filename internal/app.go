package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/skvdmt/skvdmt-back/internal/delivery"
	"github.com/skvdmt/skvdmt-back/internal/model"
)

const (
	defaultTimeout        = 10
	defaultMaxHeaderBytes = 1 << 20 // 1Mb
	get                   = "GET %s"
	url_text              = "/text/{id}"
	url_technologies      = "/technologies"
	url_examples          = "/examples"
	url_software          = "/software"
	url_libs              = "/libs"
	url_links             = "/links"
)

// App Основная структура приложения.
type App struct {
	// Канал сигналов операционной системы для отслеживания
	// сигналов прерывания работы приложения.
	interrupt chan os.Signal
	// Контекст приложения.
	ctx context.Context
	// Функция отмены контекста всего приложения.
	cancel context.CancelFunc
	// Транспортный слой.
	delivery Delivery
	// Роутер.
	router *http.ServeMux
	// Сервер.
	server *http.Server
	// Приложение запущено.
	started bool
	// Ошибки приложения.
	eg []error
	// Приложение уже закрывается.
	stopping bool
	// Корректное завершение горутин.
	wg *sync.WaitGroup
}

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Logs.Info.Info(fmt.Sprintf("%s creating", model.APP_NAME))
	var err error
	// Загрузка конфигурации.
	if err = model.LoadConfig(); err != nil {
		return nil, err
	}
	// Создаем глобальный канал ошибок.
	model.Errors = make(chan error)
	// Создание сервера.
	r := http.NewServeMux()
	a := &App{
		interrupt: make(chan os.Signal),
		router:    r,
		server: &http.Server{
			Addr:           fmt.Sprintf(":%d", model.Config.Server.Port),
			Handler:        r,
			ReadTimeout:    defaultTimeout * time.Second,
			WriteTimeout:   defaultTimeout * time.Second,
			MaxHeaderBytes: defaultMaxHeaderBytes,
		},
		wg: &sync.WaitGroup{},
	}
	// Создание контекста.
	a.ctx, a.cancel = context.WithCancel(context.Background())
	// Создание транспортного слоя из которого по
	// цепочки создаются остальные слои приложения.
	a.delivery, err = delivery.NewApp(a.ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Start Запуск приложения.
func (a *App) Start() error {
	go a.errorHandler()

	model.Logs.Info.Info(fmt.Sprintf("%s starting", model.APP_NAME))

	a.wg.Add(1)
	go func() {
		// Настройка и запуск сервера.
		defer a.wg.Done()
		a.routes()
		model.Logs.Info.Info(fmt.Sprintf("http server starting on %d port",
			model.Config.Server.Port))
		if err := a.server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			model.Errors <- err
		}
	}()

	a.wg.Add(1)
	go func() {
		// Запуск слоев приложения по цепочке.
		defer a.wg.Done()
		if err := a.delivery.Start(a.ctx); err != nil {
			model.Errors <- err
		}
		a.started = true
	}()
	return a.interruptHandler()
}

// errorHandle Обработчик глобального канала ошибок.
func (a *App) errorHandler() {
	model.Logs.Info.Info("error handler starting")
	for {
		err := <-model.Errors
		if err == nil {
			return
		}
		if errors.Is(err, context.Canceled) {
			continue
		}
		a.eg = append(a.eg, err)
		if !a.stopping {
			a.interrupt <- syscall.SIGTERM
		}
	}
}

// interruptHandler Обработчик сигналов остановки приложения.
func (a *App) interruptHandler() error {
	model.Logs.Info.Info("error handler starting")
	signal.Notify(a.interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-a.interrupt
	// close sources
	if err := a.stop(); err != nil {
		model.Errors <- err
	}
	model.Errors <- nil
	close(model.Errors)
	if len(a.eg) > 0 {
		return fmt.Errorf("%v", a.eg)
	}
	return nil
}

// stop Остановка приложения.
func (a *App) stop() error {
	a.stopping = true

	// Отмена контекста.
	a.cancel()
	model.Logs.Info.Info("context canceled")

	// Дождаться запуска приложения.
	for {
		if a.started {
			break
		}
	}

	// Остановка сервера.
	if err := a.server.Shutdown(context.Background()); err != nil {
		return err
	}
	model.Logs.Info.Info("http server shutdown")

	// Остановка транспортного слоя из которо по цепочке
	// останавливаются все остальные слои.
	if err := a.delivery.Stop(a.ctx); err != nil {
		return err
	}
	a.wg.Wait()
	// Закрытие канала отслеживающего сигналы
	// прерывания операционной системы.
	close(a.interrupt)
	model.Logs.Info.Info(fmt.Sprintf("%s stopped", model.APP_NAME))
	return nil
}

// routes Настройка маршрутов.
func (a *App) routes() error {
	bu := model.Config.Server.BaseUrl
	a.router.HandleFunc(fmt.Sprintf(get, path.Join(bu, url_text)), a.delivery.Text)
	a.router.HandleFunc(fmt.Sprintf(get, path.Join(bu, url_technologies)), a.delivery.Technologies)
	a.router.HandleFunc(fmt.Sprintf(get, path.Join(bu, url_examples)), a.delivery.Examples)
	a.router.HandleFunc(fmt.Sprintf(get, path.Join(bu, url_software)), a.delivery.Software)
	a.router.HandleFunc(fmt.Sprintf(get, path.Join(bu, url_libs)), a.delivery.Libs)
	a.router.HandleFunc(fmt.Sprintf(get, path.Join(bu, url_links)), a.delivery.Links)
	return nil
}
