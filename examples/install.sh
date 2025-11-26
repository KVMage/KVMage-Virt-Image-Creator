sudo kvmage --run \
    --install \
    --image-name debian \
    --os-var debiantesting \
    --image-size 100G \
    --install-file ./preseed.cfg \
    --install-media ./debian.iso \
    --image-dest . \
    --hostname debian \
    --firmware efi