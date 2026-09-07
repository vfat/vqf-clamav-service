package yara

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/vfat/vqf-clamav-service/internal/storage"
)

type mockReloader struct {
	reloadCalled int
	reloadErr    error
}

func (m *mockReloader) Reload(ctx context.Context) error {
	m.reloadCalled++
	return m.reloadErr
}

func setupTestManager(t *testing.T) (*Manager, *storage.DB, *mockReloader, string) {
	t.Helper()
	tmpDir := t.TempDir()
	rulesDir := filepath.Join(tmpDir, "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		t.Fatalf("failed to create rulesDir: %v", err)
	}

	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := storage.NewDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}

	mockReload := &mockReloader{}
	mgr := NewManager(rulesDir, db, mockReload)
	return mgr, db, mockReload, rulesDir
}

func TestValidateRule(t *testing.T) {
	tests := []struct {
		name        string
		ruleName    string
		content     string
		expectError bool
	}{
		{
			name:     "Valid YARA Rule",
			ruleName: "detect_webshell_php",
			content: `
rule detect_webshell_php {
    meta:
        description = "Detects simple PHP webshell"
    strings:
        $cmd = "passthru($_GET['cmd'])"
    condition:
        $cmd
}
`,
			expectError: false,
		},
		{
			name:        "Empty Content",
			ruleName:    "empty_rule",
			content:     "",
			expectError: true,
		},
		{
			name:        "Invalid Rule Name with Path Traversal",
			ruleName:    "../evil_rule",
			content:     "rule evil_rule { condition: true }",
			expectError: true,
		},
		{
			name:        "Invalid Rule Name with Spaces",
			ruleName:    "bad rule name",
			content:     "rule bad { condition: true }",
			expectError: true,
		},
		{
			name:     "Missing condition keyword",
			ruleName: "no_condition",
			content: `
rule no_condition {
    strings:
        $a = "test"
}
`,
			expectError: true,
		},
		{
			name:     "Missing rule declaration keyword",
			ruleName: "no_rule_keyword",
			content: `
webshell {
    condition:
        true
}
`,
			expectError: true,
		},
		{
			name:     "Unbalanced braces",
			ruleName: "unbalanced",
			content: `
rule unbalanced {
    condition:
        true
`,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateRule(tc.ruleName, tc.content)
			if tc.expectError && err == nil {
				t.Errorf("expected validation error, got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("expected valid rule, got error: %v", err)
			}
		})
	}
}

func TestManager_AddRule_Success(t *testing.T) {
	mgr, _, mockReload, rulesDir := setupTestManager(t)
	ctx := context.Background()

	validRule := `
rule test_eicar {
    meta:
        description = "Detect EICAR test string"
    strings:
        $eicar = "X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*"
    condition:
        $eicar
}
`
	rule, err := mgr.AddRule(ctx, "test_eicar", "Detect EICAR test string", validRule, "sec-admin")
	if err != nil {
		t.Fatalf("AddRule failed: %v", err)
	}

	if rule.ID == "" {
		t.Error("expected generated rule ID, got empty string")
	}
	if rule.RuleName != "test_eicar" {
		t.Errorf("expected rule name 'test_eicar', got '%s'", rule.RuleName)
	}

	// Verify file was written to rulesDir
	filePath := filepath.Join(rulesDir, "test_eicar.yara")
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("rule file not found on disk: %v", err)
	}
	if string(data) != validRule {
		t.Errorf("file content mismatch, expected '%s', got '%s'", validRule, string(data))
	}

	// Verify clamd reload was invoked
	if mockReload.reloadCalled != 1 {
		t.Errorf("expected clamd reload to be called once, got %d", mockReload.reloadCalled)
	}
}

func TestManager_AddRule_InvalidRuleDoesNotPersist(t *testing.T) {
	mgr, _, mockReload, rulesDir := setupTestManager(t)
	ctx := context.Background()

	invalidRule := "not a valid rule syntax"
	_, err := mgr.AddRule(ctx, "invalid_rule", "bad syntax", invalidRule, "sec-admin")
	if err == nil {
		t.Fatal("expected error for invalid rule, got nil")
	}

	// Verify no file created
	filePath := filepath.Join(rulesDir, "invalid_rule.yara")
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("expected rule file not to exist on disk")
	}

	if mockReload.reloadCalled != 0 {
		t.Errorf("expected clamd reload not to be called, got %d", mockReload.reloadCalled)
	}
}

func TestManager_AddRule_ReloadFailureRollback(t *testing.T) {
	mgr, db, mockReload, rulesDir := setupTestManager(t)
	ctx := context.Background()

	mockReload.reloadErr = errors.New("socket connection refused")

	validRule := `
rule rollback_test {
    condition:
        true
}
`
	_, err := mgr.AddRule(ctx, "rollback_test", "rollback test", validRule, "admin")
	if err == nil {
		t.Fatal("expected error due to reload failure, got nil")
	}

	// Verify file was cleaned up
	filePath := filepath.Join(rulesDir, "rollback_test.yara")
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("expected rule file to be removed upon reload failure")
	}

	// Verify DB entry was not retained
	rules, err := db.ListYARARules()
	if err != nil {
		t.Fatalf("ListYARARules error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected 0 rules in DB after rollback, got %d", len(rules))
	}
}

func TestManager_ListAndGetRule(t *testing.T) {
	mgr, _, _, _ := setupTestManager(t)
	ctx := context.Background()

	rule1 := "rule rule_one { condition: true }"
	rule2 := "rule rule_two { condition: true }"

	r1, err := mgr.AddRule(ctx, "rule_one", "desc 1", rule1, "user1")
	if err != nil {
		t.Fatalf("failed to add rule 1: %v", err)
	}
	r2, err := mgr.AddRule(ctx, "rule_two", "desc 2", rule2, "user2")
	if err != nil {
		t.Fatalf("failed to add rule 2: %v", err)
	}

	// List rules
	list, err := mgr.ListRules(ctx)
	if err != nil {
		t.Fatalf("ListRules failed: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(list))
	}

	// Get specific rule
	fetched, err := mgr.GetRule(ctx, r1.ID)
	if err != nil {
		t.Fatalf("GetRule failed: %v", err)
	}
	if fetched.RuleName != r1.RuleName {
		t.Errorf("expected ruleName '%s', got '%s'", r1.RuleName, fetched.RuleName)
	}

	_ = r2
}

func TestManager_DeleteRule(t *testing.T) {
	mgr, _, mockReload, rulesDir := setupTestManager(t)
	ctx := context.Background()

	ruleContent := "rule to_delete { condition: true }"
	r, err := mgr.AddRule(ctx, "to_delete", "will be deleted", ruleContent, "user")
	if err != nil {
		t.Fatalf("AddRule failed: %v", err)
	}

	filePath := filepath.Join(rulesDir, "to_delete.yara")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("expected file to exist before deletion")
	}

	mockReload.reloadCalled = 0 // reset counter

	// Delete rule
	if err := mgr.DeleteRule(ctx, r.ID); err != nil {
		t.Fatalf("DeleteRule failed: %v", err)
	}

	// Verify file is gone
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("expected rule file to be deleted from disk")
	}

	// Verify reload was called
	if mockReload.reloadCalled != 1 {
		t.Errorf("expected reload to be called on delete, got %d", mockReload.reloadCalled)
	}

	// Verify GetRule returns error
	_, err = mgr.GetRule(ctx, r.ID)
	if err == nil {
		t.Error("expected GetRule to fail for deleted rule, got nil")
	}
}
