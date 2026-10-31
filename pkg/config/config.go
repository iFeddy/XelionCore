package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config contiene las variables de infraestructura y conexión de XelionCore.
// Nota: La versión de cliente (FOUS.GameVersionKey) y el catálogo de mundos (tWorldServer)
// se leen obligatoriamente desde MySQL.
type Config struct {
	// Red de Login
	LoginHost string
	LoginPort uint16

	// Base de Datos MySQL
	DBHost     string
	DBPort     uint16
	DBUser     string
	DBPassword string
	DBName     string

	// Red de ZoneService
	ZoneIP   string
	ZonePort uint16
}

// Load lee el archivo .env (si existe) y las variables de entorno del sistema.
// Si .env no existe, utiliza los valores por defecto o las variables del entorno.
func Load(envPath ...string) (*Config, error) {
	filePath := ".env"
	if len(envPath) > 0 && envPath[0] != "" {
		filePath = envPath[0]
	}

	envMap := make(map[string]string)
	if f, err := os.Open(filePath); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				val = strings.Trim(val, "\"'")
				envMap[key] = val
			}
		}
	}

	getEnv := func(key, defaultVal string) string {
		if val, exists := envMap[key]; exists && val != "" {
			return val
		}
		if val := os.Getenv(key); val != "" {
			return val
		}
		return defaultVal
	}

	getEnvUint16 := func(key string, defaultVal uint16) uint16 {
		s := getEnv(key, "")
		if s == "" {
			return defaultVal
		}
		if v, err := strconv.ParseUint(s, 10, 16); err == nil {
			return uint16(v)
		}
		return defaultVal
	}

	cfg := &Config{
		LoginHost: getEnv("LOGIN_HOST", "0.0.0.0"),
		LoginPort: getEnvUint16("LOGIN_PORT", 9010),

		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnvUint16("DB_PORT", 3306),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", getEnv("DB_DATABASE", "xelion_online")),

		ZoneIP:   getEnv("ZONE_IP", "127.0.0.1"),
		ZonePort: getEnvUint16("ZONE_PORT", 9210),
	}

	return cfg, nil
}

// DSN genera la cadena de conexión canónica a MySQL garantizando soporte de
// parseTime, multiStatements y charset utf8mb4 requeridos por los Stored Procedures.
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true&charset=utf8mb4",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
