#!/bin/bash
# usage: ./generate_go_nis.sh > ${output file}
#
# copy files from container to host
# docker cp <containerId>:/file/path/within/container /host

echo '# auto:'
echo '# docker build -t go-nis:1.5 -f Dockerfile.go1.5.nis .'
echo '# run_dnis go-nis:1.5 go-nis'
echo '#'
echo '# manual: comment out ADD plugged'
echo '# docker build -t go-tmp:1.5 -f Dockerfile.go1.5.nis .'
echo '# run_dnis go-tmp:1.5 go-nis'
echo '#   open gvim'
echo '#   PlugInstall'
echo '# docker commit ${id} go-nis:1.5'
echo '# run_dnis go-nis:1.5 go-nis'
echo ""

echo 'FROM ubuntu:14.10'
echo 'MAINTAINER James Ho'
echo ""

echo '# get the needed packages/libs'
echo 'RUN apt-get update -y && apt-get install --no-install-recommends -y -q \'
echo '    openssh-client \'
echo '    curl \'
echo '    build-essential \'
echo '    bison \'
echo '    ca-certificates \'
echo '    git \'
echo '    mercurial \'
echo '    wget \'
echo '    vim-gnome \'
echo '    libpcap0.8-dev \'
echo '    firefox \'
echo '    graphviz \'
echo '    libc6-dbg \'
echo '    gdb'
echo ""

echo '# get go compiler'
echo '# RUN wget --no-check-certificate https://storage.googleapis.com/golang/go1.5rc1.linux-amd64.tar.gz'
echo '# RUN wget https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz'
echo 'COPY go1.5.1.linux-amd64.tar.gz /'
echo 'RUN tar -xzf go1.5.1.linux-amd64.tar.gz -C / && mv /go /goroot'
echo ""

echo '# Add user'
echo "ENV GROUP_ID $(id -g)"
echo "ENV GROUP $(id -ng)"
echo "ENV USER_ID $(id -u)"
echo "ENV USER $(id -nu)"
echo 'ENV HOME /home/${USER}'
echo ""

echo 'RUN groupadd -f -g ${GROUP_ID} ${GROUP} && \'
echo '    useradd -u ${USER_ID} -g ${GROUP} ${USER} && \'
echo '    usermod -a -G 0,4,27,30,46,1000 ${USER} && \'
echo '    echo "${USER} ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers && \'
echo '    mkdir -p ${HOME}'
echo ""

echo '# install vim-plug'
echo '# need to run 'PlugInstall' the first time running gvim'
echo 'RUN curl -fLo /root/.vim/autoload/plug.vim --create-dirs \'
echo '    https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim'
echo ""

echo '# copy my own files'
echo 'COPY .profile ${HOME}/.profile'
echo 'COPY .bashrc ${HOME}/.bashrc'
echo 'COPY .vimrc ${HOME}/.vimrc'
echo 'RUN cp -R /root/.vim ${HOME}'
echo '# we copy from the host for automation'
echo '# However, if we want to upgarde vim-go, we should skip this and run "PlugInstall"'
echo 'ADD plugged ${HOME}/.vim/plugged'
echo 'RUN chown -R ${USER}:${USER} ${HOME}'
echo '# COPY entrypoint.sh /entrypoint.sh'
echo ""

echo '# set the runtime env'
echo 'RUN mkdir -p /gopath/src'
echo ""

echo '# ENTRYPOINT /entrypoint.sh '
echo 'CMD ["/bin/bash"]'
