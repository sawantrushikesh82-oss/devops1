FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY . .
RUN go mod download && go build -o /out/app ./cmd/ext-sftp

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /out/app /app/app
COPY configs /app/configs
ENTRYPOINT ["/app/app"]
