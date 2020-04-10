package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/meilisearch/meilisearch-go"
	"log"
	"strings"
	"time"
)



type Exporter struct {
	Thread int
}

func (exporter Exporter) Process(client *meilisearch.Client, work Work) {

	// TODO: сделать split driver
	gormDb, err := gorm.Open(
		work.DB_DRIVER,
		work.DB_DSN,
	)

	if err != nil {
		log.Fatal("THREAD:", exporter.Thread, "ERROR", err)
	}


	defer gormDb.Close()

	haveRecord := true
	offset := exporter.Thread * work.Limit
	increment := work.Thread * work.Limit

	for {
		query := work.Query
		// Replace the "cat" with a "dog."
		query = strings.Replace(query, ":limit", fmt.Sprintf("%v", work.Limit), -1)
		query = strings.Replace(query, ":offset", fmt.Sprintf("%v", offset), -1)

		//log.Println("QUERY:", exporter.Thread, query)


		rows, errRow := gormDb.Raw(query).Rows() // (*sql.Rows, error)
		defer rows.Close()

		if errRow != nil {
			log.Println("THREAD:", exporter.Thread, " OFFSET:", offset," ERROROW:", err)
			continue
		}

		columns, _ := rows.Columns()
		colNum := len(columns)

		var values = make([]interface{}, colNum)
		for i, _ := range values {
			var ii interface{}
			values[i] = &ii
		}

		documents := []map[string]interface{}{}

		for rows.Next() {
			rows.Scan(values...)
			item := map[string]interface{}{}

			for i, colName := range columns {
				var raw_value = *(values[i].(*interface{}))

				item[colName] = string(fmt.Sprintf("%s", raw_value))
			}

			documents = append(documents, item)
		}

		haveRecord = len(documents) > 0 // До тех пор пока есть данные

		updateRes, err := client.Documents(work.Index).AddOrReplace(documents) // => { "updateId": 0 }

		if err != nil {
			log.Println("THREAD:", exporter.Thread, " OFFSET:", offset," MEILIERROR:", err)
			continue
		}
		log.Println("THREAD:", exporter.Thread, " OFFSET:", offset," UPDATES:",  updateRes.UpdateID)

		if !haveRecord {
			break
		}

		time.Sleep(time.Duration(work.Sleep) * time.Millisecond)

		// Увеличиваем, что выбрать всю базу
		offset = offset + increment
	}





}
