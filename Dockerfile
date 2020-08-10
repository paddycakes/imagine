FROM golang:1.13-buster as builder

WORKDIR /app

# Retrieve application dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN CGO_ENABLED=0 go build -mod=readonly -v -o server

# Use a Docker multi-stage build to create a lean production image.
FROM alpine:3

RUN apk add --no-cache imagemagick

# Install certificates for secure communication with network services.
# For production containers, a single RUN statement should install all system packages.
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server .

# Run the web service on container startup.
CMD ["/server"]


