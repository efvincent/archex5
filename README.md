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

## Branch `step-3-mux` API
At Bluecore we use [Gorilla MUX](https://pkg.go.dev/github.com/gorilla/mux@v1.8.0?utm_source=gopls) for APIs in Go.

1. Install in our project
``` bash
$ go get -u github.com/gorilla/mux
```

2. Add a command to start an API server
```bash
$ cobra add server
```

3. Modify the `ServerCmd` descriptions
4. Set `host` and `port` as module level string variables to support the flags.
4. Modify the `init()` function to configure flags for `host` and `port`, with logical defaults. 

At this point we could start building the server out right here in the server command. It would be cleaner however to separate the actual server from the command that starts it.

### Implementing the API Server
1. At the project root, create a folder for the server module, call it `API`, and create an `API.go` within.
2. In `API.go` create a handler for a home route. This will handle requests that go to the root of the API. For now we'll return a simple result to prove that things are working.
```go
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Gorilla!</h1>"))
}
```
3. The `Run` function (which we could have named anything) is the entry point for the server. It expects a host and port, and will instantiate a router, add a handler for the home route, and finally start the server with the `http.ListenAndServer` function.
```go
func Run(host string, port string) {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)
	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Server running. Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
```
That's it for now. Next step will stub out some routes for the sample application.