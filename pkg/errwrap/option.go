package errwrap

import (
	"fmt"
	"net/http"
	"strings"
)

// Option функциональная опция обертки ошибок.
type Option func(e *ErrorWrapper)

// CodeHTTP установка HTTP кода.
// Также меняет message ошибки на соответствующей ее коду краткое описание
// в случае, если оно было найдено.
// Рекомендуется указывать в качестве первого параметра опций.
// Значение по-умолчанию настраивается методом CodeHTTPDefault.
// Обнулить значение по-умолчанию можно методом CodeDefaultReset.
// Обнудение также обнуляет поле message.
func CodeHTTP(code int) Option {
	return func(e *ErrorWrapper) {
		e.code = code
		if m := http.StatusText(code); m != "" {
			e.message = m
		}
	}
}

// CodegRPC установка gRPC кода.
// Также меняет message ошибки на соответствующей ее коду краткое описание
// в случае, если оно было найдено.
// Рекомендуется указывать в качестве первого параметра опций.
// Значение по-умолчанию настраивается методом CodegRPCDefault.
// Обнулить значение по-умолчанию можно методом CodeDefaultReset.
// Обнудение также обнуляет поле message.
func CodegRPC(code int) Option {
	return func(e *ErrorWrapper) {
		e.code = code
		if m := gRPCStatusText(code); m != "" {
			e.message = m
		}
	}
}

// Message установка описания ошибки
func Message(message any) Option {
	return func(e *ErrorWrapper) {
		e.message = message
	}
}

// Internal установка расширенного описания ошибки для внутреннего использования.
// Принимает функциональные опции InternalOption
func Internal(options ...InternalOption) Option {
	return func(e *ErrorWrapper) {
		internal := &InternalWrapper{}
		for _, option := range options {
			option(internal)
		}
		e.internal = fmt.Errorf("%s", strings.Join(internal.messages, "; "))
	}
}
