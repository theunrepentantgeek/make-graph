# Target-specific variable assignments should be skipped as prerequisites
# but the target itself should still be tracked

.PHONY: build-cross
build-cross: build
build-cross: LDFLAGS += -extldflags "-static"
build-cross:
	CGO_ENABLED=0 gox -osarch='linux/amd64' -ldflags '$(LDFLAGS)'

build:
	go build ./...
