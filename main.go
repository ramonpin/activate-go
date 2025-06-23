package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// --- Modelo de Bubble Tea para el selector interactivo ---
type item struct {
	path string
}

func (i item) FilterValue() string { return i.path }
func (i item) Title() string       { return filepath.Base(i.path) }
func (i item) Description() string { return i.path }

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			selectedItem, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = selectedItem.path
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" || m.quitting {
		return ""
	}
	return m.list.View()
}

func main() {
	os.Setenv("CLICOLOR_FORCE", "1")

	if len(os.Args) > 1 {
		handleDirectPath(os.Args[1])
	} else {
		handleSearch()
	}
}

func handleDirectPath(path string) {
	if isVenv(path) {
		printCommand(path)
	} else {
		fmt.Fprintf(os.Stderr, "Error: Directory '%s' is not a valid virtual environment.\n", path)
		os.Exit(1)
	}
}

func handleSearch() {
	searchDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	for {
		venvs := findVenvsInDir(searchDir)

		if len(venvs) == 1 {
			printCommand(venvs[0])
			return
		}

		if len(venvs) > 1 {
			items := []list.Item{}
			for _, venv := range venvs {
				items = append(items, item{path: venv})
			}

			l := list.New(items, list.NewDefaultDelegate(), 0, 0)
			l.Title = "Multiple environments found. Please choose one:"
			l.SetShowStatusBar(false)
			l.SetShowPagination(true)

			m := model{list: l}
			p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithOutput(os.Stderr))

			finalModel, err := p.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running selection UI: %v\n", err)
				os.Exit(1)
			}

			if chosenPath := finalModel.(model).choice; chosenPath != "" {
				printCommand(chosenPath)
			}
			return
		}

		parentDir := filepath.Dir(searchDir)
		if parentDir == searchDir {
			fmt.Fprintln(os.Stderr, "No virtual environment found in this directory or parent directories.")
			os.Exit(1)
		}
		searchDir = parentDir
	}
}

func findVenvsInDir(dir string) []string {
	var venvs []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	for _, entry := range entries {
		if entry.IsDir() {
			fullPath := filepath.Join(dir, entry.Name())
			if isVenv(fullPath) {
				venvs = append(venvs, fullPath)
			}
		}
	}
	return venvs
}

func isVenv(path string) bool {
	activateScriptPath := filepath.Join(path, "bin", "activate")
	info, err := os.Stat(activateScriptPath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func printCommand(path string) {
	fmt.Fprintf(os.Stderr, "Activating environment %s...\n", path)
	fmt.Printf("source %s", filepath.Join(path, "bin", "activate"))
}
