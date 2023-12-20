FROM mcr.microsoft.com/devcontainers/go:1.21-bullseye

RUN apt update && apt install -y \
    vim \
    libusb-1.0-0-dev \
    git \
    build-essential \
    cmake \
    pkg-config \
    devscripts \
    curl \
    zsh

USER vscode

RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v1.1.5/zsh-in-docker.sh)" -- \
    -t robbyrussell

RUN go install go.uber.org/mock/mockgen@v0.3.0 && \
    go install github.com/rakyll/gotest@v0.0.6 && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2


RUN sudo mkdir -p /app/c && \
    sudo chown -R vscode:vscode /app/c && \
    git clone https://github.com/osmocom/rtl-sdr.git /app/c/rtl-sdr && \
    cd /app/c/rtl-sdr && sudo dpkg-buildpackage -us -uc && \
    sudo apt install -y /app/c/*.deb
