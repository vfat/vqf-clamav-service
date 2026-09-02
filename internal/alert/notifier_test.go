package alert

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestAlert_FormatTelegramMessage(t *testing.T) {
	alert := ThreatAlert{
		VirusName:      "Win.Trojan.Agent-9812",
		FileName:       "invoice.pdf.exe",
		FileSizeBytes:  2048576,
		FileSHA256:     "275a021bbfb6489e54d471899f7db9d1663fc695ec2fe2a2c4538aabf651fd0f",
		QuarantineID:   "Q-20260902-9812",
		SourceConsumer: "Billing-Service",
		DetectedAt:     time.Now().UTC(),
	}

	msg := formatTelegramMessage(alert)

	if !strings.Contains(msg, "Win.Trojan.Agent-9812") {
		t.Errorf("expected virus name in telegram message, got:\n%s", msg)
	}
	if !strings.Contains(msg, "invoice.pdf.exe") {
		t.Errorf("expected filename in telegram message, got:\n%s", msg)
	}
	if !strings.Contains(msg, "Q-20260902-9812") {
		t.Errorf("expected quarantine ID in telegram message, got:\n%s", msg)
	}
}

func TestAlert_FormatDiscordEmbed(t *testing.T) {
	alert := ThreatAlert{
		VirusName:      "Eicar-Test-Signature",
		FileName:       "eicar.com",
		FileSizeBytes:  68,
		FileSHA256:     "275a021b...",
		QuarantineID:   "Q-20260902-01",
		SourceConsumer: "Gateway",
		DetectedAt:     time.Now().UTC(),
	}

	embedPayload, err := formatDiscordPayload(alert)
	if err != nil {
		t.Fatalf("failed to format discord payload: %v", err)
	}

	jsonStr := string(embedPayload)
	if !strings.Contains(jsonStr, "Eicar-Test-Signature") {
		t.Errorf("expected virus name in discord embed, got:\n%s", jsonStr)
	}
	if !strings.Contains(jsonStr, "QUARANTINED") {
		t.Errorf("expected QUARANTINED in discord embed, got:\n%s", jsonStr)
	}
}

func TestAlert_ThrottlingLogic(t *testing.T) {
	notifier := NewNotifier(Config{})

	// Fire 5 alerts rapidly: all 5 should return throttled=false
	for i := 0; i < 5; i++ {
		throttled, count := notifier.recordAndCheckThrottle()
		if throttled {
			t.Fatalf("alert %d should not be throttled (under threshold 5)", i+1)
		}
		if count != i+1 {
			t.Errorf("expected count %d, got %d", i+1, count)
		}
	}

	// 6th alert must be throttled
	throttled, count := notifier.recordAndCheckThrottle()
	if !throttled {
		t.Fatalf("6th alert should be throttled")
	}
	if count != 6 {
		t.Errorf("expected count 6, got %d", count)
	}

	// DispatchThreat with disabled config should not error
	ctx := context.Background()
	alert := ThreatAlert{VirusName: "Test"}
	if err := notifier.DispatchThreat(ctx, alert); err != nil {
		t.Fatalf("DispatchThreat failed: %v", err)
	}
}
