FROM golang:1.23-alpine

ARG PORT=8080

WORKDIR /app

COPY --link go.mod go.sum ./

RUN go mod tidy

COPY --link . .

RUN go build -o main .

EXPOSE ${PORT}

ENV PORT=${PORT}

CMD ["./main"]
