FROM cypress/included:16.1.0@sha256:9618e4f3798bfc5434aeab1283813653437bbacc5d15922433b37f10b33aad49

WORKDIR /root

COPY package-lock.json package-lock.json
COPY package.json package.json
RUN npm ci --ignore-scripts


ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://local-proxy/supervision/workflow

COPY cypress.config.js .
COPY cypress cypress
