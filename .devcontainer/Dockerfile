FROM golang:1.26 AS go

FROM debian:13

ENV TZ="UTC"
ENV LC_ALL="en_US.UTF-8"
ENV LANG="en_US.UTF-8"
ENV LANGUAGE="en_US.UTF-8"

RUN apt update && apt upgrade -y \
    && apt install -y git git-lfs ca-certificates make vim postgresql-client-17 \
    bash gcc g++ zsh curl wget zip coreutils sqlite3 python3 gnupg libpcre2-dev pkg-config build-essential binutils jq b3sum \
    coreutils clang lld
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
COPY --from=go /usr/local/go /usr/local/go
RUN /usr/local/go/bin/go telemetry off

WORKDIR /
