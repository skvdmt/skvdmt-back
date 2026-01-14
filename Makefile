restart:
	make go-build && docker container restart skvdmt-back

go-build:
	go build -v -o ./build/skvdmt-back ./cmd/main.go

download:
	go get -u github.com/stretchr/testify/require
	go get -u github.com/stretchr/testify/assert
	go get -u github.com/patrickmn/go-cache
	go get -u github.com/jackc/pgx/v5
	go get -u gopkg.in/yaml.v3
	go get -u github.com/google/uuid

swag-serve:
	swagger serve ./swagger.yaml

swag-valid:
	swagger validate ./swagger.yaml
