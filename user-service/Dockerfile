FROM golang:alpine as build_container
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o user-service

FROM alpine
COPY --from=build_container /app/user-service /usr/bin
EXPOSE 8100
ENTRYPOINT["user-service"]