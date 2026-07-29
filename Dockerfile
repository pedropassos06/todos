FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN if [ -n "${TARGETOS}" ] && [ -n "${TARGETARCH}" ]; then \
			CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /server ./cmd/server; \
		else \
			CGO_ENABLED=0 go build -o /server ./cmd/server; \
		fi

FROM alpine:3.21

COPY --from=builder /server /server
RUN chmod +x /server

EXPOSE 8081

ENTRYPOINT ["/server"]
