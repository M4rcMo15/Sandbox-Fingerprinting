package signatures

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

// APTSignature representa una firma de APT conocida
type APTSignature struct {
	Name        string
	Description string
	Strings     []string
	Mutexes     []string
	RegKeys     []string
	Files       []string
	C2Domains   []string
}

// GetAPTSignatures retorna firmas de APTs conocidos para aumentar detección
func GetAPTSignatures() []APTSignature {
	return []APTSignature{
		// APT29 (Cozy Bear) - Ruso, muy activo
		{
			Name:        "APT29_CozyBear",
			Description: "Russian APT group, also known as Cozy Bear or The Dukes",
			Strings: []string{
				"SeDebugPrivilege",
				"WellMess",
				"SUNBURST",
				"TEARDROP",
			},
			Mutexes: []string{
				"Global\\{12345678-1234-1234-1234-123456789012}",
				"Local\\SM0:1234:304:WilStaging_02",
			},
			RegKeys: []string{
				"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run\\SolarWinds",
				"SYSTEM\\CurrentControlSet\\Services\\SolarWinds.Orion.Core.BusinessLayer",
			},
			C2Domains: []string{
				"avsvmcloud.com",
				"digitalcollege.org",
				"freescanonline.com",
			},
		},

		// APT28 (Fancy Bear) - Ruso, muy conocido
		{
			Name:        "APT28_FancyBear",
			Description: "Russian military APT group, also known as Fancy Bear or Sofacy",
			Strings: []string{
				"X-Agent",
				"Sofacy",
				"CHOPSTICK",
				"Zebrocy",
			},
			Mutexes: []string{
				"Global\\{A5F6D7E8-9B0C-1D2E-3F4A-5B6C7D8E9F0A}",
			},
			RegKeys: []string{
				"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run\\SecurityUpdate",
				"SYSTEM\\CurrentControlSet\\Services\\WinDefender",
			},
			C2Domains: []string{
				"microsoft-update.net",
				"windows-security.org",
			},
		},

		// Lazarus Group (APT38) - Norcoreano
		{
			Name:        "Lazarus_APT38",
			Description: "North Korean APT group, responsible for WannaCry and Sony hack",
			Strings: []string{
				"WannaCry",
				"HOPLIGHT",
				"ELECTRICFISH",
				"BISTROMATH",
			},
			Mutexes: []string{
				"Global\\MsWinZonesCacheCounterMutexA",
				"Global\\{B5F6D7E8-9B0C-1D2E-3F4A-5B6C7D8E9F0B}",
			},
			RegKeys: []string{
				"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run\\tasksche",
				"SYSTEM\\CurrentControlSet\\Services\\mssecsvc",
			},
			C2Domains: []string{
				"wowser.com",
				"iuqerfsodp9ifjaposdfjhgosurijfaewrwergwea.com",
			},
		},

		// APT41 (Double Dragon) - Chino
		{
			Name:        "APT41_DoubleDragon",
			Description: "Chinese APT group conducting espionage and cybercrime",
			Strings: []string{
				"MESSAGETAP",
				"HIGHNOON",
				"POISONPLUG",
				"CROSSWALK",
			},
			Mutexes: []string{
				"Global\\{C5F6D7E8-9B0C-1D2E-3F4A-5B6C7D8E9F0C}",
			},
			RegKeys: []string{
				"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run\\GoogleUpdate",
				"SYSTEM\\CurrentControlSet\\Services\\VMwareService",
			},
			C2Domains: []string{
				"update.microsoft-security.com",
				"download.windowsupdate.org",
			},
		},

		// Emotet - Botnet muy activo
		{
			Name:        "Emotet_Botnet",
			Description: "Emotet banking trojan and botnet infrastructure",
			Strings: []string{
				"Emotet",
				"Heodo",
				"Geodo",
				"Mealybug",
			},
			Mutexes: []string{
				"Global\\I5H_SingleInstance",
				"Global\\M5O_SingleInstance",
			},
			RegKeys: []string{
				"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run\\Emotet",
				"SYSTEM\\CurrentControlSet\\Services\\WinDefService",
			},
			C2Domains: []string{
				"185.184.25.78",
				"190.90.233.66",
			},
		},

		// Cobalt Strike (usado por múltiples APTs)
		{
			Name:        "CobaltStrike_Beacon",
			Description: "Cobalt Strike beacon, used by multiple APT groups",
			Strings: []string{
				"beacon.dll",
				"cobaltstrike",
				"ReflectiveLoader",
				"MZ......................@",
			},
			Mutexes: []string{
				"Global\\{D5F6D7E8-9B0C-1D2E-3F4A-5B6C7D8E9F0D}",
			},
			RegKeys: []string{
				"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run\\CortanaAssistant",
			},
			C2Domains: []string{
				"cdn.cloudflare.net",
				"api.github.com",
			},
		},

		// Conti Ransomware (muy activo 2021-2022)
		{
			Name:        "Conti_Ransomware",
			Description: "Conti ransomware group, very active in 2021-2022",
			Strings: []string{
				"CONTI",
				"All of your files are currently encrypted",
				"conti_v3",
			},
			Mutexes: []string{
				"Global\\kjsdhglkjhsdfg8734",
			},
			RegKeys: []string{
				"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run\\ContiRecovery",
			},
			C2Domains: []string{
				"contirecovery.info",
				"contirecovery.best",
			},
		},

		// BlackCat/ALPHV Ransomware (2022-2024)
		{
			Name:        "BlackCat_ALPHV",
			Description: "BlackCat/ALPHV ransomware, written in Rust",
			Strings: []string{
				"ALPHV",
				"BlackCat",
				"Your data is stolen and encrypted",
			},
			Mutexes: []string{
				"Global\\{E5F6D7E8-9B0C-1D2E-3F4A-5B6C7D8E9F0E}",
			},
			RegKeys: []string{
				"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run\\RecoveryKey",
			},
			C2Domains: []string{
				"alphv.onion",
			},
		},
	}
}

