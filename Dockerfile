FROM golang:1.7.1-alpine

# install curl 
RUN apk add --update curl && rm -rf /var/cache/apk/*

# copy deps
ADD vendor /go/src/

# copy sources
RUN mkdir /app 
ADD . /app/ 

# build
WORKDIR /app 
RUN go build -o gomicro . 

HEALTHCHECK CMD curl --fail http://localhost:3000/ping || exit 1

# run
CMD ["/app/gomicro"]
