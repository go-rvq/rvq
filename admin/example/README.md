# The example application

An admin of every feature of rvq/admin, on a PostgreSQL database
(`docker-compose.yaml`, `dev_env`): `go run .` (see `main.go`), the tests in
`integration/`.

## How-tos

- [A field of bits](docs/bits-field.md) — `InputDemo.Notifications1`: an
  unsigned integer whose bits are options, edited and shown as chips
  (`presets.BitsComponentFunc`), each option's label and description from the
  messages.
