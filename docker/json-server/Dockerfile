## build json-server for mock API
FROM node:20.2.0-alpine3.16 AS json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN yarn install

ENTRYPOINT ["yarn", "serve"]
