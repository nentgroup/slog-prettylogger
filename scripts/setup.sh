#!/bin/bash

echo "Setting up lefthook, commitlint, golangci-lint, and Task..."

command_exists () {
    type "$1" &> /dev/null ;
}

# Determine OS
OS="$(uname -s)"
case "${OS}" in
  Linux*)
    echo "Detected Linux"
    ;;
  Darwin*)
    echo "Detected macOS"
    ;;
  MINGW*|MSYS*|CYGWIN*)
    echo "Detected Windows"
    ;;
  *)
    echo "Unsupported OS: ${OS}"
    echo "Please install lefthook, commitlint, and golangci-lint manually."
    exit 1
    ;;
esac

# Initialize status variables
lefthook_installed=false
commitlint_installed=false
golangci_lint_installed=false
task_installed=false

# Check if lefthook is already installed
if ! command_exists lefthook; then
  echo "Installing lefthook..."

  case "${OS}" in
    Linux*|MINGW*|MSYS*|CYGWIN*)
      if command_exists go; then
        go install github.com/evilmartians/lefthook@latest
        lefthook_installed=true
      else
        echo "Please install Go or lefthook manually: https://lefthook.dev/installation/"
        exit 1
      fi
      ;;
    Darwin*)
      if command_exists brew; then
        brew install lefthook
        lefthook_installed=true
      elif command_exists go; then
        go install github.com/evilmartians/lefthook@latest
        lefthook_installed=true
      else
        echo "Please install Homebrew or Go, then run this again"
        exit 1
      fi
      ;;
  esac
else
  echo "✅ lefthook is already installed"
fi

# Check if commitlint is already installed
if ! command_exists commitlint; then
  echo "Installing commitlint..."

  if command_exists go; then
    go install github.com/conventionalcommit/commitlint@latest
    commitlint_installed=true
  else
    echo "Please install Go or commitlint manually: https://github.com/conventionalcommit/commitlint"
    exit 1
  fi
else
  echo "✅ commitlint is already installed"
fi

# Check if golangci-lint is already installed
if ! command_exists golangci-lint; then
  echo "Installing golangci-lint v2..."

  case "${OS}" in
    Linux*|MINGW*|MSYS*|CYGWIN*)
      if command_exists curl; then
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)/bin" v2.0.0
        golangci_lint_installed=true
      else
        echo "Please install curl or golangci-lint manually: https://golangci-lint.run/usage/install/"
        exit 1
      fi
      ;;
    Darwin*)
      if command_exists brew; then
        brew install golangci-lint
        brew upgrade golangci-lint
        golangci_lint_installed=true
      elif command_exists curl; then
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)/bin" v2.0.0
        golangci_lint_installed=true
      else
        echo "Please install Homebrew or curl, then run this again"
        exit 1
      fi
      ;;
  esac
else
  echo "Checking golangci-lint version..."
  VERSION=$(golangci-lint --version | awk '{print $4}')
  MAJOR_VERSION=$(echo $VERSION | cut -d. -f1)

  if [ "$MAJOR_VERSION" -lt 2 ]; then
    echo "Upgrading golangci-lint to v2..."
    case "${OS}" in
      Linux*|MINGW*|MSYS*|CYGWIN*)
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)/bin" v2.0.0
        ;;
      Darwin*)
        if command_exists brew; then
          brew upgrade golangci-lint
        else
          curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)/bin" v2.0.0
        fi
        ;;
    esac
    golangci_lint_installed=true
  else
    echo "✅ golangci-lint v2 is already installed"
  fi
fi

# Check if Task is already installed
if ! command_exists task; then
  echo "Installing Task..."

  case "${OS}" in
    Linux*|MINGW*|MSYS*|CYGWIN*)
      if command_exists curl; then
        sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b "$(go env GOPATH)/bin"
        task_installed=true
      else
        echo "Please install curl or Task manually: https://taskfile.dev/installation/"
        exit 1
      fi
      ;;
    Darwin*)
      if command_exists brew; then
        brew install go-task
        task_installed=true
      elif command_exists curl; then
        sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b "$(go env GOPATH)/bin"
        task_installed=true
      else
        echo "Please install Homebrew or curl, then run this again"
        exit 1
      fi
      ;;
  esac
else
  echo "✅ Task is already installed"
fi

# Initialize lefthook hooks
echo "Setting up git hooks..."
lefthook install

# Verify installation
echo "Verifying lefthook setup..."
lefthook run pre-commit

# Provide summary message
echo "Installation summary:"
[ "$lefthook_installed" = true ] && echo "✅ lefthook newly installed" || echo "✅ lefthook was already installed"
[ "$commitlint_installed" = true ] && echo "✅ commitlint newly installed" || echo "✅ commitlint was already installed"
[ "$golangci_lint_installed" = true ] && echo "✅ golangci-lint v2 newly installed or upgraded" || echo "✅ golangci-lint v2 was already installed"
[ "$task_installed" = true ] && echo "✅ Task newly installed" || echo "✅ Task was already installed"

echo "Your git hooks are now active. They will run automatically on git operations."