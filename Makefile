build:
	@tailwindcss -i views/css/styles.css -o public/styles.css
	@templ generate view
	@node build.js    # Add esbuild to bundle the React component
	@go build -o bin/cms main.go

test:
	@go test -v ./...

run: build
	@./bin/cms

tailwind:
	@tailwindcss -i views/css/styles.css -o public/styles.css --watch

templ:
	@templ generate -watch -proxy=http://localhost:7000

migrate-up:
	@goose -dir db/migrations sqlite3 ./cms.db up

migrate-down:
	@goose -dir db/migrations sqlite3 ./cms.db down

seed:
	@go run cmd/seed/main.go


# New target for esbuild
build-js:
	@node build.js  # Run the esbuild script

# Add a watch mode for JavaScript (if you need it for development)
watch-js:
	@npm run watch  # Watches and rebuilds JavaScript on changes

dev: tailwind templ watch-js
	@go run main.go
