.PHONY: build

$(info OS is $(shell uname))

ifeq ($(shell uname),Windows_NT)
    build:
	    if not exist dist mkdir dist
	    go build -o dist\eth-scaffolder.exe .
		copy .\config.yaml dist
		copy .\password.txt dist
else
    build:
	    mkdir -p dist
	    go build -o dist/eth-scaffolder .
		cp ./config/config.yaml dist
		cp ./password.txt dist
endif

