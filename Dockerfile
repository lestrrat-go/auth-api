FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN mkdir -p /usr/local/go/src/tripleoak/auth-api
RUN ln -s /app/* /usr/local/go/src/tripleoak/auth-api
RUN go build -o /auth-api

EXPOSE 8080
CMD [ "/auth-api" ]