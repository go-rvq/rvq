#!/bin/sh
# The ECharts of the admin's charts: dist/echarts.min.js (window.RvqECharts),
# embedded in the binary (rvq js.ECharts) and served by admin/charts. Vue
# and Vuetify are the page's, not in it. Run when src/ or its dependencies
# change:
#   sh js/echarts/build.sh
set -e
cd "$(dirname "$0")"
bun install --frozen-lockfile
bun build src/bundle.js --minify --target=browser --format=iife --outfile=dist/echarts.min.js
