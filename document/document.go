package document

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gen2brain/go-fitz"
	"github.com/nguyenthenguyen/docx"
)

type extension string

const (
	PDF  extension = ".pdf"
	TXT  extension = ".txt"
	MD   extension = ".md"
	DOCX extension = ".docx"
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

// SaveDocumentFileInBucket uploads a file to GCP Cloud Storage and returns the file URL or an error.
func SaveDocumentFileInBucket(file *multipart.FileHeader, bucket *storage.BucketHandle) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create storage client: %w", err)
	}
	defer func(client *storage.Client) {
		_ = client.Close()
	}(client)

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)

	object := bucket.Object(fmt.Sprintf("uploads/%d-%s", time.Now().Unix(), file.Filename))
	writer := object.NewWriter(ctx)
	defer func(writer *storage.Writer) {
		_ = writer.Close()
	}(writer)

	if _, err = io.Copy(writer, src); err != nil {
		return "", fmt.Errorf("failed to copy file to storage: %w", err)
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket.BucketName(), object.ObjectName())
	return url, nil
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
	case DOCX:
		return ExtractTextFromDocx(filePath)
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
	// Open the PDF file
	doc, err := fitz.New(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %v", err)
	}
	defer func(doc *fitz.Document) {
		_ = doc.Close()
	}(doc)

	// Extract text from each page
	var extractedText string
	for i := 0; i < doc.NumPage(); i++ {
		text, err := doc.Text(i)
		if err != nil {
			return "", fmt.Errorf("failed to extract text from page %d: %v", i+1, err)
		}
		extractedText += text + "\n"
	}

	return extractedText, nil
}

// ExtractTextFromDocx extracts text from a DOCX file and returns it as a string.
func ExtractTextFromDocx(filePath string) (string, error) {
	// Open the .doc file
	doc, err := docx.ReadDocxFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open .doc file: %v", err)
	}
	defer func(doc *docx.ReplaceDocx) {
		_ = doc.Close()
	}(doc)

	// Extract text from the document
	text := doc.Editable().GetContent()
	text, err = removeXMLTags(text)
	if err != nil {
		return "", fmt.Errorf("failed to remove XML tags: %v", err)
	}

	return text, nil
}

// removeXMLTags removes all XML tags and returns the text content.
func removeXMLTags(xmlContent string) (string, error) {
	// Create a decoder for the XML content
	decoder := xml.NewDecoder(strings.NewReader(xmlContent))

	// Buffer to store the extracted text
	var textContent bytes.Buffer

	// Iterate through the XML tokens
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break // End of XML content
		}
		if err != nil {
			return "", fmt.Errorf("failed to decode XML: %v", err)
		}

		// Extract text content from character data
		switch t := token.(type) {
		case xml.CharData:
			textContent.Write(t)
			textContent.WriteRune('\n')
		}
	}

	return strings.TrimSpace(textContent.String()), nil
}
