FROM cypress/included:16.0.0@sha256:2ad59498de089a8eb03b6a3923c4e2f0801099f9be044ae2825531cb585e9370

WORKDIR /root

COPY package-lock.json package-lock.json
COPY package.json package.json
RUN npm ci --ignore-scripts


ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://local-proxy/supervision/workflow

COPY cypress.config.js .
COPY cypress cypress
