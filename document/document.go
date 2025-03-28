package document

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
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

// SaveDocumentFileInBucket uploads a file to GCP Cloud Storage and returns the file URL or an error.
func SaveDocumentFileInBucket(fileHeader *multipart.FileHeader, bucket *storage.BucketHandle) (string, error) {
	ctx := context.Background()
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open fileHeader: %w", err)
	}
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)

	object := bucket.Object(fmt.Sprintf("uploads/%d-%s", time.Now().Unix(), fileHeader.Filename))
	writer := object.NewWriter(ctx)
	defer func(writer *storage.Writer) {
		_ = writer.Close()
	}(writer)

	if _, err = io.Copy(writer, src); err != nil {
		return "", fmt.Errorf("failed to copy fileHeader to storage: %w", err)
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket.BucketName(), object.ObjectName())
	return url, nil
}

// extractFilePath extracts the GCS file path from the full URL
func extractFilePath(gcsURL string, bucket *storage.BucketHandle) (string, error) {
	prefix := fmt.Sprintf("https://storage.googleapis.com/%s/", bucket.BucketName())
	if !strings.HasPrefix(gcsURL, prefix) {
		return "", fmt.Errorf("invalid GCS URL")
	}

	// Remove the prefix to get the bucket name and file path
	path := strings.TrimPrefix(gcsURL, prefix)

	return path, nil
}

// DeleteDocumentFileFromBucket deletes a file from GCP Cloud Storage and optionally returns an error
func DeleteDocumentFileFromBucket(fileURL string, bucket *storage.BucketHandle) error {
	ctx := context.Background()
	filePath, err := extractFilePath(fileURL, bucket)
	if err != nil {
		return fmt.Errorf("failed to extract file path from URL: %w", err)
	}
	object := bucket.Object(filePath)
	if err := object.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}
	return nil
}

// ExtractTextFromDocumentFile extracts text content from a document file based on its file extension.
// Supported formats include plain text (TXT, MD) and PDF. Unsupported formats return an error.
func ExtractTextFromDocumentFile(fileHeader *multipart.FileHeader) (string, error) {
	fileExtension := extension(filepath.Ext(fileHeader.Filename))

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open fileHeader: %w", err)
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)
	data, err := io.ReadAll(file)

	switch fileExtension {
	case TXT:
		return ExtractTextFromPlainText(data), nil
	case MD:
		return ExtractTextFromPlainText(data), nil
	case PDF:
		return ExtractTextFromPDF(data)
	case DOCX:
		return ExtractTextFromDocx(data)
	}
	return "", errors.New("unsupported format")
}

// ExtractTextFromPlainText reads the content of a plain text file and returns it as a string.
// It takes the file path as input and returns an error if the operation fails.
func ExtractTextFromPlainText(data []byte) string {
	return string(data)
}

// ExtractTextFromPDF extracts text from a PDF file and returns it as a string.
func ExtractTextFromPDF(data []byte) (string, error) {
	// Open the PDF file
	doc, err := fitz.NewFromReader(bytes.NewReader(data))
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

// ExtractTextFromDocx extracts text from a DOCX file given as []byte.
func ExtractTextFromDocx(data []byte) (string, error) {
	// Read DOCX from bytes
	reader := bytes.NewReader(data)

	// Read DOCX from memory
	doc, err := docx.ReadDocxFromMemory(reader, int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("failed to open DOCX from bytes: %w", err)
	}
	defer func(doc *docx.ReplaceDocx) {
		_ = doc.Close()
	}(doc)

	// Extract raw text
	text := doc.Editable().GetContent()

	// Remove XML tags (if needed)
	cleanText, err := removeXMLTags(text)
	if err != nil {
		return "", fmt.Errorf("failed to remove XML tags: %w", err)
	}

	return cleanText, nil
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
