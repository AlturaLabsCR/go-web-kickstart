# Replace with your output binary name
BIN = app

GO = go
BUILD_FLAGS = -ldflags "-s -w"
DIST_FOLDER = dist
GO_ENV = CGO_ENABLED=0

SQLC = github.com/sqlc-dev/sqlc/cmd/sqlc
TEMPL = github.com/a-h/templ/cmd/templ
AIR = github.com/air-verse/air@v1.63.6

NPM = npm
NPX = npx

GEN =

.PHONY: all
all: build

SQL = $(wildcard database/**/migrations/*.sql)
SQL += $(wildcard database/**/queries/*.sql)
GEN += database/.gen
database/.gen: $(SQL)
	(cd database/ && $(GO) tool $(SQLC) generate)
	@touch $@

.PHONY: clean/sql
clean/sql:
	rm -rf database/**/db
	rm -rf database/.gen

TEMPLATES = $(wildcard templates/*.templ)
TEMPLATES += $(wildcard templates/**/*.templ)
TEMPLATES_GEN := $(TEMPLATES:.templ=_templ.go)
GEN += $(TEMPLATES_GEN)
$(TEMPLATES_GEN) &: $(TEMPLATES)
	$(GO) tool $(TEMPL) generate
	@touch $@

.PHONY: clean/templates
clean/templates:
	rm -rf $(TEMPLATES_GEN)

node_modules: package.json package-lock.json
	$(NPM) ci $(LOG)

.PHONY: clean/node_modules
clean/node_modules:
	rm -rf node_modules

ESBUILD_IN := $(wildcard resources/ts/*.ts) $(wildcard resources/ts/*.js)
ESBUILD_OUT := $(addsuffix .js,$(basename $(patsubst resources/ts/%,assets/js/%,$(ESBUILD_IN))))
GEN += $(ESBUILD_OUT)
$(ESBUILD_OUT) &: $(ESBUILD_IN) node_modules
	@mkdir -p assets/js
	$(NPX) --yes esbuild $(ESBUILD_IN) --bundle --outdir=assets/js $(LOG)
	@touch $@

.PHONY: clean/assets/js
clean/assets/js:
	rm -rf $(ESBUILD_OUT)

GEN += assets/css/styles.css
assets/css/styles.css: resources/css/tailwind.css $(TEMPLATES) node_modules
	$(NPX) --yes @tailwindcss/cli -i ./resources/css/tailwind.css -o ./$@
	@touch $@

.PHONY: clean/assets/css
clean/assets/css:
	rm -rf assets/css/styles.css

.PHONY: clean/gen
clean/assets:
	rm -rf $(GEN)

.PHONY: gen
gen: $(GEN)

$(BIN): $(GEN)
	$(GO_ENV) $(GO) build $(BUILD_FLAGS)

.PHONY: build
build: $(BIN)

.PHONY: clean/build
clean/build:
	rm -rf $(BIN)

.PHONY: run
run: $(GEN)
	$(GO) run .

live/sql:
	@$(GO) run $(AIR) \
	--build.cmd "(cd database/ && $(GO) tool $(SQLC) generate)" \
	--build.bin "/bin/true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "database" \
	--build.include_ext "sql" \
	--log.main_only "true"

live/templ:
	@printf "\033[1mstarting templ proxy, url:\033[0m \033[32m\033[4m%s\033[0m\033[0m\n" "http://localhost:7331"
	@$(GO) tool github.com/a-h/templ/cmd/templ generate --watch --proxy="http://localhost:8080" --open-browser=false --log-level="warn"

live/server: $(GEN)
	@$(GO) run $(AIR) \
	--build.cmd "$(GO) build -o tmp/bin/main" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go,js,css,toml" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true \
	--log.main_only "true"

live/tailwind: node_modules
	@$(GO) run $(AIR) \
	--build.cmd "$(NPX) --yes @tailwindcss/cli -i ./resources/css/tailwind.css -o ./assets/css/styles.css" \
	--build.bin "/bin/true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "resources,templates" \
	--build.include_ext "js,css,templ" \
	--log.main_only "true"

live/esbuild: node_modules
	@$(NPX) --yes esbuild ./resources/ts/index.ts --bundle --outdir=assets/js --watch=forever $(LIVELOG)

live/sync_assets: assets
	@$(GO) run $(AIR) \
	--build.cmd "$(GO) tool github.com/a-h/templ/cmd/templ generate --notify-proxy" \
	--build.bin "/bin/true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "assets" \
	--build.include_ext "js,css" \
	--log.main_only "true"

# See:
# https://templ.guide/developer-tools/live-reload-with-other-tools
live: gen
	@$(MAKE) -j live/sql live/templ live/server live/esbuild live/sync_assets live/tailwind

$(DIST_FOLDER)/$(BIN)-x86_64-linux: $(GEN)
	@mkdir -p $(DIST_FOLDER)
	$(GO_ENV) GOARCH=amd64 GOOS=linux $(GO) build $(BUILD_FLAGS) -o $@

$(DIST_FOLDER)/$(BIN)-x86_64-windows.exe: $(GEN)
	@mkdir -p $(DIST_FOLDER)
	$(GO_ENV) GOARCH=amd64 GOOS=windows $(GO) build $(BUILD_FLAGS) -o $@

.PHONY: dist/x86_64
dist/amd64: $(DIST_FOLDER)/$(BIN)-x86_64-linux $(DIST_FOLDER)/$(BIN)-x86_64-windows.exe

.PHONY: dist
dist: dist/x86_64

.PHONY: clean/dist
clean/dist:
	rm -rf $(DIST_FOLDER)

.PHONY: clean
clean:
	$(MAKE) -j clean/assets/js clean/assets/css clean/node_modules clean/sql clean/templates clean/build clean/live clean/dist
