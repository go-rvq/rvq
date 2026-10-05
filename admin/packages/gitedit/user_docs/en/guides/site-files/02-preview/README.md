# Preview

**Preview the site** opens, in a new tab, the site as your draft makes it:
its templates, styles and images, with the site's real content (the pages and
posts published).

![The site as the draft makes it](images/preview.png)

- Browse at will: the links stay in the preview (`/admin/site-preview/…`).
- Only you see your draft (and whoever has the link of its public preview);
  the visitors still see the site.
- Saved a file? Reload the page of the preview.

## Public preview

To show the draft to someone with no account in the admin — a client, a
reviewer —, on the **Your draft** tab use **Make the preview public**. Choose
until when it lasts (empty: until you end it) and copy the **Link**: whoever
has it sees the site as your draft makes it, with no login.

- The link is an address nobody guesses (`/_preview/<token>/…`) and search
  engines do not index it.
- It shows the draft as it is now — its uncommitted changes too.
- **End the public preview** stops the link at once. Making it public again
  makes a new link; the former stops working.
- It asks the permission `!public_preview` (and `!preview`). Whoever sees every
  draft (`!drafts`) may end anyone's public preview.
