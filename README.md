# migrations

1. run PostgreSQL for test
```sh
docker-compose up
```

2. create 'test_database'
```sh
docker exec -i postgres psql -U postgres -c "drop database if exists test_database" && \
docker exec -i postgres psql -U postgres -c "create database test_database"
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