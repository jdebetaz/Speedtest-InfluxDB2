FROM golang:latest

#Install speedtest-cli
RUN apt-get install curl
RUN curl -s https://packagecloud.io/install/repositories/ookla/speedtest-cli/script.deb.sh | bash
RUN apt-get install speedtest

# Setup application
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . /go/src/app

# Install dependencies
RUN go get -d -v ./...
RUN go install -v ./...

# Build the application
RUN go build -o speedtest-influx .

# Run the application
CMD ["./speedtest-influx"]