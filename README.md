# KVMage: Virt Image Creator

KVMage is an image creation tool for building and customizing qcow2 images for KVM. It works similarly to HashiCorp Packer but is designed to leverage tools already available on a KVM hypervisor such as `virt-install`, `virt-customize`, and `qemu-img`. There is nothing extra to install beyond KVMage itself — it is a single compiled Go binary that works with what you already have.

## Requirements

The full list of requirements is in the `REQUIREMENTS` file in the root of the repo. KVMage performs a requirements check during install and at the beginning of command execution. You can also run a manual check using the `-R, --check-requirements` flag.

Required executables:

```
curl (or wget)
lsof
osinfo-query
qemu-img
virt-customize
virt-install
virt-resize
virsh
```

## Installation

There are several ways to install and use KVMage:

- Download a precompiled binary from the [GitLab Releases](https://gitlab.com/kvmage/kvmage/-/releases) page.
- Pull the container image from the GitLab Container Registry.
- Clone the repository and compile from source.

### Precompiled Binaries

Download the latest binary for your platform from the [Releases](https://gitlab.com/kvmage/kvmage/-/releases) page. Binaries are available for:

- `kvmage-linux-amd64`
- `kvmage-linux-arm64`
- `kvmage-darwin-amd64`
- `kvmage-darwin-arm64`

```bash
# Example: download and install on Linux amd64
curl -fsSL -o kvmage https://gitlab.com/kvmage/kvmage/-/releases/permalink/latest/downloads/kvmage-linux-amd64
chmod +x kvmage
sudo mv kvmage /usr/local/bin/
```

### Container Image

Container images are published to the GitLab Container Registry on every release. Available tags include the full version, major.minor, major, and `latest`.

```bash
docker pull registry.gitlab.com/kvmage/kvmage:latest
```

Or pin to a specific version:
```bash
docker pull registry.gitlab.com/kvmage/kvmage:2.2.11
```

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

KVMage provides a streamlined method for creating qcow2 images that is designed to feel like a natural extension to KVM by using existing tools already available on a KVM hypervisor.

### Build Modes

KVMage has two build modes:

- `install`: Creates a brand-new image from installation media (ISO or URL) and an unattended install file. Supported install file types:
  - **Kickstart** — for RHEL-based distros (Fedora, Alma, Rocky, CentOS, etc.)
  - **Preseed** — for Debian-based distros (Debian, Ubuntu, etc.)

- `customize`: Takes an existing qcow2 image as a source and modifies it. You can upload files, run scripts, resize the disk, expand partitions, and set the hostname.

### Operating Modes

KVMage supports two methods for providing configuration:

- `run`: Use the `-r, --run` flag to pass all options as CLI arguments.
- `config`: Use the `-f, --config` flag to provide a YAML config file. Config mode supports defining multiple image builds in a single file — builds are executed sequentially in the order they appear.

### Variable Substitution

When using config mode, KVMage supports variable substitution in YAML config files using `${VAR}` or `$VAR` syntax. Variables are loaded from three sources in order of precedence (highest last):

1. A `.env` file in the same directory as the config file (auto-loaded if present)
2. A file specified with `--env-file`
3. System environment variables

### Install Mode

RUN example:
```bash
kvmage \
    --run \
    --install \
    --image-name almalinux9 \
    --os-var almalinux9 \
    --image-size 100G \
    --install-file ks.cfg \
    --install-media almalinux9.5-minimal.iso \
    --image-dest . \
    --firmware hybrid
```

CONFIG example:
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
    firmware: hybrid
```

### Customize Mode

RUN example:
```bash
kvmage \
    --run \
    --customize \
    --image-name almalinux9-latest \
    --image-src almalinux9.qcow2 \
    --image-dest . \
    --upload configs/ \
    --execute setup.sh
```

CONFIG example:
```yaml
---
kvmage:
  almalinux9:
    image_name: almalinux9-latest
    virt_mode: customize
    image_src: almalinux9.qcow2
    image_dest: .
    upload:
      - configs/
      - scripts/helper.sh
    execute:
      - setup.sh
```

### Upload and Execute

In customize mode, you can upload files and directories into the guest image and execute scripts inside it:

- **Upload** (`-U, --upload`): Copies files or directories into `/tmp/kvmage/` inside the guest. Accepts local paths (relative or absolute). Multiple items can be specified.
- **Execute** (`-E, --execute`): Runs scripts inside the guest in the order specified. If an execute file was not already included in the upload list, it is automatically uploaded. Scripts run from `/tmp/kvmage/` and the directory is cleaned up after execution.

### Firmware Options

KVMage supports three firmware modes for install mode:

- `bios` (default): Standard BIOS boot.
- `efi`: UEFI boot with the Q35 machine type.
- `hybrid`: BIOS and UEFI compatible boot using Q35 with UEFI firmware. Useful for images that need to boot in both BIOS and UEFI environments, such as bootc images.

### Console Options

For install mode, you can control how the VM console is presented:

- `serial`: Headless serial console (`console=ttyS0`). Useful for automated builds without a display.
- `graphical`: VNC console on `127.0.0.1`. Useful for debugging installs visually.
- `dual`: Both serial and graphical consoles simultaneously. VNC is available for graphical access while serial console is accessible via `virsh console`. Useful when you want to monitor an install via serial but also have graphical access available.
- If unset, the default libvirt console behavior is used.

### Options Reference

```
Usage:
  kvmage [--run | --config] [flags]

Execution Modes (required):
  -r, --run                     Use CLI arguments directly
  -f, --config <file>           Use a YAML config file

Build Modes (required with --run):
  -i, --install                 Install mode (create image from ISO)
  -c, --customize               Customize mode (modify existing image)

Image Options:
  -n, --image-name <name>       Name of the image
  -o, --os-var <os>             OS variant (use osinfo-query os)
  -s, --image-size <size>       Image size (e.g., 100G), expands image in customize mode
  -P, --image-part <device>     Partition to expand (e.g., /dev/sda1)
  -k, --install-file <file>     Path to Kickstart or Preseed file
  -j, --install-media <path>    Install media path or URL (ISO or install tree)
  -S, --image-src <file>        Source QCOW2 image (customize mode)
  -D, --image-dest <file>       Destination QCOW2 image
  -H, --hostname <name>         Hostname to set inside the image (optional)
  -U, --upload <path>           Files or directories to upload (temp)
  -E, --execute <file>          Scripts to execute in order
  -W, --network <iface>         Virtual network name (optional)
  -m, --firmware <type>         Firmware type: bios (default), efi, or hybrid
      --console <type>          Console type: serial, graphical, or dual (optional)
      --env-file <file>         Path to env file for variable substitution

Global Options:
  -h, --help                    Show help and exit
  -v, --verbose                 Enable verbose output (-v/-vv/-vvv)
      --verbose-level <n>       Set verbosity level explicitly (0-3)
  -q, --quiet                   Suppress all output
  -V, --version                 Show version info for KVMage and tools
  -R, --check-requirements      Check system requirements and exit
  -u, --uninstall               Uninstall KVMage from /usr/local/bin
  -X, --cleanup                 Run cleanup mode to remove orphaned kvmage temp files
```

## Docker Container Usage

KVMage comes packaged as a container for easy image generation in CI/CD pipelines.

Since KVM requires kernel level access (hence the "K" in KVM) you need to pass through certain parameters from the container to the host.

```bash
sudo docker run --rm -it \
  --privileged \
  --device /dev/kvm \
  -v ${PWD}:/kvmage \
  -v /var/run/libvirt:/var/run/libvirt \
  -v /var/lib/libvirt:/var/lib/libvirt \
  registry.gitlab.com/kvmage/kvmage:latest \
  --config kvmage.yml
```

### Auto Build

```bash
KVMAGE_BRANCH=main bash <(curl -s https://gitlab.com/kvmage/kvmage/-/raw/main/scripts/autobuild.sh)
```