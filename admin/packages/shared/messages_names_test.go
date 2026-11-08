package shared

import "testing"

// Every event, kind of request and state has its name, in each language.
func TestNamesTranslated(t *testing.T) {
	events := []string{EventConfirmationRequested, EventOwnerConfirmed, EventOwnerRejected, EventConfirmationExpired,
		EventInvited, EventAccepted, EventRejected, EventGranted, EventRoleChangeRequested, EventRoleChanged,
		EventRevoked, EventOwnerTransferRequest, EventOwnerTransferred, EventMemberAdded, EventMemberRemoved}
	kinds := []string{RequestShare, RequestRoleChange, RequestOwnerTransfer}
	states := []string{StatePending, StateConfirmed, StateAccepted, StateRejected, StateExpired}
	for lang, m := range map[string]*Messages{"en": Messages_en_US, "pt-BR": Messages_pt_BR} {
		check := func(code, name string) {
			// in English a name may be its code ("invited")
			if name == "" || lang != "en" && name == code {
				t.Errorf("%s: %q has no name", lang, code)
			}
		}
		for _, e := range events {
			check(e, m.EventName(e))
		}
		for _, k := range kinds {
			check(k, m.KindName(k))
		}
		for _, s := range states {
			check(s, m.StateName(s))
		}
	}
	if Messages_pt_BR.EventName(EventGranted) != "acesso concedido" || Messages_en_US.EventName("other") != "other" {
		t.Error("the names")
	}
}
