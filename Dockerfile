FROM alpine:3.10

COPY /kiosk-linux-* /app/kiosk
COPY /migration /app/migration

VOLUME /app/configs

WORKDIR /app

EXPOSE 9090
EXPOSE 9091
ENTRYPOINT ["./kiosk"]