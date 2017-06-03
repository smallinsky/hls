FROM ubuntu:latest

RUN  apt-get update && apt-get install -y \
    wget \
    tar \
    git \
    ffmpeg


RUN wget https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz
RUN tar -xvf go1.8.linux-amd64.tar.gz
RUN mv go /usr/local
ENV GOROOT=/usr/local/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
RUN go get -u golang.org/x/sys/unix


RUN mkdir hls
COPY src hls
RUN cd hls && go build

RUN mkdir video
