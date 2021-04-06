# Arch Example

## Branch `step-1-scaffold`
This **ArchEX5** project has branches that show the result of doing blocks of steps. Branch `step-1-scaffold` is the result of minimally scaffolding out the project
1. From `GOROOT` which for me is `~/go`, create a folder under `github.com/[user]/[project]`. This is the project root folder
2. Initialize a git repo here
2. Touch `go.mod`, add one line `go 1.15`
3. Create hello world `main.go` in project root
4. Verify it runs. From project root: 
```bash
$ go run main.go
```

## Branch `step-2-cobra` Configuration and CLI
At Bluecore we use [cobra](https://github.com/spf13/cobra) to create CLIs. In this step we add a reference to the package and scaffold out our first commands.
1. Run the install steps from the [cobra](https://github.com/spf13/cobra) docs. Cobra implements a CLI itself that allows you to quickly add configuration and CLI abilities to your project.
2. Use the cobra CLI to initialize cobra in this project
```
$ cobra init --pkg-name github.com/[user]/[project] .
```
3. Clean up the comments and unneeded stuff cobra created, and edit the `cmd/root.go` for our use. Currently the root execution (with no command line parameters) will not do anything other than print usage.
    * Cobra init has also set up [viper](github.com/spf13/viper), which is makes it easy to get config information from the command line, config files, environment variables, and more. See the repo for details.
    * Lastly Cobra init installs [go-homedir](github.com/mitchellh/go-homedir) which is a cross platform lib to get the running process home directory
4. Add a `hello` command using cobra
    * Adds `hello.go` in your `cmd` folder
    * Wires the hello command to the root command which makes it available
    * You'll now see help for the hello when you run the program and without commands or with `--help` 
```
$ cobra add hello
```



