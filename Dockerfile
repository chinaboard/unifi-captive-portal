FROM golang AS builder
WORKDIR /go/src/app
ENV GOPROXY=https://goproxy.io,direct
COPY . .
RUN go mod tidy
RUN make build
RUN chmod +x /go/src/app/bin/unifi-captive-portal

FROM ubuntu

WORKDIR /app

EXPOSE 80

RUN apt update && apt install -y ca-certificates tzdata tini

COPY --from=builder /go/src/app/bin/unifi-captive-portal /sbin/unifi-captive-portal
COPY assets /app/assets
COPY templates /app/templates
RUN chmod +x /sbin/unifi-captive-portal

ENTRYPOINT ["tini", "--"]

CMD ["unifi-captive-portal"]