# Development

This is for development and testing before pushing to Main

## Installation

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
bash <(curl -s https://gitlab.com/kvmage/kvmage/-/raw/main/scripts/autoinstall.sh)
```

### Manually Build Docker Image

Build KVMage Docker container with `latest` tag
```bash
docker build -t kvmage:latest https://gitlab.com/kvmage/kvmage.git
```

Build KVMage Docker container with `VERSION` tag
```bash
docker build -t kvmage:$(curl -fsSL https://gitlab.com/kvmage/kvmage/-/raw/main/VERSION | tr -d '\n') https://gitlab.com/kvmage/kvmage.git
```