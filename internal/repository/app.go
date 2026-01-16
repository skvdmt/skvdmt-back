package repository

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skvdmt/skvdmt-back/init/inserts"
	"github.com/skvdmt/skvdmt-back/internal/entities"
	"github.com/skvdmt/skvdmt-back/internal/model"
	erw "github.com/skvdmt/skvdmt-back/pkg/errwrap"
)

const (
	pkg            = "repository"
	text           = "text"
	technologies   = "technologies"
	examples       = "examples"
	software       = "software"
	libs           = "libs"
	links          = "links"
	sources        = "sources"
	updateInterval = 5
)

// App Репозиторный слой.
type App struct {
	db             *pgxpool.Pool
	muTexts        *sync.RWMutex
	muTechnologies *sync.RWMutex
	muExamples     *sync.RWMutex
	muSoftware     *sync.RWMutex
	muLinks        *sync.RWMutex
	muLibs         *sync.RWMutex
	updateRunner   *time.Ticker
	exit           chan struct{}
	update         *sync.WaitGroup
	source         *sync.WaitGroup
	texts          map[string]*entities.Text
	technologies   []*entities.Technology
	examples       []*entities.Example
	software       []*entities.Software
	links          []*entities.Link
	libs           []*entities.Lib
}

const (
	postgres    = "postgres"
	DB_PASSWORD = "DB_PASSWORD"
)

// NewApp Конструктор.
func NewApp(ctx context.Context) (*App, error) {
	model.Logs.Info.Info("repository layer creating")
	model.Logs.Info.Info("database connection creating")
	pwd, ok := os.LookupEnv(DB_PASSWORD)
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
	dbpool, err := pgxpool.New(ctx, q)
	if err != nil {
		return nil, err
	}

	model.Logs.Info.Info("insert data to database")
	// Вставка данных в базу данных.
	inserts.InsertData(dbpool)

	a := &App{
		db:             dbpool,
		texts:          make(map[string]*entities.Text),
		exit:           make(chan struct{}),
		updateRunner:   time.NewTicker(time.Minute * updateInterval),
		update:         &sync.WaitGroup{},
		source:         &sync.WaitGroup{},
		muTexts:        &sync.RWMutex{},
		muTechnologies: &sync.RWMutex{},
		muExamples:     &sync.RWMutex{},
		muSoftware:     &sync.RWMutex{},
		muLinks:        &sync.RWMutex{},
		muLibs:         &sync.RWMutex{},
	}
	// Первый запуск обновлений.
	a.updateAll(ctx)

	// Запуск ресурса.
	a.source.Add(1)
	go a.handler(ctx)
	return a, nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	// Отправляем сигнал завершения.
	a.exit <- struct{}{}
	// Ожидание завершения всех ресурсов.
	a.source.Wait()
	// Закрытие соединения с базой данных.
	a.db.Close()
	model.Logs.Info.Info("database connection closed")
	model.Logs.Info.Info("repository layer stopped")
	return nil
}

// Text Репозиторий текстов.
func (a *App) Text(ctx context.Context, name string) (*entities.Text, error) {
	a.muTexts.RLock()
	t, ok := a.texts[name]
	a.muTexts.RUnlock()
	if !ok {
		return nil, erw.New(
			erw.CodeHTTP(404),
			erw.Internal(
				erw.Location(pkg, text),
				erw.Error(model.Errs[model.ErrTextNotFound]),
			))
	}
	return t, nil
}

// Technologies Репозиторий технологий.
func (a *App) Technologies(ctx context.Context) ([]*entities.Technology, error) {
	a.muTechnologies.RLock()
	defer a.muTechnologies.RUnlock()
	return a.technologies, nil
}

// Examples Репозиторий примеров.
func (a *App) Examples(ctx context.Context) ([]*entities.Example, error) {
	a.muExamples.RLock()
	defer a.muExamples.RUnlock()
	return a.examples, nil
}

// Software Репозиторий программ.
func (a *App) Software(ctx context.Context) ([]*entities.Software, error) {
	a.muSoftware.RLock()
	defer a.muSoftware.RUnlock()
	return a.software, nil
}

// Libs Репозиторий библиотек.
func (a *App) Libs(ctx context.Context) ([]*entities.Lib, error) {
	a.muLibs.RLock()
	defer a.muLibs.RUnlock()
	return a.libs, nil
}

