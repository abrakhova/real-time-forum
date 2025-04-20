FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o forum
# Using the debian:bookworm-slim image for a small runtime environment
FROM debian:bookworm-slim
WORKDIR /app
# Metadata
LABEL project="Forum"
LABEL description="A web forum"
COPY --from=builder /app/forum /app/forum
COPY --from=builder /app/data /app/data
COPY --from=builder /app/images /app/images
COPY --from=builder /app/internal /app/internal
COPY --from=builder /app/web /app/web

# Expose port 8080 to allow the app to be accessible from outside the container
EXPOSE 8080
CMD ["/app/forum"]