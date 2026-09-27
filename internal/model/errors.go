package model

import (
	"fmt"
)

// Errors Глобальный канал ошибок.
var Errors chan error

// Errs Глобальная переменная с картой описания ошибок.
var Errs errs

// Errs Карта с описанием ошибок.
type errs map[int]error

const (
	ErrTextNotFound = iota + 1
	ErrIncorrectTextId
	ErrDatabase
	ErrConversionError
	ErrConversionCache
)

// CreateErrors Создание описаний ошибок.
func CreateErrors() {
	Logs.Info.Info("Errors map creating")
	e := make(errs)
	e[ErrTextNotFound] = fmt.Errorf("text not found")
	e[ErrIncorrectTextId] = fmt.Errorf("incorrect text id")
	e[ErrDatabase] = fmt.Errorf("error database")
	e[ErrConversionError] = fmt.Errorf("can't conversion error")
	e[ErrConversionCache] = fmt.Errorf("can't conversion cache")
	Errs = e
}
