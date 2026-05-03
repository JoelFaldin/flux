# flux

A simple TCP server, written in [Go(lang)](https://go.dev/).

The idea of the project is to set up a server that listen to port `8080`, and launch a [goroutine](https://go.dev/tour/concurrency/1) for each incoming connection. Then, store variables in key/value, and provide methods to interact with them.

## Methods

1. `Set`
2. `Get`
3. `Del`
4. `PING`

## Project structure

```
flux/
├── main.go          # Main file
├── main_test.go     # Simple test file
├── internal/        # App internal logic
|   ├── handler/     # Read the request
|   ├── parser/      # Decide what to do with each action
|   ├── server/      # Configures and starts the server
|   ├── store/       # Store a key-value variable
├── go.mod
├── README.md
```

## Prerequisites

1. Golang [1.21 or higher](https://go.dev/dl/)
2. Git installed on your machine.
3. `nc` installed (or equivalent for your os, examples here use `ncat`).

## Installation

1. Clone the repo with `git clone`.
2. Run `go install .`.
3. Start the main server with `go run .\main.go`.
4. In a new terminal, set a value:
```
echo "Set test 23" | ncat -C localhost 8080
OK
```

5. Get the value:
```
echo "Get test" | ncat -C localhost 8080
23
```

6. Delete the value:
```
echo "Del test" | ncat -C localhost 8080
OK
```

7. Check health of the server:
```
echo "PING" | ncat -C localhost 8080
PONG
```

## Thats all, folks! Thanks for visiting!!
