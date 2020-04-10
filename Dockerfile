# Start from the latest golang base image
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY ./pkg ./pkg
COPY ./cmd/jorel .
COPY ./go.mod .
COPY ./config.yaml .
COPY ./service-key.json .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -o jorel .

######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/jorel .
COPY --from=builder /app/service-key.json .
COPY --from=builder /app/config.yaml .

# Command to run the executable
CMD ["./jorel"]