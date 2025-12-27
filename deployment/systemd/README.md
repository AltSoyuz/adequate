# Systemd

This directory contains the systemd service unit file for Adequate.

To deploy Adequate as a systemd service, copy the `adequate.service` file to your system's systemd directory (usually `/etc/systemd/system/`), then enable and start the service with the following commands:

```bash
sudo cp deployment/systemd/adequate.service /etc/systemd/system/
sudo systemctl enable adequate
sudo systemctl start adequate
```