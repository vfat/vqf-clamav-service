package scanner

import (
	"archive/zip"
	"bytes"
	"testing"
)

func createTestZip(files map[string][]byte) (*bytes.Reader, int64) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	for name, content := range files {
		f, _ := w.Create(name)
		f.Write(content)
	}
	w.Close()

	return bytes.NewReader(buf.Bytes()), int64(buf.Len())
}

func TestArchiveInspector_CleanZip(t *testing.T) {
	reader, size := createTestZip(map[string][]byte{
		"readme.txt": []byte("Hello World"),
		"data.csv":   []byte("id,name\n1,Alice\n2,Bob"),
	})

	inspector := NewArchiveInspector(ArchiveLimits{
		MaxRecursion: 5,
		MaxFiles:     1000,
		MaxExtractMB: 250,
	})

	res, err := inspector.InspectZip(reader, size, "")
	if err != nil {
		t.Fatalf("InspectZip failed on clean archive: %v", err)
	}

	if res.IsBomb {
		t.Errorf("expected clean zip, got IsBomb=true")
	}
	if res.TotalFiles != 2 {
		t.Errorf("expected 2 files, got %d", res.TotalFiles)
	}
	if res.IsEncrypted {
		t.Errorf("expected IsEncrypted=false")
	}
}

func TestArchiveInspector_ZipBombFilesLimit(t *testing.T) {
	// Create zip with 10 files, limit max files to 5
	files := make(map[string][]byte)
	for i := 0; i < 10; i++ {
		files[string(rune('a'+i))+".txt"] = []byte("content")
	}
	reader, size := createTestZip(files)

	inspector := NewArchiveInspector(ArchiveLimits{
		MaxRecursion: 5,
		MaxFiles:     5, // Limit 5
		MaxExtractMB: 250,
	})

	res, err := inspector.InspectZip(reader, size, "")
	if err != nil {
		t.Fatalf("InspectZip failed: %v", err)
	}

	if !res.IsBomb {
		t.Errorf("expected archive to be flagged as Zip Bomb (exceeded max files 5)")
	}
	if res.BombReason != "EXCEEDED_MAX_FILES" {
		t.Errorf("expected BombReason 'EXCEEDED_MAX_FILES', got '%s'", res.BombReason)
	}
}

func TestArchiveInspector_ZipBombSizeLimit(t *testing.T) {
	// Create file with 2 MB uncompressed, limit to 1 MB
	bigContent := bytes.Repeat([]byte("A"), 2*1024*1024)
	reader, size := createTestZip(map[string][]byte{
		"big.dat": bigContent,
	})

	inspector := NewArchiveInspector(ArchiveLimits{
		MaxRecursion: 5,
		MaxFiles:     1000,
		MaxExtractMB: 1, // Limit 1 MB
	})

	res, err := inspector.InspectZip(reader, size, "")
	if err != nil {
		t.Fatalf("InspectZip failed: %v", err)
	}

	if !res.IsBomb {
		t.Errorf("expected archive to be flagged as Zip Bomb (exceeded max extract size)")
	}
	if res.BombReason != "EXCEEDED_MAX_SIZE" {
		t.Errorf("expected BombReason 'EXCEEDED_MAX_SIZE', got '%s'", res.BombReason)
	}
}
