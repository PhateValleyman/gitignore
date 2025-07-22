BINARY_NAME=gitignore-auto
VERSION=1.0

build:
	go build -trimpath -v -x -o "$(BINARY_NAME)" ./main.go

install: build
	install -m 0755 "$(BINARY_NAME)" /data/data/com.termux/files/usr/bin/"$(BINARY_NAME)"

uninstall:
	rm -f /data/data/com.termux/files/usr/bin/"$(BINARY_NAME)"

clean:
	rm -f $(BINARY_NAME)

.PHONY: build install uninstall clean
