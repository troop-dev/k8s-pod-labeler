FROM ghcr.io/troop-dev/go-kit:1.17.0

install-deps:
	# add git to known hosts
	RUN mkdir -p /root/.ssh && \
		chmod 700 /root/.ssh && \
		ssh-keyscan github.com >> /root/.ssh/known_hosts
	RUN git config --global url."git@github.com:".insteadOf "https://github.com/"
	# add dependencies
	COPY go.mod go.sum .
	RUN --ssh go mod download -x

add-code:
	FROM +install-deps
	# add code
	COPY --dir app .

vendor:
	FROM +add-code
	RUN --ssh go mod vendor

compile:
	FROM +vendor
	ARG GOOS=linux
	ARG GOARCH
	ARG GOARM
	RUN go build -mod=vendor -o bin/server ./app/main.go
	SAVE ARTIFACT bin/server /bin/server
	SAVE IMAGE --push ghcr.io/troop-dev/k8s-pod-labeler:compile

build:
	FROM alpine:3.13.6
	RUN apk add --no-cache libc6-compat
	ARG VERSION=dev
	WORKDIR /app
	COPY +compile/bin/server .
	RUN chmod +x ./server
	USER nobody
	ENTRYPOINT ./server run
	SAVE IMAGE --push ghcr.io/troop-dev/k8s-pod-labeler:${VERSION}

build-arm-v7:
	# THIS TARGET IS ONLY FOR BUILDING ARM64 IMAGES FROM AN
	# AMD64 HOST (Github Actions)
	FROM --platform=linux/arm64 alpine:3.13.6
	RUN apk add --no-cache libc6-compat
	ARG VERSION=dev
	WORKDIR /app
	COPY --platform=linux/amd64 --build-arg GOARCH=arm --build-arg GOARM=7 +compile/bin/server .
	RUN chmod +x ./server
	USER nobody
	ENTRYPOINT ./server run
	SAVE IMAGE --push ghcr.io/troop-dev/frontend-api:${VERSION}

test:
	FROM +compile
	RUN go test -coverpkg=./app/... -mod=vendor -coverprofile=coverage.out ./app/...
	SAVE ARTIFACT coverage.out AS LOCAL coverage.out
	SAVE IMAGE --push ghcr.io/troop-dev/k8s-pod-labeler-cache:test

lint:
	FROM +vendor
	# Runs golangci-lint with settings:
	RUN golangci-lint run --timeout 10m --skip-dirs-use-default

ci:
	BUILD +lint
	BUILD +test
	BUILD +build
	BUILD +build-arm-v7
