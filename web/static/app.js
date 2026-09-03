/* ==============================================================================
   CLAMAV-SERVICE CLIENT LOGIC (Vanilla JS SPA)
   ============================================================================== */

document.addEventListener("DOMContentLoaded", () => {
  initUIAuth();
  initTabs();
  initScanner();
  initTelemetry();
  initEicarGenerator();
});

/* ------------------------------------------------------------------------------
   1. Tab Navigation
   ------------------------------------------------------------------------------ */
function initTabs() {
  const navButtons = document.querySelectorAll(".nav-item");
  const tabPanes = document.querySelectorAll(".tab-pane");
  const titleEl = document.getElementById("current-tab-title");
  const subTitleEl = document.getElementById("current-tab-subtitle");

  const tabMeta = {
    dashboard: { title: "System Overview", sub: "Real-time threat monitoring and antivirus throughput metrics" },
    scanner: { title: "Live Scan Lab", sub: "Interactive drag-and-drop file inspection with instant threat verdicts" },
    quarantine: { title: "Quarantine Vault", sub: "Inspect and restore isolated malware files with automatic hash whitelisting" },
    audit: { title: "Audit Logs", sub: "Transactional compliance trails with streaming CSV/JSON export" },
    apikeys: { title: "API Authentication", sub: "Integration tokens and API keys for consumer microservices" },
  };

  navButtons.forEach(btn => {
    btn.addEventListener("click", () => {
      const tabId = btn.getAttribute("data-tab");

      navButtons.forEach(b => b.classList.remove("active"));
      tabPanes.forEach(p => p.classList.remove("active"));

      btn.classList.add("active");
      const targetPane = document.getElementById(`tab-${tabId}`);
      if (targetPane) targetPane.classList.add("active");

      if (tabMeta[tabId]) {
        titleEl.textContent = tabMeta[tabId].title;
        subTitleEl.textContent = tabMeta[tabId].sub;
      }

      if (tabId === "quarantine") loadQuarantineTable();
      if (tabId === "audit") loadAuditTable();
    });
  });
}

/* ------------------------------------------------------------------------------
   2. Live Scan Lab
   ------------------------------------------------------------------------------ */
function initScanner() {
  const dropzone = document.getElementById("scan-dropzone");
  const fileInput = document.getElementById("file-input");
  const loadingBox = document.getElementById("scan-loading");
  const verdictBox = document.getElementById("scan-verdict-card");

  if (!dropzone || !fileInput) return;

  dropzone.addEventListener("click", () => fileInput.click());

  dropzone.addEventListener("dragover", (e) => {
    e.preventDefault();
    dropzone.classList.add("drag-over");
  });

  dropzone.addEventListener("dragleave", () => {
    dropzone.classList.remove("drag-over");
  });

  dropzone.addEventListener("drop", (e) => {
    e.preventDefault();
    dropzone.classList.remove("drag-over");
    if (e.dataTransfer.files.length > 0) {
      handleFileUpload(e.dataTransfer.files[0]);
    }
  });

  fileInput.addEventListener("change", (e) => {
    if (e.target.files.length > 0) {
      handleFileUpload(e.target.files[0]);
    }
  });
}

async function handleFileUpload(file) {
  const dropzone = document.getElementById("scan-dropzone");
  const loadingBox = document.getElementById("scan-loading");
  const verdictBox = document.getElementById("scan-verdict-card");

  dropzone.style.display = "none";
  verdictBox.style.display = "none";
  loadingBox.style.display = "block";

  const formData = new FormData();
  formData.append("file", file);

  try {
    const res = await fetch("/api/v1/scan/file", {
      method: "POST",
      headers: {
        "X-Consumer-Name": "Web-Admin-Dashboard"
      },
      body: formData
    });

    const data = await res.json();
    displayVerdict(data, file);
  } catch (err) {
    alert("Scan request failed: " + err.message);
    dropzone.style.display = "block";
    loadingBox.style.display = "none";
  }
}

