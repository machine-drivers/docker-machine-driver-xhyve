class DockerMachineDriverXhyve < Formula
  desc "Docker Machine driver for xhyve"
  homepage "https://github.com/zchee/docker-machine-driver-xhyve"
  url "https://github.com/zchee/docker-machine-driver-xhyve/archive/v0.2.0.tar.gz"
  sha256 "affba9fe235e0ceb2711a52ea15b380bdcbec8aa2c106b54f750ede3c762dcb8"

  depends_on "go" => :build
  depends_on "docker-machine"

  def install
    contents = Dir["{*,.git,.gitignore}"]
    gopath = buildpath/"gopath"
    (gopath/"src/github.com/zchee/docker-machine-driver-xhyve").install contents

    ENV["GOPATH"] = gopath
    ENV["GO15VENDOREXPERIMENT"] = "1"

    cd gopath/"src/github.com/zchee/docker-machine-driver-xhyve" do
      system "go", "build", "-o", "docker-machine-driver-xhyve", "./bin/main.go"
      bin.install "docker-machine-driver-xhyve"
    end
  end

  def caveats; <<-EOS.undent
    docker-machine-driver-xhyve requires root privileges to correctly set up
    networking. You can either use docker-machine with sudo, or change the
    driver's owner and set the setuid bit:

      sudo chown root:wheel #{bin}/docker-machine-driver-xhyve
      sudo chmod u+s #{bin}/docker-machine-driver-xhyve
    EOS
  end
end
