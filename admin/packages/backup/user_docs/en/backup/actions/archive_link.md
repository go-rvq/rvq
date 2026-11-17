# Download link (.tar.gz)

Makes a link that downloads the **whole state** as a `.tar.gz`, **with no
login**, for **24 hours**. It takes the backup elsewhere without setting up
a client: a browser, `curl` or `wget` download the file directly.

The link is shown **only once**, in the dialog, with the buttons **Copy**
and **Open** (opens it in the browser and the download starts). After that
the application keeps only a sum of it: no one gets it back. For another,
make another.

- It may be used **more than once** within the 24 hours. The `.tar.gz` is
  made as it is sent and does not resume: if the download falls, start
  again by the same link.
- Each use is recorded in **Links** (how many, the last address and the
  program), where the link can also be **revoked**.
- The link stops working when it expires, when it is revoked or when the
  state is deleted.
- The archive keeps permissions, times and symbolic links. A file changed
  on the server after the state is left out, and the `.tar.gz` ends with
  `BACKUP-SKIPPED.txt`, the list of them.

Whoever has the link has **all the data and secrets of the application**.
Send it only by a safe channel and revoke it when no longer needed.

```
curl -fLo backup.tar.gz 'https://your-site/admin/!backup/<id>/archive/<token>.tar.gz'
tar -xzf backup.tar.gz
```

## When it is available

With the state **ready**, for whoever has this page's **Download**
permission.
