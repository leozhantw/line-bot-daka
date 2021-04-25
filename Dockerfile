FROM golang:1.16.3-alpine3.13 AS build
WORKDIR /src
COPY . .
RUN go build -mod=vendor -o bin/server ./cmd/server

FROM alpine:3.13
COPY --from=build /src/bin/server .
ENV PORT 8080
EXPOSE 8080
ENTRYPOINT ./server
