# syntax=docker/dockerfile:1

FROM node:20-alpine3.17
WORKDIR /home/build
COPY package.json /home/build/package.json
RUN npm install -g pnpm
RUN pnpm install
COPY . /home/build
RUN pnpm build

FROM golang:1.20
WORKDIR /home/build
COPY . /home/build
RUN git clone https://github.com/magefile/mage && \
    cd mage && \
    go run bootstrap.go
RUN cd ..
RUN mage -v
COPY --from=0 /home/build/dist /home/build/dist

VOLUME ./dist
