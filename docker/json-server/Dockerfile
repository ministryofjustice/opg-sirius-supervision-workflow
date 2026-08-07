## build json-server for mock API
FROM node:lts-alpine3.24@sha256:5fa278c599dbba0c8f873d8717d50ecbb57c5ae6a53b7ab240c25135e0b65995 AS json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN npm ci --ignore-scripts

ENTRYPOINT ["npm", "run", "serve"]
