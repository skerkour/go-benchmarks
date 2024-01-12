FROM golang:1.21

RUN apt update && apt upgrade -y
RUN apt install make curl wget

ARG USERNAME=benchmarks
ARG USER_UID=10001
ARG USER_GID=$USER_UID

RUN addgroup --gid $USER_GID $USERNAME \
    && adduser --uid $USER_UID --ingroup $USERNAME --disabled-password --shell /bin/bash --gecos "" $USERNAME

USER $USERNAME

WORKDIR /home/${USERNAME}

COPY . ./

RUN make download_and_verify_deps

ENTRYPOINT ["/bin/bash", "-c"]
CMD ["make", "run"]
