FROM golang1.27.2-alpine3.24 AS BUILD

LABEL site="api"
LABEL stage="builder"

WORKDIR /src/

COPY go.mod ./
COPY go.sum ./
COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go generate

COPY *.go ./

RUN apk update && apk add git

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/api

EXPOSE 8081

ENTRYPOINT ["/bin/api"]