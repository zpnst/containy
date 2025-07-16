build:
	@go build -o containy

extract:
	@tar -xvzf debian.tar.gz -C bundle

install:
	@mv containy /usr/local/bin/

uninstall:
	@rm /usr/local/bin/containy

.PHRONY: build, clean, extract, install, uninstall