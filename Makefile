restart:
	make go-build && docker container restart skvdmt-back

go-build:
	go build -v -o ./build/skvdmt-back ./cmd/main.go

docker-build:
	docker build -t skvdmt-back:1.0.0 .

download:
	go get -u github.com/labstack/echo/v4
	go get -u github.com/stretchr/testify/require
	go get -u github.com/stretchr/testify/assert
	go get -u github.com/patrickmn/go-cache
	go get -u github.com/jackc/pgx/v5
	go get -u gopkg.in/yaml.v3
	go get -u github.com/google/uuid
	go get -u github.com/skvdmt/errwrap

swagger-serve:
	swagger serve ./swagger.yaml

swagger-validate:
	swagger validate ./swagger.yaml
