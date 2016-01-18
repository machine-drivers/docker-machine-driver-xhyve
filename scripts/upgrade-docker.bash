#/bin/sh
set -ex

# Upgrade docker scripts inside boot2docker
# Usage:
# $ cat ./scripts/upgrade-docker.bash | docker-machine ssh xhyve VERSION=1.10.0-rc1 sh

if [[ "$VERSION" == "master" ]]; then
  DOCKER_VERSION=master
  DOCKER_DOWNLOAD_URL="https://master.dockerproject.org/linux/amd64/docker"
else
  DOCKER_VERSION=${VERSION}
  DOCKER_DOWNLOAD_URL="https://test.docker.com/builds/Linux/x86_64/docker-${DOCKER_VERSION}"
fi

sudo sh -c "rm -f /usr/local/bin/docker && \
  curl -fL -o /usr/local/bin/docker ${DOCKER_DOWNLOAD_URL} && \
  chmod +x /usr/local/bin/docker && \
  /etc/init.d/docker stop && \
  /etc/init.d/docker start"

exit 0
