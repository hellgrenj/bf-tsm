FROM golang:1.14-alpine3.11 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .


RUN ls -la
# Build the Go app
RUN go build -o main ./cmd/webdemo/main.go

# Expose port 1337 to the outside world
EXPOSE 1337

FROM alpine:3.11 as runtime

# copy only executable from build step
COPY --from=builder ./app/main main
RUN ls -la
COPY --from=builder ./app/cmd/webdemo/static static

# Command to run the executable
CMD ["./main"]