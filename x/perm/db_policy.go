package perm

import (
	"strings"
	"time"

	"github.com/go-rvq/rvq/admin/utils/uuidkey"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type DBPolicyBuilder struct {
	db            *gorm.DB
	model         DBPolicy
	loadFrequency time.Duration
}

func NewDBPolicy(db *gorm.DB) *DBPolicyBuilder {
	uuidkey.MustRegister(db) // its records have UUID keys
	return &DBPolicyBuilder{
		db:            db,
		model:         DefaultDBPolicy{},
		loadFrequency: time.Minute,
	}
}

func (dpb *DBPolicyBuilder) Model(m DBPolicy) *DBPolicyBuilder {
	dpb.model = m
	return dpb
}

func (dpb *DBPolicyBuilder) LoadFrequency(d time.Duration) *DBPolicyBuilder {
	dpb.loadFrequency = d
	return dpb
}

type DefaultDBPolicy struct {
	uuidkey.Model

	// ReferID names what the policy belongs to: a role's key, or a free-form
	// reference ("share:<subject>", "adhoc"). Text, and no foreign key.
	ReferID   string
	Subject   string
	Effect    string
	Actions   pq.StringArray `gorm:"type:text[]"`
	Resources pq.StringArray `gorm:"type:text[]"`

	// SharedID, when set, marks the policy as belonging to a resource "share"
	// (e.g. an organization or project shared with another user). It identifies
	// the share so the sharing UI can group/manage the policies it created,
	// separately from ad-hoc grants made in the permission manager.
	SharedID *uuid.UUID `gorm:"type:uuid;index"`
}

func (p DefaultDBPolicy) LoadDBPolicies(db *gorm.DB, startFrom *time.Time) (toUpdateOrCreate []*PolicyBuilder, toDelete []*PolicyBuilder) {
	var ps []DefaultDBPolicy
	if startFrom == nil || startFrom.IsZero() {
		db.Find(&ps)
	} else {
		db.Unscoped().Where("updated_at >= ? or deleted_at >= ?", startFrom, startFrom).Find(&ps)
	}

	for _, p := range ps {
		if p.DeletedAt.Valid {
			toDelete = append(toDelete, p.ToPolicy())
		} else {
			toUpdateOrCreate = append(toUpdateOrCreate, p.ToPolicy())
		}
	}
	return
}

func (p DefaultDBPolicy) ToPolicy() *PolicyBuilder {
	var res []string
	for _, r := range p.Resources {
		res = append(res, SplitResources(r)...)
	}
	return PolicyFor(p.Subject).WhoAre(p.Effect).ToDo(p.Actions...).On(res...).ID(p.ID.String())
}

// SplitResources is the resources of s, separated by commas — out of the
// alternatives of a pattern, "{a,b}", and of a record, "<…>", whose commas
// are theirs —, trimmed, the empty ones left out.
func SplitResources(s string) (res []string) {
	depth, start := 0, 0
	add := func(end int) {
		if r := strings.TrimSpace(s[start:end]); r != "" {
			res = append(res, r)
		}
	}
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++
		case '{', '<', '[':
			depth++
		case '}', '>', ']':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				add(i)
				start = i + 1
			}
		}
	}
	add(len(s))
	return
}
