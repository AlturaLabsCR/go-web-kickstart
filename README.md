# Go Web Kickstart

Blueprint for a go, templ, sqlc fullstack web application

## Development

```sh
make gen   # make generated files
make build # build the binary
make run   # run the binary
make dist  # create dist folder with multiple archs and OSs
make live  # live reload on change of templates, schema, css classes, etc
make clean
```

## Deployments

```sh
# railway-compatible Dockerfile
docker compose up

# or try with a database
docker compose -f docker-compose.db.yml --profile [postgres/sqlite] up
```
