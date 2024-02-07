FROM golang:alpine 


WORKDIR /app

COPY . ./

RUN go mod tidy


RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

CMD ["/docker-gs-ping"]

