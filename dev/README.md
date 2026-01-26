# Development

The `dev` branch is the active development branch. All new features and fixes 
should be committed or merged here before being promoted to `main`.

Note: The `dev` branch may contain experimental or unstable features.

## Installation

### Manual Installation

Use `git clone` to download the repo locally:

``` bash
git clone --branch dev --single-branch https://gitlab.com/kvmage/kvmage.git
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
bash <(curl -s https://gitlab.com/kvmage/kvmage/-/raw/dev/scripts/autoinstall.sh)
```

### Manually Build Docker Image

Build KVMage Docker container with `latest` tag
```bash
docker build -t kvmage:latest https://gitlab.com/kvmage/kvmage.git#dev
```

Build KVMage Docker container with `VERSION` tag
```bash
docker build -t kvmage:$(curl -fsSL https://gitlab.com/kvmage/kvmage/-/raw/dev/VERSION | tr -d '\n') https://gitlab.com/kvmage/kvmage.git#dev
```