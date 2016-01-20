FROM osxcc/golang:devel
MAINTAINER zchee <k@zchee.io>

RUN set -ex \
	&& uname -a \
	&& go version \
	&& go env \
	\
	&& go get -u -v -d github.com/zchee/docker-machine-driver-xhyve \
	&& cd $GOPATH/src/github.com/zchee/docker-machine-driver-xhyve \
	&& make

CMD ["cat", "src/github.com/zchee/docker-machine-driver-xhyve/bin/docker-machine-driver-xhyve"]
