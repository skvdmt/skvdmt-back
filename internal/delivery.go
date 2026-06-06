package internal

import (
	"context"
	"net/http"
)

// Delivery Интерфейс транспортного слоя.
type Delivery interface {
	// Запуск.
	Start(ctx context.Context) error
	// Остановка.
	Stop(ctx context.Context) error
	// Текстовая информация.
	Text(http.ResponseWriter, *http.Request)
	// Техннологии.
	Technologies(http.ResponseWriter, *http.Request)
	// Примеры.
	Examples(http.ResponseWriter, *http.Request)
	// Програмное обеспечение.
	Software(http.ResponseWriter, *http.Request)
	// Библиотеки.
	Libs(http.ResponseWriter, *http.Request)
	// Ссылки.
	Links(http.ResponseWriter, *http.Request)
}
