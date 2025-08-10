package containers

import (
	"mime/multipart"
)

func ParseCardDeposits(formFile *multipart.FileHeader) ([]CardDeposit, error) {
	file, err := formFile.Open()

	if err != nil {
		return nil, err
	}

	defer file.Close()

	return nil, nil
}
