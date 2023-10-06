.PHONY: build-dump1090 build-driver clean

build-dump1090: build-dump1090.amd64 build-dump1090.arm64

build-driver: build-driver.amd64 build-driver.arm64

build-driver.%:
	$(MAKE) -C build $@

build-dump1090.%:
	$(MAKE) -C build $@

clean:
	$(MAKE) -C build clean