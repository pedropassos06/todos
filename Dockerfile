FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN if [ -n "${TARGETOS}" ] && [ -n "${TARGETARCH}" ]; then \
			CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bootstrap ./cmd/lambda; \
		else \
			CGO_ENABLED=0 go build -o /bootstrap ./cmd/lambda; \
		fi

FROM public.ecr.aws/lambda/provided:al2023

COPY --from=builder /bootstrap /var/task/bootstrap
CMD ["bootstrap"]
