FROM osxcc/golang:devel
MAINTAINER zchee <k@zchee.io>

COPY . ${GOPATH}/src/github.com/zchee/docker-machine-driver-xhyve

RUN set -ex \
	&& uname -a \
	&& go version \
	&& go env \
	\
	&& cd $GOPATH/src/github.com/zchee/docker-machine-driver-xhyve \
	&& make

CMD ["cat", "src/github.com/zchee/docker-machine-driver-xhyve/bin/docker-machine-driver-xhyve"]
