package models

type BackupArchiveRecord struct {
	ID           string
	UserID       string
	IncludeMedia bool
	SizeBytes    int64
	EntityCounts string
	ArchivePath  *string
	CreatedAt    string
	UpdatedAt    string
}

type BackupArchive struct {
	ID           string         `json:"id"`
	UserID       string         `json:"userId"`
	IncludeMedia bool           `json:"includeMedia"`
	SizeBytes    int64          `json:"sizeBytes"`
	EntityCounts map[string]int `json:"entityCounts"`
	ArchivePath  *string        `json:"archivePath,omitempty"`
	CreatedAt    string         `json:"createdAt"`
	UpdatedAt    string         `json:"updatedAt"`
}

type CreateBackupInput struct {
	IncludeMedia bool `json:"includeMedia"`
}

type BackupExportResult struct {
	DownloadID string `json:"downloadId"`
	SizeBytes  int64  `json:"sizeBytes"`
	Timestamp  string `json:"timestamp"`
}

type BackupImportSummary struct {
	SchemaVersion int            `json:"schemaVersion"`
	AppVersion    string         `json:"appVersion"`
	EntityCounts  map[string]int `json:"entityCounts"`
	MediaCount    int            `json:"mediaCount"`
	Warnings      []string       `json:"warnings"`
}

type BackupImportConfirmResult struct {
	Status       string         `json:"status"`
	EntityCounts map[string]int `json:"entityCounts"`
	MediaCount   int            `json:"mediaCount"`
}

type ImportValidationState struct {
	UserID  string
	Summary BackupImportSummary
	ZipData []byte
}

func BackupArchiveFromRecord(r *BackupArchiveRecord) *BackupArchive {
	if r == nil {
		return nil
	}
	var path *string
	if r.ArchivePath != nil {
		path = r.ArchivePath
	}
	return &BackupArchive{
		ID:           r.ID,
		UserID:       r.UserID,
		IncludeMedia: r.IncludeMedia,
		SizeBytes:    r.SizeBytes,
		EntityCounts: ParseEntityCounts(r.EntityCounts),
		ArchivePath:  path,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func ParseEntityCounts(raw string) map[string]int {
	return map[string]int{}
}
