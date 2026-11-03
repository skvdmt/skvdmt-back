package errwrap

import "net/http"

var defaultCode int = -1
var defaultMessage any
var defaultLog = true

// SetDefaultLog Установка значения по-умолчанию
// для параметра разрешающего запись в журнал.
func SetDefaultLog(value bool) {
	defaultLog = value
}

// CodeHTTPDefault устанавливает HTTP код статуса по-умолчанию для всех
// экземпляров оберток ошибок ErrorWrapper. Также изменяет message
// на текстовое описание статуса.
func CodeHTTPDefault(code int) {
	defaultCode = code
	if defaultMessage == nil {
		if m := http.StatusText(code); m != "" {
			defaultMessage = m
		}
	}
}

// CodegRPCDefault устанавливает gRPC код статуса по-умолчанию для всех
// экземпляров оберток ошибок ErrorWrapper. Также изменяет message
// на текстовое описание статуса.
func CodegRPCDefault(code int) {
	defaultCode = code
	if defaultMessage == nil {
		if m := gRPCStatusText(code); m != "" {
			defaultMessage = m
		}
	}
}

// CodeDefaultReset сброс кода статуса и сообщения по-умолчанию.
func CodeDefaultReset() {
	defaultCode = -1
	MessageDefaultReset()
}

// MessageDefault установка сообщения по-умолчанию.
func MessageDefault(message any) {
	defaultMessage = message
}

// MessageDefaultReset сброс сообщения и кода статуса по-умолчанию.
func MessageDefaultReset() {
	defaultMessage = nil
	CodeDefaultReset()
}
