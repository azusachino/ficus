# ficus

Ficus (フィークス) 无花果

## description

1. a blog management service
2. build with go and gofiber
3. should be high performance

## setup

```sh
# install podman
yum install podman
# install podman-compose
yum install python3-pip
pip3 install podman-compose
# set up alias maybe?
vim ~/.zsh
alias pc=podman-compose
# install go-task
go install github.com/go-task/task/v3/cmd/task@latest
```

### with testcontainer

- enable sock `systemctl --user start podman.socket`
  - check status
- set up `DOCKER_HOST` = `unix://$XDG_RUNTIME_DIR/podman/podman.sock`

## References

- [gofiber](https://github.com/gofiber/fiber)
- [testcontainers](https://testcontainers.com/guides/getting-started-with-testcontainers-for-go/)
