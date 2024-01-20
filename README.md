# Go Benchmarks

Comprehensive and reproducible benchmarks for Go developers.

**Motivation**: We need real-world numbers in order to design efficient applications and protocols. Here they are.


* [Encoding](encoding)
* [CGO](cgo)
* [Checksum](checksum)
* [Chunking](chunking)
* [Compression](compression)
* [Encryption](encryption)
* [Hashing](hashing)
* [Signatures](signatures)


## Usage

```shell
$ make run
```

or with `docker` (amd64, arm64):

```shell
$ make docker_build # optional
$ docker run --pull=always -ti --rm ghcr.io/skerkour/go-benchmarks:latest
```
<!--
sudo apt update && sudo apt upgrade -y && sudo apt dist-upgrade -y
curl -fsSL https://get.docker.com -o /tmp/get-docker.sh && sh /tmp/get-docker.sh
reboot
docker run --pull=always -d ghcr.io/skerkour/go-benchmarks:latest
ssh xx@xx -i xx 'docker logs xx'
-->

## Results

**Last update**: 2023-05-06

amd64:
* [AMD EPYC 7543 (Scaleway POP2-8C-32G)](results/scaleway_POP2-8C-32G.txt)
* [AMD EPYC 9R14 (AWS EC2 c7a.4xlarge)](results/aws_c7a_4xlarge.txt)
* [4th Generation Intel Xeon Scalable 8375C @ 2.90GHz (AWS EC2 c7i.4xlarge)](results/aws_c7i_4xlarge.txt)

arm64:
* [Ampere Altra Max Neoverse-N1 (Scaleway COPARM1-8C-32G)](results/scaleway_COPARM1-8C-32G.txt)
* [Graviton 3 (AWS EC2 c7g.4xlarge)](results/aws_c7g_4xlarge.txt)
