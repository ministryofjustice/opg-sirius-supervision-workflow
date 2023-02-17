## build json-server for mock API
FROM node:18.14.1-alpine3.16 as json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN yarn install

ENTRYPOINT ["yarn", "serve"]