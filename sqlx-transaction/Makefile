# Set environment variables
DB_CONTAINER_NAME=my_postgres
DB_USER=user
DB_NAME=test
MIGRATION_FILE=migration.sql

# Run migration
migrate:
	docker exec -i $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME) < $(MIGRATION_FILE)

# Copy migration file and execute (if file is not inside the container)
migrate-copy:
	docker cp $(MIGRATION_FILE) $(DB_CONTAINER_NAME):/$(MIGRATION_FILE)
	docker exec -i $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME) -f /$(MIGRATION_FILE)

# Check tables after migration
check-tables:
	docker exec -it $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME) -c "\dt"

# Run bash inside the PostgreSQL container (optional for debugging)
psql:
	docker exec -it $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME)
