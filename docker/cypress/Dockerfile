FROM cypress/included:9.7.0

WORKDIR /root

RUN npm install "cypress-failed-log@2.9.2"

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://app:8888

COPY cypress.json .
COPY cypress cypress
