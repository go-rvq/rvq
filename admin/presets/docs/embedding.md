# Embedding the admin

The admin can run inside another page — a site's "edit this page" dialog, a
dashboard's panel — in an `<iframe>`. There the room is the page's, not the
admin's: the admin starts with its **side menu closed**, and the button of its
toolbar (☰) opens it, as on a narrow screen.

## The signal: the window's name

An admin that runs in a window named `presets.EmbeddedWindowName`
(`"rvq-embedded"`) starts with the side menu closed:

```html
<iframe name="rvq-embedded" src="/admin/content/posts/42_en-US"></iframe>
```

The name is the iframe's `name` attribute — inside it, `window.name` —, not a
parameter of the URL: it stays the window's while the admin navigates in it (a
listing, a detail, a form), so every page of the admin in that iframe starts
the same. Nothing is sent to the server for it; the layout reads it in the
browser, when it sets its initial state (`vars.navDrawer`).

An admin in any other window — its own tab, an iframe with another name — starts
with the menu open, as ever.

## Loading it when it is needed

An iframe with a `src` loads the admin with the page that holds it. When the
admin is opened by the user — a button, a dialog —, keep the address in
`data-src` and set `src` on the first opening:

```html
<a href="/admin/content/posts/42_en-US" data-open-admin>Edit this page</a>
<dialog data-admin-dialog>
  <a href="/admin/content/posts/42_en-US" target="_blank" rel="noopener">Open in a new tab</a>
  <button type="button" data-close>Close</button>
  <iframe name="rvq-embedded" data-src="/admin/content/posts/42_en-US"></iframe>
</dialog>
<script>
  const dialog = document.querySelector('[data-admin-dialog]');
  const frame = dialog.querySelector('iframe');
  document.querySelector('[data-open-admin]').addEventListener('click', (e) => {
    // the middle button, Ctrl/⌘/Shift: the link's — a new tab
    if (e.button !== 0 || e.ctrlKey || e.metaKey || e.shiftKey || e.altKey) return;
    e.preventDefault();
    if (!frame.src) frame.src = frame.dataset.src;
    dialog.showModal();
  });
  dialog.querySelector('[data-close]').addEventListener('click', () => dialog.close());
</script>
```

The link stays a link: without the script, or with the middle button, the
admin opens in a tab of its own.

## The same origin

The embedding page and the admin must be of the same origin for the session to
be the user's (the login cookie) and for the page to reach into the iframe. The
presets send no header that forbids framing; an application that serves the
admin to the world should send `X-Frame-Options: SAMEORIGIN` (or
`Content-Security-Policy: frame-ancestors 'self'`) — which keeps its own pages
able to embed it and every other site not.
