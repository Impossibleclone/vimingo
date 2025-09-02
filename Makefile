APP = vmg
SRC = $(wildcard *.go)
PREFIX=/usr/bin/

all: build

install: build
	mv $(APP) $(PREFIX)

build:
	go build -o $(APP) $(SRC)

rebuild:
	rm -f $(APP)
	go build -o $(APP) $(SRC)
run: build
	./$(APP)

clean:
	rm -f $(APP)
