APP = vmg
SRC = $(wildcard *.go)

all: build

build:
	go build -o $(APP) $(SRC)

rebuild:
	rm -f $(APP)
	go build -o $(APP) $(SRC)
run: build
	./$(APP)

clean:
	rm -f $(APP)
