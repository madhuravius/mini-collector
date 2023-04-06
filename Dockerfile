FROM ubuntu:trusty

ARG active_user
ARG PYENV_VERSION=3.6

RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    git \
    libbz2-dev \
    libncurses5-dev \
    libreadline-dev \
    libssl-dev \
    libsqlite3-dev \
    liblzma-dev \
    python3 \
    python3-pip \
    python3-openssl \
    sudo \
    wget \
    zip \
    zlib1g-dev # needed for pyenv / python installations

RUN useradd -m  ${active_user} && echo "${active_user}:staff" |  chpasswd &&  usermod -aG sudo ${active_user}
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
RUN echo '${active_user} ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers

RUN mkdir -p /Users/${active_user}/work/mini-collector
RUN chown -R ${active_user}:staff /Users/${active_user}/work/mini-collector
RUN chown -R ${active_user}:staff /home/${active_user}

WORKDIR /Users/${active_user}/work/mini-collector
COPY ./scripts /Users/${active_user}/work/mini-collector/scripts

USER ${active_user}

# These are used to fully mimic CI environments with our versioned/dated libraries
RUN ./scripts/LOCAL-install-protoc.sh
RUN ./scripts/LOCAL-install-pyenv-go.sh
