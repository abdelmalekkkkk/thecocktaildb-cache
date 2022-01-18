FROM golang:1.17
COPY . /app

# Install dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Build server

COPY *.go ./

RUN go build -o /updater

CMD [ "/updater" ]