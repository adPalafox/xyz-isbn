package converter

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"xyz-isbn/models"
)

func IsbnConverter(book *models.Book) (bool, error) {
	var isUpdated = false
	if book.ISBN10 == "" && book.ISBN13 != "" {
		fmt.Println("ISBN 13 -> 10")
		isbn10, err := Isbn13To10(book.ISBN13)
		if err != nil {
			return false, fmt.Errorf("error converting ISBN13 to ISBN10: %w", err)
		}
		if isbn10 != book.ISBN10 {
			isUpdated = true
		}
		book.ISBN10 = isbn10
	} else if book.ISBN13 == "" && book.ISBN10 != "" {
		fmt.Println("ISBN 10 -> 13")
		isbn13, err := Isbn10To13(book.ISBN13)
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

func CheckIfExists(book *models.Book) error {
	if book.ISBN10 == "" || book.ISBN13 == "" {
		return errors.New("missing both ISBN10 or ISBN13")
	}
	return nil
}

func computeCheckDigitIsbn13(isbn13 string) (string, error) {
	multipliers := []int{1, 3}
	sum := 0
	for idx, char := range isbn13[:12] {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return "", fmt.Errorf("failure in type conversion")
		}

		sum += digit * multipliers[idx%2]
	}

	checkDigit := (10 - sum%10) % 10
	return strconv.Itoa(checkDigit), nil
}

func computeCheckDigitIsbn10(isbn10 string) (string, error) {
	total := 0
	for i, c := range isbn10[:9] {
		char := string(c)
		digit, err := strconv.Atoi(char)
		if err != nil {
			return "", fmt.Errorf("failure in type conversion")
		}

		total += digit * (10 - i)
	}

	checkDigit := (11 - total%11) % 11
	if checkDigit == 10 {
		return "X", nil
	}

	return strconv.Itoa(checkDigit), nil
}

func ValidateIsbns(isbn string) bool {
	var checkDigitString string
	var err error

	if len(isbn) != 13 && len(isbn) != 10 {
		return false
	}

	if len(isbn) == 13 {
		if !strings.HasPrefix(isbn, "978") {
			return false
		}
		checkDigitString, err = computeCheckDigitIsbn13(isbn)
		if err != nil {
			return false
		}
	} else if len(isbn) == 10 {
		checkDigitString, err = computeCheckDigitIsbn10(isbn)
		if err != nil {
			return false
		}
	}
	computerCheckDigit := isbn[len(isbn)-1]
	return checkDigitString == string(computerCheckDigit)
}

func Isbn10To13(isbn10 string) (string, error) {
	if len(isbn10) != 10 {
		return "", fmt.Errorf("invalid length ISBN 10")
	}

	if !ValidateIsbns(isbn10) {
		return "", fmt.Errorf("not valid ISBN 10")
	}

	first12Digits := "978" + string(isbn10[:9])
	checkDigitStr, err := computeCheckDigitIsbn13(first12Digits)
	if err != nil {
		return "", err
	}
	return first12Digits + checkDigitStr, nil
}

func Isbn13To10(isbn13 string) (string, error) {
	if len(isbn13) != 13 {
		return "", fmt.Errorf("invalid length ISBN 13")
	}

	if !strings.HasPrefix(isbn13, "978") {
		return "", fmt.Errorf("not valid ISBN 13 prefix")
	}

	if !ValidateIsbns(isbn13) {
		return "", fmt.Errorf("not valid ISBN 13")
	}

	checkDigitStr, err := computeCheckDigitIsbn10(string(isbn13[3:12]))
	if err != nil {
		return "", err
	}

	return string(isbn13[3:12]) + checkDigitStr, nil
}
