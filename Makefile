# Replace with your output binary name
BIN = app

LOGFILE = build.log
LOG = >> $(LOGFILE) 2>&1
LIVELOGFILE = live.log
LIVELOG = >> $(LIVELOGFILE) 2>&1

GO = go
SQLC = $(GO) tool github.com/sqlc-dev/sqlc/cmd/sqlc
TEMPL = $(GO) tool github.com/a-h/templ/cmd/templ

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
	$(TEMPL) generate -v $(LOG)
	@touch $@

.PHONY: templates
templates: $(TEMPLATES_GEN)

.PHONY: clean/templates
clean/templates:
	rm -rf $(TEMPLATES_GEN)

node_modules: package.json package-lock.json
	npm install $(LOG)

.PHONY: clean/node_modules
clean/node_modules:
	rm -rf node_modules

ESBUILD_IN := $(wildcard resources/ts/*.ts) $(wildcard resources/ts/*.js)
ESBUILD_OUT := $(addsuffix .js,$(basename $(patsubst resources/ts/%,assets/js/%,$(ESBUILD_IN))))
$(ESBUILD_OUT) &: $(ESBUILD_IN) node_modules
	@mkdir -p assets/js
	npx --yes esbuild $(ESBUILD_IN) --bundle --outdir=assets/js $(LOG)
	@touch $@

.PHONY: assets/js
assets/js: $(ESBUILD_OUT)

.PHONY: clean/assets/js
clean/assets/js:
	rm -rf $(ESBUILD_OUT)

assets/css/styles.css: resources/css/tailwind.css $(TEMPLATES) node_modules
	npx --yes @tailwindcss/cli -i ./resources/css/tailwind.css -o ./$@ $(LOG)
	@touch $@

.PHONY: assets/css
assets/css: assets/css/styles.css

.PHONY: clean/assets/css
clean/assets/css:
	rm -rf assets/css/styles.css

.PHONY: prep
prep:
	@$(MAKE) -j4 assets/css assets/js templates sql

$(BIN): prep
	CGO_ENABLED=0 $(GO) build $(LOG)

.PHONY: build/log-notice
build/log-notice:
	@printf "\033[34m\033[1m%s\033[0m\033[0m\n" "logs are saved into $(LOGFILE)"

.PHONY: build
build: build/log-notice $(BIN)
	@printf "\033[32m\033[1m%s\033[0m\033[0m\n" "success"

.PHONY: clean/build
clean/build:
	rm -rf $(BIN) $(LOGFILE)

.PHONY: run
run: prep
	$(GO) run .

# dist: build scripts/release.sh
# 	./scripts/release.sh

live/sql:
	@$(GO) run github.com/air-verse/air@v1.62.0 \
	--build.cmd "$(GO) tool github.com/sqlc-dev/sqlc/cmd/sqlc generate" \
	--build.bin "/bin/true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "database" \
	--build.include_ext "sql" \
	--log.main_only "true" $(LIVELOG)

live/templ:
	@printf "\033[1mstarting templ proxy, url:\033[0m \033[32m\033[4m%s\033[0m\033[0m\n" "http://localhost:7331"
	@$(GO) tool github.com/a-h/templ/cmd/templ generate --watch --proxy="http://localhost:8080" --open-browser=false --log-level="warn" $(LIVELOG)

live/server:
	@$(GO) run github.com/air-verse/air@v1.62.0 \
	--build.cmd "$(GO) build -o tmp/bin/main" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go,js,css" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true \
	--log.main_only "true" $(LIVELOG)

live/tailwind: node_modules
	@$(GO) run github.com/air-verse/air@v1.62.0 \
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
	@$(GO) run github.com/air-verse/air@v1.62.0 \
	--build.cmd "$(GO) tool github.com/a-h/templ/cmd/templ generate --notify-proxy" \
	--build.bin "/bin/true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "assets" \
	--build.include_ext "js,css" \
	--log.main_only "true" $(LIVELOG)

.PHONY: live/log-notice
live/log-notice:
	@printf "\033[1m\033[34m%s\033[0m\033[0m\n" "logs are saved into $(LIVELOGFILE)"

# See:
# https://templ.guide/developer-tools/live-reload-with-other-tools
live: live/log-notice
	@$(MAKE) -j6 live/sql live/templ live/server live/esbuild live/sync_assets live/tailwind

.PHONY: clean/live
clean/live:
	rm -rf $(LIVELOGFILE)

.PHONY: clean
clean:
	$(MAKE) -j7 clean/assets/js clean/assets/css clean/node_modules clean/sql clean/templates clean/build clean/live
