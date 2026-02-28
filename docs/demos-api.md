# Demo API Documentation

## Overview

The Demo API is a comprehensive system for uploading, storing, and retrieving game demo files with automatic compression. It's designed to handle both tournament and casual game demos with metadata management and efficient file storage using compression.

## Architecture

The Demo API follows a layered architecture pattern:

```
HTTP Layer (Gin Router)
    ↓
Handler Layer (demo.go)
    ↓
Service Layer (demo_service.go)
    ↓
Repository Layer (demo_repository.go)
    ↓
Storage Layer (MinIO + PostgreSQL)
```

### Components

1. **Handler Layer**: HTTP request/response handling and validation
2. **Service Layer**: Business logic, file compression, and coordination
3. **Repository Layer**: Database operations using GORM
4. **Storage Layer**: File storage using MinIO object storage with PostgreSQL metadata

## Data Model

### Demo Entity

```go
type Demo struct {
    ID             int       `json:"id"`
    IsTournament   bool      `json:"is_tournament"`
    TournamentID   *int      `json:"tournament_id,omitempty"`
    MatchID        *int      `json:"match_id,omitempty"`
    RoundID        *int      `json:"round_id,omitempty"`
    RawDemoName    string    `json:"raw_demo_name"`
    FileSize       int64     `json:"file_size"`
    BlueTeamName   string    `json:"blue_team_name"`
    RedTeamName    string    `json:"red_team_name"`
    ObjectName     string    `json:"object_name"`
    ContentType    string    `json:"content_type"`
    IsCompressed   bool      `json:"is_compressed"`
    CompressedSize *int64    `json:"compressed_size,omitempty"`
    UploadedAt     time.Time `json:"uploaded_at"`
}
```

### Database Schema

```sql
CREATE TABLE demos (
    id SERIAL PRIMARY KEY,
    is_tournament BOOLEAN NOT NULL DEFAULT FALSE,
    tournament_id INTEGER,
    match_id INTEGER,
    round_id INTEGER,
    raw_demo_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    blue_team_name VARCHAR(100) NOT NULL,
    red_team_name VARCHAR(100) NOT NULL,
    object_name VARCHAR(500) NOT NULL UNIQUE,
    content_type VARCHAR(100),
    is_compressed BOOLEAN DEFAULT TRUE,
    compressed_size BIGINT,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_tournament_data CHECK (
        (is_tournament = FALSE) OR 
        (is_tournament = TRUE AND tournament_id IS NOT NULL AND match_id IS NOT NULL AND round_id IS NOT NULL)
    )
);
```

### Indexing Strategy

The database includes optimized indexes for common query patterns:
- `idx_demos_match_round`: For match and round lookups
- `idx_demos_tournament`: For tournament-specific queries
- `idx_demos_tournament_match_round`: For combined tournament queries
- `idx_demos_uploaded_at`: For chronological sorting
- `idx_demos_object_name`: For file existence checks

## API Endpoints

Base path: `/api/v1/demos`

### 1. Upload Demo
**POST** `/upload`

Upload a new demo file with automatic compression.

**Authentication**: Requires `secret_password` query parameter

**Request**: `multipart/form-data`
```
demo_file: [file]                    // Required: Demo file to upload
is_tournament: boolean               // Required: Whether this is a tournament demo
tournament_id: integer              // Required if is_tournament=true
match_id: integer                   // Required if is_tournament=true
round_id: integer                   // Required if is_tournament=true
raw_demo_name: string              // Required: Original demo name
blue_team_name: string             // Required: Blue team name
red_team_name: string              // Required: Red team name
```

**Response**:
```json
{
    "status": "success",
    "message": "Demo uploaded successfully",
    "data": {
        "id": 1,
        "is_tournament": true,
        "tournament_id": 1,
        "match_id": 1,
        "round_id": 1,
        "raw_demo_name": "demo.dem",
        "file_size": 1024000,
        "blue_team_name": "Team Alpha",
        "red_team_name": "Team Beta",
        "object_name": "demos/tournament_1/match_1/round_1/demo_1709123456.dem.zst",
        "content_type": "application/octet-stream",
        "is_compressed": true,
        "compressed_size": 512000,
        "uploaded_at": "2026-02-28T10:00:00Z"
    }
}
```

### 2. Get Demo Information
**GET** `/`

Retrieve demo information by query parameters.

**Request Parameters**:
```
match_id: integer       // Optional: Match ID
round_id: integer      // Optional: Round ID  
tournament_id: integer // Optional: Tournament ID
```

**Query Priority**:
1. tournament_id + match_id + round_id (exact tournament match)
2. match_id + round_id (match without tournament context)

**Response**:
```json
{
    "status": "success",
    "message": "Demo found",
    "data": {
        // Demo object (same as upload response)
    }
}
```

### 3. Download Demo File
**GET** `/download/:id`

Download the actual demo file by demo ID. The file is automatically decompressed before streaming to the client.

