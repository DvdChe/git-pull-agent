# Stage 1: Builder
# Use a Go base image with necessary tools to build the application
FROM golang:1.25.7-alpine AS builder

# Set necessary environment variables for static compilation
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod ./
COPY go.sum ./

# Download dependencies. This step is cached if go.mod/go.sum don't change
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
# The -o flag specifies the output file name
# The -tags netgo ensures that the net package uses pure Go DNS resolution
# The -ldflags "-s -w" reduces the size of the binary by stripping debug information
RUN go build -o /git-pull-agent -tags netgo -ldflags="-s -w" main.go

# Stage 2: Runner
# Use a distroless image for the final runtime image
FROM gcr.io/distroless/static-debian11

# Set the working directory (optional, but good practice)
WORKDIR /

# Copy the compiled binary from the builder stage
COPY --from=builder /git-pull-agent /git-pull-agent

# Set the entrypoint for the container
ENTRYPOINT ["/git-pull-agent"]

# Expose any ports if your application needed them (not applicable for this agent)
# EXPOSE 8080
