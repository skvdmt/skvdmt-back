package inserts

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skvdmt/skvdmt-back/internal/model"
)

type request struct {
	query string
	args  []any
	id    uuid.UUID
}

var requests = []request{
	// texts
	{
		query: `INSERT INTO texts (name, text) VALUES ($1, $2);`,
		args:  []any{"main", "Dmitry Skidanov"},
	},
	{
		query: `INSERT INTO texts (name, text) VALUES ($1, $2);`,
		args:  []any{"tech", "Main technologies"},
	},
	{
		query: `INSERT INTO texts (name, text) VALUES ($1, $2);`,
		args:  []any{"exam", "Examples"},
	},
	{
		query: `INSERT INTO texts (name, text) VALUES ($1, $2);`,
		args:  []any{"favo", "Favorites"},
	},
	{
		query: `INSERT INTO texts (name, text) VALUES ($1, $2);`,
		args:  []any{"soft", "Software for development"},
	},
	{
		query: `INSERT INTO texts (name, text) VALUES ($1, $2);`,
		args:  []any{"libs", "Golang external libraries i like"},
	},
	{
		query: `INSERT INTO texts (name, text) VALUES ($1, $2);`,
		args:  []any{"prof", "full stack engineer."},
	},
	{
		query: `INSERT INTO texts (name, text) VALUES ($1, $2);`,
		args:  []any{"abou", "Dmitry Skidanov — full stack engineer 2026"},
	},
	{
		query: `INSERT INTO texts (name, text) VALUES ($1, $2);`,
		args:  []any{"lock", "Russian Federation, Moscow"},
	},

	// software
	{
		query: `INSERT INTO software(title, url) VALUES($1, $2);`,
		args:  []any{"Visual Studio Code", "https://code.visualstudio.com/"},
	},
	{
		query: `INSERT INTO software(title, url) VALUES($1, $2);`,
		args:  []any{"GoLand", "https://www.jetbrains.com/go/"},
	},
	{
		query: `INSERT INTO software(title, url) VALUES($1, $2);`,
		args:  []any{"WebStorm", "https://www.jetbrains.com/webstorm/"},
	},
	{
		query: `INSERT INTO software(title, url) VALUES($1, $2);`,
		args:  []any{"DataGrip", "https://www.jetbrains.com/datagrip/"},
	},
	{
		query: `INSERT INTO software(title, url) VALUES($1, $2);`,
		args:  []any{"Postman", "https://www.postman.com/"},
	},
	{
		query: `INSERT INTO software(title, url) VALUES($1, $2);`,
		args:  []any{"Swagger", "https://swagger.io/"},
	},
	{
		query: `INSERT INTO software(title, url) VALUES($1, $2);`,
		args:  []any{"Vite", "https://vite.dev/"},
	},

	// libs
	{
		query: `INSERT INTO libs(url) VALUES($1);`,
		args:  []any{"https://github.com/labstack/echo"},
	},
	{
		query: `INSERT INTO libs(url) VALUES($1);`,
		args:  []any{"https://github.com/grpc/grpc-go"},
	},
	{
		query: `INSERT INTO libs(url) VALUES($1);`,
		args:  []any{"https://github.com/patrickmn/go-cache"},
	},
	{
		query: `INSERT INTO libs(url) VALUES($1);`,
		args:  []any{"https://github.com/jackc/pgx"},
	},
	{
		query: `INSERT INTO libs(url) VALUES($1);`,
		args:  []any{"https://github.com/gorilla/websocket"},
	},

	// footer_links
	{
		query: `INSERT INTO footer_links(title, url) VALUES($1, $2);`,
		args:  []any{"GitHub repositories", "https://github.com/skvdmt"},
	},
	{
		query: `INSERT INTO footer_links(title, url) VALUES($1, $2);`,
		args:  []any{"Docker Hub registry", "https://hub.docker.com/u/skvdmt"},
	},
	{
		query: `INSERT INTO footer_links(title, url) VALUES($1, $2);`,
		args:  []any{"Telegram", "https://t.me/skidanovdima"},
	},
	{
		query: `INSERT INTO footer_links(title, url) VALUES($1, $2);`,
		args:  []any{"Email", "mailto:skvdmt@yandex.ru"},
	},
}

