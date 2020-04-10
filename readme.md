config.json.dist -> config.json

```
{
    "host": "http://0.0.0.0:7700", //host meili
    "key": null, //api_key meili
    "works": [
        {
          "index": "stopp", // индекс в meili
          "primary": "id", // primary_id в индекс
          "db_driver": "postgres", // driver postgres|mysql
          "db_dsn": "host=127.0.0.1 port=5432 user=postgres password=password dbname=test sslmode=disable", // DSN, зависит от типа базы
          "query": "SELECT id, name, key  FROM test ORDER BY id ASC LIMIT :limit OFFSET :offset", // limit, offset позволяет организовать пачки, без них зациклится
          "sleep": 1000, // время ожидания между пачками для разгрузки  системы
          "thread": 10, // количество потоков
          "limit": 1000 // размер пачки
        }
    ]
}
```

Запуск
go run main.go config.go work_control.go  exporter.go 
