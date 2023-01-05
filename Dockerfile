FROM node:latest

#Install speedtest-cli
RUN apt-get install curl
RUN curl -s https://packagecloud.io/install/repositories/ookla/speedtest-cli/script.deb.sh | bash
RUN apt-get install speedtest

WORKDIR /home/node/app
