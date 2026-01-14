package errwrap

import (
	"fmt"
	"strings"
)

// InternalOption функциональные опции внутреннего описания ошибки.
type InternalOption func(internal *InternalWrapper)

// InternalData информация функциональных опций
// внутреннего состояния обертки ошибок.
type InternalWrapper struct {
	messages []string
}

// Location указания места возникновения ошибки.
func Location(locations ...string) InternalOption {
	return func(internal *InternalWrapper) {
		internal.messages = append(internal.messages,
			fmt.Sprintf("location: %s", strings.Join(locations, ".")))
	}
}

// Error передача внутренней ошибки возвращенной сторонним пакетом для логирования.
func Error(err error) InternalOption {
	return func(internal *InternalWrapper) {
		internal.messages = append(internal.messages, fmt.Sprintf("error: %s", err))
	}
}

// SQL передача sql запроса и его аргументов для логирования.
func SQL(query string, args ...any) InternalOption {
	return func(internal *InternalWrapper) {
		if len(args) > 0 {
			a := []string{}
			for _, arg := range args {
				a = append(a, fmt.Sprintf("%v", arg))
			}
			internal.messages = append(internal.messages,
				fmt.Sprintf("query: '%s' args: [%s]", query, strings.Join(a, ", ")))
			return
		}
		internal.messages = append(internal.messages,
			fmt.Sprintf("query: %s", query))
	}
}