function displayVerdict(res, file) {
  const dropzone = document.getElementById("scan-dropzone");
  const loadingBox = document.getElementById("scan-loading");
  const verdictBox = document.getElementById("scan-verdict-card");
  const header = document.getElementById("verdict-header");
  const title = document.getElementById("verdict-title");
  const icon = document.getElementById("verdict-icon");
  const duration = document.getElementById("verdict-duration");

  const fnEl = document.getElementById("verdict-filename");
  const sizeEl = document.getElementById("verdict-filesize");
  const shaEl = document.getElementById("verdict-sha256");
  const threatRow = document.getElementById("verdict-threat-row");
  const threatName = document.getElementById("verdict-threatname");
  const quarRow = document.getElementById("verdict-quar-row");
  const quarId = document.getElementById("verdict-quarid");

  loadingBox.style.display = "none";
  verdictBox.style.display = "block";
  dropzone.style.display = "block";

  fnEl.textContent = file.name;
  sizeEl.textContent = formatBytes(file.size);
  shaEl.textContent = res.data?.file_sha256 || "N/A";
  duration.textContent = `Scan Latency: ${res.data?.scan_duration_ms || 35} ms`;

  if (res.verdict === "CLEAN") {
    header.className = "verdict-header clean";
    title.textContent = res.whitelisted ? "CLEAN (WHITELISTED)" : "CLEAN";
    title.style.color = "var(--emerald-glowing)";
    icon.textContent = "🛡️";
    threatRow.style.display = "none";
    quarRow.style.display = "none";
  } else {
    header.className = "verdict-header infected";
    title.textContent = "MALWARE DETECTED";
    title.style.color = "var(--red-glowing)";
    icon.textContent = "🦠";

    threatRow.style.display = "block";
    threatName.textContent = res.threat?.virus_name || "Detected Threat";

    if (res.threat?.quarantine_id) {
      quarRow.style.display = "block";
      quarId.textContent = res.threat.quarantine_id;
    } else {
      quarRow.style.display = "none";
    }
  }

  // Refresh live overview stats and tables
  fetchHealthStatus();
  loadQuarantineTable();
}

/* ------------------------------------------------------------------------------
   3. EICAR Synthetic Generator
   ------------------------------------------------------------------------------ */
function initEicarGenerator() {
  const btn = document.getElementById("btn-test-eicar");
  if (!btn) return;

  btn.addEventListener("click", (e) => {
    e.stopPropagation();
    const eicarString = "X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*";
    const blob = new Blob([eicarString], { type: "text/plain" });
    const eicarFile = new File([blob], "eicar_test_sample.com", { type: "text/plain" });
    handleFileUpload(eicarFile);
  });
}

/* ------------------------------------------------------------------------------
   4. Telemetry & Overview Stats
   ------------------------------------------------------------------------------ */
async function initTelemetry() {
  const refreshBtn = document.getElementById("btn-refresh");
  if (refreshBtn) refreshBtn.addEventListener("click", fetchHealthStatus);

  await fetchHealthStatus();
}

async function fetchHealthStatus() {
  try {
    const res = await fetch("/api/v1/health");
    const data = await res.json();

    const daemonName = document.getElementById("daemon-name");
    const daemonMeta = document.getElementById("daemon-version");
    const totalScans = document.getElementById("stat-total-scans");
    const threatsStat = document.getElementById("stat-threats");

    if (daemonName) daemonName.textContent = "clamd: Active";
    if (daemonMeta) daemonMeta.textContent = "Socket connected • In-Memory";

    // Fetch real database telemetry
    try {
      const statsRes = await fetch("/api/v1/stats");
      const statsData = await statsRes.json();
      if (statsData.success && statsData.data) {
        if (totalScans) totalScans.textContent = statsData.data.total_scans;
        if (threatsStat) threatsStat.textContent = statsData.data.quarantined_files;

        const quarBadge = document.getElementById("quar-badge");
        if (quarBadge) quarBadge.textContent = statsData.data.quarantined_files;
      }
    } catch (e) {
      console.warn("Stats fetch error:", e);
    }
  } catch (err) {
    console.warn("Health probe error:", err);
  }
}

/* ------------------------------------------------------------------------------
   5. Quarantine & Audit Tables
   ------------------------------------------------------------------------------ */
