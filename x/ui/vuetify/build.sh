mkdir -p dist
curl https://cdn.jsdelivr.net/npm/vuetify@3.8.1/dist/vuetify.min.js > dist/vuetify.min.js
curl https://cdn.jsdelivr.net/npm/vuetify@3.8.1/dist/vuetify.min.css > dist/vuetify.min.css
curl https://cdn.jsdelivr.net/npm/vuetify@3.8.1/dist/vuetify-labs.min.js > dist/vuetify-labs.min.js
curl https://cdn.jsdelivr.net/npm/vuetify@3.8.1/dist/vuetify-labs.min.css > dist/vuetify-labs.min.css
cd vuetifyjs && make build
