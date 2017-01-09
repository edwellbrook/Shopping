FROM node:6.9.2-alpine

RUN mkdir -p /src/app
WORKDIR /src/app

COPY package.json .
RUN npm install

COPY . /src/app

EXPOSE 8080
ENTRYPOINT npm start
