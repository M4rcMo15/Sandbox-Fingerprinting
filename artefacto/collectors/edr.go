package collectors

import (
	"os"
	"strings"
	"unsafe"

	"github.com/m4rcmo15/artefacto/models"
	"golang.org/x/sys/windows"
)

// DetectEDR identifica productos de seguridad instalados
func DetectEDR() *models.EDRInfo {
	info := &models.EDRInfo{
		DetectedProducts: []models.EDRProduct{},
		RunningProcesses: []string{},
		InstalledDrivers: []string{},
	}

	// Lista de productos EDR/AV conocidos
	edrProducts := []struct {
		name      string
		processes []string
		drivers   []string
		regKeys   []string
	}{
		{
			name:      "Windows Defender",
			processes: []string{"MsMpEng.exe", "NisSrv.exe", "SecurityHealthService.exe"},
			drivers:   []string{"WdFilter.sys", "WdNisDrv.sys"},
		},
		{
			name:      "CrowdStrike Falcon",
			processes: []string{"CSFalconService.exe", "CSFalconContainer.exe"},
			drivers:   []string{"csagent.sys", "csdevicecontrol.sys"},
		},
		{
			name:      "SentinelOne",
			processes: []string{"SentinelAgent.exe", "SentinelServiceHost.exe"},
			drivers:   []string{"SentinelMonitor.sys"},
		},
		{
			name:      "Carbon Black",
			processes: []string{"cb.exe", "RepMgr.exe", "RepUtils.exe"},
			drivers:   []string{"cbk7.sys", "parity.sys"},
		},
		{
			name:      "Cylance",
			processes: []string{"CylanceSvc.exe", "CylanceUI.exe"},
			drivers:   []string{"CylanceDrv.sys"},
		},
		{
			name:      "Symantec Endpoint Protection",
			processes: []string{"ccSvcHst.exe", "Smc.exe"},
			drivers:   []string{"SRTSP.sys", "SymEFA.sys"},
		},
		{
			name:      "McAfee Endpoint Security",
			processes: []string{"mfemms.exe", "mfevtps.exe"},
			drivers:   []string{"mfehidk.sys", "mfefirek.sys"},
		},
		{
			name:      "Kaspersky",
			processes: []string{"avp.exe", "avpui.exe"},
			drivers:   []string{"klif.sys", "kl1.sys"},
		},
		{
			name:      "Trend Micro",
			processes: []string{"TMBMSRV.exe", "TMBMServer.exe"},
			drivers:   []string{"tmcomm.sys", "tmactmon.sys"},
		},
		{
			name:      "ESET",
			processes: []string{"ekrn.exe", "egui.exe"},
			drivers:   []string{"eamonm.sys", "ehdrv.sys"},
		},
		{
			name:      "Palo Alto Traps",
			processes: []string{"cyserver.exe", "cyvera.exe"},
			drivers:   []string{"tlaworker.sys"},
		},
		{
			name:      "FireEye",
			processes: []string{"xagt.exe", "xagtnotif.exe"},
			drivers:   []string{"xagt.sys"},
		},
	}

	// Obtener procesos en ejecución
	runningProcesses := getRunningProcessNames()
	info.RunningProcesses = runningProcesses

	// Obtener drivers instalados
	installedDrivers := getInstalledDrivers()
	info.InstalledDrivers = installedDrivers

	// Verificar cada producto EDR
	for _, edr := range edrProducts {
		product := models.EDRProduct{
			Name:     edr.name,
			Type:     "EDR/AV",
			Detected: false,
		}

		// Verificar procesos
		for _, proc := range edr.processes {
			if containsIgnoreCase(runningProcesses, proc) {
				product.Detected = true
				product.Method = "process"
				break
			}
		}

		// Verificar drivers
		if !product.Detected {
			for _, driver := range edr.drivers {
				if containsIgnoreCase(installedDrivers, driver) {
					product.Detected = true
					product.Method = "driver"
					break
				}
			}
		}

		if product.Detected {
			info.DetectedProducts = append(info.DetectedProducts, product)
		}
	}

	return info
}

func getRunningProcessNames() []string {
	processes := []string{}

	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return processes
	}
	defer windows.CloseHandle(snapshot)

	var procEntry windows.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))

	err = windows.Process32First(snapshot, &procEntry)
	if err != nil {
		return processes
	}

	for {
		name := windows.UTF16ToString(procEntry.ExeFile[:])
		processes = append(processes, name)

		err = windows.Process32Next(snapshot, &procEntry)
		if err != nil {
			break
		}
	}

	return processes
}

func getInstalledDrivers() []string {
	drivers := []string{}

	// Enumerar drivers en C:\Windows\System32\drivers
	driverPath := "C:\\Windows\\System32\\drivers"
	
	// Leer archivos .sys del directorio
	entries, err := os.ReadDir(driverPath)
	if err != nil {
		return drivers
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".sys") {
			drivers = append(drivers, entry.Name())
		}
	}
	
	return drivers
}

func containsIgnoreCase(slice []string, item string) bool {
	itemLower := strings.ToLower(item)
	for _, s := range slice {
		if strings.ToLower(s) == itemLower {
			return true
		}
	}
	return false
}
