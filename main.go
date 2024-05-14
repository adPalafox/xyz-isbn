package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"xyz-isbn/constant"
	"xyz-isbn/processor"
)

func main() {
	fmt.Println("Started microservice")
	csvFile, err := os.Create("books.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	processor := processor.NewBookProcessor(csvFile)

	ctx := context.Background()

	var page = 1
	url := fmt.Sprintf("%vlist?length=100&page=1&sort=title&order=asc", constant.BASE_URL)
	for {
		books, err := processor.FetchBookList(ctx, url)
		if err != nil {
			log.Printf("Error fetching book list: %v\n", err)
			break
		}

		if err := processor.ProcessBooks(ctx, books); err != nil {
			log.Printf("Error processing books: %v\n", err)
			break
		}

		if len(books) < constant.MAX_BATCH_COUNT {
			break
		}

		page++
		url = fmt.Sprintf("%vlist?length=100&page=%d&sort=title&order=asc", constant.BASE_URL, page)
	}

	fmt.Println("Finished processing all book pages.")
}
