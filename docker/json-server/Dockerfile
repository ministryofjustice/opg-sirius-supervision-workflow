## build json-server for mock API
FROM node:14.17.5-alpine3.14 as json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN yarn install

ENTRYPOINT ["yarn", "serve"]