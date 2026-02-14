# PDF Password Removal API

A simple REST API built with Go that removes password encryption from PDF files.

## Features

- 🔓 Remove password protection from encrypted PDFs
- 🚀 Simple REST API interface
- ✅ Built with the open-source `pdfcpu` library
- 🛡️ Proper error handling for invalid passwords and files

## Prerequisites

- Go 1.25.0 or higher (for local development)
- Docker (for containerized deployment)

## Deployment Options

You can run this API either locally with Go or using Docker.

### Option 1: Docker (Recommended)

#### Using Docker Compose

The easiest way to run the API:

```bash
# Start the container
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the container
docker-compose down
```

#### Using Docker CLI

```bash
# Build the image
docker build -t pdf-decrypt-api .

# Run the container
docker run -d -p 8080:8080 --name pdf-decrypt-api pdf-decrypt-api

# View logs
docker logs -f pdf-decrypt-api

# Stop the container
docker stop pdf-decrypt-api

# Remove the container
docker rm pdf-decrypt-api
```

#### Running on a Different Port

If port 8080 is already in use:

```bash
# Run on port 8081 instead
docker run -d -p 8081:8080 --name pdf-decrypt-api pdf-decrypt-api
```

Then access the API at `http://localhost:8081/api/decrypt-pdf`

### Option 2: Local Installation

1. Clone or navigate to the project directory:
```bash
cd /Users/teammobile/dev/pdf-tools
```

2. Install dependencies:
```bash
go mod tidy
```

## Usage

### Start the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`

### API Endpoint

**POST** `/api/decrypt-pdf`

Decrypts a password-protected PDF and returns the unencrypted version.

#### Request

- **Content-Type**: `multipart/form-data`
- **Fields**:
  - `pdf` (file): The encrypted PDF file
  - `password` (string): The password to decrypt the PDF

#### Response

- **Success (200)**: Returns the decrypted PDF file
- **Error (400)**: Invalid request, missing fields, or decryption failure
- **Error (401)**: Incorrect password
- **Error (405)**: Method not allowed (use POST)
- **Error (500)**: Internal server error

### Example Usage

Using `curl`:

```bash
curl -X POST http://localhost:8080/api/decrypt-pdf \
  -F "pdf=@encrypted_document.pdf" \
  -F "password=yourpassword" \
  --output decrypted_document.pdf
```

Using JavaScript (fetch):

```javascript
const formData = new FormData();
formData.append('pdf', pdfFile); // File object
formData.append('password', 'yourpassword');

fetch('http://localhost:8080/api/decrypt-pdf', {
  method: 'POST',
  body: formData
})
  .then(response => response.blob())
  .then(blob => {
    // Download the decrypted PDF
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'decrypted.pdf';
    a.click();
  });
```

## Build

To build a standalone executable:

```bash
go build -o pdf-decrypt-api
```

Then run:

```bash
./pdf-decrypt-api
```

## CORS Configuration

The API is configured to accept requests from `http://localhost:4200`, making it compatible with Angular development servers. The CORS middleware handles:

- Preflight OPTIONS requests
- Cross-origin POST requests
- Proper CORS headers for browser compatibility

If you need to allow requests from a different origin, modify the `corsMiddleware` function in `main.go`:

```go
w.Header().Set("Access-Control-Allow-Origin", "http://your-domain:port")
```

## Error Handling

The API handles various error scenarios:

- **Missing PDF file**: Returns 400 Bad Request
- **Missing password**: Returns 400 Bad Request
- **Incorrect password**: Returns 401 Unauthorized
- **Invalid PDF file**: Returns 400 Bad Request
- **Server errors**: Returns 500 Internal Server Error

## License

This project uses the `pdfcpu` library, which is licensed under Apache License 2.0.

## Dependencies

- [pdfcpu](https://github.com/pdfcpu/pdfcpu) - PDF processing library for Go
