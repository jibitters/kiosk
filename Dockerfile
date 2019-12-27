FROM ubuntu:18.04

RUN apt update && apt install ca-certificates -y && rm -rf /var/cache/apt/*

COPY /kiosk-linux-* /app/kiosk
COPY /migration /app/migration

VOLUME /app/configs

WORKDIR /app

EXPOSE 9090
EXPOSE 8080
EXPOSE 9091

ENTRYPOINT ["./kiosk"]