var examples = []request{
	{
		query: `INSERT INTO examples (name, title, description) VALUES($1, $2, $3) RETURNING id;`,
		args: []any{
			"chess",
			"Chess game",
			`A chess server that allows two users to play chess.
The first user to connect gets white pieces, the second gets black pieces, the rest of the users get into spectators.
To start the game, two clients must be logged into the page.`,
		},
	},
	{
		query: `INSERT INTO examples (name, title, description) VALUES($1, $2, $3) RETURNING id;`,
		args: []any{
			"tgbot",
			"Telegram bot",
			`Bot that saves messages from users.
Before leaving a message, you must pass the test by entering the characters from the picture.
One user can leave no more than one message per 24 hours.`,
		},
	},
	{
		query: `INSERT INTO examples (name, title, description) VALUES($1, $2, $3) RETURNING id;`,
		args: []any{
			"home",
			"Homepage",
			`My personal homepage.
You are here.`,
		},
	},
}

var links = []request{
	{
		query: `INSERT INTO links(title, url) VALUES($1, $2) RETURNING id;`,
		args:  []any{"play", "https://chess.skvdmt.ru/"},
	},
	{
		query: `INSERT INTO links(title, url) VALUES($1, $2) RETURNING id;`,
		args:  []any{"open telegram bot", "https://t.me/skdmtr_bot"},
	},
	{
		query: `INSERT INTO links(title, url) VALUES($1, $2) RETURNING id;`,
		args:  []any{"view messages", "https://tgbot.skvdmt.ru/messages"},
	},
}

var technologies = []request{
	{
		query: `INSERT INTO technologies (title, url) VALUES($1, $2) RETURNING id;`,
		args:  []any{"Go", "https://go.dev/"},
	},
	{
		query: `INSERT INTO technologies (title, url) VALUES($1, $2) RETURNING id;`,
		args:  []any{"Postgres", "https://www.postgresql.org/"},
	},
	{
		query: `INSERT INTO technologies (title, url) VALUES($1, $2) RETURNING id;`,
		args:  []any{"Docker", "https://docker.com/"},
	},
	{
		query: `INSERT INTO technologies (title, url) VALUES($1, $2) RETURNING id;`,
		args:  []any{"REST API", "https://restfulapi.net/"},
	},
	{
		query: `INSERT INTO technologies (title, url) VALUES($1, $2) RETURNING id;`,
		args:  []any{"gRPC", "https://grpc.io/"},
	},
	{
		query: `INSERT INTO technologies (title, url) VALUES($1, $2) RETURNING id;`,
		args:  []any{"Git", "https://git-scm.com/"},
	},
	{
		query: `INSERT INTO technologies (title, url) VALUES($1, $2) RETURNING id;`,
		args:  []any{"CI/CD", "https://github.com/features/actions"},
	},
	{
		query: `INSERT INTO technologies (title, url) VALUES($1, $2) RETURNING id;`,
		args:  []any{"TypeScript", "https://www.typescriptlang.org/"},
	},
	{
		query: `INSERT INTO technologies (title, url) VALUES($1, $2) RETURNING id;`,
		args:  []any{"JavaScript", "https://developer.mozilla.org/en-US/docs/Web/JavaScript"},
	},
	{
		query: `INSERT INTO technologies (title, url) VALUES($1, $2) RETURNING id;`,
		args:  []any{"Vue", "https://vuejs.org/"},
	},
}

var sources = []request{
	{
		query: `INSERT INTO sources(url) VALUES($1) RETURNING id;`,
		args:  []any{"https://github.com/skvdmt/chess-front"},
	},
	{
		query: `INSERT INTO sources(url) VALUES($1) RETURNING id;`,
		args:  []any{"https://github.com/skvdmt/chess-back-game"},
	},
	{
		query: `INSERT INTO sources(url) VALUES($1) RETURNING id;`,
		args:  []any{"https://github.com/skvdmt/tgbot-front-messages"},
	},
	{
		query: `INSERT INTO sources(url) VALUES($1) RETURNING id;`,
		args:  []any{"https://github.com/skvdmt/tgbot-back-messages"},
	},
	{
		query: `INSERT INTO sources(url) VALUES($1) RETURNING id;`,
		args:  []any{"https://github.com/skvdmt/tgbot-back-app"},
	},
	{
		query: `INSERT INTO sources(url) VALUES($1) RETURNING id;`,
		args:  []any{"https://github.com/skvdmt/skvdmt-front"},
	},
	{
		query: `INSERT INTO sources(url) VALUES($1) RETURNING id;`,
		args:  []any{"https://github.com/skvdmt/skvdmt-back"},
	},
}

