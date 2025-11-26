FROM ghcr.io/almalinux/almalinux:9.6-minimal-20251117

ARG KVMAGE_VERSION
ARG BUILD_DATE
ARG RELEASEVER=9.6

LABEL maintainer="KVMage"
LABEL org.opencontainers.image.title="KVMage: Virt Image Builder" \
      org.opencontainers.image.description="Virt Image Builder (KVMage) container with virt and guestfs tools" \
      org.opencontainers.image.version="${KVMAGE_VERSION}" \
      org.opencontainers.image.created="${BUILD_DATE}"

RUN echo "${RELEASEVER}" > /etc/dnf/vars/releasever

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

# Configure secure_path
RUN echo 'Defaults secure_path="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"' \
    > /etc/sudoers.d/10-secure-path && \
    chmod 0440 /etc/sudoers.d/10-secure-path


# Custom tool installation script
COPY autoinstall.sh /usr/local/bin/autoinstall.sh
RUN chmod +x /usr/local/bin/autoinstall.sh && \
    /usr/local/bin/autoinstall.sh

WORKDIR /kvmage

ENV LIBGUESTFS_BACKEND=direct

CMD ["/bin/bash"]
