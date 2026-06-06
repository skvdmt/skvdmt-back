package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skvdmt/skvdmt-back/internal/model"
	"github.com/skvdmt/skvdmt-back/internal/usecase"
	erw "github.com/skvdmt/skvdmt-back/pkg/errwrap"
)

const (
	pkg = "delivery"
	app = "app"
)

// App Транспортный слой.
type App struct {
	usecase Usecase
}

// NewApp Конструктор.
func NewApp(ctx context.Context) (*App, error) {
	model.Logs.Info.Info("delivery layer creating")
	// Создание сервисного слоя.
	uc, err := usecase.NewApp(ctx)
	if err != nil {
		return nil, err
	}
	return &App{
		usecase: uc,
	}, nil
}

// Start Запуск.
func (a *App) Start(ctx context.Context) error {
	model.Logs.Info.Info("delivery layer starting")
	return a.usecase.Start(ctx)
}

// Stop Остановка транспортного слоя.
func (a *App) Stop(ctx context.Context) error {
	// Остановка сервисного слоя.
	if err := a.usecase.Stop(ctx); err != nil {
		return err
	}
	model.Logs.Info.Info("delivery layer stopped")
	return nil
}

// Text Обработчик запроса текста.
func (a *App) Text(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("id")
	txt, err := a.usecase.Text(r.Context(), name)
	if err != nil {
		a.errorHandle(w, err)
		return
	}
	model.Logs.Info.Info(fmt.Sprintf("get text with name: %s", name))
	a.sendJSON(w, http.StatusOK, txt)
}

// Technologies Обработчик запроса технологий.
func (a *App) Technologies(w http.ResponseWriter, r *http.Request) {
	tls, err := a.usecase.Technologies(r.Context())
	if err != nil {
		a.errorHandle(w, err)
		return
	}
	model.Logs.Info.Info("get technologies")
	a.sendJSON(w, http.StatusOK, tls)
}

// Examples Обработчик запроса примеров.
func (a *App) Examples(w http.ResponseWriter, r *http.Request) {
	els, err := a.usecase.Examples(r.Context())
	if err != nil {
		a.errorHandle(w, err)
		return
	}
	model.Logs.Info.Info("get examples")
	a.sendJSON(w, http.StatusOK, els)
}

// Software Обработчик запроса программ.
func (a *App) Software(w http.ResponseWriter, r *http.Request) {
	sfw, err := a.usecase.Software(r.Context())
	if err != nil {
		a.errorHandle(w, err)
		return
	}
	model.Logs.Info.Info("get software")
	a.sendJSON(w, http.StatusOK, sfw)
}

// Libs Обработчик запроса бибилиотек.
func (a *App) Libs(w http.ResponseWriter, r *http.Request) {
	lbs, err := a.usecase.Libs(r.Context())
	if err != nil {
		a.errorHandle(w, err)
		return
	}
	model.Logs.Info.Info("get libs")
	a.sendJSON(w, http.StatusOK, lbs)
}

// Links Обработчик запроса ссылкок.
func (a *App) Links(w http.ResponseWriter, r *http.Request) {
	lks, err := a.usecase.Links(r.Context())
	if err != nil {
		a.errorHandle(w, err)
		return
	}
	model.Logs.Info.Info("get links")
	a.sendJSON(w, http.StatusOK, lks)
}

// sendJSON Отправка ответа в JSON.
func (a *App) sendJSON(w http.ResponseWriter, code int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(value)
}

// errorHandle Логирование и отправка ошибок.
func (a *App) errorHandle(w http.ResponseWriter, err error) {
	const m = "errorHandle"
	var e *erw.ErrorWrapper
	var ok bool
	if e, ok = err.(*erw.ErrorWrapper); !ok {
		e = erw.New(erw.Internal(
			erw.Location(pkg, app, m),
			erw.Error(fmt.Errorf(
				"%v; %v dosent match the type *errwrap.ErrorWrapper",
				model.Errs[model.ErrConvertionError], err),
			),
		))
	}
	switch {
	case 400 >= e.Code() && e.Code() <= 499:
		model.Logs.Info.Info(fmt.Sprintf("%v", e.Detailed()))
	case 500 >= e.Code() && e.Code() <= 599:
		model.Logs.Info.Info(fmt.Sprintf("%v", e.Detailed()))
	}
	a.sendJSON(w, e.Code(), map[string]string{"message": e.Message()})
}
