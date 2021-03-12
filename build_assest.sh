#!/usr/bin/env sh

rm -rf web/static/*

mkdir -p web/static/stylesheets
npx sass --load-path=. web/assets/main.scss web/static/stylesheets/all.css

mkdir -p web/static/javascript 
npx webpack --mode production --entry ./web/assets/main.js --output web/static/javascript/all.js

mkdir -p web/static/assets/fonts
cp node_modules/govuk-frontend/govuk/assets/fonts/* web/static/assets/fonts

mkdir -p web/static/assets/images
cp node_modules/govuk-frontend/govuk/assets/images/* web/static/assets/images
cp node_modules/@ministryofjustice/frontend/moj/assets/images/* web/static/assets/images