// Links Репозиторий ссылок.
func (a *App) Links(ctx context.Context) ([]*entities.Link, error) {
	a.muLinks.RLock()
	defer a.muLinks.RUnlock()
	return a.links, nil
}

// handler Система отслеживания ресурсов.
func (a *App) handler(ctx context.Context) {
	model.Logs.Info.Info("repository handler creating")
	for {
		select {
		case <-a.exit:
			// Завершение ресурса.
			model.Logs.Info.Info("repository handler stopped")
			a.source.Done()
			return
		case <-a.updateRunner.C:
			// Запуск обновлений.
			a.updateAll(ctx)
		}
	}
}

// updateAll Обновить все.
func (a *App) updateAll(ctx context.Context) {
	a.update.Add(6)
	go a.updateTexts(ctx)
	go a.updateTechnologies(ctx)
	go a.updateExamples(ctx)
	go a.updateSoftware(ctx)
	go a.updateLibs(ctx)
	go a.updateLinks(ctx)
	a.update.Wait()
	model.Logs.Info.Info("all updated")
}

// updateTexts Обновление текстов.
func (a *App) updateTexts(ctx context.Context) {
	const query = "select id, name, text from texts"
	row, err := a.db.Query(ctx, query)
	if err != nil {
		model.Errors <- erw.New(erw.Internal(
			erw.Location(pkg, text),
			erw.Error(err),
			erw.SQL(query),
		))
		a.update.Done()
		return
	}
	cls, err := pgx.CollectRows(row, pgx.RowToStructByName[Text])
	if err != nil {
		model.Errors <- erw.New(erw.Internal(
			erw.Location(pkg, text),
			erw.Error(err),
			erw.SQL(query),
		))
		a.update.Done()
		return
	}
	texts := make(map[string]*entities.Text)
	for _, cl := range cls {
		texts[cl.Name] = &entities.Text{
			Id:   cl.Id,
			Text: cl.Text,
		}
	}
	a.muTexts.Lock()
	a.texts = texts
	a.muTexts.Unlock()
	a.update.Done()
}

// updateTechnologies Обновление технологий.
func (a *App) updateTechnologies(ctx context.Context) {
	query := fmt.Sprintf(
		"select %s, %s, %s from %s order by id",
		"id", "title", "url", technologies,
	)
	rows, err := a.db.Query(ctx, query)
	if err != nil {
		model.Errors <- erw.New(erw.Internal(
			erw.Location(pkg, technologies),
			erw.Error(err),
			erw.SQL(query),
		))
		a.update.Done()
		return
	}
	cls, err := pgx.CollectRows(rows, pgx.RowToStructByName[Technology])
	if err != nil {
		model.Errors <- erw.New(erw.Internal(
			erw.Location(pkg, technologies),
			erw.Error(err),
			erw.SQL(query),
		))
		a.update.Done()
		return
	}
	tcs := []*entities.Technology{}
	for _, cl := range cls {
		tcs = append(tcs, &entities.Technology{
			Id:    cl.Id,
			Title: cl.Title,
			Url:   cl.Url,
		})
	}
	a.muTechnologies.Lock()
	a.technologies = tcs
	a.muTechnologies.Unlock()
	a.update.Done()
}

