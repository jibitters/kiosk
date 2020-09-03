FROM ubuntu:18.04

RUN apt-get update && apt-get install -y ca-certificates
RUN DEBIAN_FRONTEND="noninteractive" apt-get install -y tzdata

COPY /kiosk-linux-* /app/kiosk
COPY /migration /app/migration

VOLUME /app/configs

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["./kiosk"]
