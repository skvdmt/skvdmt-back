package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	erw "github.com/skvdmt/errwrap"
	"github.com/skvdmt/skvdmt-back/internal/model"
	"github.com/skvdmt/skvdmt-back/internal/usecase"
)

// App trasport layer application
type App struct {
	router  *http.ServeMux
	usecase Usecase
}

// NewApp constructor trasport layer application
func NewApp() (*App, error) {
	uc, err := usecase.NewApp()
	if err != nil {
		return nil, err
	}
	a := &App{
		usecase: uc,
		router:  http.NewServeMux(),
	}
	a.setRoutes()
	return a, nil
}

// Close
func (a *App) Close() error {
	return nil
}

const (
	pkg = "delivery"
	app = "app"

	url_text         = "/text/{id}"
	url_technologies = "/technologies"
	url_examples     = "/examples"
	url_software     = "/software"
	url_libs         = "/libs"
	url_links        = "/links"
)

// SetRoutes setting up the main segment router of the application
func (a *App) setRoutes() {
	bu := model.Config.Server.BaseUrl
	a.router.HandleFunc(fmt.Sprintf("GET %s", path.Join(bu, url_text)), a.Text)
	a.router.HandleFunc(fmt.Sprintf("GET %s", path.Join(bu, url_technologies)), a.Technologies)
	a.router.HandleFunc(fmt.Sprintf("GET %s", path.Join(bu, url_examples)), a.Examples)
	a.router.HandleFunc(fmt.Sprintf("GET %s", path.Join(bu, url_software)), a.Software)
	a.router.HandleFunc(fmt.Sprintf("GET %s", path.Join(bu, url_libs)), a.Libs)
	a.router.HandleFunc(fmt.Sprintf("GET %s", path.Join(bu, url_links)), a.Links)
}

// Router
func (a *App) Router() http.Handler {
	return a.router
}

// Text handler homepage implementation
func (a *App) Text(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("id")
	txt, err := a.usecase.Text(r.Context(), name)
	if err != nil {
		a.errorHandle(w, err)
	}
	model.Logs.Info.Info(fmt.Sprintf("get text with name: %s; text: %s;", name, txt))
	a.sendJSON(w, http.StatusOK, txt)
}

// Technologies handler homepage implementation
func (a *App) Technologies(w http.ResponseWriter, r *http.Request) {
	tls, err := a.usecase.Technologies(r.Context())
	if err != nil {
		a.errorHandle(w, err)
	}
	model.Logs.Info.Info("get technologies list")
	a.sendJSON(w, http.StatusOK, tls)
}

// Examples handler homepage implementation
func (a *App) Examples(w http.ResponseWriter, r *http.Request) {
	els, err := a.usecase.Examples(r.Context())
	if err != nil {
		a.errorHandle(w, err)
	}
	model.Logs.Info.Info("get examples list")
	a.sendJSON(w, http.StatusOK, els)
}

// Software handler homepage implementation
func (a *App) Software(w http.ResponseWriter, r *http.Request) {
	sfw, err := a.usecase.Software(r.Context())
	if err != nil {
		a.errorHandle(w, err)
	}
	model.Logs.Info.Info("get software list")
	a.sendJSON(w, http.StatusOK, sfw)
}

// Libs handler homepage implementation
func (a *App) Libs(w http.ResponseWriter, r *http.Request) {
	lbs, err := a.usecase.Libs(r.Context())
	if err != nil {
		a.errorHandle(w, err)
	}
	model.Logs.Info.Info("get libs list")
	a.sendJSON(w, http.StatusOK, lbs)
}

// Links handler homepage implementation
func (a *App) Links(w http.ResponseWriter, r *http.Request) {
	lks, err := a.usecase.Links(r.Context())
	if err != nil {
		a.errorHandle(w, err)
	}
	model.Logs.Info.Info("get footer links list")
	a.sendJSON(w, http.StatusOK, lks)
}

func (a *App) sendJSON(w http.ResponseWriter, code int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(value)
}

// errorHandle sending and logging error
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
