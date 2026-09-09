#!/usr/bin/env bash
# prepare-deps.sh — build the git-dependency libraries that ship only source.
#
# `vuetify-pro-tiptap` is pinned to a GitHub commit that ships its `src/` but no
# built `lib/`, and its own `prepare` script (`pnpm build:lib`) is broken. bun
# re-extracts the git dependency fresh — with no `lib/` — every time the
# dependency graph changes, so a plain `bun install` silently leaves the vite
# build unable to resolve the package. This step rebuilds that `lib/` in place
# against whatever store hash bun currently resolved, using bun + vite directly.
#
# Idempotent: it is a no-op once `lib/` exists.
set -euo pipefail

BUN=${BUN:-bun}
here=$(cd "$(dirname "$0")/.." && pwd)

# Resolve the currently-installed package directory (follows bun's store hash).
dep=$(readlink -f "$here/node_modules/vuetify-pro-tiptap" 2>/dev/null || true)
if [ -z "$dep" ] || [ ! -f "$dep/package.json" ]; then
  echo "prepare-deps: vuetify-pro-tiptap not installed; run '$BUN install' first" >&2
  exit 1
fi

if [ -f "$dep/lib/vuetify-pro-tiptap.js" ] && [ -f "$dep/lib/vuetify-pro-tiptap.umd.cjs" ]; then
  echo "prepare-deps: vuetify-pro-tiptap/lib already built"
  exit 0
fi

echo "prepare-deps: building vuetify-pro-tiptap/lib in $dep"
(
  cd "$dep"
  # Install the package's own devDeps (rollup-plugin-pure, vite-plugin-dts, …).
  # Its `prepare` (pnpm build:lib) fails; ignore that — the deps still install.
  "$BUN" install || true
  # Build the library with vite directly, bypassing the broken prepare script.
  "$BUN" x vite build
)
echo "prepare-deps: done"
