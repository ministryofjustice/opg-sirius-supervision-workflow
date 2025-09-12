FROM cypress/included:14.5.3

WORKDIR /root

RUN npm install "cypress-failed-log@2.9.2"
RUN npm install axe-core
RUN npm install cypress-axe

ENV CYPRESS_VIDEO=true
ENV CYPRESS_baseUrl=http://workflow:8888/supervision/workflow

COPY cypress.config.js .
COPY cypress cypress
