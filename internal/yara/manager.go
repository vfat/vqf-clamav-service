package yara

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vfat/vqf-clamav-service/internal/storage"
)

// Reloader defines the interface to trigger ClamAV daemon database reload.
type Reloader interface {
	Reload(ctx context.Context) error
}

// Manager manages YARA rule validation, disk storage, database records, and clamd reload.
type Manager struct {
	rulesDir string
	storage  *storage.DB
	reloader Reloader
	mu       sync.Mutex
}

// NewManager creates a new YARA Manager instance.
func NewManager(rulesDir string, db *storage.DB, reloader Reloader) *Manager {
	return &Manager{
		rulesDir: rulesDir,
		storage:  db,
		reloader: reloader,
	}
}

// AddRule validates, saves to disk, records in DB, and reloads clamd.
func (m *Manager) AddRule(ctx context.Context, ruleName, description, content, author string) (*storage.YARARule, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 1. Validate syntax and safety
	if err := ValidateRule(ruleName, content); err != nil {
		return nil, fmt.Errorf("yara validation failed: %w", err)
	}

	// 2. Ensure rules directory exists
	if err := os.MkdirAll(m.rulesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create yara rules dir: %w", err)
	}

	rulePath := filepath.Join(m.rulesDir, ruleName+".yara")

	// Check if file already exists
	if _, err := os.Stat(rulePath); err == nil {
		return nil, fmt.Errorf("rule file '%s.yara' already exists", ruleName)
	}

	// 3. Write file to disk
	if err := os.WriteFile(rulePath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write yara rule file: %w", err)
	}

	// 4. Save metadata to DB
	now := time.Now().UTC()
	ruleRecord := storage.YARARule{
		ID:          uuid.New().String(),
		RuleName:    ruleName,
		Description: description,
		Content:     content,
		Author:      author,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := m.storage.InsertYARARule(ruleRecord); err != nil {
		_ = os.Remove(rulePath)
		return nil, fmt.Errorf("failed to store rule in database: %w", err)
	}

	// 5. Trigger clamd reload if reloader provided
	if m.reloader != nil {
		if err := m.reloader.Reload(ctx); err != nil {
			// Rollback file and DB record upon reload failure
			_ = os.Remove(rulePath)
			_ = m.storage.DeleteYARARule(ruleRecord.ID)
			return nil, fmt.Errorf("failed to reload clamd with new rule: %w", err)
		}
	}

	return &ruleRecord, nil
}

// ListRules retrieves all registered YARA rules from the database.
func (m *Manager) ListRules(ctx context.Context) ([]storage.YARARule, error) {
	return m.storage.ListYARARules()
}

// GetRule retrieves a specific YARA rule by its ID.
func (m *Manager) GetRule(ctx context.Context, id string) (*storage.YARARule, error) {
	return m.storage.GetYARARule(id)
}

// DeleteRule removes the rule file from disk, deletes DB record, and triggers clamd reload.
func (m *Manager) DeleteRule(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 1. Get existing rule metadata
	rule, err := m.storage.GetYARARule(id)
	if err != nil {
		return fmt.Errorf("rule with id '%s' not found: %w", id, err)
	}

	// 2. Remove file from disk
	rulePath := filepath.Join(m.rulesDir, rule.RuleName+".yara")
	_ = os.Remove(rulePath)

	// 3. Remove from DB
	if err := m.storage.DeleteYARARule(id); err != nil {
		return fmt.Errorf("failed to delete rule from database: %w", err)
	}

	// 4. Trigger clamd reload
	if m.reloader != nil {
		if err := m.reloader.Reload(ctx); err != nil {
			return fmt.Errorf("failed to reload clamd after rule deletion: %w", err)
		}
	}

	return nil
}
