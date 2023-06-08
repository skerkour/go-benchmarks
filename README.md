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

## Results

**Last update**: 2023-05-06

amd64:
* [AMD EPYC 7543 (Scaleway ENT1-S)](results/scaleway_ent1_s.txt)
* [AMD EPYC 7R13 (AWS EC2 c6a.4xlarge)](results/aws_c6a_4xlarge.txt)
* [Intel Xeon Platinum 8375C @ 2.90GHz (AWS EC2 c6i.4xlarge)](results/aws_c6i_4xlarge.txt)

arm64:
* [Ampere Altra Max Neoverse-N1 (Scaleway AMP2-C8)](results/scaleway_amp2_c8.txt)
* [Graviton 3 (AWS EC2 c7g.4xlarge)](results/aws_c7g_4xlarge.txt)
