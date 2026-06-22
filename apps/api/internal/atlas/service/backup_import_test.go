// FILE: apps/api/internal/atlas/service/backup_import_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for BackupImportService covering Validate and Confirm with ZIP structure, manifest, schema, data.json, media count, warnings, expiration, and user mismatch.
//   SCOPE: Success paths (no media, with media), invalid ZIP, missing manifest, invalid manifest, schema mismatch, missing data, Confirm success, expired validation, user mismatch.
//   DEPENDS: apps/api/internal/atlas/service, apps/api/internal/atlas/models.
//   LINKS: M-API / V-M-API / WAVE-09.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT

package service_test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/api/internal/atlas/models"
	"monorepo-template/apps/api/internal/atlas/service"
)

var backupCtx = context.Background()
var backupUserID = "550e8400-e29b-41d4-a716-446655440000"

func buildTestBackupZIP(t *testing.T, manifestType string, schemaVersion int, sections []string, entityData map[string][]string, mediaCount int) []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	manifest := models.BackupManifest{
		Type:             manifestType,
		SchemaVersion:    schemaVersion,
		AppVersion:       models.BackupAppVersion,
		ExportedAt:       "2026-06-22T00:00:00Z",
		IncludedSections: sections,
		MediaIncluded:    mediaCount > 0,
	}
	mData, err := json.Marshal(manifest)
	require.NoError(t, err)

	mf, err := zw.Create("manifest.json")
	require.NoError(t, err)
	_, err = mf.Write(mData)
	require.NoError(t, err)

	dataBytes, err := json.Marshal(entityData)
	require.NoError(t, err)
	df, err := zw.Create("data.json")
	require.NoError(t, err)
	_, err = df.Write(dataBytes)
	require.NoError(t, err)

	for i := 0; i < mediaCount; i++ {
		mf, err := zw.Create(fmt.Sprintf("media/file_%d", i))
		require.NoError(t, err)
		_, err = mf.Write([]byte("test media content"))
		require.NoError(t, err)
	}

	err = zw.Close()
	require.NoError(t, err)

	return buf.Bytes()
}

func TestBackupImportService_Validate_Success(t *testing.T) {
	svc := service.NewBackupImportService()

	entityData := map[string][]string{
		"exercises":  {"ex1", "ex2"},
		"dailyLogs": {"log1"},
	}
	zipData := buildTestBackupZIP(t, "full_backup", models.BackupSchemaVersion, []string{"exercises", "dailyLogs"}, entityData, 0)

	validationID, summary, err := svc.Validate(backupCtx, backupUserID, zipData)
	require.NoError(t, err)
	require.NotEmpty(t, validationID)
	require.NotNil(t, summary)
	assert.Equal(t, models.BackupSchemaVersion, summary.SchemaVersion)
	assert.Equal(t, models.BackupAppVersion, summary.AppVersion)
	assert.Equal(t, map[string]int{"exercises": 2, "dailyLogs": 1}, summary.EntityCounts)
	assert.Equal(t, 0, summary.MediaCount)
	assert.Empty(t, summary.Warnings)
}

func TestBackupImportService_Validate_WithMedia(t *testing.T) {
	svc := service.NewBackupImportService()

	entityData := map[string][]string{
		"exercises": {"ex1", "ex2", "ex3"},
	}
	zipData := buildTestBackupZIP(t, "full_backup", models.BackupSchemaVersion, []string{"exercises"}, entityData, 5)

	validationID, summary, err := svc.Validate(backupCtx, backupUserID, zipData)
	require.NoError(t, err)
	require.NotEmpty(t, validationID)
	require.NotNil(t, summary)
	assert.Equal(t, models.BackupSchemaVersion, summary.SchemaVersion)
	assert.Equal(t, models.BackupAppVersion, summary.AppVersion)
	assert.Equal(t, map[string]int{"exercises": 3}, summary.EntityCounts)
	assert.Equal(t, 5, summary.MediaCount)
	assert.Empty(t, summary.Warnings)
}

func TestBackupImportService_Validate_InvalidZIP(t *testing.T) {
	svc := service.NewBackupImportService()

	validationID, summary, err := svc.Validate(backupCtx, backupUserID, []byte("not a valid zip file"))
	assert.ErrorIs(t, err, service.ErrBackupInvalidZIP)
	assert.Empty(t, validationID)
	assert.Nil(t, summary)
}

