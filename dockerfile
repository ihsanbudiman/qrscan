# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.23.4-alpine3.20 AS build-stage

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./

# copy .env
COPY .env .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /app/app .

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8099

# Run
ENTRYPOINT ["/app/app"]
