package db_tools

import (
	"strings"
	"testing"

	db_tools "github.com/go-rvq/rvq/x/packages/db-tools"
)

func TestPersistenceFormat(t *testing.T) {
	for _, tc := range []struct {
		m    *Messages
		per  db_tools.Persistence
		want string
	}{
		{Messages_en_US, db_tools.Persistence{Enabled: true, Days: 7, Months: 2, Other: db_tools.PersistenceOtherYears},
			", keep for "},
		{Messages_pt_BR, db_tools.Persistence{Enabled: true, Days: 7, Months: 2},
			", manter por "},
		{Messages_en_US, db_tools.Persistence{}, "Do not keep</span>."},
		{Messages_pt_BR, db_tools.Persistence{}, "Não manter</span>."},
	} {
		got, err := tc.m.Persistence.Format(&tc.per)
		if err != nil || !strings.Contains(string(got), tc.want) {
			t.Errorf("%+v: %q, %v; want it to hold %q", tc.per, got, err, tc.want)
		}
		t.Logf("%s", got)
	}
}
