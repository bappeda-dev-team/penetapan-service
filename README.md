# PENETAPAN SERVICE

Snapshot data perencanaan saat ditetapkan

digunakan untuk realisasi dan laporan

# Struktur folder

### cmd
untuk routes dan handler (controller)

### http
digunakan untuk enpoint testing

# Running

``` sh
make run
```
atau
``` sh
go run ./cmd/api
```

# Database
Postgresql

## Migrations
Folder: db/migration

``` sh
flyway -user=postgres -password=root -url=jdbc:postgresql://localhost:5432/penetapan_local -locations="filesystem:./db/migration" migrate
```

# Generate Docs
Swagger

``` sh
swag init -g ./cmd/api/main.go
```
atau

``` sh
make swagger
```

# Test

``` sh
go test ./...
```
atau

``` sh
make test
```
