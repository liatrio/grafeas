package storage

// createTables represents the initial query grafeas uses to populate the database schema.
// for backwards compatibility, this query will remain here. this will act as the "base" migration point.
const createTables = `
BEGIN;
CREATE TABLE IF NOT EXISTS projects (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS notes (
	id SERIAL PRIMARY KEY,
	project_name TEXT NOT NULL,
	note_name TEXT NOT NULL,
	data TEXT,
	UNIQUE (project_name, note_name)
);
CREATE TABLE IF NOT EXISTS occurrences (
	id SERIAL PRIMARY KEY,
	project_name TEXT NOT NULL,
	occurrence_name TEXT NOT NULL,
	data TEXT,
	note_id int REFERENCES notes NOT NULL,
	UNIQUE (project_name, occurrence_name)
);
CREATE TABLE IF NOT EXISTS operations (
	id SERIAL PRIMARY KEY,
	project_name TEXT NOT NULL,
	operation_name TEXT NOT NULL,
	data TEXT,
	UNIQUE (project_name, operation_name)
);
COMMIT;
`

// upMigrations is a slice of strings that contains each database migration from the base.
// this is an ordered slice, so new migrations should be placed at the end
var upMigrations = []string{
	createTables,
	alterOccurrenceTableAddFilterableColumns,
}

// alterOccurrenceTableAddFilterableColumns alters the existing table for occurrences to add columns for each filterable field
const alterOccurrenceTableAddFilterableColumns = `
BEGIN;
	ALTER TABLE occurrences
		ADD COLUMN resource_name TEXT,
		ADD COLUMN resource_uri TEXT,
		ADD COLUMN kind TEXT,
		ADD COLUMN create_time timestamp,
		ADD COLUMN update_time timestamp;
COMMIT;
`
