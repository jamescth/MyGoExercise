package main

const templ = `# auto:
# docker build -t go-nis:{{.Gover}} -f Dockerfile.go{{.Gover}}.nis .
# run_dnis go-nis:{{.Gover}} go-nis
#
# manual: comment out ADD plugged
# docker build -t go-tmp:{{.Gover}} -f Dockerfile.go{{.Gover}}.nis .
# run_dnis go-tmp:{{.Gover}} go-nis
#   open gvim
#   PlugInstall
# docker commit ${id} go-nis:{{.Gover}}
# run_dnis go-nis:{{.Gover}} go-nis

FROM {{.BaseImg}}
MAINTAINER {{.Maintainer}}

# get the needed packages/libs
RUN apt-get update -y && apt-get install --no-install-recommends -y -q \
    openssh-client \
    curl \
    build-essential \
    bison \
    ca-certificates \
    git \
    mercurial \
    wget \
    vim-gnome \
    libpcap0.8-dev \
    firefox \
    graphviz \
    libc6-dbg \
    gdb

##### Add user #####
ENV GROUP_ID {{.Gid}}
ENV GROUP {{.Username}}
ENV USER_ID {{.Uid}}
ENV USER {{.Username}}
ENV HOME /home/{{.Username}}

RUN groupadd -f -g ${GROUP_ID} ${GROUP}
RUN useradd -u ${USER_ID} -g ${GROUP} ${USER}
RUN usermod -a -G 0,4,27,30,46,1000 ${USER}
RUN echo "${USER} ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers 
RUN mkdir -p ${HOME}

##### Golang #####
# 1. Get golang compiler
# RUN wget --no-check-certificate https://storage.googleapis.com/golang/go{{.Gover}}.linux-amd64.tar.gz
RUN wget https://storage.googleapis.com/golang/go{{.Gover}}.linux-amd64.tar.gz
# COPY go{{.Gover}}.linux-amd64.tar.gz /
RUN tar -C /usr/local -xzf go{{.Gover}}.linux-amd64.tar.gz

# 2. Set up Golang runtime env
RUN mkdir -p ${HOME}/GoWorkspace/src

ENV GOROOT /usr/local/go
ENV GOPATH ${HOME}/GoWorkspace
ENV PATH $GOROOT/bin:$GOPATH/bin:.:$PATH

##### vim-plug #####
# https://github.com/junegunn/vim-plug
# files copying into the container have to be in this dir
# 1. download plug.vim
RUN curl -fLo /root/.vim/autoload/plug.vim --create-dirs \
    https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
# 2. Copy .vimrc from the host
COPY .vimrc ${HOME}/.vimrc

# 3. Run 'PlugInstall' or copy it from host
ADD plugged ${HOME}/.vim/plugged
ADD autoload ${HOME}/.vim/autoload

##### copy my own files #####
COPY .profile ${HOME}/.profile
COPY .bashrc ${HOME}/.bashrc

RUN chown -R ${USER}:${USER} ${HOME}

# COPY entrypoint.sh /entrypoint.sh
# ENTRYPOINT /entrypoint.sh
CMD ["/bin/bash"]

# in order to run
# docker run -ti --rm -u $(id -nu) -h jhodocker -e DISPLAY=$DISPLAY -v /tmp/.X11-unix:/tmp/.X11-unix -v /home/james/GoWorkspace:/home/james/GoWorkspeac $1
`
