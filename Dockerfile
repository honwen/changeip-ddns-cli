FROM alpine:edge
MAINTAINER CHENHW2 <https://github.com/chenhw2>

ARG BIN_URL=https://github.com/chenhw2/changeip-ddns-cli/releases/download/v20170420/changeip_linux-amd64-20170420.tar.gz
ARG TZ=Asia/Hong_Kong

RUN apk add --update --no-cache wget supervisor ca-certificates tzdata \
    && update-ca-certificates \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime \
    && rm -rf /var/cache/apk/*

RUN mkdir -p /opt \
    && cd /opt \
    && wget -qO- ${BIN_URL} | tar xz \
    && mv changeip_* changeip

ENV Username=1234567890 \
    Password=abcdefghijklmn \
    Domain=ddns.changeip.com \
    Redo=0

ADD Docker_entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
