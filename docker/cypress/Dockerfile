FROM cypress/included:5.2.0

RUN npm install axe-core
RUN npm install cypress-axe

WORKDIR /root

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://app:8888

COPY cypress.json .
COPY cypress cypress