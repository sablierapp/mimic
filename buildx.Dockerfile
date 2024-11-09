# syntax=docker/dockerfile:1.2
FROM golang:1-alpine as builder

RUN apk --no-cache --no-progress add git tzdata make \
    && rm -rf /var/cache/apk/*

# syntax=docker/dockerfile:1.2
FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY mimic .

ENTRYPOINT ["/mimic"]
EXPOSE 80