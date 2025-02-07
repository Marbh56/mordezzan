# Mordezzan Setup Guide
Prerequisites

Go 1.23 or higher
SQLite 3
Git

Installation

Clone the repository:

```bash
git clone https://github.com/marbh56/mordezzan.git
cd mordezzan
```

Install dependencies:

```bash
go mod download
```
Set up the database:

# Create a new SQLite database
```bash
touch mordezzan.db
```
# Run migrations (you'll need goose installed)
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir sql/migrations sqlite3 ./mordezzan.db up

Build the project:

bashCopygo build -o mordezzan ./cmd/web
Running the Server

Start the server:

Copy./mordezzan

Access the web interface at http://localhost:8080

Development

Database queries are managed using sqlc. After modifying SQL files, regenerate the Go code:

bashCopy# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate Go code
sqlc generate
Project Structure

cmd/web/: Main application entry point
internal/: Application code

database/: Database connection
db/: Generated SQL code
rules/: Game rules implementation
server/: HTTP server and handlers


sql/: SQL migrations and queries
templates/: HTML templates
static/: Static assets

License
This project is dual-licensed:

Source code: MIT License
Game content (rules directory and certain SQL migrations): OGL License
