FROM cypress/included:16.1.1@sha256:e652cda84399aaf99c68d1d0322085dc3faa7ca9ace88d9bb4f0f045d2ce6764

WORKDIR /root

COPY package-lock.json package-lock.json
COPY package.json package.json
RUN npm ci --ignore-scripts


ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://local-proxy/supervision/workflow

COPY cypress.config.js .
COPY cypress cypress
