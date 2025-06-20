package main

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
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
