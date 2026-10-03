# Permissions

Each action of the page is a permission of it, given to the roles in **User
Roles** (see the policies and permissions):

| Action | Permission |
| --- | --- |
| See the page and the files | `@get` |
| Change the files of the draft | `!edit` |
| Commit | `!commit` |
| Update | `!update` |
| Publish | `!publish` |
| Discard | `!discard` |
| Start over | `!reset` |
| Preview | `!preview` |

So publishing can stay with fewer people than editing: whoever edits commits,
and someone else, who reviews, publishes.
