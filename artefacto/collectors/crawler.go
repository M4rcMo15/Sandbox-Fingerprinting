package collectors

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/m4rcmo15/artefacto/models"
)

// CrawlFiles busca archivos específicos en el sistema
func CrawlFiles(patterns []string, maxFiles int) *models.CrawlerInfo {
	info := &models.CrawlerInfo{
		ScannedPaths: []string{},
		FoundFiles:   []string{},
		TotalFiles:   0,
	}

	// Obtener unidades montadas
	drives := getLogicalDrives()

	for _, drive := range drives {
		if info.TotalFiles >= maxFiles {
			break
		}
		
		info.ScannedPaths = append(info.ScannedPaths, drive)
		crawlDirectory(drive, patterns, info, maxFiles)
	}

	return info
}

func getLogicalDrives() []string {
	drives := []string{}
	
	// En Windows, verificar de A: a Z:
	for i := 'A'; i <= 'Z'; i++ {
		drive := string(i) + ":\\"
		if _, err := os.Stat(drive); err == nil {
			drives = append(drives, drive)
		}
	}

	return drives
}

func crawlDirectory(root string, patterns []string, info *models.CrawlerInfo, maxFiles int) {
	// Directorios a evitar para no perder tiempo
	skipDirs := map[string]bool{
		"Windows":         true,
		"Program Files":   true,
		"Program Files (x86)": true,
		"$Recycle.Bin":   true,
		"System Volume Information": true,
	}

	filepath.Walk(root, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return nil // Ignorar errores y continuar
		}

		if info.TotalFiles >= maxFiles {
			return filepath.SkipDir
		}

		// Saltar directorios del sistema
		if fileInfo.IsDir() {
			if skipDirs[fileInfo.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		// Verificar si el archivo coincide con algún patrón
		for _, pattern := range patterns {
			if matchPattern(fileInfo.Name(), pattern) {
				info.FoundFiles = append(info.FoundFiles, path)
				info.TotalFiles++
				break
			}
		}

		return nil
	})
}

func matchPattern(filename, pattern string) bool {
	// Búsqueda simple por extensión o nombre
	filename = strings.ToLower(filename)
	pattern = strings.ToLower(pattern)

	if strings.HasPrefix(pattern, "*.") {
		// Patrón de extensión
		ext := strings.TrimPrefix(pattern, "*")
		return strings.HasSuffix(filename, ext)
	}

	// Búsqueda por substring
	return strings.Contains(filename, pattern)
}
