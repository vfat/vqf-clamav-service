package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/vfat/vqf-clamav-service/internal/alert"
	"github.com/vfat/vqf-clamav-service/internal/api"
	"github.com/vfat/vqf-clamav-service/internal/clamd"
	"github.com/vfat/vqf-clamav-service/internal/crypto"
	"github.com/vfat/vqf-clamav-service/internal/quarantine"
	"github.com/vfat/vqf-clamav-service/internal/ratelimit"
	"github.com/vfat/vqf-clamav-service/internal/storage"
	"github.com/vfat/vqf-clamav-service/internal/supervisor"
)

func main() {
	fmt.Println("==================================================================")
	fmt.Println("       🛡️  CLAMAV-SERVICE (High-Performance Antivirus API)        ")
	fmt.Println("==================================================================")

	envPath := getEnv("ENV_FILE", ".env")

	// 1. Zero-Touch Master Key Generation
	keyStr, _, created, err := crypto.EnsureMasterKey(envPath)
	if err != nil {
		log.Fatalf("[SECURITY FATAL] Failed to ensure master encryption key: %v", err)
	}

	if created {
		fmt.Println("------------------------------------------------------------------")
		fmt.Println("🔑 [SECURITY NOTICE] A new Master Encryption Key was generated:")
		fmt.Printf("   %s\n", keyStr)
		fmt.Println("   Auto-injected into .env. Keep a secure offline backup for recovery!")
		fmt.Println("------------------------------------------------------------------")
	}

	// 2. Storage & SQLite Initialization
	dbPath := getEnv("DB_PATH", "/data/clamav-service.db")
	db, err := storage.NewDB(dbPath)
	if err != nil {
		log.Fatalf("[STORAGE FATAL] Failed to initialize SQLite WAL database at %s: %v", dbPath, err)
	}
	defer db.Close()
	log.Printf("[STORAGE] SQLite initialized at %s (WAL mode)", dbPath)

	// 3. Quarantine Vault Initialization
	quarDir := getEnv("QUARANTINE_DIR", "/data/quarantine")
	vault := quarantine.NewVault(quarDir, db)
	log.Printf("[VAULT] Quarantine Vault initialized at %s", quarDir)

	// 4. Alert Dispatcher Initialization
	alertCfg := alert.Config{
		TelegramEnabled:   getEnvBool("ALERT_TELEGRAM_ENABLED", false),
		TelegramBotToken:  getEnv("ALERT_TELEGRAM_BOT_TOKEN", ""),
		TelegramChatID:    getEnv("ALERT_TELEGRAM_CHAT_ID", ""),
		DiscordEnabled:    getEnvBool("ALERT_DISCORD_ENABLED", false),
		DiscordWebhookURL: getEnv("ALERT_DISCORD_WEBHOOK_URL", ""),
		SMTPEnabled:       getEnvBool("ALERT_SMTP_ENABLED", false),
		SMTPHost:          getEnv("ALERT_SMTP_HOST", ""),
		SMTPPort:          getEnvInt("ALERT_SMTP_PORT", 587),
		SMTPUser:          getEnv("ALERT_SMTP_USER", ""),
		SMTPPassword:      getEnv("ALERT_SMTP_PASSWORD", ""),
		SMTPFrom:          getEnv("ALERT_SMTP_FROM", ""),
	}
	notifier := alert.NewNotifier(alertCfg)

	// 5. Rate Limiter Initialization
	limiter := ratelimit.NewLimiter()
	rateLimitRPM := getEnvInt("RATE_LIMIT_DEFAULT_RPM_PER_KEY", 100)
	rateLimitEnabled := getEnvBool("RATE_LIMIT_ENABLED", true)

	// 6. ClamAV Native Process Supervisor (PID 1 supervises clamd & freshclam)
	clamdSocket := getEnv("CLAMD_SOCKET", "/var/run/clamav/clamd.ctl")
	clamdConf := getEnv("CLAMD_CONFIG", "/etc/clamav/clamd.conf")
	freshclamConf := getEnv("FRESHCLAM_CONFIG", "/etc/clamav/freshclam.conf")
	dataSigDir := "/var/lib/clamav"

	sup := supervisor.NewSupervisor()
	_ = sup.Start(context.Background(), clamdConf, freshclamConf, dataSigDir, clamdSocket)

	clamdClient := clamd.NewClient("unix", clamdSocket)
	log.Printf("[CLAMD] Unix Socket Client targeting %s", clamdSocket)

	// 7. Background Auto-Purge Worker (3d Logs, 7d Quarantine)
	logRetentionDays := getEnvInt("LOG_RETENTION_DAYS", 3)
	quarRetentionDays := getEnvInt("QUARANTINE_RETENTION_DAYS", 7)
	startBackgroundCleaner(db, vault, logRetentionDays)

	// 8. HTTP REST API Server
	maxScanMB := int64(getEnvInt("MAX_SCAN_SIZE_MB", 100))
	server := api.NewServer(api.ServerConfig{
		DB:               db,
		Vault:            vault,
		Notifier:         notifier,
		Limiter:          limiter,
		Clamd:            clamdClient,
		RequireAPIKey:    getEnvBool("REQUIRE_API_KEY", false),
		MaxScanSizeMB:    maxScanMB,
		RateLimitRPM:     rateLimitRPM,
		RateLimitEnabled: rateLimitEnabled,
		LogRetention:     logRetentionDays,
		QuarRetention:    quarRetentionDays,
		AuthMode:         getEnv("AUTH_MODE", "none"),
		BasicUser:        getEnv("AUTH_BASIC_USER", "admin"),
		BasicPass:        getEnv("AUTH_BASIC_PASS", ""),
		BearerToken:      getEnv("AUTH_BEARER_TOKEN", ""),
		UIPassword:       getEnv("UI_PASSWORD", "123456"),
	})

	port := getEnv("PORT", "8080")
	addr := ":" + port
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      server.Router(),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	// Graceful Shutdown Channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("[API] Server listening on http://0.0.0.0:%s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[API FATAL] Server failed: %v", err)
		}
	}()

	<-stop
	log.Println("[SHUTDOWN] Received termination signal. Initiating graceful shutdown...")

	sup.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("[SHUTDOWN ERROR] HTTP shutdown failed: %v", err)
	}

	log.Println("[SHUTDOWN] Service stopped gracefully.")
}

func startBackgroundCleaner(db *storage.DB, vault *quarantine.Vault, logRetentionDays int) {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			ctx := context.Background()
			purgedLogs, err := db.PurgeAuditLogs(logRetentionDays)
			if err == nil && purgedLogs > 0 {
				log.Printf("[CLEANER] Auto-purged %d audit logs (> %d days)", purgedLogs, logRetentionDays)
			}

			purgedFiles, err := vault.PurgeExpired(ctx)
			if err == nil && purgedFiles > 0 {
				log.Printf("[CLEANER] Auto-purged %d expired quarantine files", purgedFiles)
			}
		}
	}()
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}
	return b
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return i
}