func TestBackupImportService_Validate_MissingManifest(t *testing.T) {
	svc := service.NewBackupImportService()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	df, err := zw.Create("data.json")
	require.NoError(t, err)
	_, err = df.Write([]byte(`{}`))
	require.NoError(t, err)
	err = zw.Close()
	require.NoError(t, err)

	validationID, summary, err := svc.Validate(backupCtx, backupUserID, buf.Bytes())
	assert.ErrorIs(t, err, service.ErrBackupManifestMissing)
	assert.Empty(t, validationID)
	assert.Nil(t, summary)
}

func TestBackupImportService_Validate_InvalidManifest(t *testing.T) {
	svc := service.NewBackupImportService()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	mf, err := zw.Create("manifest.json")
	require.NoError(t, err)
	_, err = mf.Write([]byte(`not valid json`))
	require.NoError(t, err)
	err = zw.Close()
	require.NoError(t, err)

	validationID, summary, err := svc.Validate(backupCtx, backupUserID, buf.Bytes())
	assert.ErrorIs(t, err, service.ErrBackupManifestInvalid)
	assert.Empty(t, validationID)
	assert.Nil(t, summary)
}

func TestBackupImportService_Validate_SchemaMismatch(t *testing.T) {
	svc := service.NewBackupImportService()

	zipData := buildTestBackupZIP(t, "full_backup", 999, []string{"exercises"}, map[string][]string{"exercises": {"ex1"}}, 0)

	validationID, summary, err := svc.Validate(backupCtx, backupUserID, zipData)
	assert.ErrorIs(t, err, service.ErrBackupSchemaMismatch)
	assert.Empty(t, validationID)
	assert.Nil(t, summary)
}

func TestBackupImportService_Validate_MissingData(t *testing.T) {
	svc := service.NewBackupImportService()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	mf, err := zw.Create("manifest.json")
	require.NoError(t, err)
	_, err = mf.Write([]byte(`{"type":"full_backup","schemaVersion":1,"appVersion":"1.0.0","exportedAt":"2026-06-22T00:00:00Z","includedSections":[],"mediaIncluded":false}`))
	require.NoError(t, err)
	err = zw.Close()
	require.NoError(t, err)

	validationID, summary, err := svc.Validate(backupCtx, backupUserID, buf.Bytes())
	assert.ErrorIs(t, err, service.ErrBackupDataMissing)
	assert.Empty(t, validationID)
	assert.Nil(t, summary)
}

func TestBackupImportService_Confirm_Success(t *testing.T) {
	svc := service.NewBackupImportService()

	entityData := map[string][]string{
		"exercises":  {"ex1", "ex2"},
		"dailyLogs": {"log1"},
	}
	zipData := buildTestBackupZIP(t, "full_backup", models.BackupSchemaVersion, []string{"exercises", "dailyLogs"}, entityData, 3)

	validationID, summary, err := svc.Validate(backupCtx, backupUserID, zipData)
	require.NoError(t, err)
	require.NotEmpty(t, validationID)
	require.NotNil(t, summary)
	assert.Equal(t, 3, summary.MediaCount)

	result, err := svc.Confirm(backupCtx, backupUserID, validationID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "confirmed", result.Status)
	assert.Equal(t, map[string]int{"exercises": 2, "dailyLogs": 1}, result.EntityCounts)
	assert.Equal(t, 3, result.MediaCount)
}

func TestBackupImportService_Confirm_Expired(t *testing.T) {
	svc := service.NewBackupImportService()

	result, err := svc.Confirm(backupCtx, backupUserID, "550e8400-e29b-41d4-a716-446655449999")
	assert.ErrorIs(t, err, service.ErrBackupValidationExpired)
	assert.Nil(t, result)
}

func TestBackupImportService_Confirm_UserMismatch(t *testing.T) {
	svc := service.NewBackupImportService()

	zipData := buildTestBackupZIP(t, "full_backup", models.BackupSchemaVersion, []string{"exercises"}, map[string][]string{"exercises": {"ex1"}}, 0)

	validationID, _, err := svc.Validate(backupCtx, backupUserID, zipData)
	require.NoError(t, err)
	require.NotEmpty(t, validationID)

	result, err := svc.Confirm(backupCtx, "different-user-id", validationID)
	assert.ErrorIs(t, err, service.ErrBackupValidationUserMismatch)
	assert.Nil(t, result)
}
