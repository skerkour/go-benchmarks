DOCKER_IMAGE = ghcr.io/skerkour/go-benchmarks:latest

.PHONY: run
run:
	go run tools/system_info/main.go
	go test -benchmem -bench=. github.com/skerkour/go-benchmarks/hashing
	go test -benchmem -bench=. github.com/skerkour/go-benchmarks/checksum
	go test -benchmem -bench=. github.com/skerkour/go-benchmarks/chunking
	go test -benchmem -bench=. github.com/skerkour/go-benchmarks/encryption
	go test -timeout 1h -benchmem -bench=. github.com/skerkour/go-benchmarks/compression
	go test -benchmem -bench=. github.com/skerkour/go-benchmarks/signatures
# disable inlining
	go test -benchmem  -bench=. -gcflags '-l' github.com/skerkour/go-benchmarks/cgo
	go test -benchmem -bench=. github.com/skerkour/go-benchmarks/encoding

.PHONY: run_docker
run_docker:
	docker run -ti --rm $(DOCKER_IMAGE)


# Docker
.PHONY: docker_build
docker_build:
	docker build -t $(DOCKER_IMAGE) . -f Dockerfile

.PHONY: docker_push
docker_push:
	docker push $(DOCKER_IMAGE)


# Other
.PHONY: download_and_verify_deps
download_and_verify_deps:
	go mod download
	go mod verify

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: update_deps
update_deps:
	go get -u ./...
	go mod tidy
	go mod tidy
