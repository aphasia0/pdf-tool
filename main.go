package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	allowedOrigin := os.Getenv("CORS_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:4200"
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	http.HandleFunc("/api/decrypt-pdf", corsMiddleware(DecryptPDFHandler))

	port := ":8080"
	log.Printf("Starting PDF decryption API server on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// DecryptPDFHandler handles the PDF decryption request
func DecryptPDFHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 50MB)
	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get the password from form
	password := r.FormValue("password")
	if password == "" {
		http.Error(w, "Password field is required", http.StatusBadRequest)
		return
	}

	// Get the PDF file from form
	file, header, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "PDF file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create temporary directory for processing
	tempDir, err := os.MkdirTemp("", "pdf-decrypt-*")
	if err != nil {
		log.Printf("Error creating temp directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)

	// Save uploaded file to temp location
	inputPath := filepath.Join(tempDir, "input.pdf")
	outputPath := filepath.Join(tempDir, "output.pdf")

	inputFile, err := os.Create(inputPath)
	if err != nil {
		log.Printf("Error creating input file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(inputFile, file)
	inputFile.Close()
	if err != nil {
		log.Printf("Error saving uploaded file: %v", err)
		http.Error(w, "Failed to save uploaded file", http.StatusInternalServerError)
		return
	}

	// Decrypt the PDF using pdfcpu
	// Create configuration with the password
	conf := model.NewDefaultConfiguration()
	conf.UserPW = password
	conf.OwnerPW = password

	err = api.DecryptFile(inputPath, outputPath, conf)
	if err != nil {
		log.Printf("Error decrypting PDF: %v", err)
		// Check if it's a password error
		errMsg := err.Error()
		if errMsg == "pdfcpu: please provide the correct password" ||
			errMsg == "pdfcpu: wrong password" ||
			errMsg == "pdfcpu: please provide the owner password" {
			http.Error(w, "Incorrect password provided", http.StatusUnauthorized)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to decrypt PDF: %v", err), http.StatusBadRequest)
		return
	}

	// Read the decrypted PDF
	decryptedPDF, err := os.ReadFile(outputPath)
	if err != nil {
		log.Printf("Error reading decrypted PDF: %v", err)
		http.Error(w, "Failed to read decrypted PDF", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"decrypted_%s\"", header.Filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(decryptedPDF)))

	// Send the decrypted PDF
	_, err = w.Write(decryptedPDF)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}

	log.Printf("Successfully decrypted PDF: %s", header.Filename)
}
