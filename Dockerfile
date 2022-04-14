FROM ubuntu:20.04

RUN apt-get update -y
RUN apt-get install -y vim wget git build-essential

WORKDIR /stockfish
RUN git clone https://github.com/official-stockfish/Stockfish.git

WORKDIR /stockfish/Stockfish/src
RUN make net
RUN make -j clean build ARCH=x86-64
RUN ln -s /stockfish/Stockfish/src/stockfish /usr/bin/stockfish

RUN wget https://dl.google.com/go/go1.16.5.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.16.5.linux-amd64.tar.gz

WORKDIR /mongosh
RUN wget https://downloads.mongodb.com/compass/mongodb-mongosh_1.3.1_amd64.deb
RUN dpkg -i mongodb-mongosh_1.3.1_amd64.deb

WORKDIR /go/src/github.com/kaushikc92/cagliostro
COPY cmd cmd
COPY go.mod .
COPY go.sum .
COPY main.go .
COPY pkg pkg

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

CMD ["tail", "-f", "/dev/null"]
#CMD ["sh", "-c", "/bin/bash"]
