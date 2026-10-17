package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// WorldConfig representa la configuración de un mundo/servidor disponible en el Login.
type WorldConfig struct {
	ID   uint8
	Name string
	IP   string
	Port uint16
}

// Config contiene todas las variables de configuración de LoginService.
type Config struct {
	// Red de Login
	LoginHost     string
	LoginPort     uint16
	ClientVersion string

	// Base de Datos MySQL
	DBHost     string
	DBPort     uint16
	DBUser     string
	DBPassword string
	DBName     string

	// Mundos configurados
	Worlds []WorldConfig
}

// Load lee el archivo .env (si existe) y las variables de entorno del sistema.
// Si .env no existe, utiliza los valores por defecto o las variables del entorno.
func Load(envPath ...string) (*Config, error) {
	filePath := ".env"
	if len(envPath) > 0 && envPath[0] != "" {
		filePath = envPath[0]
	}

	// Cargar archivo .env en el mapa de entorno si existe
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

	// Función auxiliar para obtener un valor del .env o del SO o un default
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
		LoginHost:     getEnv("LOGIN_HOST", "0.0.0.0"),
		LoginPort:     getEnvUint16("LOGIN_PORT", 9010),
		ClientVersion: getEnv("CLIENT_VERSION", "10022024000000"),

		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnvUint16("DB_PORT", 3306),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", getEnv("DB_DATABASE", "xelion_online")),
	}

	// Cargar mundos indexados dinámicamente: WORLD_0_NAME, WORLD_1_NAME, etc.
	for i := 0; i < 255; i++ {
		nameKey := fmt.Sprintf("WORLD_%d_NAME", i)
		name := getEnv(nameKey, "")
		if name == "" {
			if i == 0 {
				legacyName := getEnv("WORLD_NAME", "")
				if legacyName != "" {
					cfg.Worlds = append(cfg.Worlds, WorldConfig{
						ID:   0,
						Name: legacyName,
						IP:   getEnv("WORLD_IP", "127.0.0.1"),
						Port: getEnvUint16("WORLD_PORT", 9014),
					})
				}
			}
			break
		}

		ipKey := fmt.Sprintf("WORLD_%d_IP", i)
		portKey := fmt.Sprintf("WORLD_%d_PORT", i)

		cfg.Worlds = append(cfg.Worlds, WorldConfig{
			ID:   uint8(i),
			Name: name,
			IP:   getEnv(ipKey, "127.0.0.1"),
			Port: getEnvUint16(portKey, 9014),
		})
	}

	// Si no se configuró ningún mundo, registrar al menos un mundo por defecto
	if len(cfg.Worlds) == 0 {
		cfg.Worlds = append(cfg.Worlds, WorldConfig{
			ID:   0,
			Name: "Xelion",
			IP:   "127.0.0.1",
			Port: 9014,
		})
	}

	return cfg, nil
}
