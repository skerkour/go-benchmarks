FROM golang:alpine

RUN apk update && \
    apk upgrade --no-cache && \
    apk add --no-cache ca-certificates make curl wget git

ARG USERNAME=benchmarks
ARG USER_UID=10001
ARG USER_GID=$USER_UID

ENV TZ="Europe/London"
ENV LC_ALL="en_US.UTF-8"
ENV LANG="en_US.UTF-8"
ENV LANGUAGE="en_US:en"

RUN echo "${TZ}" > /etc/timezone
RUN update-ca-certificates

RUN adduser \
    --disabled-password \
    --gecos "" \
    --shell "/bin/sh" \
    --uid "${USER_UID}" \
    "${USERNAME}"

USER $USERNAME

WORKDIR /home/${USERNAME}/benchmarks

COPY /home/runner/work/go-benchmarks/go-benchmarks/.git ./.git
COPY . ./
RUN git config --global --add safe.directory /home/benchmarks/benchmarks

RUN make download_and_verify_deps

ENTRYPOINT ["/bin/sh", "-c"]
CMD ["make", "run"]
