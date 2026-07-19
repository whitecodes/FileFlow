FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

COPY fileflow /usr/local/bin/fileflow

EXPOSE 8080

ENTRYPOINT ["fileflow"]
