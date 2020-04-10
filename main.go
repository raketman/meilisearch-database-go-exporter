package main

import (
	"github.com/meilisearch/meilisearch-go"
	"log"
)

func main() {
	var config Config

	config = Config{}

	config.Read("config.json")

	log.Println("Config:", config)

	var client = meilisearch.NewClient(meilisearch.Config{
		Host: config.Host,
		APIKey: config.Key,
	})

	log.Println(client)

	// Создадим канал, для обмена информацией
	works := make(chan int, len(config.Works))
	done := make(chan bool)

	// Запустим контроль работ
	WorkControl(works, done, "BASE WORKS", len(config.Works))

	// Для каждого запрос создадим индекс primary
	for _, itemWork := range config.Works {
		go func(work Work) {
			// Create an index if your index does not already exist
			index, _ := client.Indexes().Get(work.Index)

			if index == nil {
				_, err := client.Indexes().Create(meilisearch.CreateIndexRequest{
					UID:        work.Index,
					PrimaryKey: work.Primary,
				})
				if err != nil {
					log.Fatal("Create index: ", err)
				}
			}

			threads := make(chan int, len(config.Works))
			threadDone := make(chan bool)

			WorkControl(threads, threadDone, "THREAD WORK " + work.Index, work.Thread)

			for thread := 0; thread < work.Thread; thread++ {
				go func(thread int) {
					exporter := Exporter{
						Thread: thread,
					}
					// TODO: Сделать в виде интерфейса, чтобы любой exporter могу это сделать
					exporter.Process(client, work)

					threads <- 1
				}(thread)
			}

			<- threadDone
			// Отправим, что заверешнили
			works <- 1
		} (itemWork)
	}

	<- done
}
