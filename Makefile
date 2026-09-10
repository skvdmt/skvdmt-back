go-mod-install:
	go get -u github.com/jackc/pgx/v5
	go get -u gopkg.in/yaml.v3
	go get -u github.com/google/uuid

swag-serve:
	swagger serve ./swagger.yaml

swag-valid:
	swagger validate ./swagger.yaml
