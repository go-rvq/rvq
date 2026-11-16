# File system

The files of the server the site keeps: browse the folders, upload files and
folders (drag them onto the list), create, edit, rename, copy and delete;
download a file, or a folder as zip; extract an archive.

The same files may be reached from a computer through **WebDAV**, with the
address and the account shown on the page — as a network drive: in the
**WebDAV protocol** panel above the explorer, or by the button at the end of
its bar, which also says which folders are read-only.

## Permissions

The WebDAV asks the permissions of this page: seeing and downloading the
files ask `@get`; sending, changing, moving and deleting, `!write`.

Beside the folders of the server, the WebDAV shows folders of each user: the
draft of the **Site files** (`site-files`), with the permissions of editing
it in the admin.

## Terminal

The terminal button of the bar (or **F4**) opens a shell below the list, in
the open folder, and follows it when the folder changes. It sees only the
folders of the file system — in `$FILES_ROOT` —, each read or written as the
user may; nothing else of the server, nor the network (unless the
configuration allows it).

It asks this page's `!terminal` permission, in no set of permissions: give
it to who needs it. The server may turn it off — the WebDAV panel says why to
who has the permission.
