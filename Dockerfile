FROM golang:latest 

RUN mkdir -p /go/src/github.com/Unanoc/InstaFollower/
WORKDIR /go/src/github.com/Unanoc/InstaFollower/
COPY . .

RUN go get -u github.com/go-telegram-bot-api/telegram-bot-api 
RUN go get -u -v gopkg.in/ahmdrz/goinsta.v2

RUN go build -o main .
CMD ["./main", "-cfg", "config.json"]