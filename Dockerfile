FROM golang:1-bullseye

ARG active_user

RUN apt-get update && apt-get install -y \
    build-essential \
    git \
    golang-goprotobuf-dev \
    protobuf-compiler \
    sudo \
    wget

RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN apt-get -y install apt-transport-https ca-certificates curl gnupg2 software-properties-common
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
RUN add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable"
RUN apt-get update && \
      apt-cache policy docker-ce && \
      apt-get -y install docker-ce

RUN useradd -m  ${active_user} && echo "${active_user}:staff" |  chpasswd &&  usermod -aG sudo ${active_user}
RUN usermod -aG docker ${active_user}
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
RUN echo '${active_user} ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers

RUN mkdir -p /Users/${active_user}/work/mini-collector
RUN chown -R ${active_user}:staff /Users/${active_user}/work/mini-collector
RUN chown -R ${active_user}:staff /home/${active_user}
RUN gpasswd -a ${active_user} daemon

WORKDIR /Users/${active_user}/work/mini-collector
COPY ./scripts /Users/${active_user}/work/mini-collector/scripts

USER ${active_user}
