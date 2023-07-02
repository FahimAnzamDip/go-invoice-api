package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	u "github.com/fahimanzamdip/go-invoice-api/utils"
)

var uploadDir string = "./public/uploads"

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // Max file size 32MB
	if err != nil {
		u.Respond(w, u.Message(false, "Failed to parse multipart form"))
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		u.Respond(w, u.Message(false, "Failed to retrieve file"))
		return
	}

	defer file.Close()

	// Generate a unique filename
	filename := generateUniqueFilename(handler.Filename)

	// Save the file to a specific location on server
	path, err := saveFile(file, filename)
	if err != nil {
		u.Respond(w, u.Message(false, "Failed to save image file"))
		return
	}

	u.Respond(w, map[string]interface{}{"status": true, "image_url": path})
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.FormValue("file_name")

	if fileName == "" {
		u.Respond(w, u.Message(false, "File name is missing"))
		return
	}

	filePath := uploadDir + "/" + fileName

	err := os.Remove(filePath)
	if err != nil {
		u.Respond(w, u.Message(false, "Failed to delete file"))
		return
	}

	u.Respond(w, map[string]interface{}{"status": true, "file_name": fileName})
}

func saveFile(file io.Reader, filename string) (string, error) {
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create upload directory: %s", err.Error())
	}

	filePath := filepath.Join(uploadDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %s", err.Error())
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %s", err.Error())
	}

	return filePath, nil
}

func generateUniqueFilename(originalFilename string) string {
	// Timestmap based filename
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
