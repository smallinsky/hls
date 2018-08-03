FROM golang:1.10.3-alpine3.7
COPY . /go/src/github.com/smallinsky/hls
RUN go install github.com/smallinsky/hls/cmd/hlsrun

FROM alpine:3.7
RUN apk --no-cache add ffmpeg
COPY --from=0 /go/bin/hlsrun /hlsrun
RUN mkdir /tmp/video && mkdir /tmp/segments
RUN /hlsrun -d /tmp/video -s /tmp/segments -p 8080
