package main

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestIsVenv(t *testing.T) {
	tempDir := t.TempDir()

	// --- Caso 1: Un entorno virtual VÁLIDO ---
	validVenvPath := filepath.Join(tempDir, "valid_venv")
	os.MkdirAll(filepath.Join(validVenvPath, "bin"), 0755)
	os.WriteFile(filepath.Join(validVenvPath, "bin", "activate"), []byte(""), 0644)

	if !isVenv(validVenvPath) {
		t.Errorf("isVenv() failed to identify a valid virtual environment at %s", validVenvPath)
	}

	// --- Caso 2: Un directorio que NO es un entorno virtual ---
	notVenvPath := filepath.Join(tempDir, "not_a_venv")
	os.Mkdir(notVenvPath, 0755)

	if isVenv(notVenvPath) {
		t.Errorf("isVenv() incorrectly identified a plain directory as a virtual environment")
	}

	// --- Caso 3: Un directorio con 'bin' pero sin 'activate' ---
	noActivatePath := filepath.Join(tempDir, "no_activate_venv")
	os.MkdirAll(filepath.Join(noActivatePath, "bin"), 0755)

	if isVenv(noActivatePath) {
		t.Errorf("isVenv() incorrectly identified a directory with only a 'bin' folder as a venv")
	}

	// --- Caso 4: Una ruta que no existe ---
	if isVenv("path/that/does/not/exist") {
		t.Errorf("isVenv() should return false for a non-existent path")
	}
}

func TestFindVenvsInDir(t *testing.T) {
	tempDir := t.TempDir()

	// Creamos una estructura de directorios para el test.
	// Dos entornos válidos, uno oculto y un directorio normal.
	venv1 := filepath.Join(tempDir, "venv1")
	os.MkdirAll(filepath.Join(venv1, "bin"), 0755)
	os.WriteFile(filepath.Join(venv1, "bin", "activate"), []byte(""), 0644)

	venv2 := filepath.Join(tempDir, ".hidden_venv") // Un venv oculto
	os.MkdirAll(filepath.Join(venv2, "bin"), 0755)
	os.WriteFile(filepath.Join(venv2, "bin", "activate"), []byte(""), 0644)

	os.Mkdir(filepath.Join(tempDir, "not_a_venv"), 0755) // Directorio normal

	// Llamamos a la función que queremos testear.
	foundVenvs := findVenvsInDir(tempDir)

	// Definimos el resultado que esperamos.
	expectedVenvs := []string{venv1, venv2}

	sort.Strings(foundVenvs)
	sort.Strings(expectedVenvs)

	if !reflect.DeepEqual(foundVenvs, expectedVenvs) {
		t.Errorf("findVenvsInDir() failed.\nExpected: %v\nGot:      %v", expectedVenvs, foundVenvs)
	}
}

func TestShellEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple path - no escaping needed",
			input:    "/home/user/venv",
			expected: "/home/user/venv",
		},
		{
			name:     "path with hyphen - no escaping needed",
			input:    "/home/user/my-venv",
			expected: "/home/user/my-venv",
		},
		{
			name:     "path with underscore - no escaping needed",
			input:    "/home/user/my_venv",
			expected: "/home/user/my_venv",
		},
		{
			name:     "path with space",
			input:    "/home/user/my venv",
			expected: "'/home/user/my venv'",
		},
		{
			name:     "path with single quote",
			input:    "/home/user/it's venv",
			expected: "'/home/user/it'\"'\"'s venv'",
		},
		{
			name:     "injection: semicolon command",
			input:    "/path; rm -rf /",
			expected: "'/path; rm -rf /'",
		},
		{
			name:     "injection: backtick command substitution",
			input:    "/path/`whoami`/venv",
			expected: "'/path/`whoami`/venv'",
		},
		{
			name:     "injection: dollar command substitution",
			input:    "/path/$(rm -rf /)/venv",
			expected: "'/path/$(rm -rf /)/venv'",
		},
		{
			name:     "injection: pipe",
			input:    "/path | cat /etc/passwd",
			expected: "'/path | cat /etc/passwd'",
		},
		{
			name:     "injection: ampersand background",
			input:    "/path & echo hacked &",
			expected: "'/path & echo hacked &'",
		},
		{
			name:     "injection: redirect",
			input:    "/path > /tmp/hacked",
			expected: "'/path > /tmp/hacked'",
		},
		{
			name:     "injection: variable expansion",
			input:    "/path/$HOME/venv",
			expected: "'/path/$HOME/venv'",
		},
		{
			name:     "multiple single quotes",
			input:    "it's john's venv",
			expected: "'it'\"'\"'s john'\"'\"'s venv'",
		},
		{
			name:     "injection: newline",
			input:    "/path\nrm -rf /",
			expected: "'/path\nrm -rf /'",
		},
		{
			name:     "injection: tab character",
			input:    "/path/venv\tmalicious",
			expected: "'/path/venv\tmalicious'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shellEscape(tt.input)
			if result != tt.expected {
				t.Errorf("shellEscape(%q)\n  got:  %q\n  want: %q",
					tt.input, result, tt.expected)
			}
		})
	}
}

// TestShellEscapeIntegration verifica que rutas con nombres peligrosos se escapan correctamente
func TestShellEscapeIntegration(t *testing.T) {
	tempDir := t.TempDir()

	// Crear venv con nombre potencialmente peligroso
	dangerousName := "venv; echo HACKED > /tmp/pwned"
	venvPath := filepath.Join(tempDir, dangerousName)
	os.MkdirAll(filepath.Join(venvPath, "bin"), 0755)
	activateScript := filepath.Join(venvPath, "bin", "activate")
	os.WriteFile(activateScript, []byte("#!/bin/bash\necho SAFE"), 0755)

	// Generar comando escapado
	escaped := shellEscape(activateScript)

	// Verificar que está entre comillas
	if !strings.HasPrefix(escaped, "'") || !strings.HasSuffix(escaped, "'") {
		t.Errorf("Expected dangerous path to be quoted, got: %s", escaped)
	}

	// Verificar que el comando completo contiene el nombre peligroso escapado
	// El punto y coma debe estar dentro de las comillas (literal, no ejecutable)
	if !strings.Contains(escaped, "venv; echo HACKED") {
		t.Errorf("Full dangerous path should be present in escaped output, got: %s", escaped)
	}

	// Verificar que no hay comillas sueltas que permitan inyección
	// Todo entre la primera y última comilla debe ser literal
	if strings.Count(escaped, "'") > 2 {
		// Puede tener más comillas si hay comillas simples en el path,
		// pero deben estar correctamente escapadas con '"'"'
		// Verificar que las comillas extras están en el patrón de escape
		if !strings.Contains(escaped, "'\"'\"'") {
			t.Errorf("Path contains unescaped quotes, got: %s", escaped)
		}
	}
}
