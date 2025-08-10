package containers

import (
	"mime/multipart"
)

func ParseCardRequests(formFile *multipart.FileHeader) ([]CardRequest, error) {
	file, err := formFile.Open()

	if err != nil {
		return nil, err
	}

	defer file.Close()

	return nil, nil
}
