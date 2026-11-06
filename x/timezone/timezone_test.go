package timezone

import (
	"context"
	"testing"
	"time"
)

func TestLocation(t *testing.T) {
	if Location(context.Background()) != time.UTC {
		t.Error("without one: UTC")
	}
	sp, err := Load("America/Sao_Paulo")
	if err != nil {
		t.Fatal(err)
	}
	if got := Location(WithLocation(context.Background(), sp)); got != sp {
		t.Errorf("the one carried: %v", got)
	}
	if loc, err := Load(" "); err != nil || loc != time.UTC {
		t.Errorf("an empty name: %v %v", loc, err)
	}
	if _, err := Load("Terra/Media"); err == nil {
		t.Error("a name not known")
	}
	if LoadOr("Terra/Media", sp) != sp || LoadOr("", nil) != time.UTC {
		t.Error("the default")
	}
}
