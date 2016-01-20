FROM osxcc/golang:devel
MAINTAINER zchee <k@zchee.io>

RUN set -ex \
	&& uname -a \
	&& go version \
	&& go env \
	\
	&& go get -u -v -d github.com/zchee/docker-machine-driver-xhyve \
	&& cd $GOPATH/src/github.com/zchee/docker-machine-driver-xhyve \
	&& go build -v -o bin/docker-machine-driver-xhyve -ldflags "-w -s -X `go list ./xhyve`.GitCommit=`git rev-parse --short HEAD 2>/dev/null`" github.com/zchee/docker-machine-driver-xhyve

CMD ["cat", "src/github.com/zchee/docker-machine-driver-xhyve/bin/docker-machine-driver-xhyve"]
