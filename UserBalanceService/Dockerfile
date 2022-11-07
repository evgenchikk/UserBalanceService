FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . ./

RUN /go/bin/swag init -d ./cmd,./app/handlers,./api,./app/models
RUN go build -o /userbalanceservice cmd/main.go

# CMD [ "/userbalanceservice" ]

