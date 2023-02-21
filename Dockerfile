FROM golang:1.19
ENV TELEGRAM_APITOKEN ""
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o bot
CMD [ "./bot" ]