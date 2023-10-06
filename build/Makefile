image_name=driver-rtl-sdr
PWD:=$(shell pwd)

.PHONY: clean

build-dump1090.%: docker.% deb.%
	docker run \
	    --mount type=bind,source=$(PWD)/$(subst build-dump1090,deb,$@),target=/deb \
		--volume $(PWD)/..:/application:ro \
		$(image_name)-$(subst build-dump1090.,,$@) \
		build-dump1090

build-driver.%: docker.% deb.%
	docker run \
	    --mount type=bind,source=$(PWD)/$(subst build-driver,deb,$@),target=/deb \
		$(image_name)-$(subst build-driver.,,$@) \
		build-driver

docker.%:
	docker build -t $(image_name)-$(subst docker.,,$@) -f $(subst docker,Dockerfile,$@) .

clean:
	docker image rm -f $(image_name)-arm64
	docker image rm -f $(image_name)-amd64
	rm -rf deb.*

deb.%:
	mkdir -p $@

