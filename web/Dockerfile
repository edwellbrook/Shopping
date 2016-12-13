FROM node:6.9.2-alpine

RUN mkdir -p /src/app
WORKDIR /src/app

VOLUME /src/app/node_modules
COPY package.json .
RUN npm install

COPY . .

EXPOSE 8080
ENTRYPOINT npm start
