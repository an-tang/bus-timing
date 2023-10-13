# build stage
FROM golang:1.21.3-alpine AS build

ENV APP bus-timing
ENV GO111MODULE=on

RUN apk add git

COPY . /go/src/$APP/
WORKDIR /go/src/$APP

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/$APP main.go

###
FROM golang:1.21.3-alpine
ENV APP bus-timing

RUN apk add git; \
    apk add build-base; \
    apk add --no-cache tzdata;

# COPY ./database/migrations/ /database/migrations/
COPY ./entrypoint.sh /entrypoint.sh
COPY ./configuration /go/configuration
COPY --from=build /go/bin/$APP /go/bin/$APP

EXPOSE 8080
ENTRYPOINT ["sh","/entrypoint.sh"]