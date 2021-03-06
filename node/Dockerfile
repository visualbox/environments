FROM node:11-alpine
LABEL maintainer="Pontus Aurdal <pontus.edvard@gmail.com>"

ARG UID=1000
ARG GID=1000

USER root
RUN apk add --no-cache shadow sudo && \
    if [ -z "`getent group $GID`" ]; then \
      addgroup -S -g $GID visualbox; \
    else \
      groupmod -n visualbox `getent group $GID | cut -d: -f1`; \
    fi && \
    if [ -z "`getent passwd $UID`" ]; then \
      adduser -S -u $UID -G visualbox -s /bin/sh visualbox; \
    else \
      usermod -l visualbox -g $GID -d /home/visualbox -m `getent passwd $UID | cut -d: -f1`; \
    fi && \
    echo "visualbox ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/visualbox && \
    chmod 0440 /etc/sudoers.d/visualbox

WORKDIR /home/visualbox

COPY ./bootstrap/bootstrap-linux-amd64 .
COPY ./node/visualbox /opt/node_modules/visualbox/

RUN chown -Rv visualbox:visualbox /home/visualbox
RUN chmod -Rv a+r /opt/node_modules
USER visualbox

ENV NODE_PATH=/opt/node_modules:/home/visualbox/app/node_modules

CMD ./bootstrap-linux-amd64
