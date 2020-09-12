FROM debian:buster
RUN apt-get update -y && apt-get upgrade -y && apt-get autoremove -y
RUN apt-get install -y wget
RUN apt-get install -y git
# RUN apt-get install -y curl

# Install golang
RUN wget https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz
RUN tar -xvf go1.14.1.linux-amd64.tar.gz
RUN mv go /usr/local
ENV GOROOT=/usr/local/go
ENV GOPATH=/home/go
RUN mkdir $GOPATH && mkdir $GOPATH/bin && mkdir $GOPATH/pkg && mkdir $GOPATH/src
ENV PATH="$GOPATH/bin:$GOROOT/bin:${PATH}"

# Install node.js
RUN apt-get install -y curl
RUN apt-get install -y build-essential
RUN curl -sL https://deb.nodesource.com/setup_12.x | bash -
RUN apt-get -f install -y nodejs
RUN apt-get -f install -y aptitude
RUN aptitude install -y npm

# Install Typescript and other utilities
RUN npm install -g typescript
RUN npm install -g esbuild-linux-64

# Install golang tools
RUN go get github.com/mdempsky/gocode
RUN go get github.com/uudashr/gopkgs/v2/cmd/gopkgs
RUN go get github.com/ramya-rao-a/go-outline
RUN go get github.com/acroca/go-symbols
RUN go get golang.org/x/tools/cmd/guru
RUN go get golang.org/x/tools/cmd/gorename
RUN go get github.com/cweill/gotests/...
RUN go get github.com/fatih/gomodifytags
RUN go get github.com/josharian/impl
RUN go get github.com/davidrjenni/reftools/cmd/fillstruct
RUN go get github.com/haya14busa/goplay/cmd/goplay
RUN go get github.com/godoctor/godoctor
RUN go get github.com/go-delve/delve/cmd/dlv
RUN go get github.com/stamblerre/gocode
RUN go get github.com/rogpeppe/godef
RUN go get golang.org/x/tools/cmd/goimports
RUN go get golang.org/x/lint/golint
RUN go get golang.org/x/tools/gopls