// EmbedAPTSignatures embebe firmas de APT en el binario para aumentar detección
func EmbedAPTSignatures() {
	signatures := GetAPTSignatures()

	// Crear strings embebidos (no se usan, solo para detección estática)
	var embeddedStrings []string

	for _, sig := range signatures {
		embeddedStrings = append(embeddedStrings, sig.Strings...)
	}

	// Los strings quedan embebidos en el binario
	_ = embeddedStrings
}

// CreateAPTArtifacts crea artefactos de APT en el sistema (archivos, mutex, etc)
func CreateAPTArtifacts() {
	signatures := GetAPTSignatures()

	// Seleccionar 2-3 APTs aleatorios para mezclar firmas
	selectedAPTs := []APTSignature{
		signatures[0], // APT29
		signatures[2], // Lazarus
		signatures[5], // Cobalt Strike
	}

	for _, apt := range selectedAPTs {
		// Crear archivos temporales con nombres de APT
		for _, filename := range apt.Files {
			createTempFile(filename)
		}

		// Los mutex y regkeys se crean en otros módulos si es necesario
	}
}

// createTempFile crea un archivo temporal con contenido de APT
func createTempFile(filename string) {
	tempDir := os.TempDir()
	filePath := filepath.Join(tempDir, filename)

	// Contenido que parece de APT
	content := []byte("MZ\x90\x00\x03\x00\x00\x00\x04\x00\x00\x00\xFF\xFF\x00\x00")

	os.WriteFile(filePath, content, 0644)
}

// GenerateAPTHashes genera hashes conocidos de APTs (para YARA rules)
func GenerateAPTHashes() map[string]string {
	hashes := make(map[string]string)

	// Hashes MD5 de malware conocido (ejemplos ficticios pero realistas)
	hashes["APT29_SUNBURST"] = "d0d626deb3f9484e649294a8dfa814c5"
	hashes["APT28_XAgent"] = "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
	hashes["Lazarus_WannaCry"] = "db349b97c37d22f5ea1d1841e3c89eb4"
	hashes["Emotet_Loader"] = "4a17a9d2e5f6b7c8d9e0f1a2b3c4d5e6"
	hashes["CobaltStrike_Beacon"] = "5b28c9e3f6a7d8e9f0a1b2c3d4e5f6a7"

	return hashes
}

