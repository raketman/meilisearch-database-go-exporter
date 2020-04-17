config.json.dist -> config.json

```
{
    "host": "http://0.0.0.0:7700", //host meili
    "key": null, //api_key meili
    "works": [
        {
          "index": "stopp", // индекс в meili
          "primary": "id", // primary_id в индекс
          "delete_before": true, // нужно ли удалять индекс перед созданием
          "searchable_attributes": [ 
            "parentguid",
            "housenum",
            "buildnum",
            "structnum"
          ],
          "displayed_attributes": [],
          "db_driver": "postgres", // driver postgres|mysql  | can use env variables in db_dsn, example env:DBDSN
          "db_dsn": "host=127.0.0.1 port=5432 user=postgres password=password dbname=test sslmode=disable", // DSN, зависит от типа базы
          "query": "SELECT id, name, key  FROM test ORDER BY id ASC LIMIT :limit OFFSET :offset", // limit, offset позволяет организовать пачки, без них зациклится
          "sleep": 1000, // время ожидания между пачками для разгрузки  системы
          "thread": 10, // количество потоков
          "limit": 1000, // размер пачки
          "offset": 0 // Стартовый оffset
        },
        {
          "index": "my_second_index",
          "primary": "aoguid",
          "db_driver": "mysql",
          "db_dsn": "env:DBDSN", // example for ENV
          "query": "SELECT aoguid, parentguid, offname FROM address_list ORDER BY aoguid ASC LIMIT :offset, :limit",
          "sleep": 1000,
          "thread": 10,
          "limit": 1000,
          "offset": 0 // Стартовый оffset
        }
    ]
}
```

Запуск
go run main.go config.go  exporter.go 
