package delivery

import (
	"context"

	"github.com/skvdmt/skvdmt-back/internal/entities"
)

// Usecase application businer logic interface
type Usecase interface {
	Text(c context.Context, name string) (*entities.Text, error)
	Technologies(c context.Context) (*[]entities.Technology, error)
	Examples(c context.Context) (*[]entities.Example, error)
	Software(c context.Context) (*[]entities.Software, error)
	Libs(c context.Context) (*[]entities.Lib, error)
	Links(c context.Context) (*[]entities.Link, error)
}
