FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure

RUN go build -ldflags '-w -s' -a -installsuffix cgo -o app

CMD ["./app"]
EXPOSE 3000
