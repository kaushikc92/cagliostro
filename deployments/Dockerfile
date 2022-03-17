FROM ubuntu:18.04

RUN apt-get update -y
RUN apt-get install -y vim wget git

RUN wget https://dl.google.com/go/go1.16.5.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.16.5.linux-amd64.tar.gz

WORKDIR /go/src/github.com/kaushikc92/chesstrocity

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

CMD ["tail", "-f", "/dev/null"]
#CMD ["sh", "-c", "/bin/bash"]
