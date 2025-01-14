FROM golang:1.22-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy all files to the working directory
COPY . .

# Build the binary and name it 'task-one'
RUN go build -o task-one

FROM alpine:3.20 AS target

# Copy the built binary from builder stage to the target stage
COPY --from=builder /app/task-one /usr/local/bin/task-one

# Set the default command to run the binary
CMD [ "/usr/local/bin/task-one" ]