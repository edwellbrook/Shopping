FROM node:6.9.2-alpine

# Add build tools for node-gyp
RUN apk add --no-cache make gcc g++ python

RUN mkdir -p /src/app
WORKDIR /src/app

COPY package.json .
RUN npm install

COPY . /src/app

EXPOSE 8080
ENTRYPOINT npm start
