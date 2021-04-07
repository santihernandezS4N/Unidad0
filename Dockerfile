FROM golang

ENV GO111MODULE=on

# Add Maintainer Info
LABEL maintainer="Sanhernandezmon <sanhernandezmon@unal.edu.co>"

RUN apt-get update
RUN apt-get install -y git
RUN git clone https://github.com/santihernandezS4N/Unidad0.git

WORKDIR /Unidad0

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]