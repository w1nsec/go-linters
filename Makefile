all: clean build checkExit

build:
	go build -o staticlinter staticlinter.go
checkExit:
	./staticlinter linters/testdata/pkg3_osExit.go
clean:
	rm -rf staticlinter
