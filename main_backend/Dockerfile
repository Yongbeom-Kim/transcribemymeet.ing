FROM golang:1.23-alpine

ARG PORT=8080
ARG TF_VAR_backend_identity_key
ARG TF_VAR_resource_name

WORKDIR /app

COPY --link go.mod go.sum ./

RUN go mod tidy

COPY --link . .

RUN go build -o main .

EXPOSE ${PORT}

ENV PORT=${PORT}
ENV TF_VAR_backend_identity_key=${TF_VAR_backend_identity_key}
ENV TF_VAR_resource_name=${TF_VAR_resource_name}

CMD ["./main"]
