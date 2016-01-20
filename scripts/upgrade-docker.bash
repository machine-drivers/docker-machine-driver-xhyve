#/bin/sh
set -ex

# Upgrade docker scripts inside docker-machine vm
# Required /etc/init.d/docker script inside vm
# Usage:
# $ cat ./scripts/upgrade-docker.bash | docker-machine ssh xhyve VERSION=1.10.0-rc1 sh

if [[ "$VERSION" == "master" ]]; then
  DOCKER_VERSION=master
  DOCKER_DOWNLOAD_URL="https://master.dockerproject.org/linux/amd64/docker"
else
  DOCKER_VERSION=${VERSION}
  DOCKER_DOWNLOAD_URL="https://test.docker.com/builds/Linux/x86_64/docker-${DOCKER_VERSION}"
fi

sudo sh -c "rm -f /usr/local/bin/docker \
  && curl -fL -o /usr/local/bin/docker ${DOCKER_DOWNLOAD_URL} \
  && chmod +x /usr/local/bin/docker \
  && /etc/init.d/docker restart \
  && printf '\033[0;36mWait for restart docker server...\033[0m\n\n' \
  && sleep 3 \
  && printf '\033[0;36mUpgraded.\033[0m\n' \
  && printf '\033[0;36mdocker version\033[0m\n' \
  && docker version"

exit 0
