package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/patrickmn/go-cache"
	erw "github.com/skvdmt/errwrap"
	"github.com/skvdmt/skvdmt-back/init/inserts"
	"github.com/skvdmt/skvdmt-back/internal/entities"
	"github.com/skvdmt/skvdmt-back/internal/model"
)

const (
	pkg          = "repository"
	text         = "text"
	technologies = "technologies"
	examples     = "examples"
	software     = "software"
	libs         = "libs"
	links        = "links"
	sources      = "sources"
)

// App
type App struct {
	cache *cache.Cache
	db    *pgx.Conn
}

const (
	postgres    = "postgres"
	DB_PASSWORD = "DB_PASSWORD"
)

// NewApp
func NewApp() (*App, error) {
	pwd, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		return nil, fmt.Errorf("env %s unset", DB_PASSWORD)
	}
	q := fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s",
		postgres,
		model.Config.Postgres.User,
		pwd,
		model.Config.Postgres.Host,
		model.Config.Postgres.Port,
		model.Config.Postgres.Database,
	)
	db, err := pgx.Connect(context.Background(), q)
	if err != nil {
		return nil, err
	}

	// insert default data to database tables
	inserts.InsertData(db)

	return &App{
		cache: cache.New(2*time.Minute, 20*time.Second),
		db:    db,
	}, nil
}

func (a *App) Close() error {
	return nil
}

// Text repository homepage implementation
func (a *App) Text(c context.Context, name string) (*entities.Text, error) {
	data, ok := a.cache.Get(fmt.Sprintf("%s_%s", text, name))
	if !ok {
		const query = "select id, text from texts where name = $1 order by id"
		row, err := a.db.Query(c, query, name)
		if err != nil {
			return nil, erw.New(erw.Internal(
				erw.Location(pkg, text),
				erw.Error(err),
				erw.SQL(query, name),
			))
		}
		cl, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Text])
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, erw.New(
					erw.CodeHTTP(404),
					erw.Internal(
						erw.Location(pkg, text),
						erw.Error(err),
						erw.SQL(query, name),
					))
			}
			return nil, erw.New(erw.Internal(
				erw.Location(pkg, text),
				erw.Error(err),
				erw.SQL(query, name),
			))
		}
		txt := &entities.Text{
			Id:   cl.Id,
			Text: cl.Text,
		}
		// setting to cache
		a.cache.Set(fmt.Sprintf("%s_%s", text, name), txt, cache.DefaultExpiration)
		return txt, nil
	}
	txt, ok := data.(*entities.Text)
	if !ok {
		return nil, erw.New(erw.Internal(
			erw.Location(pkg, text),
			erw.Error(model.Errs[model.ErrConvertionCache]),
		))
	}
	return txt, nil
}

// Technologies repository homepage implementation
func (a *App) Technologies(c context.Context) (*[]entities.Technology, error) {
	data, ok := a.cache.Get(technologies)
	if !ok {
		query := fmt.Sprintf(
			"select %s, %s, %s from %s order by id",
			"id", "title", "url", technologies,
		)
		rows, err := a.db.Query(c, query)
		if err != nil {
			return nil, erw.New(erw.Internal(
				erw.Location(pkg, technologies),
				erw.Error(err),
				erw.SQL(query),
			))
		}
		cls, err := pgx.CollectRows(rows, pgx.RowToStructByName[Technology])
		if err != nil {
			return nil, erw.New(erw.Internal(
				erw.Location(pkg, technologies),
				erw.Error(err),
				erw.SQL(query),
			))
		}
		tls := new([]entities.Technology)
		for _, cl := range cls {
			*tls = append(*tls, entities.Technology{
				Id:    cl.Id,
				Title: cl.Title,
				Url:   cl.Url,
			})
		}
		// setting to cache
		a.cache.Set(technologies, tls, cache.DefaultExpiration)
		return tls, nil
	}
	tls, ok := data.(*[]entities.Technology)
	if !ok {
		return nil, erw.New(erw.Internal(
			erw.Location(pkg, technologies),
			erw.Error(model.Errs[model.ErrConvertionCache]),
		))
	}
	return tls, nil
}

