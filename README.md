# Tic-Tac-Toe Game

This is a Tic-Tac-Toe game implemented in Go with support for both local and TCP multiplayer modes. Follow the instructions below to clone, build, and run the project.

## Libraries Used
- [github.com/charmbracelet/bubbles v0.18.0](https://pkg.go.dev/github.com/charmbracelet/bubbles@v0.18.0)
- [github.com/charmbracelet/bubbletea v0.26.4](https://pkg.go.dev/github.com/charmbracelet/bubbletea@v0.26.4)
- [github.com/charmbracelet/lipgloss v0.11.0](https://pkg.go.dev/github.com/charmbracelet/lipgloss@v0.11.0)
- [github.com/spf13/cobra v1.8.0](https://pkg.go.dev/github.com/spf13/cobra@v1.8.0)


## Table of Contents

- Libraries Used
- Cloning the Repository
- Installing the Dependencies
- Running the Project
- CI/CD with GoReleaser
- Installation via Homebrew



## Cloning the Repository

To clone the repository, run the following command in your terminal:

```sh
git clone https://github.com/dziedzicgrzegorz/Tic-Tac-Toe.git
cd Tic-Tac-Toe
```
## Installing the dependencies
```go
go mod vendor
```
go mod vendor is a command used in Go to create a vendor directory in your project, which will include all the dependencies specified in your go.mod file. This command copies the dependencies from the module cache to the vendor directory, making them available for version control and ensuring that your project can be built with the exact versions of dependencies you specified.

## Building the Project

To build the project, use the provided Makefile:

```sh
make build
```

This will compile the Go code and create the binary.

## Running the Project

To run the project, use the following command:

```sh
make all
```

This will run all the "build install run "

## Linting the Code

To lint the code, use:

```sh
make lint
```
## CI/CD with GoReleaser
This project uses GoReleaser to automate the release process. When you push a new tag to the repository, GoReleaser will create a new release and publish it to Homebrew
Please make sure you have filled in all the details (repository name and GitHub Token)

This will run the linter to check for code quality issues.

## Installation via Homebrew

You can also install the Tic-Tac-Toe game using Homebrew. First, make sure you have Homebrew installed. Then run the following command:

```sh
brew tap dziedzicgrzegorz/tic-tac-toe
brew install dziedzicgrzegorz
```

Once installed, you can run the game with:

```sh
Tic-Tac-Toe
```
![brew-install](https://github.com/DziedzicGrzegorz/Tic-Tac-Toe/assets/110931212/b77e4b71-9866-4cc8-b74f-cef9eec1d091)

Here's a screenshot from the game:

![TCP](https://github.com/DziedzicGrzegorz/Tic-Tac-Toe/assets/110931212/0a224d6a-d34c-41d9-afeb-cab44f1ba44b)

## Additional Information



## Official golang Website
- Ensure you have Go installed on your machine. You can download it from the [official Go website](https://golang.org/dl/).

## Feel free to open issues or contribute to the project by submitting pull requests.

Happy coding!
