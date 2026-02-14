package main

import (
	"fmt"
	"log"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func main() {
	// First, let's download or use a simple PDF
	// For testing, we'll create a basic PDF using pdfcpu's import images feature
	// Or we can use a simpler approach - just encrypt an existing PDF

	inputPDF := "test_sample.pdf"
	encryptedPDF := "test_encrypted.pdf"
	password := "testpassword123"

	// Create a simple text-based PDF
	// Using ImportImages with a simple approach
	err := api.ImportImagesFile([]string{}, inputPDF, nil, nil)
	if err != nil {
		// If that doesn't work, let's try a different approach
		log.Printf("Note: Could not create PDF with images, trying alternative method")

		// Alternative: Create a very basic PDF structure manually
		// For now, let's just inform the user to provide a sample PDF
		fmt.Println("Please provide a sample PDF file named 'test_sample.pdf' in this directory")
		fmt.Println("Or download one from: https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf")
		return
	}

	// Encrypt the PDF with password
	conf := model.NewAESConfiguration(password, password, 256)
	err = api.EncryptFile(inputPDF, encryptedPDF, conf)
	if err != nil {
		log.Fatalf("Failed to encrypt PDF: %v", err)
	}

	fmt.Printf("✓ Created encrypted test PDF: %s\n", encryptedPDF)
	fmt.Printf("✓ Password: %s\n", password)
}
