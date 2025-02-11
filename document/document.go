package document

import (
	"io"
	"mime/multipart"
	"os"
)

func SaveDocumentFile(file *multipart.FileHeader) (string, error) {
	filePath := "./uploads/" + file.Filename
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)

	defer func(dst *os.File) {
		_ = dst.Close()
	}(dst)

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return filePath, nil
}
