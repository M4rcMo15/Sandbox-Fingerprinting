package models

import "time"

// Payload principal que se envía al servidor
type Payload struct {
	Timestamp   time.Time            `json:"timestamp"`
	Hostname    string               `json:"hostname"`
	PublicIP    string               `json:"public_ip"`
	BinarySize  int64                `json:"binary_size_bytes"`
	RawData     *RawData             `json:"raw_data"`
	SystemInfo  *SystemInfo          `json:"system_info"`
	HookInfo    *HookInfo            `json:"hook_info"`
	CrawlerInfo *CrawlerInfo         `json:"crawler_info"`
	XSSPayloads []XSSPayloadMetadata `json:"xss_payloads,omitempty"`
}

// XSSPayloadMetadata contiene la metadata de los payloads XSS inyectados
type XSSPayloadMetadata struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Vector string `json:"vector"`
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

type MousePoint struct {
	X    int32 `json:"x"`
	Y    int32 `json:"y"`
	Time int64 `json:"time"`
}

// HookInfo - Información de hooks detectados
type HookInfo struct {
	HookedFunctions []HookedFunction `json:"hooked_functions"`
	SuspiciousDLLs  []string         `json:"suspicious_dlls"`
}

type HookedFunction struct {
	Module     string `json:"module"`
	Function   string `json:"function"`
	IsHooked   bool   `json:"is_hooked"`
	FirstBytes string `json:"first_bytes"`
}

// CrawlerInfo - Información del crawler de archivos
type CrawlerInfo struct {
	ScannedPaths []string `json:"scanned_paths"`
	FoundFiles   []string `json:"found_files"`
	TotalFiles   int      `json:"total_files"`
}

// RawData - Datos en bruto sin procesamiento
// El análisis se realiza en el servidor
type RawData struct {
	VMFiles           []string       `json:"vm_files"`
	RegistryKeys      []RegistryKey  `json:"registry_keys"`
	SecurityProcesses []string       `json:"security_processes"`
	Drivers           []string       `json:"drivers"`
	DiskInfo          DiskInfo       `json:"disk_info"`
	CPUInfo           CPUInfo        `json:"cpu_info"`
	WindowCount       int            `json:"window_count"`
	TimingDiscrepancy float64        `json:"timing_discrepancy"`
	MouseHistory      []MousePoint   `json:"mouse_history"`
	CPUIDHypervisor   bool           `json:"cpuid_hypervisor_bit"`
	MACOUI            string         `json:"mac_address_oui"`
	ClipboardPreview  string         `json:"clipboard_content_preview"`
}

type RegistryKey struct {
	Path   string            `json:"path"`
	Name   string            `json:"name"`
	Exists bool              `json:"exists"`
	Values map[string]string `json:"values,omitempty"`
}

type DiskInfo struct {
	Identifier   string `json:"identifier"`
	SerialNumber string `json:"serial_number"`
	TotalBytes   int64  `json:"total_bytes"`
	FreeBytes    int64  `json:"free_bytes"`
}

type CPUInfo struct {
	ProcessorName string  `json:"processor_name"`
	Vendor        string  `json:"vendor"`
	Identifier    string  `json:"identifier"`
	Temperature   float64 `json:"temperature"`
}
