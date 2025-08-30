# See:
# https://templ.guide/developer-tools/live-reload-with-other-tools

all: build

deps:
	npm install

live/sql:
	@go run github.com/air-verse/air@v1.62.0 \
	--build.cmd "go tool github.com/sqlc-dev/sqlc/cmd/sqlc generate" \
	--build.bin "/bin/true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "database" \
	--build.include_ext "sql" \
	--log.main_only "true"

live/templ:
	@printf "\033[1mstarting templ proxy, url: \033[0m\033[32m%s\033[0m\n" "http://localhost:7331"
	@go tool github.com/a-h/templ/cmd/templ generate --watch --proxy="http://localhost:8080" --open-browser=false --log-level="warn"

# run air to detect any go file changes to re-build and re-run the server.
live/server:
	@go run github.com/air-verse/air@v1.62.0 \
	--build.cmd "go build -o tmp/bin/main" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go,js,css" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true \
	--log.main_only "true"

# run tailwindcss to generate the styles.css bundle in watch mode.
live/tailwind: deps
	@go run github.com/air-verse/air@v1.62.0 \
	--build.cmd "npx --yes @tailwindcss/cli -i ./resources/css/tailwind.css -o ./assets/css/styles.css" \
	--build.bin "/bin/true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "resources,templates" \
	--build.include_ext "js,css,templ" \
	--log.main_only "true"

# run esbuild to generate the index.js bundle in watch mode.
live/esbuild: deps
	@npx --yes esbuild ./resources/ts/index.ts --bundle --outdir=assets/js --watch=forever

# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
live/sync_assets:
	@go run github.com/air-verse/air@v1.62.0 \
	--build.cmd "go tool github.com/a-h/templ/cmd/templ generate --notify-proxy" \
	--build.bin "/bin/true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "assets" \
	--build.include_ext "js,css" \
	--log.main_only "true"

# start all 5 watch processes in parallel.
live:
	$(MAKE) -j6 live/sql live/templ live/server live/esbuild live/sync_assets live/tailwind

build/sql: database/schema.sql database/queries.sql
	go tool github.com/sqlc-dev/sqlc/cmd/sqlc generate

build/templ:
	go tool github.com/a-h/templ/cmd/templ generate -v

build/tailwind: deps
	npx --yes @tailwindcss/cli -i ./resources/css/tailwind.css -o ./assets/css/styles.css

build/esbuild: deps
	npx --yes esbuild ./resources/ts/index.ts --bundle --outdir=assets/js

build:
	make -j4 build/sql build/templ build/tailwind build/esbuild
	CGO_ENABLED=0 go build

run: build
	go run .

dist: build scripts/release.sh
	./scripts/release.sh

clean:
	rm -rf node_modules tmp app dist

.PHONY: all deps \
	live live/sql live/templ live/server live/tailwind live/esbuild live/sync_assets \
	build build/sql build/templ build/tailwind build/esbuild \
	run dist clean
