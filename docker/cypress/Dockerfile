FROM cypress/included:15.20.1@sha256:327d62ee6d1c9bd7b4b490bb2772c1b70d37d56647953a4ec3105a939b7c61ff

WORKDIR /root

COPY package-lock.json package-lock.json
COPY package.json package.json
RUN npm ci --ignore-scripts


ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://workflow:8888/supervision/workflow

COPY cypress.config.js .
COPY cypress cypress
