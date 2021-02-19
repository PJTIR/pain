FROM ubuntu:18.04

WORKDIR /root

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install --no-install-recommends -y \
        libguestfs-tools \
        qemu-utils \
        linux-image-generic \
        python3

ENV LIBGUESTFS_BACKEND=direct \
    HOME=/root


WORKDIR /
ADD exporter / 
ADD iou_export /

#ENTRYPOINT [ "/bin/bash" ]

CMD [ "./exporter" ]
