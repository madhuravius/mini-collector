#!/usr/bin/env bash

CPU=$(dpkg --print-architecture)
GO_VERSION=1.20.1
PYENV_VERSION=3.6
ACTIVE_USER=$(whoami)

echo "Running for $ACTIVE_USER"

# copy protoc to path
sudo cp /tmp/protoc

# install go
echo "Installing go"
wget https://go.dev/dl/go${GO_VERSION}.linux-${CPU}.tar.gz
tar -xvf go${GO_VERSION}.linux-${CPU}.tar.gz
sudo cp -r go/ /usr/local
sudo cp go/bin/go /usr/local/bin
sudo mkdir -p /usr/local/go
sudo mkdir -p ~/go
sudo chown -R $ACTIVE_USER ~/go
rm -rf ./go

echo 'export GOROOT=/usr/local/go' >> ~/.bashrc
echo "export GOPATH=/home/$ACTIVE_USER/go" >> ~/.bashrc
echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.bashrc

# install pyenv
git clone https://github.com/pyenv/pyenv.git ~/.pyenv
echo "export PYENV_ROOT=\"/home/$ACTIVE_USER/.pyenv\"" >> ~/.bashrc
echo 'export PATH="$PYENV_ROOT/bin:$PATH"' >> ~/.bashrc
echo 'eval "$(pyenv init -)"' >> ~/.bashrc

PYENV_ROOT="/home/$ACTIVE_USER/.pyenv"
export PATH="$PYENV_ROOT/bin:$PATH:"

eval "$(pyenv init -)"
source ~/.bashrc

pyenv install -v $PYENV_VERSION
pyenv local 3.6
pyenv rehash

echo "Confirm that following output is correctly set: $PYENV_VERSION - is expected for python version"
python --version