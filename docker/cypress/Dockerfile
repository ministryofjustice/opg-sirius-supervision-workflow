FROM cypress/included:13.8.0

WORKDIR /root

RUN npm install "cypress-failed-log@2.9.2"
RUN npm install axe-core
RUN npm install cypress-axe

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://workflow:8888

COPY cypress.config.js .
COPY cypress cypress
