# migrations

#### up commands

1. migrate up by one
```sh
go run ./cmd/up
```

2. migrate up all the pendings
```sh
go run ./cmd/up all
```

#### down commands

1. migrate down by one
```sh
go run ./cmd/down
```

2. migrate down all the migrations
```sh
go run ./cmd/down all
```

#### useful commands

```
go install golang.org/x/vuln/cmd/govulncheck@latest
~/go/bin/govulncheck ./...
```