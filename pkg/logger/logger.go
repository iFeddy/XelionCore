package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// InitSessionLogger inicializa la persistencia de logs de sesión en la carpeta logs/.
// Crea el archivo con el formato requerido: <serviceName>_YYYYMMDD_HHMM.log
// y configura el logger estándar para escribir concurrentemente a os.Stdout y al archivo de sesión.
func InitSessionLogger(serviceName string) (*os.File, error) {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_1504")
	fileName := fmt.Sprintf("%s_%s.log", serviceName, timestamp)
	filePath := filepath.Join(logDir, fileName)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file '%s': %w", filePath, err)
	}

	multi := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multi)

	return file, nil
}
