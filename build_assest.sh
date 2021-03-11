#!/usr/bin/env bash

rm -rf web/static/*

mkdir -p web/static/{javascript,stylesheets}
npx sass --load-path=. web/assets/main.scss web/static/stylesheets/all.css
npx webpack --entry ./web/assets/main.js --output web/static/javascript/all.js

mkdir -p web/static/assets/{fonts,images}
cp node_modules/govuk-frontend/govuk/assets/fonts/* web/static/assets/fonts
cp node_modules/govuk-frontend/govuk/assets/images/* web/static/assets/images
cp node_modules/@ministryofjustice/frontend/moj/assets/images/* web/static/assets/images