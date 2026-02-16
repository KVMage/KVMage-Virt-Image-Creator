# KVMage: Virt Image Creator

KVMage is an image creation software similar to something like "HashiCorp Packer" that is used to assist with and even fully automate the creation of qcow2 image for use with KVM. What makes KVMage uniique is that it is designed to leverage tools you may already have installed on a KVM hypervisor like virt-install and virt-customize. This is huge when you dont want to install new packages or programs and just want to work with what you have and since its written in Go, its super fast and runs and a single compiled binary on your system. 

## Requirements

The requirements are listed in the `REQUIREMENTS` file located in the root of the repo. KVMage performs a requirements check during install and at the beginning of command execution. You can also optionally perform a manual check using the `-R, --check-requirements` flag options.

There are three primary required executables:

```bash
virt-customize
virt-install
qemu-img
```

## Installation

There are currently two different officially supported methods for use.

- Clone the repository and compile the code. (We include a script that will do this automatically)
- Clone the repository and create the Docker image.

Precompiled binaries and readily available container images in GHCR and Docker Hub will be available in the future.

### Manual Installation

Use `git clone` to download the repo locally:

``` bash
git clone https://gitlab.com/kvmage/kvmage.git
```

``` bash
cd kvmage
mkdir -p dist
bash build.sh
bash install.sh
cd ..
rm -rf kvmage
```

Autoinstall script
```
KVMAGE_BRANCH=main bash <(curl -s https://gitlab.com/kvmage/kvmage/-/raw/main/scripts/autoinstall.sh)
```

### Manually Build Docker Image

Build KVMage Docker container with `latest` tag
```bash
docker build -t kvmage:latest https://gitlab.com/kvmage/kvmage.git
```

Build KVMage Docker container with `VERSION` tag
```bash
docker build \
  -t kvmage:latest \
  -t "kvmage:$(curl -fsSL https://gitlab.com/kvmage/kvmage/-/raw/main/VERSION | tr -d '\n')" \
  https://gitlab.com/kvmage/kvmage.git
```

## How to Use KVMage

KVMage provides a streamlined method for creating qcow images that is designed to feel like a natural extension to KVM by using existing commands and features already available to users with a deployed KVM hypervisor.

### KVMage Build Methods

KVMage has two operating modes:

- `install`: creates a brand-new image using an installation media (an ISO or URL) and a startup script to perform the automated unattended installation of the system.

> **NOTE**
> The only supported methods for install currently are using a Kickstart file with RHEL-based distros such as Fedora, Alma, Rocky, etc...

- `customize`: creates an image using an existing qcow2 image as a source and modifies it with an identified script file (such as bash)

### KVMage Operating Modes

KVMage supports two different methods for operating:

- `run`: Use the `-r, --run` option with the `kvmage` command to perform setup using command line arguments and options.
- `config`: Use the `-f, --config` option with the `kvmage` command to perform setup using a config file (YAML) where the options are defined. Config mode is particularly useful for managing multiple image builds as you can stack builds.

### KVMage Install

RUN Example:
```bash
kvmage \
    --run \
    --install \
    --image-name almalinux9 \
    --os-var almalinux9 \
    --image-size 100G \
    --install-file ks.cfg \
    --install-media almalinux9.5-minimal.iso \
    --image-dest .
```

CONFIG Example:
```yaml
---
kvmage:
  almalinux9:
    image_name: almalinux9
    virt_mode: install
    os_var: almalinux9
    image_size: 100G
    install_file: ks.cfg
    install_media: almalinux9-minimal.iso
    image_dest: .
```

### KVMage Customize

RUN Example:

```bash
kvmage \
    --run \
    --customize \
    --image-name almalinux9-latest \
    --os-var almalinux9 \
    --image-src almalinux9.qcow2 \
    --image-dest $PWD \
    --execute script.sh
```

```yaml
---
kvmage:
  almalinux9:
    image_name: almalinux9-latest
    virt_mode: customize
    os_var: almalinux9
    image_src: almalinux9.qcow2
    img_dest: .
    execute: script.sh
```
### KVMage Options


```cfg
Usage:
  kvmage [--run | --config] [flags]

Execution Modes (required):
  -r, --run                     Use CLI arguments directly
  -f, --config <file>           Use a YAML config file

Installation Methods (required):
  -i, --install                 Install mode (create image from ISO)
  -c, --customize               Customize mode (modify existing image)

Image Options:
  -n, --image-name <name>       Name of the image
  -o, --os-var <os>             OS variant (use `osinfo-query os`)
  -s, --image-size <size>       Image size (e.g., 100G), expands image in customize mode
  -P, --image-part <device>     Partition to expand (e.g., /dev/sda1)
  -k, --install-file <file>     Path to Install file
  -l, --install-media <path>    Install media path or URL
  -S, --image-src <file>        Source QCOW2 image (customize mode)
  -D, --image-dest <file>       Destination QCOW2 image
  -H, --hostname <name>         Hostname to set inside the image (optional)
  -U, --upload <path>              Files or directories to upload (temp)
  -E, --execute <file>             Files to execute scripts (in order)
  -W, --network <iface>         Virtual network name (optional)
  -m, --firmware <type>         Firmware type: bios (default) or efi

Global Options:
  -h, --help                    Show help and exit
  -v, --verbose                 Enable verbose output (-v/-vv/-vvv)
      --verbose-level <n>       Set verbosity level explicitly (0-3)
  -q, --quiet                   Suppress all output
  -V, --version                 Show version info for KVMage and tools
  -u, --uninstall               Uninstall KVMage from /usr/local/bin
  -X, --cleanup                 Run cleanup mode to remove orphaned kvmage temp files
```

## Docker Container Usage

KVMage comes packaged as a container for easy image generation in CI/CD pipelines.

Since KVM requires kernel level access (hence the "K" in KVM) you need to pass through certain parameters from the container to the host.

Below is an example of what you want to execute:

```bash
sudo docker run --rm -it \
  --privileged \
  --device /dev/kvm \
  -v kvmage:/kvmage \
  -v /var/run/libvirt:/var/run/libvirt \
  -v /var/lib/libvirt:/var/lib/libvirt \
  kvmage:latest \
  install --run --customize kvmage.yml
```

### Auto Build

```bash
KVMAGE_BRANCH=main bash <(curl -s https://gitlab.com/kvmage/kvmage/-/raw/main/scripts/autobuild.sh)
```