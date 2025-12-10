# db-tools

This module provides admin interface to manage database backup tool.

Features:

- Manual backup creation with custom message;
- Auto backup creator (using Cron rules);
- List backups;
- Remove backups;
- Download backup files;
- Configure backup history preservation; 

<style>
figure {
  border: thin silver solid;
  display: flex;
  flex-flow: column;
  padding: 5px;
}

figcaption {
  background-color: #222222;
  color: white;
  font: italic smaller sans-serif;
  padding: 3px;
  text-align: center;
}
</style>

<figure>
    <img src="./docs/screen.png" alt="Sample Screen" />
    <figcaption>Sample Screen</figcaption>
</figure>

<figure>
    <img src="./docs/backup-preservation.png" alt="Configure Backup Preservation" />
    <figcaption>Configure Backup Preservation</figcaption>
</figure>

<figure>
    <img src="./docs/create-new-backup.png" alt="Manual backup creation with custom message" />
    <figcaption>Manual backup creation with custom message</figcaption>
</figure>