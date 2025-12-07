package model

import (
	"fmt"
)

var Errs errs

// errs
type errs map[int]error

const (
	ErrTextNotFound = iota + 1
	ErrIncorrectTextId
	ErrDatabase
	ErrConvertionError
	ErrConvertionCache
	ErrConvertionResponse
)

// LoadErrors
func LoadErrors() error {
	Errs := make(errs)
	Errs[ErrTextNotFound] = fmt.Errorf("text not found")
	Errs[ErrIncorrectTextId] = fmt.Errorf("incorrect text id")
	Errs[ErrDatabase] = fmt.Errorf("error database")
	Errs[ErrConvertionError] = fmt.Errorf("can't conversion error")
	Errs[ErrConvertionCache] = fmt.Errorf("can't conversion cache")
	Errs[ErrConvertionResponse] = fmt.Errorf("can't conversion response")
	return nil
}
