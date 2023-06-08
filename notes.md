## Resources

```bash
sudo su
apt update && apt upgrade -y && apt install -y tmux
curl -fsSL https://get.docker.com | sh
reboot
tmux new -s bench
docker run --pull=always -ti --rm ghcr.io/skerkour/go-benchmarks:latest > result.txt
```

if SSH connection drop:
```bash
tmux a -t bench
```

* https://easyperf.net/blog/2019/08/02/Perf-measurement-environment-on-Linux
* https://tmuxcheatsheet.com/
