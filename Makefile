HEADHASH := $(shell git rev-parse --short HEAD)
DT := $(shell date +%Y-%m-%d_%H%M%S)
TAG := $(shell git describe --tags --abbrev=0)
IMAGENAME = "jirm/gwc-server:$(HEADHASH)"
IMAGELATEST = "jirm/gwc-server:latest"
IMAGENAMETAG = "jirm/gwc-server:$(TAG)"
PWD := $(shell pwd)

.PHONY: run
run:
	@printf "\033[32m--> Running GWC Server...\n\033[0m"
	go run ./app

.PHONY: build
build:
	@printf "\033[32m--> Building GWC Server\n\033[0m"
	go mod download
	go mod verify
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s -extldflags "-static"' -a -o ./gwc-server app/main.go

.PHONY: clean
clean:
	@printf "\033[32m--> Cleaning\n\033[0m"
	rm -r ./gwc-server

# Building GWC Server docker image
.PHONY: build_docker
build-docker:
	@printf "\033[32m--> Building GWC Server docker image\n\033[0m"
	docker build --tag --tag $(IMAGENAME) --tag $(IMAGELATEST) .

.PHONY: build-docker-no-cache-tag
build-docker-no-cache-tag:
	@printf "\033[32m--> Building GWC Server docker image\n\033[0m"
	docker build --no-cache --tag $(IMAGENAMETAG) --tag $(IMAGENAME) --tag $(IMAGELATEST) .

# Running GWC Server docker container
.PHONY: run_docker
run-docker:
	@printf "\033[32m--> Run the GWC Server server\n\033[0m"
	docker run --rm -v $(PWD)/config.yml:/gwc/config.yml -v $(PWD)/ed25519_test.key:/gwc/ed25519_test.key -v $(PWD)/known_hosts:/gwc/known_hosts -p 8080:8080 jirm/gwc-server