// Examples repository homepage implementation
func (a *App) Examples(c context.Context) (*[]entities.Example, error) {
	data, ok := a.cache.Get(examples)
	if !ok {
		query := fmt.Sprintf(
			"select %s, %s, %s, %s from %s order by id",
			"id", "name", "title", "description", examples,
		)
		rows, err := a.db.Query(c, query)
		if err != nil {
			return nil, erw.New(erw.Internal(
				erw.Location(pkg, examples),
				erw.Error(err),
				erw.SQL(query),
			))
		}
		cls, err := pgx.CollectRows(rows, pgx.RowToStructByName[Example])
		if err != nil {
			return nil, erw.New(erw.Internal(
				erw.Location(pkg, examples),
				erw.Error(err),
				erw.SQL(query),
			))
		}
		exs := new([]entities.Example)
		for _, cl := range cls {
			lks, err := a.exampleLinks(c, cl.Id)
			if err != nil {
				return nil, err
			}
			tks, err := a.exampleTechnologies(c, cl.Id)
			if err != nil {
				return nil, err
			}
			srs, err := a.exampleSources(c, cl.Id)
			if err != nil {
				return nil, err
			}
			ex := entities.Example{
				Id:           cl.Id,
				Name:         cl.Name,
				Title:        cl.Title,
				Description:  cl.Description,
				Technologies: tks,
				Sources:      srs,
			}
			if len(*lks) > 0 {
				ex.Links = lks
			}
			*exs = append(*exs, ex)
		}
		a.cache.Set(examples, exs, cache.DefaultExpiration)
		return exs, nil
	}
	exs, ok := data.(*[]entities.Example)
	if !ok {
		return nil, erw.New(erw.Internal(
			erw.Location(pkg, examples),
			erw.Error(model.Errs[model.ErrConvertionCache]),
		))
	}
	return exs, nil

}

// Software repository homepage implementation
func (a *App) Software(c context.Context) (*[]entities.Software, error) {
	data, ok := a.cache.Get(software)
	if !ok {
		query := fmt.Sprintf(
			"select %s, %s, %s from %s order by id",
			"id", "title", "url", software,
		)
		rows, err := a.db.Query(c, query)
		if err != nil {
			return nil, erw.New(erw.Internal(
				erw.Location(pkg, software),
				erw.Error(err),
				erw.SQL(query),
			))
		}
		cls, err := pgx.CollectRows(rows, pgx.RowToStructByName[Software])
		if err != nil {
			return nil, erw.New(erw.Internal(
				erw.Location(pkg, software),
				erw.Error(err),
				erw.SQL(query),
			))
		}
		sfw := new([]entities.Software)
		for _, cl := range cls {
			*sfw = append(*sfw, entities.Software{
				Id:    cl.Id,
				Title: cl.Title,
				Url:   cl.Url,
			})
		}
		// setting to cache
		a.cache.Set(software, sfw, cache.DefaultExpiration)
		return sfw, nil
	}
	sfw, ok := data.(*[]entities.Software)
	if !ok {
		return nil, erw.New(erw.Internal(
			erw.Location(pkg, software),
			erw.Error(model.Errs[model.ErrConvertionCache]),
		))
	}
	return sfw, nil
}

// Libs repository homepage implementation
func (a *App) Libs(c context.Context) (*[]entities.Lib, error) {
	data, ok := a.cache.Get(libs)
	if !ok {
		query := fmt.Sprintf(
			"select %s, %s from %s order by id",
			"id", "url", libs,
		)
		rows, err := a.db.Query(c, query)
		if err != nil {
			return nil, erw.New(erw.Internal(
				erw.Location(pkg, libs),
				erw.Error(err),
				erw.SQL(query),
			))
		}
		cls, err := pgx.CollectRows(rows, pgx.RowToStructByName[Lib])
		if err != nil {
			return nil, erw.New(erw.Internal(
				erw.Location(pkg, libs),
				erw.Error(err),
				erw.SQL(query),
			))
		}
		lbs := new([]entities.Lib)
		for _, cl := range cls {
			*lbs = append(*lbs, entities.Lib{
				Id:  cl.Id,
				Url: cl.Url,
			})
		}
		// setting to cache
		a.cache.Set(libs, lbs, cache.DefaultExpiration)
		return lbs, nil
	}
	lbs, ok := data.(*[]entities.Lib)
	if !ok {
		return nil, erw.New(erw.Internal(
			erw.Location(pkg, libs),
			erw.Error(model.Errs[model.ErrConvertionCache]),
		))
	}
	return lbs, nil
}

