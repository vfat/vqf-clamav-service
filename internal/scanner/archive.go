package scanner

import (
	"archive/zip"
	"fmt"
	"io"
)

// ArchiveLimits defines boundaries for inspecting archives to prevent DoS attacks.
type ArchiveLimits struct {
	MaxRecursion int
	MaxFiles     int
	MaxExtractMB int64
}

// ArchiveResult contains metadata and security verdict from archive inspection.
type ArchiveResult struct {
	IsBomb                 bool   `json:"is_bomb"`
	BombReason             string `json:"bomb_reason,omitempty"`
	IsEncrypted            bool   `json:"is_encrypted"`
	TotalFiles             int    `json:"total_files"`
	TotalUncompressedBytes int64  `json:"total_uncompressed_bytes"`
}

// ArchiveInspector analyzes zip and compressed archives for zip bombs and encryption.
type ArchiveInspector struct {
	limits ArchiveLimits
}

// NewArchiveInspector initializes the archive security inspector.
func NewArchiveInspector(limits ArchiveLimits) *ArchiveInspector {
	if limits.MaxRecursion <= 0 {
		limits.MaxRecursion = 5
	}
	if limits.MaxFiles <= 0 {
		limits.MaxFiles = 1000
	}
	if limits.MaxExtractMB <= 0 {
		limits.MaxExtractMB = 250
	}

	return &ArchiveInspector{limits: limits}
}

// InspectZip checks a zip archive against file limits, extraction size limits, and password encryption.
func (ai *ArchiveInspector) InspectZip(r io.ReaderAt, size int64, password string) (*ArchiveResult, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, fmt.Errorf("failed to parse zip archive: %w", err)
	}

	res := &ArchiveResult{
		TotalFiles: len(zr.File),
	}

	// 1. Check Max Files Limit
	if len(zr.File) > ai.limits.MaxFiles {
		res.IsBomb = true
		res.BombReason = "EXCEEDED_MAX_FILES"
		return res, nil
	}

	maxBytes := ai.limits.MaxExtractMB * 1024 * 1024
	var totalUncompressed int64

	for _, f := range zr.File {
		// Check for encryption (Bit 0 of General Purpose Bit Flag indicates encryption)
		if f.Flags&0x1 != 0 {
			res.IsEncrypted = true
		}

		totalUncompressed += int64(f.UncompressedSize64)

		// 2. Check Total Uncompressed Size Limit
		if totalUncompressed > maxBytes {
			res.IsBomb = true
			res.BombReason = "EXCEEDED_MAX_SIZE"
			res.TotalUncompressedBytes = totalUncompressed
			return res, nil
		}

		// 3. Check suspicious compression ratio (e.g. > 100x ratio)
		if f.CompressedSize64 > 0 {
			ratio := f.UncompressedSize64 / f.CompressedSize64
			if ratio > 100 && f.UncompressedSize64 > 10*1024*1024 {
				res.IsBomb = true
				res.BombReason = "SUSPICIOUS_COMPRESSION_RATIO"
				return res, nil
			}
		}
	}

	res.TotalUncompressedBytes = totalUncompressed
	return res, nil
}
