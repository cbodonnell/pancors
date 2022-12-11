# Build stage
FROM golang:1-alpine as builder

RUN mkdir /pancors
WORKDIR /pancors

COPY go.mod go.mod
RUN go mod download

COPY . .

RUN go build -o ./bin/pancors cmd/pancors/main.go

# Production stage
FROM alpine:latest as prod

RUN mkdir /pancors
WORKDIR /pancors

COPY --from=builder /pancors/bin/pancors ./

EXPOSE 8080

CMD [ "./pancors" ]