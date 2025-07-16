build:
	@go build -o bundle/containy

run: build
	@./bundle/containy

clean:
	@ rm bundle/containyc
