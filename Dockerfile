FROM chenhw2/alpine:base
MAINTAINER CHENHW2 <https://github.com/chenhw2>

ARG VER=20170715
ARG URL=https://github.com/chenhw2/changeip-ddns-cli/releases/download/v$VER/changeip_linux-amd64-$VER.tar.gz

RUN mkdir -p /usr/bin \
    && cd /usr/bin \
    && wget -qO- ${URL} | tar xz \
    && mv changeip_* changeip

USER nobody

ENV USERNAME=1234567890 \
    PASSWORD=abcdefghijklmn \
    DOMAIN=ddns.changeip.com \
    REDO=0

CMD changeip \
    --username ${USERNAME} \
    --password ${PASSWORD} \
    auto-update \
    --domain ${DOMAIN} \
    --redo ${REDO}
