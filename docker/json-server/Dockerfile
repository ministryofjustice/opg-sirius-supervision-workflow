## build json-server for mock API
FROM node:24.1.0-alpine3.20 AS json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN yarn install

ENTRYPOINT ["yarn", "serve"]
