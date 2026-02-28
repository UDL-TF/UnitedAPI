package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/UDL-TF/UnitedAPI/internal/model"
	"github.com/UDL-TF/UnitedAPI/internal/repository"
	"github.com/UDL-TF/UnitedAPI/internal/storage"
	"github.com/klauspost/compress/zstd"
)

// DemoService handles demo business logic
type DemoService struct {
	demoRepo    *repository.DemoRepository
	minioClient *storage.MinIOClient
}

// NewDemoService creates a new demo service
func NewDemoService(demoRepo *repository.DemoRepository, minioClient *storage.MinIOClient) *DemoService {
	return &DemoService{
		demoRepo:    demoRepo,
		minioClient: minioClient,
	}
}

// UploadDemo handles demo file upload with compression
func (s *DemoService) UploadDemo(req *model.DemoUploadRequest, file multipart.File, header *multipart.FileHeader) (*model.Demo, error) {
	// Validate tournament data if needed
	if err := req.ValidateTournamentData(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Generate unique object name
	objectName := s.generateObjectName(req, header.Filename)

	// Check if object already exists
	exists, err := s.demoRepo.ExistsByObjectName(objectName)
	if err != nil {
		return nil, fmt.Errorf("failed to check object existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("demo with object name %s already exists", objectName)
	}

	// Compress the file
	compressedReader, compressedSize, err := s.compressFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to compress file: %w", err)
	}

	// Upload compressed file to MinIO
	contentType := "application/zstd"
	if err := s.minioClient.UploadFile(objectName, compressedReader, compressedSize, contentType); err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	// Create demo record
	demo := &model.Demo{
		IsTournament:   req.IsTournament,
		TournamentID:   req.TournamentID,
		MatchID:        req.MatchID,
		RoundID:        req.RoundID,
		RawDemoName:    req.RawDemoName,
		FileSize:       header.Size,
		BlueTeamName:   req.BlueTeamName,
		RedTeamName:    req.RedTeamName,
		ObjectName:     objectName,
		ContentType:    contentType,
		IsCompressed:   true,
		CompressedSize: &compressedSize,
	}

	// Save to database
	if err := s.demoRepo.Create(demo); err != nil {
		// Try to cleanup uploaded file if database save fails
		s.minioClient.DeleteFile(objectName)
		return nil, fmt.Errorf("failed to save demo record: %w", err)
	}

	return demo, nil
}

// GetDemo retrieves demo information by various criteria
func (s *DemoService) GetDemo(req *model.DemoGetRequest) (*model.Demo, error) {
	// Priority: Tournament + Match + Round → Match + Round
	if req.TournamentID != nil && req.MatchID != nil && req.RoundID != nil {
		return s.demoRepo.GetByTournamentMatchRound(*req.TournamentID, *req.MatchID, *req.RoundID)
	}

	if req.MatchID != nil && req.RoundID != nil {
		return s.demoRepo.GetByMatchAndRound(*req.MatchID, *req.RoundID)
	}

	return nil, fmt.Errorf("insufficient parameters: need either (match_id + round_id) or (tournament_id + match_id + round_id)")
}

// DownloadDemo gets the demo file from storage and decompresses it
func (s *DemoService) DownloadDemo(demoID int) (io.ReadCloser, *model.Demo, error) {
	// Get demo record
	demo, err := s.demoRepo.GetByID(demoID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get demo: %w", err)
	}

	// Download from MinIO
	object, err := s.minioClient.DownloadFile(demo.ObjectName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download file from storage: %w", err)
	}
	defer object.Close()

	// Decompress the file
	decompressedReader, err := s.decompressFile(object)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decompress file: %w", err)
	}

	return decompressedReader, demo, nil
}

// ListDemos retrieves a list of demos with pagination
func (s *DemoService) ListDemos(page, limit int) ([]*model.Demo, error) {
	offset := (page - 1) * limit
	return s.demoRepo.List(offset, limit)
}

// DeleteDemo removes a demo and its file
func (s *DemoService) DeleteDemo(demoID int) error {
	// Get demo record first
	demo, err := s.demoRepo.GetByID(demoID)
	if err != nil {
		return fmt.Errorf("failed to get demo: %w", err)
	}

	// Delete from storage
	if err := s.minioClient.DeleteFile(demo.ObjectName); err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Delete from database
	if err := s.demoRepo.Delete(demoID); err != nil {
		return fmt.Errorf("failed to delete demo record: %w", err)
	}

	return nil
}

// compressFile compresses the file using zstd
func (s *DemoService) compressFile(file multipart.File) (io.Reader, int64, error) {
	// Reset file pointer to beginning
	if _, err := file.Seek(0, 0); err != nil {
		return nil, 0, fmt.Errorf("failed to seek file: %w", err)
	}

	// Create zstd encoder
	encoder, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedDefault))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create zstd encoder: %w", err)
	}

	// Read file content
	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read file: %w", err)
	}

	// Compress the data
	compressedData := encoder.EncodeAll(fileData, make([]byte, 0, len(fileData)))
	encoder.Close()

	// Return compressed data as reader
	compressedReader := strings.NewReader(string(compressedData))
	return compressedReader, int64(len(compressedData)), nil
}

// decompressFile decompresses a zstd compressed file and returns a reader
func (s *DemoService) decompressFile(compressedReader io.Reader) (io.ReadCloser, error) {
	// Create zstd decoder
	decoder, err := zstd.NewReader(compressedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create zstd decoder: %w", err)
	}

	// Read and decompress all data
	decompressedData, err := io.ReadAll(decoder)
	decoder.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to decompress data: %w", err)
	}

	// Return decompressed data as ReadCloser
	return io.NopCloser(strings.NewReader(string(decompressedData))), nil
}

// generateObjectName creates a unique object name for the demo
func (s *DemoService) generateObjectName(req *model.DemoUploadRequest, originalFilename string) string {
	timestamp := time.Now().Unix()
	ext := filepath.Ext(originalFilename)

	if req.IsTournament {
		return fmt.Sprintf("demos/tournament_%d/match_%d/round_%d/%s_%d%s.zst",
			*req.TournamentID, *req.MatchID, *req.RoundID,
			strings.TrimSuffix(req.RawDemoName, ext), timestamp, ext)
	}

	return fmt.Sprintf("demos/casual/%s_%d%s.zst",
		strings.TrimSuffix(req.RawDemoName, ext), timestamp, ext)
}
