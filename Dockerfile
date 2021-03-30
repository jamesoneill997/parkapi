FROM golang:1.15-alpine AS build

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    mongosrv=mongodb+srv://james:Paul&James2021@cluster0.dznpg.mongodb.net/parkai?retryWrites=true&w=majority \
    secret=dcusoc2021

WORKDIR /.

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Export necessary port
EXPOSE 8080

# Command to run when starting the container
CMD ["/main"]