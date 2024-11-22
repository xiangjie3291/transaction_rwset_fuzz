
build:
	cd standard-dfa && ./build.sh dfa
	cd standard-nfa && ./build.sh nfa
	cd standard-identity && ./build.sh identity
	cd standard-evidence && ./build.sh evidence
	cd encdata && ./build.sh encdata
	cd exchange && ./build.sh exchange
	cd fact && ./build.sh fact
	cd itinerary && ./build.sh itinerary
	cd raffle && ./build.sh raffle
	cd trace && ./build.sh trace
	cd vote && ./build.sh vote
	ls -laht */ |grep .7z

lint:
	cd standard-dfa && golangci-lint run ./...
	cd standard-nfa && golangci-lint run ./...
	cd standard-identity && golangci-lint run ./...
	cd standard-evidence && golangci-lint run ./...
	cd encdata && golangci-lint run ./...
	cd exchange && golangci-lint run ./...
	cd fact && golangci-lint run ./...
	cd raffle && golangci-lint run ./...
	cd trace && golangci-lint run ./...
	cd vote && golangci-lint run ./...

gomod:
	cd standard-dfa && go mod tidy
	cd standard-nfa && go mod tidy
	cd standard-identity && go mod tidy
	cd standard-evidence && go mod tidy
	cd encdata && go mod tidy
	cd exchange && go mod tidy
	cd fact && go mod tidy
	cd itinerary && go mod tidy
	cd raffle && go mod tidy
	cd trace && go mod tidy
	cd vote && go mod tidy

all: gomod lint build