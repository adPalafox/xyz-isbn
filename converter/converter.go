package converter

import (
	"errors"
	"fmt"
	"xyz-isbn/models"
)

func IsbnConverter(book *models.Book) (bool, error) {
	var isUpdated = false
	if book.ISBN10 == "" && book.ISBN13 != "" {
		fmt.Println("ISBN 13 -> 10")
		isbn10, err := ISBN13To10(book.ISBN13)
		if err != nil {
			return false, fmt.Errorf("error converting ISBN13 to ISBN10: %w", err)
		}
		if isbn10 != book.ISBN10 {
			isUpdated = true
		}
		book.ISBN10 = isbn10
	} else if book.ISBN13 == "" && book.ISBN10 != "" {
		fmt.Println("ISBN 10 -> 13")
		isbn13, err := ISBN10To13(book.ISBN13)
		if err != nil {
			return false, fmt.Errorf("error converting ISBN10 to ISBN13: %w", err)
		}
		if isbn13 != book.ISBN13 {
			isUpdated = true
		}
		book.ISBN13 = isbn13
	}
	if isUpdated {
		fmt.Printf("Successfully converted: %v\n", book.ISBN13)
	}
	return isUpdated, nil
}

func ValidateISBNs(book *models.Book) error {
	if book.ISBN10 == "" && book.ISBN13 == "" {
		return errors.New("missing both ISBN10 and ISBN13")
	}
	return nil
}

func ISBN13To10(isbn13 string) (string, error) {
	return isbn13, nil
}

func ISBN10To13(isbn10 string) (string, error) {
	return isbn10, nil
}
