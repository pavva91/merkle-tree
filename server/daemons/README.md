# Create systemd Unit

How to run the service as a systemd unit (daemon)

## Unit file location

../daemons/gotask.service
should be in:

```bash
/etc/systemd/system/<service_name>.service
```

## Create Symlink

```bash
sudo ln -s ~/work/task/daemons/gotask.service /etc/systemd/system/gotask.service
```

## Start the service

```bash
systemctl start gotask.service
```

### Start the service at boot automatically
Start the script on boot, enable the service with systemd

```bash
systemctl enable gotask.service
```

## Verify Status

```bash
sudo systemctl status gotask.service
```

## Try with cURL

```bash
curl --location --request GET 'http://localhost:8080/task/1' \
--header 'Content-Type: application/json' | jq
```

## Guide

[https://fedoraproject.org/wiki/Packaging:Systemd](https://fedoraproject.org/wiki/Packaging:Systemd)
