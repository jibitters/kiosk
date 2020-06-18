FROM ubuntu:18.04

RUN apt-get update && apt-get install -y ca-certificates

COPY /kiosk-linux-* /app/kiosk
COPY /migration /app/migration

VOLUME /app/configs

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["./kiosk"]
