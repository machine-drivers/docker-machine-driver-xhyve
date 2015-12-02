TODO
===
Current TODO is listed in [README.md](https://github.com/zchee/docker-machine-driver-xhyve/blob/master/README.md#todo).

## archived
- [x] Support(Ensure) `kill`, `ls`, `restart`, `status`, `stop` command

- [x] Daemonize xhyve use `syscall` or `go execute external process myself` or `OS X launchd daemon` or other daemonize method
    - Since it is was hard to exec itself that embedded xhyve.Exec, for the time being, is separated into [xhyve-bindings](https://github.com/zchee/xhyve-bindings/tree/daemonize).

- [x] Replace exec uuid2mac binary to standalone `vmnet.go`, `dhcp.go`

- [x] Update xhyve source to unofficial edge branch
    - See [update-xhyve-to-edge](https://github.com/zchee/docker-machine-driver-xhyve/tree/update-xhyve-to-edge)
    - Replace `Grand Central Dispatch` instead of `pthreads` , and etc...
    - See https://github.com/AntonioMeireles/xhyve/tree/edgy
      - Separated [xhyve-bindings](https://github.com/zchee/xhyve-bindings/tree/daemonize)

- [x] Occasionally fail convert UUID to IP
    - Fixed [1960629b3c8683aec193631a0e9573c5143832ab](https://github.com/zchee/docker-machine-driver-xhyve/commit/1960629b3c8683aec193631a0e9573c5143832ab)

- [x] Crash on boot because of `prltoolsd`
    - Crash it's not an empty disk.img?
    - See https://github.com/ailispaw/boot2docker-xhyve/pull/16
      - Solved on `boot2docker v1.8.3`

