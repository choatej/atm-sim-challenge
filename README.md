# atm-sim
___


## Setup

The application requires Golang 1.20+

After cloning the repo, the following should be run
```sh
go mod tidy
```

In order to run all of the build targets, also run:
```sh
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Using the application

To run the application from the source code:
```bash
go run atm-sim.go
```

To create, install and run the application as a binary
```bash
make clean install
atm-sim
```

To run the application as a docker container (assuming you have a docker daemon running)
```bash
make clean docker
docker run -it atm-sim:latest
```

## Development
The application uses the cobra library for creating CLI commands.

To install the cobra CLI, run:
```sh
go install github.com/spf13/cobra
```

### Design and Assumptions
- The auth data (account-> pin) should be kept separate from balance data (account-> balance)
- All pins are 4 digits
- The balances are stored in US dollars
- The source data in csv is clean
- Balance and history checks do not need to be logged
- All logins will be logged, wheteher they fail or succeed
- all transactions (deposit / withdrawal) will be logged

### Unit tests
The goal of the unit tests was not to achieve 100% coverage but to ensure that the
critical areas of the code were being tested.

### Some Potential enhancements
- store the data from the csv in a database
- encrypt the pins in the csv
- store the encrypted pins and their salt separately
- set a limit on the overdraft amount
- create a remote application for the application logic that is secure, HA and allows for multiple clients