# docker build . -t notionallabs/eve:latest
# docker run --rm -it notionallabs/eve /bin/sh

FROM archlinux

# install packages
RUN pacman -Sy --noconfirm go 
RUN pacman -Sy --noconfirm archlinux-keyring 
RUN pacman -Sy --noconfirm make gcc base jq

# set working directory
WORKDIR /app

# copy the current directory contents into the container at /usr/src/app
COPY . .

ENV PATH "$PATH:/root/go/bin"

RUN make install