# Use the official Go image as the base image
FROM golang:1.20 as builder

# Set the working directory
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go dependencies
RUN go mod download

# Copy the Go source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a minimal alpine image for the final container
FROM alpine:3.14

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Set the environment variables
ENV PLAYLISTIFY_SPOTIFY_SECRET=$PLAYLISTIFY_SPOTIFY_SECRET
ENV PLAYLISTIFY_SPOTIFY_ID=$PLAYLISTIFY_SPOTIFY_ID
ENV PLAYLISTIFY_DBHOST=localhost
ENV PLAYLISTIFY_DBPORT=5433
ENV PLAYLISTIFY_DBUSERNAME=dev
ENV PLAYLISTIFY_DBPASSWORD=dev
ENV PLAYLISTIFY_DBNAME=plailist
ENV PLAYLISTIFY_SPOTIFY_REDIR="http://127.0.0.1:8888/spotify/auth"
ENV PLAYLISTIFY_AUTH_REDIR="http://127.0.0.1:5173/"

# Expose necessary ports
EXPOSE 8888 5173

# Run the binary
CMD ["./main"]