FROM golang:1.23 AS builder
WORKDIR /src
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/tinykv ./cmd/tinykv

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /out/tinykv /tinykv
USER nonroot:nonroot
ENTRYPOINT ["/tinykv"]