// Links repository homepage implementation
func (a *App) Links(c context.Context) (*[]entities.Link, error) {
	data, ok := a.cache.Get(links)
	if !ok {
		query := fmt.Sprintf(
			"select %s, %s, %s from %s order by id",
			"id", "title", "url", "footer_links",
		)
		rows, err := a.db.Query(c, query)
		if err != nil {
			return nil, erw.New(erw.Internal(
				erw.Location(pkg, links),
				erw.Error(err),
				erw.SQL(query),
			))
		}
		cls, err := pgx.CollectRows(rows, pgx.RowToStructByName[Link])
		if err != nil {
			return nil, erw.New(erw.Internal(
				erw.Location(pkg, links),
				erw.Error(err),
				erw.SQL(query),
			))
		}
		lks := new([]entities.Link)
		for _, cl := range cls {
			*lks = append(*lks, entities.Link{
				Id:    cl.Id,
				Title: cl.Title,
				Url:   cl.Url,
			})
		}
		// setting to cache
		a.cache.Set(links, lks, cache.DefaultExpiration)
		return lks, nil
	}

	lks, ok := data.(*[]entities.Link)
	if !ok {
		return nil, erw.New(erw.Internal(
			erw.Location(pkg, links),
			erw.Error(model.Errs[model.ErrConvertionCache]),
		))
	}
	return lks, nil
}

// exampleLinks getting links of example by id
func (a *App) exampleLinks(c context.Context, exampleID uuid.UUID) (*[]entities.Link, error) {
	query := fmt.Sprintf(`SELECT %s, %s, %s
			FROM %s
			LEFT JOIN %s
			ON %s = %s
			WHERE %s = $1 order by links.id`,
		"links.id", "links.title", "links.url",
		"examples_links", "links",
		"examples_links.link_id", "links.id",
		"example_id",
	)
	rows, err := a.db.Query(c, query, exampleID)
	if err != nil {
		return nil, erw.New(erw.Internal(
			erw.Location(pkg, examples, links),
			erw.Error(err),
			erw.SQL(query, exampleID),
		))
	}
	cls, err := pgx.CollectRows(rows, pgx.RowToStructByName[Link])
	if err != nil {
		return nil, erw.New(erw.Internal(
			erw.Location(pkg, examples, links),
			erw.Error(err),
			erw.SQL(query, exampleID),
		))
	}
	lks := new([]entities.Link)
	for _, cl := range cls {
		*lks = append(*lks, entities.Link{
			Id:    cl.Id,
			Title: cl.Title,
			Url:   cl.Url,
		})
	}
	return lks, nil
}

// exampleTechnologies getting technologies of example by id
func (a *App) exampleTechnologies(c context.Context, exampleID uuid.UUID) (*[]entities.Technology, error) {
	query := fmt.Sprintf(`SELECT %s, %s, %s
			FROM %s
			LEFT JOIN %s
			ON %s = %s
			WHERE %s = $1 order by technologies.id`,
		"technologies.id", "technologies.title", "technologies.url",
		"examples_technologies", "technologies",
		"examples_technologies.technology_id", "technologies.id",
		"example_id",
	)
	rows, err := a.db.Query(c, query, exampleID)
	if err != nil {
		return nil, erw.New(erw.Internal(
			erw.Location(pkg, examples, technologies),
			erw.Error(err),
			erw.SQL(query, exampleID),
		))
	}
	cls, err := pgx.CollectRows(rows, pgx.RowToStructByName[Technology])
	if err != nil {
		return nil, erw.New(erw.Internal(
			erw.Location(pkg, examples, technologies),
			erw.Error(err),
			erw.SQL(query, exampleID),
		))
	}
	tks := new([]entities.Technology)
	for _, cl := range cls {
		*tks = append(*tks, entities.Technology{
			Id:    cl.Id,
			Title: cl.Title,
			Url:   cl.Url,
		})
	}
	return tks, nil
}

// exampleSources getting sources of example by id
func (a *App) exampleSources(c context.Context, exampleID uuid.UUID) (*[]entities.Source, error) {
	query := fmt.Sprintf(`SELECT %s, %s
			FROM %s
			LEFT JOIN %s
			ON %s = %s
			WHERE %s = $1 order by sources.id`,
		"sources.id", "sources.url",
		"examples_sources", "sources",
		"examples_sources.source_id", "sources.id",
		"example_id",
	)
	rows, err := a.db.Query(c, query, exampleID)
	if err != nil {
		return nil, erw.New(erw.Internal(
			erw.Location(pkg, examples, sources),
			erw.Error(err),
			erw.SQL(query, exampleID),
		))
	}
	cls, err := pgx.CollectRows(rows, pgx.RowToStructByName[Source])
	if err != nil {
		return nil, erw.New(erw.Internal(
			erw.Location(pkg, examples, sources),
			erw.Error(err),
			erw.SQL(query, exampleID),
		))
	}
	srs := new([]entities.Source)
	for _, cl := range cls {
		*srs = append(*srs, entities.Source{
			Id:  cl.Id,
			Url: cl.Url,
		})
	}
	return srs, nil
}
