package collectors

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/m4rcmo15/artefacto/models"
)

// DetectTools detecta herramientas de análisis, debugging y reversing
func DetectTools() *models.ToolsInfo {
	info := &models.ToolsInfo{
		ReversingTools:      []string{},
		DebuggingTools:      []string{},
		MonitoringTools:     []string{},
		VirtualizationTools: []string{},
		AnalysisTools:       []string{},
	}

	// Herramientas de reversing
	reversingTools := map[string]string{
		"ida.exe":        "IDA Pro",
		"ida64.exe":      "IDA Pro 64-bit",
		"idaq.exe":       "IDA Pro",
		"idaq64.exe":     "IDA Pro 64-bit",
		"ghidra.exe":     "Ghidra",
		"radare2.exe":    "Radare2",
		"r2.exe":         "Radare2",
		"cutter.exe":     "Cutter",
		"hopper.exe":     "Hopper Disassembler",
		"binary ninja.exe": "Binary Ninja",
		"pe-bear.exe":    "PE-bear",
		"pestudio.exe":   "PEStudio",
		"die.exe":        "Detect It Easy",
		"exeinfope.exe":  "ExeinfoPE",
	}

	// Herramientas de debugging
	debuggingTools := map[string]string{
		"x64dbg.exe":     "x64dbg",
		"x32dbg.exe":     "x32dbg",
		"windbg.exe":     "WinDbg",
		"ollydbg.exe":    "OllyDbg",
		"immunity debugger.exe": "Immunity Debugger",
		"gdb.exe":        "GDB",
		"devenv.exe":     "Visual Studio",
		"dnspy.exe":      "dnSpy",
		"dotpeek.exe":    "dotPeek",
	}

	// Herramientas de monitoreo
	monitoringTools := map[string]string{
		"procmon.exe":    "Process Monitor",
		"procmon64.exe":  "Process Monitor 64",
		"procexp.exe":    "Process Explorer",
		"procexp64.exe":  "Process Explorer 64",
		"wireshark.exe":  "Wireshark",
		"fiddler.exe":    "Fiddler",
		"tcpview.exe":    "TCPView",
		"autoruns.exe":   "Autoruns",
		"regshot.exe":    "Regshot",
		"apimonitor.exe": "API Monitor",
	}

	// Herramientas de virtualización
	virtualizationTools := map[string]string{
		"vmware.exe":         "VMware",
		"vmware-vmx.exe":     "VMware",
		"virtualbox.exe":     "VirtualBox",
		"vboxmanage.exe":     "VirtualBox",
		"qemu.exe":           "QEMU",
		"qemu-system-x86_64.exe": "QEMU",
		"hyper-v.exe":        "Hyper-V",
	}

	// Herramientas de análisis
	analysisTools := map[string]string{
		"sandboxie.exe":  "Sandboxie",
		"cuckoo.exe":     "Cuckoo Sandbox",
		"regshot.exe":    "Regshot",
		"fakenet.exe":    "FakeNet",
		"inetsim.exe":    "INetSim",
	}

	// Buscar en procesos
	processes := getProcessList()
	for _, proc := range processes {
		procNameLower := strings.ToLower(proc.Name)
		
		if tool, found := reversingTools[procNameLower]; found {
			info.ReversingTools = append(info.ReversingTools, tool)
		}
		if tool, found := debuggingTools[procNameLower]; found {
			info.DebuggingTools = append(info.DebuggingTools, tool)
		}
		if tool, found := monitoringTools[procNameLower]; found {
			info.MonitoringTools = append(info.MonitoringTools, tool)
		}
		if tool, found := virtualizationTools[procNameLower]; found {
			info.VirtualizationTools = append(info.VirtualizationTools, tool)
		}
		if tool, found := analysisTools[procNameLower]; found {
			info.AnalysisTools = append(info.AnalysisTools, tool)
		}
	}

	// Buscar en directorios comunes
	searchPaths := []string{
		"C:\\Program Files",
		"C:\\Program Files (x86)",
		os.Getenv("USERPROFILE") + "\\Desktop",
		os.Getenv("USERPROFILE") + "\\Downloads",
	}

	for _, basePath := range searchPaths {
		checkToolsInPath(basePath, reversingTools, &info.ReversingTools)
		checkToolsInPath(basePath, debuggingTools, &info.DebuggingTools)
		checkToolsInPath(basePath, monitoringTools, &info.MonitoringTools)
		checkToolsInPath(basePath, virtualizationTools, &info.VirtualizationTools)
		checkToolsInPath(basePath, analysisTools, &info.AnalysisTools)
	}

	// Eliminar duplicados
	info.ReversingTools = removeDuplicates(info.ReversingTools)
	info.DebuggingTools = removeDuplicates(info.DebuggingTools)
	info.MonitoringTools = removeDuplicates(info.MonitoringTools)
	info.VirtualizationTools = removeDuplicates(info.VirtualizationTools)
	info.AnalysisTools = removeDuplicates(info.AnalysisTools)

	return info
}

func checkToolsInPath(basePath string, tools map[string]string, found *[]string) {
	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		
		fileNameLower := strings.ToLower(info.Name())
		if tool, exists := tools[fileNameLower]; exists {
			*found = append(*found, tool)
		}
		
		return nil
	})
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if _, exists := keys[item]; !exists {
			keys[item] = true
			result = append(result, item)
		}
	}
	
	return result
}
