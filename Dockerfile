FROM bitnami/minideb:bullseye

WORKDIR /app

EXPOSE 80

RUN apt update && apt install -y ca-certificates tzdata tini

ENV TZ=Asia/Shanghai

COPY bin/unifi-captive-portal /sbin/unifi-captive-portal
COPY assets /app/assets
COPY templates /app/templates
RUN chmod +x /sbin/unifi-captive-portal
# install
ENTRYPOINT ["tini", "--"]

CMD ["unifi-captive-portal"]