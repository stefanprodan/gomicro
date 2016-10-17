FROM golang:latest 

ADD vendor /go/src/
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o gomicro . 

CMD ["/app/gomicro"]
