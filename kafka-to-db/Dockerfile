FROM golang:latest

LABEL maintainer="roma makhanov"

WORKDIR /app

# ARG LOG_DIR=/app/logs

# RUN mkdir -p ${LOG_DIR}

# ENV LOG_FILE_LOCATION=${LOG_DIR}/app.log 

COPY . . 

RUN apk add librdkafka-dev pkgconf
RUN go get -d github.com/gorilla/mux
RUN go build -o main .


EXPOSE 8080

VOLUME ${LOG_DIR}

CMD ["./main"]