Needed for RHEL/Fedora/Alma/Rocky etc.....

```
echo 'Defaults secure_path="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"' | sudo tee /etc/sudoers.d/10-secure-path >/dev/null
```