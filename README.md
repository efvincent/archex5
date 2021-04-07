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

## Branch `step-3-mux` Gorilla Mux
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

## Branch `step-4-models` Models, Commands, Events
In this step we implement models, commands, events, and stub out the functions in the command processor.Since one of the goals is to build and demonstrate an event sourced system, it would be good to pick something where we can demonstrate handling contention, different read models, replay, and time travel.

I choose products to model, a simplified example product model, but it is interesting enough to have several commands and events to implement, and potentially we can show different reducers and projectors.

### Models
This example has a single model (at this time), the `ProductModel`. The attributes are as you'd expect, with the exception of `SequenceNum` which is the last event in the stream that makes up this state of the `ProductModel`.

In an event sourced system there are any number of possible ways to combine the events in a stream into a model, but you typically see a "canonical" model used by the command processors. We'll see more of how event streams are turned into models when we get to the reducer step.

### Commands
Commands aren't a necessary part of an event sourced solution; rather you see them typically in [CQRS](https://docs.microsoft.com/en-us/azure/architecture/patterns/cqrs) or _Command Query Responsibility Separation_ pattern. The way to think of this is a command represents a request or attempt to take some action, often that action is to change the state of some durable entity in the system (in our case the main entities are Products).

A Key characteristic of a command is that it may fail for any number of legitimate reasons. For example, the action that is being requested might not be allowed under the constraints (business rules). Our product system may have a rule that disallows a product to be created in a namespace if that SKU already exists.

In our system the commands are found in `./commands/productCommands.go`

Also note that these commands are not the same as the `cobra.Command` struct type used by the Cobra CLI and configuration package which is unfortunately a naming collision and is completely unrelated to CQRS Commands.

### Command Processor
Commands are passed to command processors, in our case `./processor/productCmdProc.go`. The command processor typically has one function that takes any valid command instance, determines the type, and dispatches the call to a function made specifically to handle that type of command:
```go
func ProcessProductCommand(cmd interface{}) error {
	switch c := cmd.(type) {
	case *commands.CreateProductCmd:
		return ProcessCreateProduct(c)
	case *commands.HeadCheckCmd:
		return PerformHeadCheck(c)
...
```

### Events
Finally we have the events, located in `./events/ProductEvents.go`. Do not confuse these events with the "raw events" that are generated by user interaction with a web page (and other raw events) received by the API at Bluecore. In this context, events are very specifically "event sourcing" events.

Each event is a record of something **that has happened** and has been recorded in our system. That's why events are named in the past tense, for example:
```go
type ProductCreated struct {
	Event
	Source  string              `json:"source"`
	Product models.ProductModel `json:"product"`
}
```
the `ProductCreated` event records the fact that the product was created. The command processor has the logic and opportunity to access resources such as databases or other APIs to determine if a command is valid when it is received, and what event or events should occur as a result of processing the command. Once the events are created and recorded, it is a permanent part of the history of the system. Events are **immutable**, once they are created and successfully written to the event store, they cannot be deleted during the normal course of events.

As a practical matter most event sourced systems allow for compaction and/or removal of events as part of a retention scheme, but that's out of scope of this exercise.