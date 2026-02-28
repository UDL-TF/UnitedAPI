# MinIO Integration Documentation

## Overview
This API now includes MinIO integration for file storage operations including upload, download, delete, and listing files.

## Environment Variables
Set the following environment variables to configure MinIO:

```bash
MINIO_ENDPOINT=your-minio-endpoint.com:9000
MINIO_ACCESS_KEY_ID=your-username
MINIO_SECRET_ACCESS_KEY=your-password
MINIO_USE_SSL=false
MINIO_BUCKET_NAME=your-bucket-name
```

## API Endpoints

### Upload File
- **Method**: `POST`
- **Endpoint**: `/api/v1/storage/upload`
- **Content-Type**: `multipart/form-data`
- **Parameters**:
  - `file` (required): The file to upload
  - `object_name` (optional): Custom name for the object (defaults to original filename)

**Example using curl:**
```bash
curl -X POST http://localhost:8080/api/v1/storage/upload \
  -F "file=@/path/to/your/file.pdf" \
  -F "object_name=my-document.pdf"
```

### Download File
- **Method**: `GET`
- **Endpoint**: `/api/v1/storage/download/{objectName}`

**Example using curl:**
```bash
curl -X GET http://localhost:8080/api/v1/storage/download/my-document.pdf \
  --output downloaded-file.pdf
```

### Delete File
- **Method**: `DELETE`
- **Endpoint**: `/api/v1/storage/delete/{objectName}`

**Example using curl:**
```bash
curl -X DELETE http://localhost:8080/api/v1/storage/delete/my-document.pdf
```

### List Files
- **Method**: `GET`
- **Endpoint**: `/api/v1/storage/files`
- **Query Parameters**:
  - `prefix` (optional): Filter files by prefix

**Example using curl:**
```bash
curl -X GET "http://localhost:8080/api/v1/storage/files?prefix=documents/"
```

### Generate Presigned URL
- **Method**: `GET`
- **Endpoint**: `/api/v1/storage/presigned/{objectName}`
- **Query Parameters**:
  - `method` (optional): GET or PUT (default: GET)
  - `expiry` (optional): Duration like "1h", "30m", "24h" (default: 1h)

**Example using curl:**
```bash
curl -X GET "http://localhost:8080/api/v1/storage/presigned/my-document.pdf?method=GET&expiry=2h"
```

## Example Response Format

### Success Response
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "data": {
    "object_name": "my-document.pdf",
    "size": 1024000,
    "content_type": "application/pdf"
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Failed to upload file",
  "error": "detailed error message"
}
```

## Features

1. **Automatic Bucket Creation**: The system automatically creates the configured bucket if it doesn't exist
2. **File Upload**: Upload files via multipart form data
3. **File Download**: Download files with proper headers and streaming
4. **File Deletion**: Remove files from storage
5. **File Listing**: List all files with optional prefix filtering
6. **Presigned URLs**: Generate time-limited direct access URLs
7. **File Metadata**: Get file information including size, content type, and modification date

## Error Handling
The API includes comprehensive error handling for:
- MinIO connection issues
- File not found errors
- Invalid parameters
- Upload/download failures
- Authorization issues

## Security Notes
- Store MinIO credentials securely using environment variables
- Use HTTPS in production (set `MINIO_USE_SSL=true`)
- Implement proper authentication/authorization for file operations
- Consider adding file size limits and type restrictions
- Use presigned URLs for direct client uploads to reduce server load