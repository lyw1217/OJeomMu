# OJEOMMU
ARG ARCH=
ARG BUILD_IMAGE=lyw1217/ojeommu-base:0.2
ARG BASE_IMAGE=${ARCH}debian:buster-slim

FROM ${BUILD_IMAGE} AS builder

FROM ${BASE_IMAGE}
LABEL AUTHOR Youngwoo Lee (mvl100d@gmail.com)
ENV GIN_MODE=debug \
    PORT=30000 \
    TZ=Asia/Seoul

ARG ssh_pub_key

RUN set -x \
&& apt-get update -y \
&& apt-get install -y --no-install-recommends \
  openssh-server \
&& rm -rf /var/lib/apt/lists/* \
&& mkdir /var/run/sshd \
&& echo 'root:P@ssw0rd' | chpasswd \
&& sed 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' -i /etc/ssh/sshd_config \
&& sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd \
&& mkdir -p ~/.ssh/ \
&& echo ${ssh_pub_key} > ~/.ssh/authorized_keys \
&& chmod 0700 ~/.ssh \
&& chmod 0600 ~/.ssh/authorized_keys \
&& mkdir -p /run/openrc \
&& touch /run/openrc/softlevel
ENV NOTVISIBLE "in users profile"
RUN echo "export VISIBLE=now" >> /etc/profile

EXPOSE 22

WORKDIR /usr/src/app

COPY --chown=0:0 --from=builder /usr/src/app/assets					      /usr/src/app/assets
COPY --chown=0:0 --from=builder /usr/src/app/static					      /usr/src/app/static
COPY --chown=0:0 --from=builder /usr/src/app/templates				    /usr/src/app/templates

COPY --chown=0:0 --from=builder /usr/src/app/scripts				      /usr/src/app/scripts
COPY --chown=0:0 --from=builder /usr/src/app/config/keys.json		  /usr/src/app/config/keys.json
COPY --chown=0:0 --from=builder /usr/src/app/config/logging.json	/usr/src/app/config/logging.json
COPY --chown=0:0 --from=builder /usr/src/app/log					        /usr/src/app/log
COPY --chown=0:0 --from=builder /usr/src/app/go.mod					      /usr/src/app/go.mod
COPY --chown=0:0 --from=builder /usr/src/app/go.sum					      /usr/src/app/go.sum
COPY --chown=0:0 --from=builder /usr/src/app/ojeommu				      /usr/src/app/ojeommu

RUN ln -s /usr/src/app/ojeommu /usr/local/bin/ojeommu

EXPOSE ${PORT}

ENTRYPOINT /bin/sh -e /etc/init.d/ssh start; ojeommu; tail -f /dev/null