FROM golang:1.22 AS go

FROM ubuntu:latest

ENV TZ="Europe/London"
ENV LC_ALL="en_US.UTF-8"
ENV LANG="en_US.UTF-8"
ENV LANGUAGE="en_US.UTF-8"

RUN apt update && apt upgrade -y
RUN apt install -y git git-lfs ca-certificates make vim postgresql-client-14 \
    bash gcc g++ zsh curl wget zip coreutils sqlite3 python3 gnupg xz-utils bzip2
    #  \
    # gcc-multilib gcc-mingw-w64

RUN update-ca-certificates

# Setup dev user
ARG USERNAME=dev
ARG USER_UID=10001
ARG USER_GID=$USER_UID

RUN addgroup --gid $USER_GID $USERNAME \
    && adduser --uid $USER_UID --ingroup $USERNAME --disabled-password --shell /bin/bash --gecos "" $USERNAME

USER $USERNAME

COPY bashrc /home/$USERNAME/.bashrc
COPY vimrc /home/$USERNAME/.vimrc

# Setup go
RUN mkdir -p /home/$USERNAME/.local/gopath

COPY --from=go /usr/local/go /home/$USERNAME/.local/go

ENV GOROOT /home/$USERNAME/.local/go
ENV GOPATH /home/$USERNAME/.local/gopath
ENV PATH $PATH:$GOPATH/bin:$GOROOT/bin

WORKDIR /
