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

## Deploy

```sh
docker compose up

# or sidecar a database
docker compose -f docker-compose.db.yml --profile [postgres/sqlite] up
```

## Install from binary

### Download app

```sh
curl -L "https://github.com/AlturaLabsCR/go-web-kickstart/releases/latest/download/app-amd64-linux" -o /usr/local/bin/app
chmod 755 /usr/local/bin/app
```

### Add user `app`

```sh
useradd -s /bin/bash -m -d /var/lib/app app
```

### App's configuration

```sh
curl -L "https://raw.githubusercontent.com/AlturaLabsCR/go-web-kickstart/refs/heads/main/contrib/config.toml" -o /etc/app/config.toml
```

### App's data (only if using sqlite and/or storage Type=local)

```sh
mkdir /var/lib/app/data
chown -R app:app /var/lib/app/data
```

### Systemd service

```sh
curl -L "https://raw.githubusercontent.com/AlturaLabsCR/go-web-kickstart/refs/heads/main/contrib/systemd/app.service" -o /etc/systemd/system/app.service
systemctl daemon-reload
systemctl enable --now app.service
```

### Done

Visit [localhost:8080](http://localhost:8080)
