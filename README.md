# flux

A simple TCP server, written in [Go(lang)](https://go.dev/).

The idea of the project is to set up a server that listen to port `8080`, and launch a [goroutine](https://go.dev/tour/concurrency/1) for each incoming connection. Then, store variables in key/value, and provide methods to interact with them.

## Methods

1. `Set`
2. `Get`
3. `Del`
4. `PING`
5. `Set name val EX 10`
6. `LPUSH name val`

## Project structure

```
flux/
├── main.go          # Main file
├── main_test.go     # Simple test file
├── internal/        # App internal logic
|   ├── handler/     # Read the request
|   ├── loader/      # Load yaml data
|   ├── models/      # Structs used across the project
|   ├── parser/      # Decide what to do with each action
|   ├── server/      # Configures and starts the server
|   ├── store/       # Store a key-value variable
├── data.yaml
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
    * Use a custom port with the flag `--port=3000`.
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

8. Set a value with expiration date:
```
echo "Set test1 23 EX 10" | ncat -C localhost 8080
```

9. Add items to slice:
```
echo "LPUSH test1 val" | ncat -C localhost 8080
```

This will store `val` in an slice. Run the same command changing `val` to add more items to the slice.
Then, you can just run the `Get` command, and you will get the full slice.

## Thats all, folks! Thanks for visiting!!
