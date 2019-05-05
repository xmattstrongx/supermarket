deploy: build run

build:
	echo building image
	docker build . -t supermarket:latest || exit 1;

run:
	docker run --rm -it -p 8080:8080 supermarket:latest

install:
	go install github.com/xmattstrongx/supermarket
