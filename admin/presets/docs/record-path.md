# Record paths and the page of an action

A **record path** says "this record, under these parents" — or this listing,
this page, this action — as text. It is kept in a field (the origin of a mail
sent, of a notification) and turned into the admin's address later, when the
admin that reads it may have moved the model in its menu.

## The format

It is the resource of the permissions by the unique names, the groups of the
menu left out (`perm.Verifier.PreferredResource`, see [Menu](menu.md)), made
by the same generator (`ModelPermissioner`), with the **whole** id of each
record between `<…>` (`ID.String`, path-escaped: a `>` in it is `%3E`):

```text
admin:posts:                          the listing
admin:posts:<01a0…_en-US>:            a record
admin:surveys:<1>:places:<2>:         a record under its parent's
admin:mail_sender:                    a singleton
admin:posts:<01a0…_en-US>:/preview:   a page of the record
admin:posts:/import:                  a page of the listing
admin:/site-files:                    a page of the admin
admin:mail_sender:!test_send_mail     an action (ActionPerm)
```

The parents go root first; a nested model is named by its own unique name.

## Making one

```go
mb := b.GetModel(&models.Post{})
id, _ := mb.ParseRecordID("01a0…_en-US")

mb.RecordPath(nil, id)                         // admin:posts:<01a0…_en-US>:
mb.ListingPath(nil)                            // admin:posts:
presets.PagePath(mb.RecordPath(nil, id), "/preview")
presets.ActionPath(mb.RecordPath(nil, id), "Publish") // …:!publish
b.PagePath("/site-files")                      // admin:/site-files:
```

`ModelInfo` has the same `RecordPath` and `ListingPath`. A child model takes
its parents' ids root first, one for each parent — zero for a singleton
parent —, as `ListingHref` takes them. Without a permissions builder
(`Builder.Permission`) there is no generator, and the path is `""`.

## Reading one

```go
t, err := b.ParseRecordPath(p)  // *RecordPathTarget{Models, Listing, Page, Action}
href := t.Href(b)               // "/admin/posts/01a0…_en-US"
href, err = b.RecordPathHref(p) // both at once
```

Each model is found by its unique name, wherever the menu put it. A path that
names no model, record, page or action there is `ErrRecordPath`
(`errors.Is`): a model removed, an action renamed — show the text, do not link.

`Href` is the detail of a record (under its parents), the listing, the
singleton's page, a page (`…/preview`), or **the page of an action**.

## The page of an action

The address of an action is the address of what it is of — a record's detail,
a listing, a page — and the suffix `/action:<name>` (`ActionSuffix`):

```text
/admin/users/01a0…/action:ResetPassword
/admin/mail-sender/action:TestSendMail
```

`Builder.ServeHTTP` takes the suffix off and draws the page it is of with the
action's form in its body (`ActionBuilder.PageComponent`), the form it shows
in its dialog; its button does it (`ParamActionPage` on the event):

- validation errors draw the form again in place (`ActionPagePortalName`);
- a flash of color `"error"` (`IsErrorFlash`) leaves the page as it is: the
  action did not happen;
- done, the browser goes back to the page the action is of (`PushState`),
  with the action's message.

So a link sent by mail ("Send the confirmation again") opens the action ready
to be done, and nothing more.

## Actions of the menu of a detail page

An action of the detail's menu (⋮) is done in its own dialog. When the detail
is a page — not a drawer or a dialog — its click says so
(`ParamToolbarAction` `"page"`, as a button of the toolbar does), and once done
the page is drawn again and the action's dialog closed: what the action
changed shows without reloading.
