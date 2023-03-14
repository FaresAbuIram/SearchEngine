# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.18-buster as builder

# Create and change to the app directory.
WORKDIR /app

# Copy go.mod and go.sum files to the workspace.
COPY go.mod go.sum ./

# Download dependencies.
RUN go mod download

# Copy the rest of the application source code.
COPY . .

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /go/bin/app

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM gcr.io/distroless/static-debian10

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/bin/app /go/bin/app

# Run the web service on container startup.
CMD ["/go/bin/app"]