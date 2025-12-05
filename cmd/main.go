package main

import (
	"os"

	"github.com/skvdmt/skvdmt-back/internal"
	"github.com/skvdmt/skvdmt-back/internal/model"
)

func main() {
	// making logger
	err := model.LoadLogger()
	if err != nil {
		panic(err)
	}

	model.Logs.Error.Error("test error")

	// making app
	hp, err := internal.NewApp()
	if err != nil {
		model.Logs.Error.Error(err.Error())
		os.Exit(1)
	}
	// starting app
	if err = hp.Start(); err != nil {
		model.Logs.Error.Error(err.Error())
		os.Exit(1)
	}
}
