FROM --platform=$BUILDPLATFORM golang:1.19-alpine3.16 AS build_base

RUN apk add --no-cache git

WORKDIR /tmp/app

COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go mod vendor

COPY . .
ARG TARGETOS TARGETARCH
RUN env GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o ./out/main ./main.go

FROM alpine:3.16 
RUN apk add ca-certificates

COPY --from=build_base /tmp/app/out/main /app/main

ENTRYPOINT ["/app/main"]