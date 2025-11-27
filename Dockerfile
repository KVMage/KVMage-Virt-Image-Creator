ARG CONTAINER_NAME=9-minimal
ARG ALMA_VERSION=9.6
ARG ALMA_TAG_DATE=20251117

FROM ghcr.io/almalinux/${CONTAINER_NAME}:${ALMA_VERSION}-${ALMA_TAG_DATE}

ARG ALMA_VERSION
ARG KVMAGE_VERSION
ARG BUILD_DATE
ARG RELEASEVER

LABEL maintainer="KVMage"
LABEL org.opencontainers.image.title="KVMage: Virt Image Builder" \
      org.opencontainers.image.description="Virt Image Builder (KVMage) container with virt and guestfs tools" \
      org.opencontainers.image.version="${KVMAGE_VERSION}" \
      org.opencontainers.image.created="${BUILD_DATE}"

RUN : echo "${RELEASEVER:=${ALMA_VERSION}}" && \
    echo "${RELEASEVER}" > /etc/dnf/vars/releasever

# System packages
RUN microdnf update -y && \
    microdnf install -y epel-release && \
    microdnf install -y \
        cloud-utils-growpart \
        dosfstools \
        edk2-ovmf \
        e2fsprogs \
        git \
        go \
        guestfs-tools \
        libosinfo \
        libvirt-client \
        osinfo-db-tools \
        parted \
        qemu-kvm \
        rsync \
        virt-install \
        virt-top \
        virt-viewer \
        vim \
        which \
        xz && \
    microdnf clean all

# Install KVMage
RUN bash <(curl -s https://gitlab.com/kvmage/kvmage/-/raw/main/scripts/autoinstall.sh)

# Configure secure_path
RUN echo 'Defaults secure_path="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"' \
    > /etc/sudoers.d/10-secure-path && \
    chmod 0440 /etc/sudoers.d/10-secure-path

WORKDIR /kvmage

ENV LIBGUESTFS_BACKEND=direct

CMD ["/bin/bash"]
