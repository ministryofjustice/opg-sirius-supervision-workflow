FROM cypress/included:15.18.1@sha256:9d3e6ac7d4a725e802080fb45f4c8d08d1a0eb78b3907232ffb3a42cd1f05e33

WORKDIR /root

COPY yarn.lock yarn.lock
COPY package.json package.json
RUN yarn install --frozen-lockfile --ignore-scripts


ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://workflow:8888/supervision/workflow

COPY cypress.config.js .
COPY cypress cypress
