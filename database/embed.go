package database

import "embed"

// MigrationsFS embeds all database migration files into the binary.
//
//go:embed migrations/*.sql
var MigrationsFS embed.FS
