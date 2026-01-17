# kickstart-go

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

## Install from binary

### Download app

```sh
curl -L "https://github.com/AlturaLabsCR/kickstart-go/releases/latest/download/app-amd64-linux" -o /usr/local/bin/app
chmod 755 /usr/local/bin/app
```

### Add user `app`

```sh
useradd -s /bin/bash -m -d /var/lib/app app
```

### App's configuration

```sh
curl -L "https://raw.githubusercontent.com/AlturaLabsCR/kickstart-go/refs/heads/main/contrib/config.toml" -o /etc/app/config.toml
```

### App's data (only if using sqlite and/or storage Type=local)

```sh
mkdir /var/lib/app/data
chown -R app:app /var/lib/app/data
```

### Systemd service

```sh
curl -L "https://raw.githubusercontent.com/AlturaLabsCR/kickstart-go/refs/heads/main/contrib/systemd/app.service" -o /etc/systemd/system/app.service
systemctl daemon-reload
systemctl enable --now app.service
```

### Done

Visit [localhost:8080](http://localhost:8080)

## Roadmap

- [X] Logout button
- [X] Light user management
- [X] SessionData is supposed to identify the login EVENT, meaning: when the session started, the device from which it originated, MAYBE last login, not stuff like permissions and username, those should be queried on every request. (Idea, use a view to get user info such as name and ~perms~ id and created time in one query)
- [ ] Sidebar layout in protected view
- [ ] Tests
