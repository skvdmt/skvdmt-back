package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
	"github.com/skvdmt/skvdmt-back/init/inserts"
	"github.com/skvdmt/skvdmt-back/internal/entities"
	"github.com/skvdmt/skvdmt-back/internal/model"
	erw "github.com/skvdmt/skvdmt-back/pkg/errwrap"
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

// App Репозиторный слой.
type App struct {
	cache *cache.Cache
	db    *pgxpool.Pool
}

const (
	postgres    = "postgres"
	DB_PASSWORD = "DB_PASSWORD"
)

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Logs.Info.Info("repository layer creating")
	model.Logs.Info.Info("database connection creating")
	pwd, ok := os.LookupEnv(DB_PASSWORD)
	if !ok {
		return nil, fmt.Errorf("env %s unset", DB_PASSWORD)
	}
	pwd, err := url.QueryUnescape(pwd)
	if err != nil {
		return nil, err
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
	dbpool, err := pgxpool.New(context.Background(), q)
	if err != nil {
		return nil, err
	}

	model.Logs.Info.Info("insert data to database")
	// Вставка данных в базу данных.
	inserts.InsertData(dbpool)

	model.Logs.Info.Info("cache creating")
	c := cache.New(2*time.Minute, 20*time.Second)

	return &App{
		cache: c,
		db:    dbpool,
	}, nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	// Очистка кэша.
	a.cache.Flush()
	model.Logs.Info.Info("cache flushed")
	// Закрытие соединения с базой данных.
	a.db.Close()
	model.Logs.Info.Info("database connection closed")
	model.Logs.Info.Info("repository layer stopped")
	return nil
}

// Text Репозиторий текстов.
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
		// Установка кэша.
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

// Technologies Репозиторий технологий.
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
		// Установка кэша.
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

// Examples Репозиторий примеров.
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
		// Установка кэша.
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

// Software Репозиторий программ.
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
		// Установка кэша.
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

// Libs Репозиторий библиотек.
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
		// Установка кэша.
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

// Links Репозиторий ссылок.
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
		// Установка кэша.
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

// exampleLinks Ссылки примера.
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

// exampleTechnologies Технологии примера.
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

// exampleSources Исходники примера.
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