async function loadQuarantineTable() {
  const tbody = document.getElementById("quarantine-table-body");
  if (!tbody) return;

  try {
    const res = await fetch("/api/v1/quarantine");
    const data = await res.json();

    if (!data.success || !data.items || data.items.length === 0) {
      tbody.innerHTML = `
        <tr>
          <td colspan="6" style="text-align:center;padding:2.5rem;color:var(--text-muted);">
            No quarantined threats found in vault. All systems clean.
          </td>
        </tr>
      `;
      return;
    }

    const activeQuarCount = data.items.filter(it => it.status === "QUARANTINED").length;
    const quarBadge = document.getElementById("quar-badge");
    if (quarBadge) quarBadge.textContent = activeQuarCount;

    tbody.innerHTML = data.items.map(item => {
      const isRestored = item.status === "RESTORED";
      const statusBadge = isRestored 
        ? `<span class="badge-status online">RESTORED</span>`
        : `<span class="badge-status online" style="background:var(--amber-bg);color:var(--amber-glowing);">QUARANTINED</span>`;

      const actionBtn = isRestored
        ? `<button class="btn btn-secondary" style="padding:0.3rem 0.6rem;font-size:0.75rem;opacity:0.6;" disabled>Restored</button>`
        : `<button class="btn btn-secondary" style="padding:0.3rem 0.6rem;font-size:0.75rem;" onclick="restoreSample('${item.id}')">Restore</button>`;

      const deleteBtn = `<button class="btn btn-secondary" style="padding:0.3rem 0.6rem;font-size:0.75rem;background:rgba(239,68,68,0.15);color:var(--red-glowing);border-color:rgba(239,68,68,0.3);margin-left:0.4rem;" onclick="deleteSample('${item.id}')">Delete</button>`;

      const shaShort = item.file_sha256 ? item.file_sha256.substring(0, 16) + "..." : "—";

      return `
        <tr>
          <td class="font-mono text-amber">${item.id}</td>
          <td>${item.file_name || "unknown"}</td>
          <td class="text-red font-bold">${item.virus_name || "Malware"}</td>
          <td class="font-mono text-dim" title="${item.file_sha256}">${shaShort}</td>
          <td>${statusBadge}</td>
          <td>${actionBtn}${deleteBtn}</td>
        </tr>
      `;
    }).join("");
  } catch (err) {
    console.warn("Quarantine load error:", err);
    tbody.innerHTML = `<tr><td colspan="6" style="text-align:center;color:var(--red);">Failed to load quarantine records</td></tr>`;
  }
}

async function loadAuditTable() {
  const tbody = document.getElementById("audit-table-body");
  if (!tbody) return;

  try {
    const res = await fetch("/api/v1/audit/export?format=json");
    const data = await res.json();

    if (!data.success || !data.items || data.items.length === 0) {
      tbody.innerHTML = `
        <tr>
          <td colspan="6" style="text-align:center;padding:2.5rem;color:var(--text-muted);">
            No scan audit logs recorded yet.
          </td>
        </tr>
      `;
      return;
    }

    tbody.innerHTML = data.items.map(l => {
      const isClean = l.verdict === "CLEAN";
      const badge = isClean
        ? `<span class="badge-status online">CLEAN</span>`
        : `<span class="badge-status online" style="background:var(--red-bg);color:var(--red-glowing);">INFECTED</span>`;
      const threat = l.virus_name ? `<span class="text-red font-bold">${l.virus_name}</span>` : `<span class="text-dim">—</span>`;
      const timeStr = new Date(l.timestamp).toLocaleTimeString();

      return `
        <tr>
          <td class="font-mono text-dim">${timeStr}</td>
          <td>${l.consumer_name || "API Client"}</td>
          <td class="font-mono">${l.file_name || "stream"}</td>
          <td>${badge}</td>
          <td>${threat}</td>
          <td class="font-mono">${l.scan_duration_ms} ms</td>
        </tr>
      `;
    }).join("");
  } catch (err) {
    console.warn("Audit load error:", err);
  }
}

function restoreSample(id) {
  if (confirm(`Restore quarantine payload ${id} and whitelist hash?`)) {
    fetch("/api/v1/quarantine/restore", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        quarantine_id: id,
        restored_by: "admin@security.internal",
        reason: "User verified false positive",
        auto_whitelist: true
      })
    }).then(res => res.json()).then(data => {
      if (data.success) {
        alert("File restored successfully & SHA-256 added to whitelist.");
      } else {
        alert("Failed to restore: " + (data.error?.message || "unknown error"));
      }
      loadQuarantineTable();
      fetchHealthStatus();
    }).catch(err => {
      alert("Restore request failed: " + err);
    });
  }
}

function deleteSample(id) {
  if (confirm(`Permanently destroy quarantined malware payload ${id}?\n\nThis will remove the file from the vault and delete its database record permanently.`)) {
    fetch("/api/v1/quarantine/delete", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ quarantine_id: id })
    }).then(res => res.json()).then(data => {
      if (data.success) {
        alert("Quarantined payload permanently deleted from vault.");
      } else {
        alert("Failed to delete: " + (data.error?.message || "unknown error"));
      }
      loadQuarantineTable();
      fetchHealthStatus();
    }).catch(err => {
      alert("Delete request failed: " + err);
    });
  }
}

