FROM cypress/included:15.17.0@sha256:9378699679f7d6ddb63c572866d4f11ec822abacb099345e89f5b86381efd53f

WORKDIR /root

COPY yarn.lock yarn.lock
COPY package.json package.json
RUN yarn install --frozen-lockfile --ignore-scripts


ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://workflow:8888/supervision/workflow

COPY cypress.config.js .
COPY cypress cypress