**Response**: Binary file stream with headers:
```
Content-Type: application/octet-stream
Content-Disposition: attachment; filename="demo.dem" 
Content-Length: [original_file_size]
```

**Note**: Files are stored compressed using zstd in the storage layer, but are automatically decompressed when downloaded, so clients receive the original uncompressed demo file.

### 4. List Demos  
**GET** `/list`

Retrieve paginated list of demos.

**Request Parameters**:
```
page: integer   // Optional: Page number (default: 1)
limit: integer  // Optional: Items per page (default: 10, max: 100)
```

**Response**:
```json
{
    "status": "success", 
    "message": "Demos retrieved successfully",
    "data": {
        "demos": [/* Array of demo objects */],
        "page": 1,
        "limit": 10,
        "count": 10
    }
}
```

### 5. Check Demo Availability
**HEAD** `/`

Check if a demo exists without returning data (same parameters as GET).

## File Storage & Compression

### Storage Strategy

**Object Storage**: MinIO is used for storing compressed demo files
- Bucket-based organization  
- Automatic file compression using Zstandard (zstd)
- Unique object naming prevents conflicts

### Object Naming Convention

**Tournament Demos**:
```
demos/tournament_{tournament_id}/match_{match_id}/round_{round_id}/{filename}_{timestamp}.{ext}.zst
```

**Casual Demos**:
```
demos/casual/{filename}_{timestamp}.{ext}.zst
```

### Storage Details

- **Algorithm**: Zstandard (zstd) with default speed level
- **Process**: Files compressed automatically during upload
- **Storage**: Only compressed versions stored in MinIO for efficiency
- **Download**: Files automatically decompressed when downloaded to clients
- **Metadata**: Both original and compressed sizes tracked
- **Format**: Files stored with `.zst` suffix, served with original names

## Authentication & Security

### Upload Security
- **Secret Password**: Required for upload operations via query parameter
- **Middleware**: `SecretPasswordAuth` middleware validates credentials
- **Public Access**: Read operations (GET, HEAD, LIST, DOWNLOAD) are public

### File Validation
- **Duplicate Prevention**: Object name uniqueness enforced at database level
- **Tournament Validation**: Tournament demos require tournament_id, match_id, and round_id
- **Size Tracking**: Both original and compressed sizes monitored

## Error Handling

### Common HTTP Status Codes

- **200 OK**: Successful operation
- **400 Bad Request**: Invalid request data or missing required fields
- **401 Unauthorized**: Missing or invalid secret_password for uploads
- **404 Not Found**: Demo not found or doesn't exist
- **500 Internal Server Error**: Server-side errors (database, storage issues)

### Error Response Format

```json
{
    "status": "error",
    "message": "Error description"
}
```

## Usage Examples

### Upload Tournament Demo

```bash
curl -X POST "https://api.udl.tf/v1/demos/upload?secret_password=your_secret" \
  -F "demo_file=@match1_round1.dem" \
  -F "is_tournament=true" \
  -F "tournament_id=1" \
  -F "match_id=1" \
  -F "round_id=1" \
  -F "raw_demo_name=match1_round1.dem" \
  -F "blue_team_name=Team Alpha" \
  -F "red_team_name=Team Beta"
```

### Get Demo by Tournament Match

```bash
curl "https://api.udl.tf/v1/demos?tournament_id=1&match_id=1&round_id=1"
```

### Download Demo File

```bash
curl -O "https://api.udl.tf/v1/demos/download/1"
# Downloads the original demo file (automatically decompressed)
# File saved as: original_name.dem
```

### List Recent Demos

```bash
curl "https://api.udl.tf/v1/demos/list?page=1&limit=20"
```

## Implementation Details

### Service Layer Responsibilities

1. **File Compression**: Automatic zstd compression during upload
2. **File Decompression**: Automatic decompression during download 
3. **Object Naming**: Generates unique object names with timestamps
4. **Validation**: Ensures tournament data consistency
5. **Storage Coordination**: Handles both database and object storage operations
6. **Cleanup**: Removes storage files if database operations fail

### Repository Layer Operations

- **Create**: Insert new demo records with validation
- **Read**: Query by ID, tournament context, or match context  
- **List**: Paginated demo retrieval with sorting
- **Check**: Object name existence verification
- **Delete**: Record removal (used by service for cleanup)

### Concurrency & Consistency

- **Unique Constraints**: Database enforces object name uniqueness
- **Transactional Safety**: Database operations fail fast if storage succeeds but DB fails
- **Cleanup Strategy**: Orphaned storage files cleaned up on DB save failure

## Performance Considerations

### Database Optimization
- Strategic indexing for common query patterns
- Offset-based pagination for large datasets
- Timestamp-based ordering for chronological access

### Storage Optimization
- Zstd compression reduces storage size and transfer time
- Object name structure enables efficient bucket organization
- Presigned URL capability (via storage service) for direct client access

### Scalability Features
- Stateless service design enables horizontal scaling
- Separate storage and database layers allow independent scaling
- Consistent object naming supports load balancing and caching
