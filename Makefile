all: clean build

check: checkExit checkErr

build:
	go build -o staticlinter staticlinter.go
clean:
	rm -rf staticlinter

checkExit:
	go test -v linters/osexit_linter.go linters/osexit_linter_test.go
checkErr:
	go test -v linters/err_linter.go linters/err_linter_test.go

