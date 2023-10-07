# Dump1090

Dump 1090 is a Mode S decoder specifically designed for RTLSDR devices.

This is a fork from the project https://github.com/antirez/dump1090 from Salvatore Sanfilippo.

The code has been cleaned to remove all net features, and has been wrapped in goland code.

## Prerequisites

* You must have a sane installation of Docker.
* Install qemu for other architecture support through docker:

```bash
sudo apt-get install -y qemu qemu-user-static
```

## Build

1. Build the drivers
2. Build the application

```bash
make build-driver
make build-dump1090
```

This will generate debian packages in folder `build/deb*`.
