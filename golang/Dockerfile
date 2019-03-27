FROM alpine:3.9
LABEL maintainer="Pontus Aurdal <pontus.edvard@gmail.com>"

ARG UID=1000
ARG GID=1000

USER root
RUN apk add --no-cache shadow sudo curl && \
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

RUN apk add --no-cache musl-dev go

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /opt/go
ENV GOBIN /go/bin
ENV PATH $GOBIN:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin ${GOBIN}
RUN curl https://glide.sh/get | sh

WORKDIR /home/visualbox

COPY ./bootstrap/bootstrap-linux-amd64 .
COPY ./golang/visualbox /opt/go/src/

RUN chown -Rv visualbox:visualbox /home/visualbox
RUN chmod -Rv a+r /opt/go
USER visualbox

CMD ./bootstrap-linux-amd64
