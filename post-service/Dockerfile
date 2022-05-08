FROM golang:alpine as build_container
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o post-service

FROM alpine
COPY --from=build_container /app/post-service /usr/bin
EXPOSE 8300
ENTRYPOINT["post-service"]