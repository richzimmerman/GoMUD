test:
	echo "This seems like a bad way to run tests..."
	go test -cover ./src/users/... ./src/telnet ./src/output ./src/mobs

build:
	go build -o bin/GameServer main.go

run:
	go run main.go

compile-linux:
	GOOS=linux GOARCH=386 go build -o bin/GameServer main.go

compile-linux-64bit:
	GOOS=linux GOARCH=amd64 go build -o bin/GameServer main.go

compile-windows:
	GOOS=windows GOARCH=386 go build -o bin/GameServer.exe main.go

compile-windows-64bit:
	GOOS=windows GOARCH=amd64 go build -o bin/GameServer.exe main.go