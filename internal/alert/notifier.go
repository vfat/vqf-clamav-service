package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// ThreatAlert encapsulates details of a detected malware payload.
type ThreatAlert struct {
	VirusName      string    `json:"virus_name"`
	FileName       string    `json:"file_name"`
	FileSizeBytes  int64     `json:"file_size_bytes"`
	FileSHA256     string    `json:"file_sha256"`
	QuarantineID   string    `json:"quarantine_id"`
	SourceConsumer string    `json:"source_consumer"`
	DetectedAt     time.Time `json:"detected_at"`
}

// Config specifies credentials for all notification channels.
type Config struct {
	TelegramEnabled  bool
	TelegramBotToken string
	TelegramChatID   string

	DiscordEnabled    bool
	DiscordWebhookURL string

	SMTPEnabled    bool
	SMTPHost       string
	SMTPPort       int
	SMTPUser       string
	SMTPPassword   string
	SMTPFrom       string
	SMTPRecipients []string
}

// Notifier dispatches alerts across configured channels with flood throttling.
type Notifier struct {
	cfg        Config
	httpClient *http.Client
	mu         sync.Mutex
	alertTimes []time.Time
}

// NewNotifier initializes the alert dispatcher.
func NewNotifier(cfg Config) *Notifier {
	return &Notifier{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		alertTimes: make([]time.Time, 0),
	}
}

// DispatchThreat broadcasts the threat alert to active channels if not throttled.
func (n *Notifier) DispatchThreat(ctx context.Context, alert ThreatAlert) error {
	throttled, _ := n.recordAndCheckThrottle()
	if throttled {
		// Throttling active: flood protection engaged
		return nil
	}

	if n.cfg.TelegramEnabled && n.cfg.TelegramBotToken != "" && n.cfg.TelegramChatID != "" {
		_ = n.sendTelegram(ctx, alert)
	}

	if n.cfg.DiscordEnabled && n.cfg.DiscordWebhookURL != "" {
		_ = n.sendDiscord(ctx, alert)
	}

	return nil
}

func (n *Notifier) recordAndCheckThrottle() (bool, int) {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-1 * time.Minute)

	// Filter out alerts older than 1 minute
	recent := make([]time.Time, 0)
	for _, t := range n.alertTimes {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}

	recent = append(recent, now)
	n.alertTimes = recent

	// Throttle if more than 5 alerts in the last minute
	if len(recent) > 5 {
		return true, len(recent)
	}

	return false, len(recent)
}

func (n *Notifier) sendTelegram(ctx context.Context, alert ThreatAlert) error {
	text := formatTelegramMessage(alert)
	payload := map[string]interface{}{
		"chat_id":    n.cfg.TelegramChatID,
		"text":       text,
		"parse_mode": "HTML",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.cfg.TelegramBotToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (n *Notifier) sendDiscord(ctx context.Context, alert ThreatAlert) error {
	body, err := formatDiscordPayload(alert)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.cfg.DiscordWebhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func formatTelegramMessage(alert ThreatAlert) string {
	return fmt.Sprintf(
		"🚨 <b>[CLAMAV ALERT] Malware Detected!</b>\n"+
			"───────────────────────────────\n"+
			"🦠 <b>Threat:</b> <code>%s</code>\n"+
			"📁 <b>File:</b> <code>%s</code> (%.2f MB)\n"+
			"🔑 <b>Hash:</b> <code>%s</code>\n"+
			"🛡️ <b>Vault:</b> <code>%s</code> (QUARANTINED)\n"+
			"🏷️ <b>Source:</b> %s\n"+
			"⏰ <b>Time:</b> %s",
		alert.VirusName,
		alert.FileName,
		float64(alert.FileSizeBytes)/(1024*1024),
		alert.FileSHA256,
		alert.QuarantineID,
		alert.SourceConsumer,
		alert.DetectedAt.Format("2006-01-02 15:04:05 MST"),
	)
}

func formatDiscordPayload(alert ThreatAlert) ([]byte, error) {
	payload := map[string]interface{}{
		"username": "ClamAV Security Guard",
		"embeds": []map[string]interface{}{
			{
				"title": "🚨 Malware Threat Detected!",
				"color": 15158332, // Red color
				"fields": []map[string]interface{}{
					{"name": "Threat Name", "value": fmt.Sprintf("`%s`", alert.VirusName), "inline": true},
					{"name": "File Name", "value": fmt.Sprintf("`%s`", alert.FileName), "inline": true},
					{"name": "Status", "value": "🛡️ **QUARANTINED**", "inline": true},
					{"name": "File SHA256", "value": fmt.Sprintf("`%s`", alert.FileSHA256), "inline": false},
					{"name": "Source Consumer", "value": alert.SourceConsumer, "inline": true},
					{"name": "Vault ID", "value": fmt.Sprintf("`%s`", alert.QuarantineID), "inline": true},
				},
				"footer": map[string]string{
					"text": fmt.Sprintf("ClamAV Security Telemetry • %s", alert.DetectedAt.Format(time.RFC3339)),
				},
			},
		},
	}

	return json.Marshal(payload)
}
