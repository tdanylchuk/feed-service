FROM golang:alpine as build-env
RUN apk update && apk add git && apk --no-cache add gcc g++ make ca-certificates
WORKDIR /app
ADD . /app
RUN cd /app && go build -o feed-service

FROM alpine
WORKDIR /app
COPY --from=build-env /app/feed-service /app

EXPOSE 8000
ENTRYPOINT ./feed-service