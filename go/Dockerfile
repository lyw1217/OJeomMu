# OJEOMMU
ARG ARCH=
ARG BUILD_IMAGE=${ARCH}golang:1.18-buster
ARG BASE_IMAGE=${ARCH}debian:buster-slim

FROM ${BUILD_IMAGE} AS builder

WORKDIR /usr/src/app

COPY . .

RUN go mod download && go mod verify

RUN go build -v -o /usr/src/app/ojeommu .

WORKDIR /usr/src
RUN tar -cvf app.tar ./app

FROM ${BASE_IMAGE}
LABEL AUTHOR Youngwoo Lee (mvl100d@gmail.com)
ENV GIN_MODE=debug \
    HTTP_PORT=80 \
    HTTPS_PORT=443 \
    REDIRECT_PORT=9090 \
    TZ=Asia/Seoul

WORKDIR /usr/src

COPY --chown=0:0 --from=builder /usr/src/app.tar        /usr/src
COPY --chown=0:0 --from=builder /usr/src/app/run.sh     /usr/src
COPY --chown=0:0 --from=builder /usr/src/app/ojeommu    /usr/src/app/ojeommu

RUN ln -s /usr/src/app/ojeommu /usr/local/bin/ojeommu

EXPOSE ${HTTP_PORT} ${HTTPS_PORT} ${REDIRECT_PORT}

ENTRYPOINT ["/usr/src/run.sh"]
