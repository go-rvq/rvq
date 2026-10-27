# A field of bits — `InputDemo.Notifications1`

**Featured Models Management › Input Demos**, any record: the field
**Notifications1** is a field of bits. It is an unsigned integer whose bits are
options — here, the channels a notification goes by —, drawn by
`presets.BitsComponentFunc` as chips:

- **to edit**, one chip an option: a click marks or unmarks it, and the field's
  value is the sum of the bits marked;
- **in the detail**, the set ones marked (a check), the others outlined;
- each chip has a `?` at its right: on hover or on click it opens the option's
  **description** (the click does not toggle the chip);
- the field's label above the chips, its hint below.

The test `integration/input_demo_bits_test.go` edits, shows and saves it.

## 1. The type and its bits

`models/input_demo.go` — a named unsigned integer and its constants, one bit
each:

```go
// NotifyChannels are the channels a notification goes by, a bit each.
type NotifyChannels uint8

const (
	NotifyEmail NotifyChannels = 1 << iota // 1
	NotifySMS                              // 2
	NotifyPush                             // 4
)

type InputDemo struct {
	// …
	Notifications1 NotifyChannels `gorm:"not null;default:0"`
}
```

The column is a small integer; `0` is no option. Test a bit with
`v&NotifySMS == NotifySMS`.

## 2. The label and the description of each option: messages

`admin/messages.go` — the texts are messages of the module, in every language it
has (en, ja, zh here), never written in the code:

```go
type Messages struct {
	DemoNotifyEmail     string
	DemoNotifyEmailDesc string
	DemoNotifySMS       string
	DemoNotifySMSDesc   string
	DemoNotifyPush      string
	DemoNotifyPushDesc  string
	// …
}

var Messages_en_US = &Messages{
	DemoNotifyEmail:     "Email",
	DemoNotifyEmailDesc: "Notifies by email, to the address of the account.",
	DemoNotifySMS:       "SMS",
	DemoNotifySMSDesc:   "Notifies by a text message, to the phone of the account.",
	DemoNotifyPush:      "Push",
	DemoNotifyPushDesc:  "Notifies on the devices where the account is signed in.",
	// …
}
```

## 3. The options and the component

`admin/input_demo_config.go` — a `presets.BitOption` a bit: its value, its label
and its description, each read from the messages of the request's language:

```go
func notifyChannelOptions() []presets.BitOption {
	msgs := func(ctx context.Context) *Messages {
		return i18n.MustGetModuleMessages(ctx, I18nExampleKey, Messages_en_US).(*Messages)
	}
	opt := func(bit models.NotifyChannels, label, desc func(*Messages) string) presets.BitOption {
		return presets.BitOption{Bit: uint64(bit),
			Label:       func(ctx context.Context) string { return label(msgs(ctx)) },
			Description: func(ctx context.Context) string { return desc(msgs(ctx)) }}
	}
	return []presets.BitOption{
		opt(models.NotifyEmail, func(m *Messages) string { return m.DemoNotifyEmail }, func(m *Messages) string { return m.DemoNotifyEmailDesc }),
		opt(models.NotifySMS, func(m *Messages) string { return m.DemoNotifySMS }, func(m *Messages) string { return m.DemoNotifySMSDesc }),
		opt(models.NotifyPush, func(m *Messages) string { return m.DemoNotifyPush }, func(m *Messages) string { return m.DemoNotifyPushDesc }),
	}
}
```

The same component edits and shows — give it to the editing and to the detail:

```go
ed := inputDemo.Editing( /* … */ "Notifications1")

notifications := presets.BitsComponentFunc(notifyChannelOptions()...)
ed.Field("Notifications1").ComponentFunc(notifications)
inputDemo.Detailing().Field("Notifications1").ComponentFunc(notifications)
```

Nothing else: the form posts the sum of the bits as the field's value, and the
default setter writes it.

## What to keep in mind

- **Bits, not values**: each constant is `1 << n`; a combination is their sum
  (`NotifyEmail | NotifyPush` is `5`).
- **The order of the chips** is the order of the options given.
- **A hint for the field** (`InputDemoNotifications1_Hint` in the messages of the
  models) is shown below the chips; the descriptions are the options'.
- **Hiding other fields by a bit** — as the mail sender hides the IMAP user when
  it takes the SMTP's — is a `v-if` on the form's value:
  `((Number(form["Inbound.Reuse"]) || 0) & 1) === 0`.

See also: [presets: a field of bits](../../presets/docs/fields.md#a-field-of-bits--bitscomponentfunc),
and the page **Editing Customizations › A field of bits** of the documentation
site (`admin/docs/docsrc`), with a live demo.
