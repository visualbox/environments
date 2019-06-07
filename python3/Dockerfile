FROM python:alpine
#FROM frolvlad/alpine-python3
LABEL maintainer="Pontus Aurdal <pontus.edvard@gmail.com>"

ARG UID=1000
ARG GID=1000

USER root
RUN apk add --no-cache shadow sudo build-base libffi-dev && \
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
COPY ./python3/visualbox.py /opt/python_modules/

RUN chown -Rv visualbox:visualbox /home/visualbox
RUN chmod -Rv a+r /opt/python_modules
USER visualbox

ENV PYTHONPATH=/opt/python_modules

CMD ./bootstrap-linux-amd64
