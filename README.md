# migrations

1. run PostgreSQL for test
```sh
docker-compose up
```

2. create 'testdb'
```sh
docker exec -i fuel-management-postgres psql -U postgres -c "drop database if exists testdb" && \
docker exec -i fuel-management-postgres psql -U postgres -c "create database testdb"
```

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