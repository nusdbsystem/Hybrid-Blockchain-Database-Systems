FROM ubuntu:20.04
ENV DEBIAN_FRONTEND=noninteractive
RUN apt update && apt-get -y install tzdata
RUN apt update && apt -y install build-essential git iperf wget iputils-ping net-tools libcurl4-openssl-dev libtool m4 automake openssh-server iproute2 psmisc vim python dstat default-jdk unzip
RUN ssh-keygen -t rsa -N "" -f ~/.ssh/id_rsa && cd ~/.ssh && cp id_rsa.pub authorized_keys
ADD id_rsa.pub /
RUN cat /id_rsa.pub >> ~/.ssh/authorized_keys
RUN echo "StrictHostKeyChecking no" > ~/.ssh/config
RUN wget https://dl.google.com/go/go1.15.6.linux-amd64.tar.gz && tar xf go1.15.6.linux-amd64.tar.gz
RUN mkdir gopath
ENV GOROOT /go
ENV GOPATH /gopath
ENV PATH $PATH:/go/bin
RUN mkdir /git
RUN wget https://download.redis.io/releases/redis-6.0.9.tar.gz && tar -xzf redis-6.0.9.tar.gz && cd redis-6.0.9 && make && make install
RUN wget https://archive.apache.org/dist/kafka/2.7.0/kafka_2.12-2.7.0.tgz && tar -xzf kafka_2.12-2.7.0.tgz
RUN wget https://github.com/RedBeardLab/rediSQL/releases/download/v1.1.1/RediSQL_v1.1.1_9b110f__release.so && mv RediSQL_v1.1.1_9b110f__release.so redisql.so && chmod u+x redisql.so
RUN wget https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-ubuntu2004-4.4.4.tgz && tar xf mongodb-linux-x86_64-ubuntu2004-4.4.4.tgz
COPY tendermint /usr/local/bin
COPY bin /bin
RUN echo "PermitUserEnvironment yes" >> /etc/ssh/sshd_config
CMD ["bash", "service ssh start"]
