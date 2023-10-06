#!/bin/bash


for arg in "$@"
do
    case $arg in
        "build-driver" )
            mkdir -p /project/c && cd /project/c
            git clone https://github.com/osmocom/rtl-sdr.git
            cd rtl-sdr
            dpkg-buildpackage -us -uc
            cp /project/c/*.deb /deb/.
           exit 0
           ;;
        "build-dump1090" )
            cd /deb
            apt install -y \
                ./librtlsdr*.deb \
                ./librtlsdr-dev*.deb \
                ./rtl-sdr*.deb \
                ./rtl-sdr-dbgsym*.deb 
            cd /application

            case $(arch) in
                "x86_64" )
                    ARCH=amd64
                    ;;
                "aarch64" )
                    ARCH=arm64
                    ;;
            esac

            mkdir -p /dump1090-${ARCH}/usr/local/bin

            cp -r /application/build/DEBIAN /dump1090-${ARCH}/.
            cp -r /application/build/etc /dump1090-${ARCH}/.

            cat /dump1090-${ARCH}/DEBIAN/control | sed -e "s/amd64/${ARCH}/g">/tmp/control
            cp /tmp/control /dump1090-${ARCH}/DEBIAN/control
            
            go build -o /dump1090-${ARCH}/usr/local/bin/dumb1090 ./main.go

            dpkg-deb --build --root-owner-group /dump1090-${ARCH}
            cp /dump1090-${ARCH}.deb /deb/.
            exit 0
           ;;
   esac
done



