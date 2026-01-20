#!/bin/bash
# Apply database schema

DATABASE_URL=${DATABASE_URL:-"postgres://postgres:postgres@localhost/notifications?sslmode=disable"}

echo "Applying schema to database..."
psql "$DATABASE_URL" < schema.sql

if [ $? -eq 0 ]; then
    echo "Schema applied successfully!"
else
    echo "Failed to apply schema. Please check your DATABASE_URL and database connection."
    exit 1
fi
