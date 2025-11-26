sudo kvmage --run \
    --customize \
    --image-name debian2 \
    --os-var debiantesting \
    --image-src ./debian.qcow2 \
    --image-dest . \
    --hostname debian \
    --custom-script script.sh \
    --firmware efi