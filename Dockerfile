FROM golang:1.20-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o main

FROM gcr.io/distroless/static-debian11
WORKDIR /app
USER nonroot
COPY --from=builder /src/main ./
EXPOSE 8080
CMD ["./main"]

