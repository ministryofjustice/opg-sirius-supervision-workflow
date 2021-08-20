FROM node:14.13.1-alpine3.12 as build

WORKDIR /app
EXPOSE 3000

COPY json-server .

RUN yarn install

ENTRYPOINT ["yarn", "serve"]