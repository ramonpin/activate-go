# activate-go: A Smart Virtual Environment Activator

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A command-line tool written in Go to quickly find and help activate Python
virtual environments (`virtualenv`, `venv`). It's a modern, self-contained
replacement for complex shell scripts.

## Demo

**Use Case 1: Automatic Search.** When run with no arguments, it searches
parent directories automatically. If multiple environments are found, it
displays an interactive selector for you to choose.

```plain
$ activate
? Multiple environments found. Please choose one:

sample
sample2
.hidden_venv
(q) quit
```

**Use Case 2: Direct Activation.** If a path is provided, it activates the
specified environment directly (if valid).

```plain
activate path/to/my-env
```

## Features

* **Automatic Search:** Scans the current and parent directories for virtual
  environments.
* **Direct Activation:** Activate a specific environment by providing its path
  as an argument.
* **Interactive Selector:** If multiple environments are found during a search,
  a TUI (Terminal User Interface) is displayed to let you choose.
* **Self-Contained:** Compiles to a single, dependency-free executable.
* **Fast:** Written in Go for optimal performance.

## Installation

### Prerequisites

* Go toolchain (version 1.18 or higher) installed.

### Steps

1. **Clone the repository** (or simply place the code in a folder):

    ```bash
    # Replace the URL with your repository's URL if you upload it to GitHub
    git clone [https://github.com/your-username/activate-go.git](https://github.com/your-username/activate-go.git)
    cd activate-go
    ```

2. **Build the executable**:

    ```bash
    go build
    ```

    This will create a binary file named `activate-go` in the current directory.

3. **Move the executable to your PATH**: To call it from anywhere, move the binary to a directory included in your `$PATH`.

    ```bash
    # Example:
    mv activate-go /usr/local/bin/
    ```

## Usage

**Important!** A child process (like this tool) cannot modify the environment of its parent shell. Therefore, `activate-go` does not activate the environment directly. Instead, it prints the `source` command that the shell must execute.

To make it work seamlessly, you must add a small wrapper function to your shell's configuration file (`~/.zshrc` for Zsh or `~/.bashrc` for Bash).

1. Add the following function to your `~/.zshrc`:

    ```zsh
    # ~/.zshrc

    function activate() {
      # Call the Go executable and capture its text output
      local command_to_run
      command_to_run=$(activate-go "$@")
      
      # If the command was successful, use `eval` to execute its output in the current shell
      if [ $? -eq 0 ]; then
        eval "$command_to_run"
      fi
    }
    ```

2. Reload your shell configuration:

    ```bash
    source ~/.zshrc
    ```

3. Done! You can now use the `activate` command:

    ```bash
    # Automatic upward search
    activate

    # Activate a specific environment by path
    activate ~/projects/my-virtual-env
    ```

## Development

This project is written in Go and uses Go Modules for dependency management.

* **Running Tests**: To ensure everything is working correctly, you can run the
  test suite:

    ```bash
    go test -v
    ```

### Dependencies

#### Runtime Dependencies

* None! The compiled binary is self-contained.

#### Development Dependencies

* [Go](https://golang.org/) (1.18+)
* [github.com/charmbracelet/bubbles](https://github.com/charmbracelet/bubbles)
* [github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

