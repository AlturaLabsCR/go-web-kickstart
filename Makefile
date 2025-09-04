# See:
# https://templ.guide/developer-tools/live-reload-with-other-tools

BIN = app

LOG = >> build.log 2>&1
LIVELOG = >> live.log 2>&1
SQLC = go tool github.com/sqlc-dev/sqlc/cmd/sqlc
TEMPL = go tool github.com/a-h/templ/cmd/templ

.PHONY: all
all: build

SQL=$(wildcard database/*.sql)
internal/db: $(SQL)
	$(SQLC) generate $(LOG)
	@touch $@

.PHONY: clean/sql
clean/sql:
	rm -rf internal/db

.PHONY: sql
sql: internal/db

TEMPLATES := $(wildcard templates/*.templ)
TEMPLATES_GEN := $(TEMPLATES:.templ=_templ.go)
$(TEMPLATES_GEN) &: $(TEMPLATES)
	$(TEMPL) generate -v >> make.log 2>&1
	@touch $@

.PHONY: templates
templates: $(TEMPLATES_GEN)

.PHONY: clean/templates
clean/templates:
	rm -rf $(TEMPLATES_GEN)

node_modules: package.json package-lock.json
	npm install >> make.log 2>&1

.PHONY: clean/node_modules
clean/node_modules:
	rm -rf node_modules

ESBUILD_IN := $(wildcard resources/ts/*.ts) $(wildcard resources/ts/*.js)
ESBUILD_OUT := $(addsuffix .js,$(basename $(patsubst resources/ts/%,assets/js/%,$(ESBUILD_IN))))
$(ESBUILD_OUT) &: $(ESBUILD_IN) node_modules
	@mkdir -p assets/js
	npx --yes esbuild $(ESBUILD_IN) --bundle --outdir=assets/js >> make.log 2>&1
	@touch $@

.PHONY: assets/js
assets/js: $(ESBUILD_OUT)

.PHONY: clean/assets/js
clean/assets/js:
	rm -rf $(ESBUILD_OUT)

assets/css/styles.css: resources/css/tailwind.css $(TEMPLATES) node_modules
	npx --yes @tailwindcss/cli -i ./resources/css/tailwind.css -o ./$@ >> make.log 2>&1
	touch $@

.PHONY: assets/css
assets/css: assets/css/styles.css

.PHONY: clean/assets/css
clean/assets/css:
	rm -rf assets/css/styles.css

.PHONY: prep
prep:
	@$(MAKE) -j4 assets/css assets/js templates sql

.PHONY: build
build: prep
	CGO_ENABLED=0 go build

.PHONY: clean/build
clean/build:
	rm -rf $(BIN) make.log

.PHONY: run
run: prep
	go run .

.PHONY: clean
clean:
	$(MAKE) -j6 clean/assets/js clean/assets/css clean/node_modules clean/sql clean/templates clean/build

# dist: build scripts/release.sh
# 	./scripts/release.sh

live/sql:
	@go run github.com/air-verse/air@v1.62.0 \
	--build.cmd "go tool github.com/sqlc-dev/sqlc/cmd/sqlc generate" \
	--build.bin "/bin/true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "database" \
	--build.include_ext "sql" \
	--log.main_only "true" $(LIVELOG)

live/templ:
	@printf "\033[1mstarting templ proxy, url: \033[0m\033[32m%s\033[0m\n" "http://localhost:7331"
	@go tool github.com/a-h/templ/cmd/templ generate --watch --proxy="http://localhost:8080" --open-browser=false --log-level="warn" $(LIVELOG)

live/server:
	@go run github.com/air-verse/air@v1.62.0 \
	--build.cmd "go build -o tmp/bin/main" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go,js,css" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true \
	--log.main_only "true" $(LIVELOG)

live/tailwind: node_modules
	@go run github.com/air-verse/air@v1.62.0 \
	--build.cmd "npx --yes @tailwindcss/cli -i ./resources/css/tailwind.css -o ./assets/css/styles.css" \
	--build.bin "/bin/true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "resources,templates" \
	--build.include_ext "js,css,templ" \
	--log.main_only "true" $(LIVELOG)

live/esbuild: node_modules
	@npx --yes esbuild ./resources/ts/index.ts --bundle --outdir=assets/js --watch=forever $(LIVELOG)

live/sync_assets:
	@go run github.com/air-verse/air@v1.62.0 \
	--build.cmd "go tool github.com/a-h/templ/cmd/templ generate --notify-proxy" \
	--build.bin "/bin/true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "assets" \
	--build.include_ext "js,css" \
	--log.main_only "true" $(LIVELOG)

live:
	$(MAKE) -j6 live/sql live/templ live/server live/esbuild live/sync_assets live/tailwind