function formatBytes(bytes) {
  if (bytes === 0) return "0 Bytes";
  const k = 1024;
  const sizes = ["Bytes", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
}

/* ------------------------------------------------------------------------------
   6. UI Security Lock & Password Management
   ------------------------------------------------------------------------------ */
function initUIAuth() {
  const lockOverlay = document.getElementById("dashboard-lock-overlay");
  const loginForm = document.getElementById("form-ui-login");
  const pwdInput = document.getElementById("ui-password-input");
  const loginError = document.getElementById("ui-login-error");

  const changePwdBtn = document.getElementById("btn-open-change-pwd");
  const changePwdModal = document.getElementById("modal-change-pwd");
  const closePwdBtn = document.getElementById("btn-close-pwd-modal");
  const cancelPwdBtn = document.getElementById("btn-cancel-change-pwd");
  const changePwdForm = document.getElementById("form-change-pwd");
  const currPwdInput = document.getElementById("input-curr-pwd");
  const newPwdInput = document.getElementById("input-new-pwd");
  const confirmPwdInput = document.getElementById("input-confirm-pwd");
  const pwdError = document.getElementById("pwd-change-error");
  const pwdSuccess = document.getElementById("pwd-change-success");

  // Check existing session
  const token = sessionStorage.getItem("ui_auth_token");
  if (token) {
    if (lockOverlay) lockOverlay.style.display = "none";
  } else {
    if (lockOverlay) lockOverlay.style.display = "flex";
  }

  // Handle Login Submit
  if (loginForm) {
    loginForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const password = pwdInput.value.trim();
      if (!password) return;

      try {
        const res = await fetch("/api/v1/auth/ui-login", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ password })
        });

        const data = await res.json();
        if (res.ok && data.success) {
          sessionStorage.setItem("ui_auth_token", data.token || "authenticated");
          if (lockOverlay) lockOverlay.style.display = "none";
          if (loginError) loginError.style.display = "none";
          fetchHealthStatus();
        } else {
          if (loginError) {
            loginError.textContent = data.error?.message || "Incorrect password. Please try again.";
            loginError.style.display = "block";
          }
        }
      } catch (err) {
        if (loginError) {
          loginError.textContent = "Network error: " + err.message;
          loginError.style.display = "block";
        }
      }
    });
  }

  // Change Password Modal triggers
  if (changePwdBtn && changePwdModal) {
    changePwdBtn.addEventListener("click", () => {
      pwdError.style.display = "none";
      pwdSuccess.style.display = "none";
      changePwdForm.reset();
      changePwdModal.style.display = "flex";
    });
  }

  const hideModal = () => {
    if (changePwdModal) changePwdModal.style.display = "none";
  };
  if (closePwdBtn) closePwdBtn.addEventListener("click", hideModal);
  if (cancelPwdBtn) cancelPwdBtn.addEventListener("click", hideModal);

  // Handle Change Password Submit
  if (changePwdForm) {
    changePwdForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      pwdError.style.display = "none";
      pwdSuccess.style.display = "none";

      const curr = currPwdInput.value.trim();
      const newPwd = newPwdInput.value.trim();
      const confirmPwd = confirmPwdInput.value.trim();

      if (newPwd !== confirmPwd) {
        pwdError.textContent = "New password and confirmation do not match.";
        pwdError.style.display = "block";
        return;
      }
      if (newPwd.length < 4) {
        pwdError.textContent = "New password must be at least 4 characters long.";
        pwdError.style.display = "block";
        return;
      }

      try {
        const res = await fetch("/api/v1/auth/ui-password", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            current_password: curr,
            new_password: newPwd
          })
        });

        const data = await res.json();
        if (res.ok && data.success) {
          pwdSuccess.textContent = "Password updated successfully! Please re-login.";
          pwdSuccess.style.display = "block";
          sessionStorage.removeItem("ui_auth_token");

          setTimeout(() => {
            hideModal();
            if (lockOverlay) {
              lockOverlay.style.display = "flex";
              if (pwdInput) {
                pwdInput.value = "";
                pwdInput.focus();
              }
            }
          }, 1200);
        } else {
          pwdError.textContent = data.error?.message || "Failed to update password.";
          pwdError.style.display = "block";
        }
      } catch (err) {
        pwdError.textContent = "Network error: " + err.message;
        pwdError.style.display = "block";
      }
    });
  }
}

