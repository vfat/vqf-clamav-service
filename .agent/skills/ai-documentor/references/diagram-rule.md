# Diagram Rule

## C4 Component Diagram

- `C4-Component-Diagrams.md` harus menggunakan Mermaid
- Untuk component diagram, gunakan syntax `C4Component`
- Nama komponen harus mengikuti nama artefak atau komponen aktual di repo, bukan label generik yang tidak traceable
- Relationship pada diagram harus tetap evidence-based; jika sebuah relasi belum terverifikasi dari code, tandai sebagai asumsi atau jangan digambar sebagai fakta implementasi

## Contoh Baseline Mermaid untuk C4 Component Diagram

```mermaid
C4Component
title Component Diagram for Desktop Apps Trader Panel

Container(webTraderPanel, "Web Trader Panel", "Web App")
Container(userApiServices, "User API Services", "API Service")
Container(mqlApiServices, "MQL API Services", "API Service")
ContainerDb(database, "Database", "Database")
ContainerDb(databaseLocal, "Database Local", "Local Database")

System_Ext(awsSes, "AWS SES", "Email Service")
System_Ext(googleServices, "Google Services", "OAuth Service")
System_Ext(fileStorage, "File Storage", "Storage Service")
System_Ext(metaTraderExt, "MetaTrader (MT5/MT4)", "Trading Platform")
System_Ext(windowsOs, "Windows OS", "Host OS")
System_Ext(siem, "SIEM", "Monitoring")

Container_Boundary(desktopApps, "Desktop Apps Trader Panel") {
	Component(login, "Login", "WPF Component")
	Component(signUp, "Sign Up", "WPF Component")
	Component(dashboard, "Dashboard", "WPF Component")
	Component(metaTraderPage, "MetaTrader", "WPF Component")
	Component(chart, "Chart", "WPF Component")
	Component(trade, "Trade", "WPF Component")
	Component(accountHistory, "Account History", "WPF Component")
	Component(createAccount, "Create Account", "WPF Component")
	Component(notification, "Notification", ".NET Component")
	Component(feedback, "Feedback", "WPF Component")
	Component(brokerInstaller, "Broker Installer", ".NET Component")
	Component(setting, "Setting", "WPF Component")
	Component(dataCollector, "Data Collector", ".NET Component")
	Component(help, "Help", "WPF Component")
	Component(appsInstaller, "Apps Installer", ".NET Component")

	Rel(login, userApiServices, "Uses")
	Rel(signUp, userApiServices, "Uses")
	Rel(dashboard, userApiServices, "Uses")
	Rel(createAccount, userApiServices, "Uses")
	Rel(notification, userApiServices, "Uses")
	Rel(feedback, userApiServices, "Uses")

	Rel(metaTraderPage, mqlApiServices, "Uses")
	Rel(brokerInstaller, mqlApiServices, "Uses")

	Rel(login, googleServices, "Uses")
	Rel(login, awsSes, "Uses")
	Rel(signUp, awsSes, "Uses")

	Rel(appsInstaller, fileStorage, "Downloads")
	Rel(brokerInstaller, fileStorage, "Downloads")
	Rel(brokerInstaller, metaTraderExt, "Installs")
	Rel(dataCollector, metaTraderExt, "Collects data from")
	Rel(dataCollector, siem, "Sends logs")

	Rel(dashboard, databaseLocal, "Reads")
	Rel(setting, databaseLocal, "Writes")
	Rel(dataCollector, databaseLocal, "Writes")
}

Rel(webTraderPanel, userApiServices, "Uses")
Rel(userApiServices, database, "Reads/Writes")
Rel(mqlApiServices, database, "Reads/Writes")
Rel(databaseLocal, database, "Syncs")
```
