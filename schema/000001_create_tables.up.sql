CREATE TABLE IF NOT EXISTS texts (
    id UUID NOT NULL DEFAULT uuidv7(),
    name VARCHAR(4) NOT NULL,
    text VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY(id),
    UNIQUE(name)
);

COMMENT ON TABLE texts IS 'Таблица текстов';
COMMENT ON COLUMN texts.id IS 'Уникальный идентификатор текста';
COMMENT ON COLUMN texts.name IS 'Уникальное название текста';
COMMENT ON COLUMN texts.text IS 'Текст';
COMMENT ON COLUMN texts.created_at IS 'Дата и время создания записи о тексте';

CREATE TABLE IF NOT EXISTS technologies (
    id UUID NOT NULL DEFAULT uuidv7(),
    title VARCHAR(20) NOT NULL,
    url VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY(id),
    UNIQUE(title)
);

COMMENT ON TABLE technologies IS 'Таблица технологий';
COMMENT ON COLUMN technologies.id IS 'Уникальный идентификатор технологии';
COMMENT ON COLUMN technologies.title IS 'Уникальное название технологии';
COMMENT ON COLUMN technologies.url IS 'Ссылка на документацию по технологии';
COMMENT ON COLUMN technologies.created_at IS 'Дата и время создания записи о технологии';

CREATE TABLE IF NOT EXISTS links (
    id UUID NOT NULL DEFAULT uuidv7(),
    title VARCHAR(32) NOT NULL,
    url VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY(id),
    UNIQUE(title)
);

COMMENT ON TABLE links IS 'Таблица ссылок';
COMMENT ON COLUMN links.id IS 'Уникальный идентификатор ссылки';
COMMENT ON COLUMN links.title IS 'Уникальное название ссылки';
COMMENT ON COLUMN links.url IS 'Адрес ссылки';
COMMENT ON COLUMN links.created_at IS 'Дата и время создания записи о ссылки';

CREATE TABLE IF NOT EXISTS sources (
    id UUID NOT NULL DEFAULT uuidv7(),
    url VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY(id),
    UNIQUE(url)
);

COMMENT ON TABLE sources IS 'Таблица источников';
COMMENT ON COLUMN sources.id IS 'Уникальный идентификатор источника';
COMMENT ON COLUMN sources.url IS 'Уникальный адрес источника';
COMMENT ON COLUMN sources.created_at IS 'Дата и время создания записи об источнике';

CREATE TABLE IF NOT EXISTS examples (
    id UUID NOT NULL DEFAULT uuidv7(),
    name VARCHAR(20) NOT NULL,
    title VARCHAR(32) NOT NULL,
    description VARCHAR(1024) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY(id),
    UNIQUE(name)
);

COMMENT ON TABLE examples IS 'Таблица примеров';
COMMENT ON COLUMN examples.id IS 'Уникальный идентификатор примера';
COMMENT ON COLUMN examples.name IS 'Уникальное имя примера';
COMMENT ON COLUMN examples.title IS 'Название примера';
COMMENT ON COLUMN examples.description IS 'Описаное примера';
COMMENT ON COLUMN examples.created_at IS 'Дата и время создания записи о примере';

CREATE TABLE IF NOT EXISTS examples_links (
    id UUID NOT NULL DEFAULT uuidv7(),
    example_id UUID NOT NULL,
    link_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    CONSTRAINT fk_examples FOREIGN KEY (example_id) REFERENCES examples (id),
    CONSTRAINT fk_link FOREIGN KEY (link_id) REFERENCES links (id),
    PRIMARY KEY(id),
    UNIQUE(example_id, link_id)
);

COMMENT ON TABLE examples_links IS 'Таблица соединений примеров и ссылок';
COMMENT ON COLUMN examples_links.id IS 'Уникальный идентификатор соединения';
COMMENT ON COLUMN examples_links.example_id IS 'Идентификатор примера';
COMMENT ON COLUMN examples_links.link_id IS 'Идентификатор ссылки';
COMMENT ON COLUMN examples_links.created_at IS 'Дата и время создания записи о соединении';

