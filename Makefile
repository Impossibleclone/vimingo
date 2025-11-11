APP = vmg
PREFIX = /usr/bin/

# The main package to build
MAIN_PKG = ./cmd/vmg

all: build

install: build
	@echo "Installing $(APP) to $(PREFIX)..."
	@mv $(APP) $(PREFIX)
	@echo "Done! You may need sudo."

build:
	@echo "Building $(APP)..."
	@go build -o $(APP) $(MAIN_PKG)
	@echo "Build complete: ./$(APP)"

rebuild: clean build

run: build
	@./$(APP)

clean:
	@echo "Cleaning..."
	@rm -f $(APP)