const dublicate = "duplicate key value violates unique constraint"

// InsertData insert default data to database tables
func InsertData(db *pgxpool.Pool) {
	for _, req := range requests {
		_, err := db.Exec(context.Background(), req.query, req.args...)
		if err != nil {
			if !strings.Contains(err.Error(), dublicate) {
				model.Logs.Error.Error(err.Error())
			}
		}
	}
	insertDataSetID(db, &examples)
	insertDataSetID(db, &links)
	insertDataSetID(db, &technologies)
	insertDataSetID(db, &sources)

	createLinks(db, "examples_links", "link_id", 1, 1)
	createLinks(db, "examples_links", "link_id", 2, 2)
	createLinks(db, "examples_links", "link_id", 2, 3)

	createLinks(db, "examples_technologies", "technology_id", 1, 1)
	createLinks(db, "examples_technologies", "technology_id", 1, 9)
	createLinks(db, "examples_technologies", "technology_id", 1, 3)
	createLinks(db, "examples_technologies", "technology_id", 1, 6)
	createLinks(db, "examples_technologies", "technology_id", 2, 1)
	createLinks(db, "examples_technologies", "technology_id", 2, 2)
	createLinks(db, "examples_technologies", "technology_id", 2, 3)
	createLinks(db, "examples_technologies", "technology_id", 2, 4)
	createLinks(db, "examples_technologies", "technology_id", 2, 6)
	createLinks(db, "examples_technologies", "technology_id", 2, 9)
	createLinks(db, "examples_technologies", "technology_id", 3, 1)
	createLinks(db, "examples_technologies", "technology_id", 3, 2)
	createLinks(db, "examples_technologies", "technology_id", 3, 3)
	createLinks(db, "examples_technologies", "technology_id", 3, 4)
	createLinks(db, "examples_technologies", "technology_id", 3, 6)
	createLinks(db, "examples_technologies", "technology_id", 3, 9)

	createLinks(db, "examples_sources", "source_id", 1, 1)
	createLinks(db, "examples_sources", "source_id", 1, 2)
	createLinks(db, "examples_sources", "source_id", 2, 3)
	createLinks(db, "examples_sources", "source_id", 2, 4)
	createLinks(db, "examples_sources", "source_id", 2, 5)
	createLinks(db, "examples_sources", "source_id", 3, 6)
	createLinks(db, "examples_sources", "source_id", 3, 7)
}

// insertDataSetID remember returning id after inserting data
func insertDataSetID(db *pgxpool.Pool, rs *[]request) {
	for i, req := range *rs {
		if err := db.QueryRow(context.Background(), req.query,
			req.args...).Scan(&(*rs)[i].id); err != nil {
			if !strings.Contains(err.Error(), dublicate) {
				model.Logs.Error.Error(err.Error())
			}
		}
	}
}

const defaultID = "00000000-0000-0000-0000-000000000000"

// createLinks creating examples links to links, technologies and sources
func createLinks(db *pgxpool.Pool, table, field string, exampleID, targetID int) {
	var id uuid.UUID
	switch table {
	case "examples_links":
		id = links[targetID-1].id
	case "examples_technologies":
		id = technologies[targetID-1].id
	case "examples_sources":
		id = sources[targetID-1].id
	}
	if examples[exampleID-1].id.String() != defaultID &&
		id.String() != defaultID {
		db.Exec(
			context.Background(),
			fmt.Sprintf(
				"INSERT INTO %s(example_id, %s) VALUES($1, $2);",
				table, field,
			),
			examples[exampleID-1].id, id,
		)
	}
}