CREATE TABLE IF NOT EXISTS examples_sources (
    id UUID NOT NULL DEFAULT uuidv7(),
    example_id UUID NOT NULL,
    source_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    CONSTRAINT fk_examples FOREIGN KEY (example_id) REFERENCES examples (id),
    CONSTRAINT fk_sources FOREIGN KEY (source_id) REFERENCES sources (id),
    PRIMARY KEY(id),
    UNIQUE(example_id, source_id)
);

COMMENT ON TABLE examples_sources IS 'Таблица соединений примеров и источников';
COMMENT ON COLUMN examples_sources.id IS 'Уникальный идентификатор соединения';
COMMENT ON COLUMN examples_sources.example_id IS 'Идентификатор примера';
COMMENT ON COLUMN examples_sources.source_id IS 'Идентификатор источника';
COMMENT ON COLUMN examples_sources.created_at IS 'Дата и время создания записи о соединении';

CREATE TABLE IF NOT EXISTS examples_technologies (
    id UUID NOT NULL DEFAULT uuidv7(),
    example_id UUID NOT NULL,
    technology_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    CONSTRAINT fk_examples FOREIGN KEY (example_id) REFERENCES examples (id),
    CONSTRAINT fk_tecknology FOREIGN KEY (technology_id) REFERENCES technologies (id),
    PRIMARY KEY(id),
    UNIQUE(example_id, technology_id)
);

COMMENT ON TABLE examples_technologies IS 'Таблица соединений примеров и технологий';
COMMENT ON COLUMN examples_technologies.id IS 'Уникальный идентификатор соединения';
COMMENT ON COLUMN examples_technologies.example_id IS 'Идентификатор примера';
COMMENT ON COLUMN examples_technologies.technology_id IS 'Идентификатор технологии';
COMMENT ON COLUMN examples_technologies.created_at IS 'Дата и время создания записи о соединении';

CREATE TABLE IF NOT EXISTS software (
    id UUID NOT NULL DEFAULT uuidv7(),
    title VARCHAR(32) NOT NULL,
    url VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY(id),
    UNIQUE(title)
);

COMMENT ON TABLE software IS 'Таблица приложений';
COMMENT ON COLUMN software.id IS 'Уникальный идентификатор приложения';
COMMENT ON COLUMN software.title IS 'Уникальное название приложения';
COMMENT ON COLUMN software.url IS 'Ссылка на приложение';
COMMENT ON COLUMN software.created_at IS 'Дата и время создания записи о приложении';

CREATE TABLE IF NOT EXISTS libs (
    id UUID NOT NULL DEFAULT uuidv7(),
    url VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY(id),
    UNIQUE(url)
);

COMMENT ON TABLE libs IS 'Таблица библиотек';
COMMENT ON COLUMN libs.id IS 'Уникальный идентификатор библиотеки';
COMMENT ON COLUMN libs.url IS 'Уникальная ссылка на библиотеку';
COMMENT ON COLUMN libs.created_at IS 'Дата и время создания записи о библиотеке';

CREATE TABLE IF NOT EXISTS footer_links (
    id UUID NOT NULL DEFAULT uuidv7(),
    title VARCHAR(32) NOT NULL,
    url VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY(id),
    UNIQUE(title)
);

COMMENT ON TABLE footer_links IS 'Таблица подвальных ссылок';
COMMENT ON COLUMN footer_links.id IS 'Уникальный идентификатор подвальной ссылки';
COMMENT ON COLUMN footer_links.title IS 'Уникальное название подвальной ссылки';
COMMENT ON COLUMN footer_links.url IS 'Адрес подвальной ссылки';
COMMENT ON COLUMN footer_links.created_at IS 'Дата и время создания записи о подвальной ссылки';
