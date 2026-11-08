package user

import (
	"slices"
	"testing"

	"github.com/google/uuid"
)

// The members of a group as the form says: the others out, the new ones in,
// the ones kept untouched (their place among their groups kept).
func TestSetMembers(t *testing.T) {
	db := groupsDB(t)
	g, _ := CreateGroup(db, "G", "")
	other, _ := CreateGroup(db, "Other", "")
	ana, bia, caio := uuid.New(), uuid.New(), uuid.New()
	_ = AddMember(db, other.ID, bia, "root")
	if err := setMembers(db, g.ID, []uuid.UUID{ana, bia}, "root"); err != nil {
		t.Fatal(err)
	}
	if err := setMembers(db, g.ID, []uuid.UUID{bia, caio}, "root"); err != nil {
		t.Fatal(err)
	}
	ids, _ := MembersOf(db, g.ID)
	if len(ids) != 2 || !slices.Contains(ids, bia) || !slices.Contains(ids, caio) {
		t.Fatalf("the members: %v", ids)
	}
	if gs, _ := GroupsOf(db, bia); len(gs) != 2 || gs[0].Name != "Other" || gs[1].String() != "G" {
		t.Errorf("the groups of bia: %v", gs)
	}
}

func TestParseIDs(t *testing.T) {
	a, b := uuid.New(), uuid.New()
	got := parseIDs([]string{b.String(), "x", a.String(), b.String()})
	if !slices.Equal(got, []uuid.UUID{b, a}) {
		t.Errorf("%v", got)
	}
}
