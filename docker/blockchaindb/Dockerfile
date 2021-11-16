#Dockerfile for BlockchainDB
FROM ubuntu:20.04
RUN apt update && apt-get -y install tzdata
RUN apt -y install build-essential git jq iperf wget iputils-ping net-tools libcurl4-openssl-dev libtool m4 automake openssh-server iproute2 psmisc vim python dstat default-jdk unzip
RUN ssh-keygen -t rsa -N "" -f ~/.ssh/id_rsa && cd ~/.ssh && cp id_rsa.pub authorized_keys
ADD id_rsa.pub /
RUN cat /id_rsa.pub >> ~/.ssh/authorized_keys
RUN echo "StrictHostKeyChecking no" > ~/.ssh/config
RUN wget https://download.redis.io/releases/redis-6.0.9.tar.gz && tar -xzf redis-6.0.9.tar.gz && cd redis-6.0.9 && make && make install
RUN wget https://dl.google.com/go/go1.16.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.linux-amd64.tar.gz 
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $PATH:/usr/local/go/bin
RUN mkdir -p /go/src/github.com/ethereum
WORKDIR /go/src/github.com/ethereum
RUN git clone https://github.com/ethereum/go-ethereum.git
WORKDIR /go/src/github.com/ethereum/go-ethereum
RUN git checkout v1.8.23
ENV GO111MODULE=off
RUN make geth
RUN cp build/bin/geth /usr/local/go/bin/
RUN geth version
RUN mkdir -p /Data/
COPY .bin /bin
RUN echo "PermitUserEnvironment yes" >> /etc/ssh/sshd_config
CMD ["bash", "service ssh start"]
