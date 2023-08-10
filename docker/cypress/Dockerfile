FROM cypress/included:11.2.0

WORKDIR /root

RUN npm install "cypress-failed-log@2.9.2"

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://workflow:8888

COPY cypress.config.js .
COPY cypress cypress
