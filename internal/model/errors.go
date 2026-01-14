package model

import (
	"fmt"
)

// Errors Глобальный канал ошибок.
var Errors chan error

// Errs Глобальная переменная с картой описания ошибок.
var Errs errs

// errs Карта с описание ошибок.
type errs map[int]error

const (
	ErrTextNotFound = iota + 1
	ErrIncorrectTextId
	ErrDatabase
	ErrConvertionError
	ErrConvertionCache
)

// LoadErrors загрузка описания ошибок.
func LoadErrors() error {
	e := make(errs)
	e[ErrTextNotFound] = fmt.Errorf("text not found")
	e[ErrIncorrectTextId] = fmt.Errorf("incorrect text id")
	e[ErrDatabase] = fmt.Errorf("error database")
	e[ErrConvertionError] = fmt.Errorf("can't conversion error")
	e[ErrConvertionCache] = fmt.Errorf("can't conversion cache")
	Errs = e
	return nil
}
