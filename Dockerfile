FROM golang:1.23.2-bullseye

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /stock-ticker

EXPOSE 8080

CMD [ "/stock-ticker" ]