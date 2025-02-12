
gocriticDS:
	gocritic check -enableAll -disable='#style' ./...

gocritic:
	gocritic check -enableAll ./...

gosec:
	gosec ./...

migration:
	@goose -dir db/migrations up

down:
	@goose -dir db/migrations down

# goose -dir db/migrations create "$(MIGRATION_NAME)" sql