// GetYARASignature genera una regla YARA que detectaría este binario
func GetYARASignature() string {
	return `
rule APT_Artefacto_Research {
    meta:
        description = "Detects Artefacto research tool with APT signatures"
        author = "Security Researcher"
        date = "2024-12-05"
        
    strings:
        $apt29_1 = "SUNBURST" ascii
        $apt29_2 = "TEARDROP" ascii
        $apt28_1 = "X-Agent" ascii
        $apt28_2 = "Sofacy" ascii
        $lazarus_1 = "WannaCry" ascii
        $lazarus_2 = "HOPLIGHT" ascii
        $emotet_1 = "Emotet" ascii
        $cobalt_1 = "cobaltstrike" ascii
        $cobalt_2 = "ReflectiveLoader" ascii
        
    condition:
        2 of them
}
`
}

// EmbedFakeC2Communication simula comunicación con C2 (solo strings, no real)
func EmbedFakeC2Communication() []string {
	return []string{
		"POST /api/v1/beacon HTTP/1.1",
		"Host: microsoft-update.net",
		"User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		"Content-Type: application/octet-stream",
		"X-Session-ID: {12345678-1234-1234-1234-123456789012}",
	}
}

// GenerateFakeMutexes crea nombres de mutex similares a APTs
func GenerateFakeMutexes() []string {
	return []string{
		"Global\\{A5F6D7E8-9B0C-1D2E-3F4A-5B6C7D8E9F0A}",
		"Global\\SM0:1234:304:WilStaging_02",
		"Local\\MsWinZonesCacheCounterMutexA",
	}
}

// GetPEMetadata retorna metadata del PE que parece de APT
func GetPEMetadata() map[string]string {
	return map[string]string{
		"CompanyName":      "Microsoft Corporation",
		"FileDescription":  "Windows Security Update",
		"FileVersion":      "10.0.19041.1234",
		"InternalName":     "SecurityUpdate.exe",
		"LegalCopyright":   "© Microsoft Corporation. All rights reserved.",
		"OriginalFilename": "SecurityUpdate.exe",
		"ProductName":      "Microsoft® Windows® Operating System",
		"ProductVersion":   "10.0.19041.1234",
	}
}

// CalculateFileHash calcula el hash del binario actual
func CalculateFileHash() (string, string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", "", err
	}

	data, err := os.ReadFile(exePath)
	if err != nil {
		return "", "", err
	}

	// MD5
	md5Hash := md5.Sum(data)
	md5Str := hex.EncodeToString(md5Hash[:])

	// SHA256
	sha256Hash := sha256.Sum256(data)
	sha256Str := hex.EncodeToString(sha256Hash[:])

	return md5Str, sha256Str, nil
}

// PrintAPTInfo imprime información de las firmas embebidas
func PrintAPTInfo() {
}

// GetIOCs retorna Indicators of Compromise para documentación
func GetIOCs() map[string][]string {
	return map[string][]string{
		"MD5": {
			"d0d626deb3f9484e649294a8dfa814c5",
			"a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
		},
		"SHA256": {
			"ce77d116a074dab7a22a0fd4f2c1ab475f16eec42e1ded3c0b0aa8211fe858d6",
			"32519b85c0b422e4656de6e6c41878e95fd95026267daab4215ee59c107d6c77",
		},
		"Domains": {
			"microsoft-update.net",
			"windows-security.org",
			"avsvmcloud.com",
		},
		"IPs": {
			"185.184.25.78",
			"190.90.233.66",
		},
		"Mutexes": {
			"Global\\{A5F6D7E8-9B0C-1D2E-3F4A-5B6C7D8E9F0A}",
			"Global\\SM0:1234:304:WilStaging_02",
		},
	}
}
