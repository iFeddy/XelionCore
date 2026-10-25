package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type IndexDef struct {
	Name  string
	Title string
	Dirs  []string
	Note  string
}

func main() {
	indices := []IndexDef{
		{"Data", "pkg/data (cargadores de datos)", []string{"pkg/data", "pkg/script"}, "Datos JSON en data/."},
		{"Zone", "ZoneService", []string{"pkg/server/zone", "cmd/ZoneService"}, "Protocolos en pkg/protocol."},
		{"Character", "WorldService / Character", []string{"pkg/server/world", "pkg/database", "cmd/WorldService"}, "Equivale a WorldManager+Character."},
		{"Protocol", "pkg/protocol", []string{"pkg/protocol"}, ""},
		{"Login", "LoginService", []string{"pkg/server/login", "cmd/LoginService"}, ""},
		{"Core", "config/crypto/logger/network", []string{"pkg/config", "pkg/crypto", "pkg/logger", "pkg/network"}, ""},
	}

	_ = os.MkdirAll("docs/Indices", 0755)
	sigRegex := regexp.MustCompile(`^(func|type)\s`)
	trimBrace := regexp.MustCompile(`\s*\{\s*$`)

	for _, idx := range indices {
		var out []string
		out = append(out, fmt.Sprintf("# Índice: %s", idx.Title), "")
		out = append(out, fmt.Sprintf("> Generado automáticamente por scripts/gen_indices. Solo firmas con n° de línea: abrir el archivo con StartLine/EndLine. No reemplaza a docs/Servicios/. %s", idx.Note), "")

		var files []string
		for _, dir := range idx.Dirs {
			_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				if err == nil && !d.IsDir() && strings.HasSuffix(d.Name(), ".go") {
					files = append(files, filepath.ToSlash(path))
				}
				return nil
			})
		}
		sort.Strings(files)

		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				continue
			}
			out = append(out, fmt.Sprintf("## %s", file))
			scanner := bufio.NewScanner(f)
			lineNum := 0
			for scanner.Scan() {
				lineNum++
				line := scanner.Text()
				if sigRegex.MatchString(line) {
					s := trimBrace.ReplaceAllString(line, "")
					if len(s) > 140 {
						s = s[:140]
					}
					out = append(out, fmt.Sprintf("- L%d: `%s`", lineNum, s))
				}
			}
			_ = f.Close()
			out = append(out, "")
		}

		targetPath := filepath.Join("docs/Indices", idx.Name+".md")
		content := strings.Join(out, "\n")
		_ = os.WriteFile(targetPath, []byte(content), 0644)
		fmt.Printf("%s.md generado (%d bytes)\n", idx.Name, len(content))
	}
}
