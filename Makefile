GOOSE_MIGRATION_DIR=./cmd/migrations

.PHONY: migration-create
migration-create:
	@goose -s create -dir $(GOOSE_MIGRATION_DIR) $(filter-out $@,$(MAKECMDGOALS)) sql

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt