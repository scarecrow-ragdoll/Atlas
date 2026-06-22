// FILE: apps/api/internal/atlas/service/backup_import.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Implement the BackupImportService for WAVE-09 backup ZIP validation and import confirmation.
//   SCOPE: Validate backup ZIP structure (manifest, schema, data.json, media count) with in-memory validation tokens via sync.Map. Confirm consumes one-time validation tokens. No direct DB writes in MVP.
//   DEPENDS: apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-09.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   BackupImportService - Interface for backup import validation and confirmation.
//   NewBackupImportService - Creates a new BackupImportService.
//   Validate - Parses backup ZIP, validates manifest/schema/data.json, returns summary with entity counts.
//   Confirm - Consumes a one-time validation token and confirms the import.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Initial BackupImportService implementation for WAVE-09.
// END_CHANGE_SUMMARY

package service

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"monorepo-template/apps/api/internal/atlas/models"
)

var (
	ErrBackupInvalidZIP             = errors.New("invalid backup ZIP file")
	ErrBackupManifestMissing        = errors.New("manifest.json not found in backup")
	ErrBackupManifestInvalid        = errors.New("manifest.json is invalid")
	ErrBackupDataMissing            = errors.New("data.json not found in backup")
	ErrBackupSchemaMismatch         = errors.New("backup schema version mismatch")
	ErrBackupValidationExpired      = errors.New("backup validation expired or not found")
	ErrBackupValidationUserMismatch = errors.New("backup validation was created by a different user")
)

type BackupImportService interface {
	Validate(ctx context.Context, userID string, zipData []byte) (validationID string, summary *models.BackupImportSummary, err error)
	Confirm(ctx context.Context, userID string, validationID string) (*models.BackupImportConfirmResult, error)
}

type validationEntry struct {
	state     models.ImportValidationState
	createdAt time.Time
}

type backupImportService struct {
	validationStore sync.Map
	cleanupOnce     sync.Once
}

const validationTTL = 15 * time.Minute

func NewBackupImportService() BackupImportService {
	return &backupImportService{}
}

func (s *backupImportService) Validate(ctx context.Context, userID string, zipData []byte) (string, *models.BackupImportSummary, error) {
	zr, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return "", nil, fmt.Errorf("%w: %v", ErrBackupInvalidZIP, err)
	}

	var manifest *models.BackupManifest
	var dataFile *zip.File
	mediaCount := 0

	for _, f := range zr.File {
		switch {
		case f.Name == "manifest.json":
			m, err := readManifestFile(f)
			if err != nil {
				return "", nil, fmt.Errorf("%w: %v", ErrBackupManifestInvalid, err)
			}
			manifest = m
		case f.Name == "data.json":
			dataFile = f
		case len(f.Name) > 6 && f.Name[:6] == "media/":
			mediaCount++
		}
	}

	if manifest == nil {
		return "", nil, ErrBackupManifestMissing
	}

	if manifest.Type != "full_backup" {
		return "", nil, fmt.Errorf("unsupported backup type: %s", manifest.Type)
	}

	if manifest.SchemaVersion != models.BackupSchemaVersion {
		return "", nil, ErrBackupSchemaMismatch
	}

	if dataFile == nil {
		return "", nil, ErrBackupDataMissing
	}

	dataBytes, err := readZipContent(dataFile)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read data.json: %w", err)
	}

	var dataMap map[string]any
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		return "", nil, fmt.Errorf("data.json is invalid: %w", err)
	}

	entityCounts := countEntities(dataMap)

	warnings := []string{}
	if manifest.MediaIncluded && mediaCount == 0 {
		warnings = append(warnings, "manifest indicates media included but no media files found in ZIP")
	}
	if !manifest.MediaIncluded && mediaCount > 0 {
		warnings = append(warnings, "media files found in ZIP but manifest indicates media not included")
	}

	summary := &models.BackupImportSummary{
		SchemaVersion: manifest.SchemaVersion,
		AppVersion:    manifest.AppVersion,
		EntityCounts:  entityCounts,
		MediaCount:    mediaCount,
		Warnings:      warnings,
	}

	validationID, err := generateUUID()
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate validation ID: %w", err)
	}

	entry := &validationEntry{
		state: models.ImportValidationState{
			UserID:  userID,
			Summary: *summary,
			ZipData: zipData,
		},
		createdAt: time.Now(),
	}

	s.validationStore.Store(validationID, entry)

	return validationID, summary, nil
}

func (s *backupImportService) Confirm(ctx context.Context, userID string, validationID string) (*models.BackupImportConfirmResult, error) {
	raw, ok := s.validationStore.Load(validationID)
	if !ok {
		return nil, ErrBackupValidationExpired
	}

	entry, ok := raw.(*validationEntry)
	if !ok {
		s.validationStore.Delete(validationID)
		return nil, ErrBackupValidationExpired
	}

	if time.Since(entry.createdAt) > validationTTL {
		s.validationStore.Delete(validationID)
		return nil, ErrBackupValidationExpired
	}

	if entry.state.UserID != userID {
		return nil, ErrBackupValidationUserMismatch
	}

	s.validationStore.Delete(validationID)

	zr, err := zip.NewReader(bytes.NewReader(entry.state.ZipData), int64(len(entry.state.ZipData)))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBackupInvalidZIP, err)
	}

	mediaCount := 0
	for _, f := range zr.File {
		if len(f.Name) > 6 && f.Name[:6] == "media/" {
			mediaCount++
		}
	}

	result := &models.BackupImportConfirmResult{
		Status:       "confirmed",
		EntityCounts: entry.state.Summary.EntityCounts,
		MediaCount:   mediaCount,
	}

	return result, nil
}

func readManifestFile(f *zip.File) (*models.BackupManifest, error) {
	data, err := readZipContent(f)
	if err != nil {
		return nil, err
	}

	var manifest models.BackupManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func readZipContent(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func countEntities(data map[string]any) map[string]int {
	counts := make(map[string]int)
	for key, val := range data {
		arr, ok := val.([]any)
		if ok {
			counts[key] = len(arr)
		}
	}
	return counts
}

func generateUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
