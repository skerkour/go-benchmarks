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

LABEL maintainer="Sylvain Kerkour <https://kerkour.com>"
LABEL homepage=https://github.com/skerkour/go-benchmarks
LABEL org.opencontainers.image.name=go-benchmarks
LABEL repository=https://github.com/skerkour/go-benchmarks
LABEL org.opencontainers.image.source=https://github.com/skerkour/go-benchmarks
