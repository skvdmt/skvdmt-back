package internal

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	erw "github.com/skvdmt/errwrap"
	"github.com/skvdmt/skvdmt-back/internal/delivery"
	"github.com/skvdmt/skvdmt-back/internal/model"
)

const (
	defaultParseFlags = true
	name              = "skvdmt-back"
	version           = "1.0.0"
)

// App facade application struct
type App struct {
	exit     chan os.Signal
	errors   chan error
	server   *http.Server
	delivery Delivery
}

// NewHome app constructor
func NewApp(options ...Option) (*App, error) {
	a := &App{}
	data, err := a.optionsProcessing(options...)
	if err != nil {
		return nil, err
	}
	if err := model.LoadConfig(); err != nil {
		return nil, err
	}
	if err := model.LoadErrors(); err != nil {
		return nil, err
	}

	// set default error status code
	erw.CodeHTTPDefault(500)

	// flags processing
	if data.parseFlags {
		var flagVersion bool
		flag.BoolVar(&flagVersion, "version", false, "print application version")
		flag.Parse()
		switch {
		case flagVersion:
			fmt.Printf("%s v%s\n", name, version)
			if err := model.Logs.Close(); err != nil {
				return nil, err
			}
			os.Exit(0)
		}
	}

	a.delivery, err = delivery.NewApp()
	if err != nil {
		return nil, err
	}

	a.exit = make(chan os.Signal)
	a.errors = make(chan error)

	// making http server
	const (
		defaultTimeout        = 10
		defaultMaxHeaderBytes = 1 << 20 // 1Mb
	)
	a.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", model.Config.Server.Port),
		Handler:        a.delivery.Router(),
		ReadTimeout:    defaultTimeout * time.Second,
		WriteTimeout:   defaultTimeout * time.Second,
		MaxHeaderBytes: defaultMaxHeaderBytes,
	}
	return a, nil
}

// Start launching the app
func (a *App) Start() error {
	// start rest api server
	go func() {
		model.Logs.Info.Info(fmt.Sprintf(
			"skvdmt-back application REST api server starting on %d port",
			model.Config.Server.Port,
		))
		a.errors <- a.server.ListenAndServe()
	}()
	return a.play()
}

// optionsProcessing processing functional options
func (a *App) optionsProcessing(options ...Option) (*OptionData, error) {
	d := &OptionData{
		// setting default options
		parseFlags: defaultParseFlags,
	}
	for _, option := range options {
		if err := option(d); err != nil {
			return nil, err
		}
	}
	return d, nil
}

// play replaying an app until critical errors or interrupt signals
func (a *App) play() error {
	signal.Notify(a.exit, syscall.SIGTERM)
	select {
	case err := <-a.errors:
		return err
	case <-a.exit:
		return a.stop()
	}
}

// stop homepage app
func (a *App) stop() error {
	// stop rest server
	if err := a.server.Shutdown(context.Background()); err != nil {
		return err
	}
	model.Logs.Info.Info("skvdmt-back application http server stopped")
	// close resources on all layers application
	if err := a.delivery.Close(); err != nil {
		return err
	}
	if err := model.Logs.Close(); err != nil {
		return err
	}
	model.Logs.Info.Info("skvdmt-back all resources closed")
	// exit
	os.Exit(0)
	return nil
}
