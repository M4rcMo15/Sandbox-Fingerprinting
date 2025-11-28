package models

import "time"

// Payload principal que se envía al servidor
type Payload struct {
	Timestamp    time.Time        `json:"timestamp"`
	Hostname     string           `json:"hostname"`
	PublicIP     string           `json:"public_ip"`
	BinarySize   int64            `json:"binary_size_bytes"`
	GeoLocation  *GeoLocation     `json:"geo_location"`
	SandboxInfo  *SandboxInfo     `json:"sandbox_info"`
	SystemInfo   *SystemInfo      `json:"system_info"`
	HookInfo     *HookInfo        `json:"hook_info"`
	CrawlerInfo  *CrawlerInfo     `json:"crawler_info"`
	EDRInfo      *EDRInfo         `json:"edr_info"`
	ToolsInfo    *ToolsInfo       `json:"tools_info"`
}

// SandboxInfo - Información de detección de sandbox
type SandboxInfo struct {
	IsVM              bool     `json:"is_vm"`
	VMIndicators      []string `json:"vm_indicators"`
	RegistryIndicators []string `json:"registry_indicators"`
	DiskIndicators    []string `json:"disk_indicators"`
	CPUTemperature    float64  `json:"cpu_temperature"`
	WindowCount       int      `json:"window_count"`
	HasDebugPrivilege bool     `json:"has_debug_privilege"`
}

// SystemInfo - Información del sistema
type SystemInfo struct {
	OS              string            `json:"os"`
	Architecture    string            `json:"architecture"`
	Language        string            `json:"language"`
	Timezone        string            `json:"timezone"`
	CPUCount        int               `json:"cpu_count"`
	TotalRAM        uint64            `json:"total_ram_mb"`
	TotalDisk       int64             `json:"total_disk_bytes"`
	BIOS            string            `json:"bios"`
	Processes       []ProcessInfo     `json:"processes"`
	Users           []string          `json:"users"`
	Groups          []string          `json:"groups"`
	NetworkConns    []NetworkConn     `json:"network_connections"`
	Services        []string          `json:"services"`
	EnvVars         map[string]string `json:"environment_variables"`
	Pipes           []string          `json:"pipes"`
	Screenshot      string            `json:"screenshot_base64,omitempty"`
	MousePosition   Point             `json:"mouse_position"`
	InstalledApps   []string          `json:"installed_apps"`
	RecentFiles     []string          `json:"recent_files"`
	UptimeSeconds   int64             `json:"uptime_seconds"`
}

type ProcessInfo struct {
	PID   uint32 `json:"pid"`
	Name  string `json:"name"`
	Owner string `json:"owner"`
	Path  string `json:"path"`
}

type NetworkConn struct {
	Protocol    string `json:"protocol"`
	LocalAddr   string `json:"local_addr"`
	RemoteAddr  string `json:"remote_addr"`
	State       string `json:"state"`
}

type Point struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

// HookInfo - Información de hooks detectados
type HookInfo struct {
	HookedFunctions []HookedFunction `json:"hooked_functions"`
	SuspiciousDLLs  []string         `json:"suspicious_dlls"`
}

type HookedFunction struct {
	Module       string `json:"module"`
	Function     string `json:"function"`
	IsHooked     bool   `json:"is_hooked"`
	FirstBytes   string `json:"first_bytes"`
}

// CrawlerInfo - Información del crawler de archivos
type CrawlerInfo struct {
	ScannedPaths  []string `json:"scanned_paths"`
	FoundFiles    []string `json:"found_files"`
	TotalFiles    int      `json:"total_files"`
}

// EDRInfo - Información de EDR/AV detectados
type EDRInfo struct {
	DetectedProducts []EDRProduct `json:"detected_products"`
	RunningProcesses []string     `json:"running_processes"`
	InstalledDrivers []string     `json:"installed_drivers"`
}

type EDRProduct struct {
	Name     string `json:"name"`
	Type     string `json:"type"` // "EDR", "AV", "Sandbox"
	Detected bool   `json:"detected"`
	Method   string `json:"method"` // "process", "driver", "registry"
}

// GeoLocation - Información de geolocalización
type GeoLocation struct {
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Region      string  `json:"region"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ISP         string  `json:"isp"`
	Organization string `json:"organization"`
}

// ToolsInfo - Información de herramientas de análisis detectadas
type ToolsInfo struct {
	ReversingTools []string `json:"reversing_tools"`
	DebuggingTools []string `json:"debugging_tools"`
	MonitoringTools []string `json:"monitoring_tools"`
	VirtualizationTools []string `json:"virtualization_tools"`
	AnalysisTools []string `json:"analysis_tools"`
}
