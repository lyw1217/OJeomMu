# OJEOMMU
ARG ARCH=
ARG BUILD_IMAGE=${ARCH}golang:1.18-alpine3.15
ARG BASE_IMAGE=${ARCH}alpine:3.15

FROM ${BUILD_IMAGE} AS builder
    
WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
COPY /etc/letsencrypt/live/mumeog.site/fullchain.pem /etc/letsencrypt/live/mumeog.site/fullchain.pem
RUN go build -v -o /usr/local/bin/ojeommu .

FROM ${BASE_IMAGE}
LABEL AUTHOR Youngwoo Lee (mvl100d@gmail.com)
ENV GIN_MODE=debug \
    PORT=30000 \
    TZ=Asia/Seoul

RUN apk --no-cache add tzdata && \
	cp /usr/share/zoneinfo/$TZ /etc/localtime && \
	echo $TZ > /etc/timezone \
	apk del tzdata

RUN set -x \
&& apk add --no-cache openssh \
&& echo 'root:P@ssw0rd' | chpasswd \
&& ssh-keygen -f /etc/ssh/ssh_host_rsa_key -N '' -t rsa \
&& mkdir -p ~/.ssh/ \
&& echo 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCshioTzYDPuFt1VM1PkBrqRxrpLRTmkG3MhKdqtooy75+2UF3Eld3R/hWl49U8w41wa+VSL5lmjAHBxIRaHjcFrHNIoRYWraqI6chHr8I82AvBnGNo+KJgxp1KcXLrPg3iy82UJZX9+XwhHjeh7g3c0SHTgp4Fiyu/mieiXioHi4NKra2ESNhrdPdg0el4oBsRxtT1Fyv+MsZuT0R0er+LFR0LbwtZ0B65QBZD6jj0YZvhuB1O4cld8gLJ9nJZAHkO+BuwQFnq6uEuTnSjWqGaBFmmW9lBBHbF5w0XZaG2uWaQ+HAkrM5rNdN3tmp+kyjcq5hd173uBxioGOtQ4DcPSDMHa5fkqK3zNuTXtU0S9PqaZK/RtJU64cS2sY6J6BjvN3b90bjDfWimQFdlgd+s+S0RK15xPS7VpQSrCb160TRHTWlJYreqd+86j10zwj42vmm2BQxq2jGlk8jV6YMB1ig9cZ/qnfR7YDWJXQLVkt30B17idw13YxTk6iUQSBE= jenkins@c753384856db' > ~/.ssh/authorized_keys
EXPOSE 22

WORKDIR /usr/src/app

COPY /etc/letsencrypt/live/mumeog.site/fullchain.pem /etc/letsencrypt/live/mumeog.site/fullchain.pem

COPY --chown=0:0 --from=builder /usr/src/app/assets					/usr/src/app/assets
COPY --chown=0:0 --from=builder /usr/src/app/static					/usr/src/app/static
COPY --chown=0:0 --from=builder /usr/src/app/templates				/usr/src/app/templates

COPY --chown=0:0 --from=builder /usr/src/app/config/keys.json		/usr/src/app/config/keys.json
COPY --chown=0:0 --from=builder /usr/src/app/log					/usr/src/app/log
COPY --chown=0:0 --from=builder /usr/src/app/go.mod					/usr/src/app/go.mod
COPY --chown=0:0 --from=builder /usr/src/app/go.sum					/usr/src/app/go.sum
COPY --chown=0:0 --from=builder /usr/local/bin/ojeommu				/usr/local/bin/ojeommu

EXPOSE ${PORT}

CMD ["ojeommu"]