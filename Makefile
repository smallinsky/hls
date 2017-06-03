all: buildimage run

buildimage:
	docker build -t hls .

run:
	[ -d video ] || mkdir video
	docker run -v $(shell pwd)/video:/video -p 8888:8888 -it hls /hls/hls -d /video
