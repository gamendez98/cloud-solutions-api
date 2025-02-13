package document

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

type extension string

const (
	PDF extension = ".pdf"
	TXT extension = ".txt"
	MD  extension = ".md"
	DOC extension = ".doc"
)

// SaveDocumentFile saves an uploaded file to the server and returns the file path or an error if the operation fails.
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

// ExtractTextFromDocumentFile extracts text content from a document file based on its file extension.
// Supported formats include plain text (TXT, MD) and PDF. Unsupported formats return an error.
func ExtractTextFromDocumentFile(filePath string) (string, error) {
	fileExtension := extension(filepath.Ext(filePath))
	switch fileExtension {
	case TXT:
		return ExtractTextFromPlainText(filePath)
	case MD:
		return ExtractTextFromPlainText(filePath)
	case PDF:
		return ExtractTextFromPDF(filePath)
	case DOC:
		return "", errors.New("unsupported format")
	}
	return "", errors.New("unsupported format")
}

// ExtractTextFromPlainText reads the content of a plain text file and returns it as a string.
// It takes the file path as input and returns an error if the operation fails.
func ExtractTextFromPlainText(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// ExtractTextFromPDF extracts text from a PDF file and returns it as a string.
func ExtractTextFromPDF(filePath string) (string, error) {
	// Open the PDF file.
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	pdfReader, err := model.NewPdfReader(file)
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %v", err)
	}
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", fmt.Errorf("failed to get number of pages: %v", err)
	}

	var extractedText string
	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			return "", fmt.Errorf("failed to get page %d: %v", i, err)
		}

		ex, err := extractor.New(page)
		if err != nil {
			return "", fmt.Errorf("failed to create extractor for page %d: %v", i, err)
		}

		text, err := ex.ExtractText()
		if err != nil {
			return "", fmt.Errorf("failed to extract text from page %d: %v", i, err)
		}

		extractedText += text + "\n"
	}

	return extractedText, nil
}
