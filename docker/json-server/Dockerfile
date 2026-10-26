## build json-server for mock API
FROM node:24-alpine3.24@sha256:06d7adf3e0ef61a13f1d28682726366f97a15aa1888c934d3ecec86c33918b53 AS json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN npm ci --ignore-scripts

ENTRYPOINT ["npm", "run", "serve"]