// updateExamples Обновление примеров.
func (a *App) updateExamples(ctx context.Context) {
	query := fmt.Sprintf(
		"select %s, %s, %s, %s from %s order by id",
		"id", "name", "title", "description", examples,
	)
	rows, err := a.db.Query(ctx, query)
	if err != nil {
		model.Errors <- erw.New(erw.Internal(
			erw.Location(pkg, examples),
			erw.Error(err),
			erw.SQL(query),
		))
		a.update.Done()
		return
	}
	cls, err := pgx.CollectRows(rows, pgx.RowToStructByName[Example])
	if err != nil {
		model.Errors <- erw.New(erw.Internal(
			erw.Location(pkg, examples),
			erw.Error(err),
			erw.SQL(query),
		))
		a.update.Done()
		return
	}
	exs := []*entities.Example{}
	for _, cl := range cls {
		lks, err := a.exampleLinks(ctx, cl.Id)
		if err != nil {
			model.Errors <- err
			a.update.Done()
			return
		}
		tks, err := a.exampleTechnologies(ctx, cl.Id)
		if err != nil {
			model.Errors <- err
			a.update.Done()
			return
		}
		srs, err := a.exampleSources(ctx, cl.Id)
		if err != nil {
			model.Errors <- err
			a.update.Done()
			return
		}
		ex := &entities.Example{
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
		exs = append(exs, ex)
	}
	a.muExamples.Lock()
	a.examples = exs
	a.muExamples.Unlock()
	a.update.Done()
}

// updateSoftware Обновление программ.
func (a *App) updateSoftware(ctx context.Context) {
	query := fmt.Sprintf(
		"select %s, %s, %s from %s order by id",
		"id", "title", "url", software,
	)
	rows, err := a.db.Query(ctx, query)
	if err != nil {
		model.Errors <- erw.New(erw.Internal(
			erw.Location(pkg, software),
			erw.Error(err),
			erw.SQL(query),
		))
		a.update.Done()
		return
	}
	cls, err := pgx.CollectRows(rows, pgx.RowToStructByName[Software])
	if err != nil {
		model.Errors <- erw.New(erw.Internal(
			erw.Location(pkg, software),
			erw.Error(err),
			erw.SQL(query),
		))
		a.update.Done()
		return
	}
	sfw := []*entities.Software{}
	for _, cl := range cls {
		sfw = append(sfw, &entities.Software{
			Id:    cl.Id,
			Title: cl.Title,
			Url:   cl.Url,
		})
	}
	a.muSoftware.Lock()
	a.software = sfw
	a.muSoftware.Unlock()
	a.update.Done()
}

// updateLibs Обновление библиотек.
func (a *App) updateLibs(ctx context.Context) {
	query := fmt.Sprintf(
		"select %s, %s from %s order by id",
		"id", "url", libs,
	)
	rows, err := a.db.Query(ctx, query)
	if err != nil {
		model.Errors <- erw.New(erw.Internal(
			erw.Location(pkg, libs),
			erw.Error(err),
			erw.SQL(query),
		))
		a.update.Done()
		return
	}
	cls, err := pgx.CollectRows(rows, pgx.RowToStructByName[Lib])
	if err != nil {
		model.Errors <- erw.New(erw.Internal(
			erw.Location(pkg, libs),
			erw.Error(err),
			erw.SQL(query),
		))
		a.update.Done()
		return
	}
	lbs := []*entities.Lib{}
	for _, cl := range cls {
		lbs = append(lbs, &entities.Lib{
			Id:  cl.Id,
			Url: cl.Url,
		})
	}
	a.muLibs.Lock()
	a.libs = lbs
	a.muLibs.Unlock()
	a.update.Done()
}

// updateLinks Основление ссылок.
func (a *App) updateLinks(ctx context.Context) {
	query := fmt.Sprintf(
		"select %s, %s, %s from %s order by id",
		"id", "title", "url", "footer_links",
	)
	rows, err := a.db.Query(ctx, query)
	if err != nil {
		model.Errors <- erw.New(erw.Internal(
			erw.Location(pkg, links),
			erw.Error(err),
			erw.SQL(query),
		))
		a.update.Done()
		return
	}
	cls, err := pgx.CollectRows(rows, pgx.RowToStructByName[Link])
	if err != nil {
		model.Errors <- erw.New(erw.Internal(
			erw.Location(pkg, links),
			erw.Error(err),
			erw.SQL(query),
		))
		a.update.Done()
		return
	}
	lks := []*entities.Link{}
	for _, cl := range cls {
		lks = append(lks, &entities.Link{
			Id:    cl.Id,
			Title: cl.Title,
			Url:   cl.Url,
		})
	}
	a.muLinks.Lock()
	a.links = lks
	a.muLinks.Unlock()
	a.update.Done()
}

// exampleLinks Ссылки примера.
func (a *App) exampleLinks(ctx context.Context, exampleID uuid.UUID) (*[]entities.Link, error) {
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
	rows, err := a.db.Query(ctx, query, exampleID)
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
func (a *App) exampleTechnologies(ctx context.Context, exampleID uuid.UUID) (*[]entities.Technology, error) {
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
	rows, err := a.db.Query(ctx, query, exampleID)
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
func (a *App) exampleSources(ctx context.Context, exampleID uuid.UUID) (*[]entities.Source, error) {
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
	rows, err := a.db.Query(ctx, query, exampleID)
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
