package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/meilisearch/meilisearch-go"
	"log"
	"os"
	"strings"
	"time"
)



type Exporter struct {
	Thread int
	Work int
}

func (exporter Exporter) Process(client *meilisearch.Client, work Work) {

	// work.DB_DSN get from env
	if strings.Contains(work.DB_DSN, "env:") {

		work.DB_DSN = strings.Replace(work.DB_DSN, "env:","", 1)

		work.DB_DSN, _ = os.LookupEnv(work.DB_DSN)

		if len(work.DB_DSN) == 0 {
			log.Fatal("Work: ", exporter.Work, " THREAD: ", exporter.Thread, " NO DB_DSN")
		}

	}

	gormDb, err := gorm.Open(
		work.DB_DRIVER,
		work.DB_DSN,
	)

	if err != nil {
		log.Fatal("Work:", exporter.Work, "THREAD:", exporter.Thread, "ERROR", err)
	}


	defer gormDb.Close()

	haveRecord := true
	offset := work.Offset + exporter.Thread * work.Limit
	increment := work.Thread * work.Limit

	for {
		query := work.Query
		// Replace the "cat" with a "dog."
		query = strings.Replace(query, ":limit", fmt.Sprintf("%v", work.Limit), -1)
		query = strings.Replace(query, ":offset", fmt.Sprintf("%v", offset), -1)

		//log.Println("QUERY:", exporter.Thread, query)


		rows, errRow := gormDb.Raw(query).Rows() // (*sql.Rows, error)

		if errRow != nil {
			log.Println("Work:", exporter.Work, "THREAD:", exporter.Thread, " OFFSET:", offset," ERROROW:", errRow)
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

				if raw_value != nil {
					raw_value = string(fmt.Sprintf("%s", raw_value))
				}
				item[colName] = raw_value
			}

			documents = append(documents, item)
		}

		haveRecord = len(documents) > 0 // До тех пор пока есть данные

		rows.Close()

		if !haveRecord {
			break
		}

		updateRes, err := client.Documents(work.Index).AddOrReplace(documents) // => { "updateId": 0 }

		if err != nil {
			log.Println("Work:", exporter.Work, "THREAD:", exporter.Thread, " OFFSET:", offset," MEILIERROR:", err)
			continue
		}
		log.Println("Work:", exporter.Work, "THREAD:", exporter.Thread, " OFFSET:", offset," UPDATES:",  updateRes.UpdateID)

		time.Sleep(time.Duration(work.Sleep) * time.Millisecond)

		// Увеличиваем, что выбрать всю базу
		offset = offset + increment
	}





}
