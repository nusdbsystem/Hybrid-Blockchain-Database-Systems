#Dockerfile for BlockchainDB
FROM ubuntu:21.04
RUN apt update && apt-get -y install tzdata
RUN apt update && apt -y install build-essential git jq libc6 python3-pip \
    iperf wget iputils-ping net-tools libcurl4-openssl-dev libtool m4 automake openssh-server iproute2 psmisc \
    vim python dstat unzip nodejs npm curl autoconf make g++  protobuf-compiler \
    && pip3 install solc-select \
    && solc-select install 0.4.24 \
    && solc-select use 0.4.24
ENV HOME /root
#RUN ssh-keygen -f /etc/ssh/ssh_host_rsa_key -N ""
RUN ssh-keygen -t rsa -N "" -f /root/.ssh/id_rsa && cd /root/.ssh && cp id_rsa.pub authorized_keys
ADD id_rsa.pub /
RUN cat /id_rsa.pub >> ~/.ssh/authorized_keys
RUN echo "StrictHostKeyChecking no" > /root/.ssh/config
RUN echo "PermitUserEnvironment yes" >> /etc/ssh/sshd_config
RUN echo "root:newpass" | chpasswd
RUN service ssh start
CMD ["/usr/sbin/sshd","-D"]

RUN wget https://download.redis.io/releases/redis-6.0.9.tar.gz && tar -xzf redis-6.0.9.tar.gz && cd redis-6.0.9 && make && make install
RUN wget https://dl.google.com/go/go1.16.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.linux-amd64.tar.gz 
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $PATH:/usr/local/go/bin
RUN go version

RUN mkdir -p /go/src/github.com/ethereum
WORKDIR /go/src/github.com/ethereum
RUN git clone https://github.com/ethereum/go-ethereum.git
WORKDIR /go/src/github.com/ethereum/go-ethereum
RUN git checkout v1.8.23
ENV GO111MODULE=off
RUN make geth
RUN cp build/bin/geth /usr/local/go/bin/
RUN geth version
RUN make devtools

RUN mkdir -p /Data/
RUN mkdir -p ${HOME}/BlockchainDB/scripts/
RUN mkdir -p ${HOME}/BlockchainDB/config/
RUN mkdir -p ${HOME}/BlockchainDB/logs/
RUN mkdir -p ${HOME}/BlockchainDB/bin/
# COPY .scripts ${HOME}/BlockchainDB/scripts
COPY .bin ${HOME}/BlockchainDB/bin
WORKDIR ${HOME}/BlockchainDB

