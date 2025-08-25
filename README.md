# Go Web Kickstart

Kickstart repo for go, templ, sqlc fullstack web application.
It also includes example tailwindcss and esbuild usage.

## Build

```sh
go mod tidy
go tool github.com/a-h/templ/cmd/templ generate
go tool github.com/sqlc-dev/sqlc/cmd/sqlc generate
go build
```

## Live reload

```sh
make live
```
