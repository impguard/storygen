FROM golang:1.13.4-stretch

RUN mkdir /storygen
ADD . /storygen/
WORKDIR /storygen

RUN go build -o storygen ./service

CMD ["./storygen"]