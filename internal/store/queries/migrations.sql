-- name: GetLastMigrationVersion :one
SELECT version
FROM schema_migrations
ORDER BY version DESC
LIMIT 1;