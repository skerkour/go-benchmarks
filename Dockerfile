FROM golang:latest

RUN apt update && apt upgrade -y && apt install -y make curl wget git

ARG USERNAME=benchmarks
ARG USER_UID=10001
ARG USER_GID=$USER_UID

RUN addgroup --gid $USER_GID $USERNAME \
    && adduser --uid $USER_UID --ingroup $USERNAME --disabled-password --shell /bin/bash --gecos "" $USERNAME

USER $USERNAME

WORKDIR /home/${USERNAME}

COPY . ./
COPY ./.git /home/${USERNAME}/benchmarks/.git
RUN git config --global --add safe.directory /home/${USERNAME}/benchmarks

RUN make download_and_verify_deps

ENTRYPOINT ["/bin/sh", "-c"]
CMD ["make", "run"]
