package shared

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/perms"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// The sharing in the admin: a model shared (Service.Install) gets the
// actions of its detail — Share, Accesses (who reaches it: change a role,
// revoke), Access history, Transfer ownership — and, with ReadScope, its
// listing and detail limited to the records the user sees; the page where
// the owner confirms a request (InstallConfirm) and the user's invites
// (NotificationContent, NotificationCount).

// Choice is a user or a group to choose.
type Choice struct{ Value, Title string }

// Directory lists the users and the groups to choose from (a Users that
// is one: UsersOf).
type Directory interface {
	ListUsers(ctx context.Context) []Choice
	ListGroups(ctx context.Context) []Choice
}

// ModelOptions say how a model is shared in the admin.
type ModelOptions struct {
	// Resource is the model's (its Name, else the model's id).
	Resource Resource
	// Roles are the roles offered (seeded once, with None).
	Roles []RoleSeed
	// ReadScope limits the listing and the detail to the records the user
	// sees.
	ReadScope bool
	// Bypass says who sees every record (the administrator).
	Bypass func(r *http.Request) bool
}

// The actions of the detail of a model shared.
const (
	ActionShare         = "share_access"
	ActionAccesses      = "accesses"
	ActionHistory       = "access_history"
	ActionTransferOwner = "transfer_owner"
)

const (
	evShareRole   = "shared_role"
	evShareRevoke = "shared_revoke"
	paramSubject  = "subject"
	paramRole     = "role"
	paramRecord   = "record"
	fieldSubject  = "shared_subject"
	fieldRole     = "shared_role"
	fieldKeep     = "shared_keep"
	fieldNewOwner = "shared_new_owner"
)

// originOf is where the request came from.
func originOf(r *http.Request) Origin {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	return Origin{IP: ip, UserAgent: r.UserAgent()}
}

// currentUser is the id of the user of the request ("" none).
func currentUser(r *http.Request) string {
	u := user.GetCurrentUser(r)
	if u == nil || u.Anonymous() {
		return ""
	}
	return u.GetID().String()
}

// Install shares the model mb in the admin, as o says.
func (s *Service) Install(mb *presets.ModelBuilder, o ModelOptions) error {
	registerMessages(mb.Builder().I18n())
	res := o.Resource
	if res.Name == "" {
		res.Name = mb.Info().URIName()
	}
	s.Register(res)
	if err := SeedRoles(s.o.DB, uuid.Nil, res.Name, o.Roles...); err != nil {
		return err
	}
	// the policies of the permissions of the admin kept with the accesses:
	// the subject (user:<id>, group:<id>) the role's actions on the record
	// (a model on top: under another, its parent's permissions decide)
	if mb.Parent() == nil {
		s.OnAccess(res.Name, func(tx *gorm.DB, a *Access, role *ShareTemplate, removed bool) error {
			ref := "sharedaccess:" + a.ID.String()
			if removed || role == nil || role.None {
				return perms.Revoke(tx, ref)
			}
			rid, err := mb.ParseRecordID(a.RecordID)
			if err != nil {
				return err
			}
			_ = perms.Revoke(tx, ref)
			_, err = perms.GrantShared(tx, ref, a.Subject, perms.RecordResource(mb, rid)+"*", []string(role.Actions), a.ID)
			if err == nil {
				if pb := mb.Builder().GetPermission(); pb != nil {
					from := s.o.Now().Add(-time.Minute)
					pb.LoadDBPoliciesToMemory(tx, &from)
				}
			}
			return err
		})
	}
	recordOf := func(id string) string {
		rid, err := mb.ParseRecordID(id)
		if err != nil || rid.IsZero() {
			return ""
		}
		return strings.Join(rid.StringValues(), "_")
	}
	bypass := func(r *http.Request) bool { return o.Bypass != nil && o.Bypass(r) }

	if o.ReadScope {
		mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
			return do.(*gorm2op.DataOperatorBuilder).
				WithReadCallbacks(func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
					cb.Pre(func(state *gorm2op.CallbackState) error {
						if state.Ctx == nil || bypass(state.Ctx.R) {
							return nil
						}
						db, err := s.Scope(state.Ctx.R.Context(), state.DB, currentUser(state.Ctx.R), res.Name)
						if err != nil {
							return err
						}
						state.DB = db
						return nil
					})
				})
		})
	}

	verifier := func(ctx *web.EventContext) *perm.Verifier {
		return mb.Permissioner().ReqListDo(ctx.R, PermShare)
	}
	dt := mb.Detailing()

	// share: a user or a group, a role — asked to the owner by e-mail
	dt.Action(ActionShare).Icon("mdi-account-multiple-plus").
		SetI18nLabel(func(ctx context.Context) string { return msgs(ctx).ShareTitle }).
		SetVerifier(verifier).
		ComponentFunc(func(id string, ctx *web.EventContext) (h.HTMLComponent, error) {
			m := msgs(ctx.Context())
			return h.Div(
				v.VAutocomplete().Label(m.Subject).Items(s.choices(ctx.R.Context(), m)).ItemTitle("title").ItemValue("value").
					Attr(web.VField(fieldSubject, "")...).Attr("data-shared-subject", true),
				v.VSelect().Label(m.Role).Items(s.roleChoices(ctx.R.Context(), res.Name)).ItemTitle("title").ItemValue("value").
					Attr(web.VField(fieldRole, "")...).Attr("data-shared-role", true),
			), nil
		}).
		UpdateFunc(func(id string, ctx *web.EventContext) error {
			roleID, err := uuid.Parse(ctx.R.FormValue(fieldRole))
			if err != nil {
				return ErrRole
			}
			if _, err := s.Share(ctx.R.Context(), res.Name, recordOf(id), ctx.R.FormValue(fieldSubject), roleID,
				currentUser(ctx.R), originOf(ctx.R)); err != nil {
				return err
			}
			ctx.Flash = msgs(ctx.Context()).RequestSent
			return nil
		})

	// who reaches the record: the owner, the accesses (own: change, revoke;
	// inherited: shown)
	dt.Action(ActionAccesses).Icon("mdi-account-key-outline").
		SetI18nLabel(func(ctx context.Context) string { return msgs(ctx).AccessesTitle }).
		ComponentFunc(func(id string, ctx *web.EventContext) (h.HTMLComponent, error) {
			return s.accessesBody(ctx, mb, res.Name, id, recordOf(id))
		})

	mb.RegisterEventFunc(evShareRole, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		if ver := verifier(ctx); ver != nil && ver.Denied() {
			return r, perm.PermissionDenied
		}
		roleID, err := uuid.Parse(ctx.R.FormValue(paramRole))
		if err != nil {
			return r, ErrRole
		}
		m := msgs(ctx.Context())
		req, err := s.ChangeRole(ctx.R.Context(), res.Name, recordOf(ctx.R.FormValue(paramRecord)), ctx.R.FormValue(paramSubject),
			roleID, currentUser(ctx.R), originOf(ctx.R))
		switch {
		case err != nil:
			presets.ShowMessage(&r, err.Error(), "error")
		case req != nil:
			presets.ShowMessage(&r, m.RequestSent, "info")
		default:
			presets.ShowMessage(&r, m.RoleChanged, "success")
		}
		return r, nil
	})
	mb.RegisterEventFunc(evShareRevoke, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		if ver := verifier(ctx); ver != nil && ver.Denied() {
			return r, perm.PermissionDenied
		}
		if err := s.Revoke(ctx.R.Context(), res.Name, recordOf(ctx.R.FormValue(paramRecord)), ctx.R.FormValue(paramSubject),
			currentUser(ctx.R), originOf(ctx.R)); err != nil {
			presets.ShowMessage(&r, err.Error(), "error")
			return r, nil
		}
		presets.ShowMessage(&r, msgs(ctx.Context()).Revoked, "success")
		r.Reload = true
		return r, nil
	})

	// the history of the record's sharing
	dt.Action(ActionHistory).Icon("mdi-history").
		SetI18nLabel(func(ctx context.Context) string { return msgs(ctx).HistoryTitle }).
		ComponentFunc(func(id string, ctx *web.EventContext) (h.HTMLComponent, error) {
			return s.historyBody(ctx, res.Name, recordOf(id))
		})

	// a new owner (of a resource with its own owner)
	if res.OwnerColumn != "" {
		dt.Action(ActionTransferOwner).Icon("mdi-account-switch-outline").
			SetI18nLabel(func(ctx context.Context) string { return msgs(ctx).TransferOwner }).
			SetVerifier(verifier).
			ComponentFunc(func(id string, ctx *web.EventContext) (h.HTMLComponent, error) {
				m := msgs(ctx.Context())
				var users []map[string]string
				for _, c := range s.listUsers(ctx.R.Context()) {
					users = append(users, map[string]string{"title": c.Title, "value": c.Value})
				}
				return h.Div(
					v.VAutocomplete().Label(m.Owner).Items(users).ItemTitle("title").ItemValue("value").
						Attr(web.VField(fieldNewOwner, "")...),
					v.VSelect().Label(m.KeepRole).Items(s.roleChoices(ctx.R.Context(), res.Name)).ItemTitle("title").
						ItemValue("value").Attr(web.VField(fieldKeep, "")...),
				), nil
			}).
			UpdateFunc(func(id string, ctx *web.EventContext) error {
				keep, _ := uuid.Parse(ctx.R.FormValue(fieldKeep))
				if _, err := s.TransferOwner(ctx.R.Context(), res.Name, recordOf(id), ctx.R.FormValue(fieldNewOwner), keep,
					currentUser(ctx.R), originOf(ctx.R)); err != nil {
					return err
				}
				ctx.Flash = msgs(ctx.Context()).RequestSent
				return nil
			})
	}
	return nil
}

// listUsers are the users to choose (the Directory's).
func (s *Service) listUsers(ctx context.Context) []Choice {
	if d, ok := s.o.Users.(Directory); ok {
		return d.ListUsers(ctx)
	}
	return nil
}

// choices are the users (user:<id>) and the groups (group:<id>) to choose.
func (s *Service) choices(ctx context.Context, m *Messages) (out []map[string]string) {
	d, ok := s.o.Users.(Directory)
	if !ok {
		return nil
	}
	for _, c := range d.ListUsers(ctx) {
		out = append(out, map[string]string{"title": c.Title + " (" + m.SubjectUser + ")", "value": UserSubject(c.Value)})
	}
	for _, c := range d.ListGroups(ctx) {
		out = append(out, map[string]string{"title": c.Title + " (" + m.SubjectGroup + ")", "value": GroupSubject(c.Value)})
	}
	return
}

// roleChoices are the roles of the resource to choose.
func (s *Service) roleChoices(ctx context.Context, resource string) (out []map[string]string) {
	roles, _ := RolesOf(s.o.DB.WithContext(ctx), uuid.Nil, resource)
	for _, r := range roles {
		out = append(out, map[string]string{"title": r.Name, "value": r.ID.String()})
	}
	return
}

// accessesBody is who reaches the record: its owner, its accesses (own:
// a role to change, revoke) and the ones of the records above it.
func (s *Service) accessesBody(ctx *web.EventContext, mb *presets.ModelBuilder, resource, id, record string) (h.HTMLComponent, error) {
	m := msgs(ctx.Context())
	c := ctx.R.Context()
	accs, err := s.Subjects(c, resource, record)
	if err != nil {
		return nil, err
	}
	rows := h.Tbody()
	if owner, err := s.Owner(c, resource, record); err == nil {
		rows.AppendChildren(h.Tr(h.Td(h.Text(s.o.Users.Name(c, owner))), h.Td(h.Text(m.Owner)), h.Td()))
	}
	roles := s.roleChoices(c, resource)
	for _, a := range accs {
		name := s.subjectName(c, a.Subject)
		if strings.HasPrefix(a.Subject, SubjectGroup) {
			name += " (" + m.SubjectGroup + ")"
		}
		roleName := ""
		if a.Role != nil {
			roleName = a.Role.Name
		}
		if a.Resource != resource || a.RecordID != record {
			rows.AppendChildren(h.Tr(h.Td(h.Text(name)), h.Td(h.Text(roleName+" ("+m.Inherited+")")), h.Td()).
				Attr("data-shared-access", a.Subject))
			continue
		}
		change := web.Plaid().URL(mb.Info().ListingHrefCtx(ctx)).EventFunc(evShareRole).
			Query(paramRecord, id).Query(paramSubject, a.Subject)
		revoke := web.Plaid().URL(mb.Info().ListingHrefCtx(ctx)).EventFunc(evShareRevoke).
			Query(paramRecord, id).Query(paramSubject, a.Subject)
		rows.AppendChildren(h.Tr(
			h.Td(h.Text(name)),
			h.Td(v.VSelect().Items(roles).ItemTitle("title").ItemValue("value").ModelValue(a.RoleID.String()).
				Density(v.DensityCompact).HideDetails(true).
				Attr("@update:model-value", fmt.Sprintf("(r) => %s", change.FieldValue(paramRole, web.Var("r")).Go()))),
			h.Td(v.VBtn("").Icon("mdi-account-remove-outline").Variant(v.VariantText).Size(v.SizeSmall).
				Attr("title", m.Revoke).Attr("@click", revoke.Go())),
		).Attr("data-shared-access", a.Subject))
	}
	return v.VTable(
		h.Thead(h.Tr(h.Th(m.Subject), h.Th(m.Role), h.Th(""))),
		rows,
	).Density(v.DensityCompact), nil
}

// historyBody is the history of the record's sharing, the newest first.
func (s *Service) historyBody(ctx *web.EventContext, resource, record string) (h.HTMLComponent, error) {
	m := msgs(ctx.Context())
	c := ctx.R.Context()
	evs, err := s.History(c, resource, record)
	if err != nil {
		return nil, err
	}
	if len(evs) == 0 {
		return h.Div(h.Text(m.NoHistory)), nil
	}
	sort.SliceStable(evs, func(i, j int) bool { return evs[i].At.After(evs[j].At) })
	roleNames := map[string]string{}
	var all []ShareTemplate
	s.o.DB.WithContext(c).Find(&all)
	for _, r := range all {
		roleNames[r.ID.String()] = r.Name
	}
	rows := h.Tbody()
	for _, e := range evs {
		role := ""
		if e.RoleAfter != nil {
			role = roleNames[e.RoleAfter.String()]
		}
		if e.RoleBefore != nil {
			role = roleNames[e.RoleBefore.String()] + " → " + role
		}
		subject := ""
		if e.Subject != "" {
			subject = s.subjectName(c, e.Subject)
		}
		by := ""
		if e.By != "" {
			by = s.o.Users.Name(c, e.By)
		}
		rows.AppendChildren(h.Tr(h.Td(h.Text(e.At.Local().Format("2006-01-02 15:04"))), h.Td(h.Text(e.Event)),
			h.Td(h.Text(subject)), h.Td(h.Text(role)), h.Td(h.Text(by))).Attr("data-shared-event", e.Event))
	}
	return v.VTable(
		h.Thead(h.Tr(h.Th(m.When), h.Th(m.Event), h.Th(m.Subject), h.Th(m.Role), h.Th(m.Who))),
		rows,
	).Density(v.DensityCompact), nil
}

// InstallConfirm puts the page where the owner confirms or declines a
// request (its token in ?t=) at path of the admin: the user of the page
// must be the owner.
func (s *Service) InstallConfirm(pb *presets.Builder, path string) {
	registerMessages(pb.I18n())
	const evConfirm, evDecline = "shared_confirm", "shared_decline"
	answer := func(confirm bool) web.EventFunc {
		return func(ctx *web.EventContext) (r web.EventResponse, err error) {
			m := msgs(ctx.Context())
			tok := ctx.R.FormValue("t")
			if confirm {
				err = s.Confirm(ctx.R.Context(), tok, currentUser(ctx.R), originOf(ctx.R))
			} else {
				err = s.Reject(ctx.R.Context(), tok, currentUser(ctx.R), originOf(ctx.R))
			}
			if err != nil {
				presets.ShowMessage(&r, err.Error(), "error")
				return r, nil
			}
			msg := m.Declined
			if confirm {
				msg = m.Confirmed
			}
			presets.ShowMessage(&r, msg, "success")
			r.Reload = true
			return r, nil
		}
	}
	// opened by the link of the mail alone: not in the menu
	page := pb.PagesRegistrator().New(presets.HttpPage(path).InMenu(false).
		TitleFunc(func(ctx context.Context) string { return msgs(ctx).ConfirmTitle })).
		Layout(func(ctx *web.EventContext) (r web.PageResponse, err error) {
			m := msgs(ctx.Context())
			r.PageTitle = m.ConfirmTitle
			tok := ctx.R.FormValue("t")
			req, err := s.RequestOf(ctx.R.Context(), tok)
			if err != nil {
				r.Body = h.Div(h.Text(ErrNotPending.Error()))
				return r, nil
			}
			c := ctx.R.Context()
			label := req.RecordID
			if rs, e := s.resource(req.Resource); e == nil && rs.Label != nil {
				label = rs.Label(c, s.o.DB.WithContext(c), req.RecordID)
			}
			role := ""
			if req.RoleID != nil {
				var t ShareTemplate
				if s.o.DB.WithContext(c).Limit(1).Find(&t, "id = ?", *req.RoleID).RowsAffected > 0 {
					role = t.Name
				}
			}
			info := v.VTable(h.Tbody(
				h.Tr(h.Th(m.Record), h.Td(h.Text(label))),
				h.Tr(h.Th(m.Subject), h.Td(h.Text(s.subjectName(c, req.Subject)))),
				h.Tr(h.Th(m.Role), h.Td(h.Text(role))),
				h.Tr(h.Th(m.RequestedBy), h.Td(h.Text(s.o.Users.Name(c, req.By)))),
				h.Tr(h.Th(m.Event), h.Td(h.Text(req.Kind+" — "+req.State))),
			)).Density(v.DensityCompact)
			var buttons h.HTMLComponent
			if req.State == StatePending {
				buttons = h.Div(
					v.VBtn(m.Confirm).Color("primary").Class("me-2").Attr("data-shared-confirm", true).
						Attr("@click", web.Plaid().EventFunc(evConfirm).Query("t", tok).Go()),
					v.VBtn(m.Decline).Variant(v.VariantTonal).Attr("data-shared-decline", true).
						Attr("@click", web.Plaid().EventFunc(evDecline).Query("t", tok).Go()),
				).Class("mt-4")
			}
			r.Body = v.VContainer(v.VCard(v.VCardTitle(h.Text(m.ConfirmTitle)), v.VCardText(info, buttons)))
			return r, nil
		})
	page.EventFunc(evConfirm, answer(true)).EventFunc(evDecline, answer(false))
	page.Build()
}

// NotificationCount is the number of the user's invites pending (the old
// ones' too).
func (s *Service) NotificationCount(ctx *web.EventContext) int {
	uid := currentUser(ctx.R)
	if uid == "" {
		return 0
	}
	invs, _ := s.Invites(ctx.R.Context(), uid)
	return len(invs) + PendingCount(s.o.DB, uid)
}

// InstallInvites puts the user's invites — accept, decline — in the menu
// of the top (the old ones' too).
func (s *Service) InstallInvites(pb *presets.Builder) {
	registerMessages(pb.I18n())
	const evAccept, evDecline, portal = "shared_invite_accept", "shared_invite_decline", "sharedAccessInvites"
	body := func(ctx *web.EventContext) h.HTMLComponent {
		uid := currentUser(ctx.R)
		if uid == "" {
			return h.Div()
		}
		invs, _ := s.Invites(ctx.R.Context(), uid)
		if len(invs) == 0 {
			return invitesBody(s.o.DB, ctx)
		}
		m := msgs(ctx.Context())
		c := ctx.R.Context()
		items := []h.HTMLComponent{v.VListSubheader(h.Text(m.MyInvites))}
		for _, inv := range invs {
			label := inv.RecordID
			if rs, e := s.resource(inv.Resource); e == nil && rs.Label != nil {
				label = rs.Label(c, s.o.DB.WithContext(c), inv.RecordID)
			}
			id := inv.ID.String()
			items = append(items, v.VListItem(
				v.VListItemTitle(h.Text(label)),
				v.VListItemSubtitle(h.Text(m.InviteFrom+": "+s.o.Users.Name(c, inv.By))),
				web.Slot(
					v.VBtn(m.AcceptInvite).Size(v.SizeSmall).Color("primary").Variant(v.VariantTonal).Class("mr-1").
						Attr("@click", web.Plaid().EventFunc(evAccept).Query("id", id).Go()),
					v.VBtn(m.RejectInvite).Size(v.SizeSmall).Variant(v.VariantText).
						Attr("@click", web.Plaid().EventFunc(evDecline).Query("id", id).Go()),
				).Name("append"),
			).Attr("data-shared-invite", id))
		}
		return h.Div(v.VList(items...).Density(v.DensityCompact), invitesBody(s.o.DB, ctx))
	}
	answer := func(accept bool) web.EventFunc {
		return func(ctx *web.EventContext) (r web.EventResponse, err error) {
			if id, e := uuid.Parse(ctx.R.FormValue("id")); e == nil {
				if accept {
					err = s.Accept(ctx.R.Context(), id, currentUser(ctx.R), originOf(ctx.R))
				} else {
					err = s.Decline(ctx.R.Context(), id, currentUser(ctx.R), originOf(ctx.R))
				}
				if err != nil {
					presets.ShowMessage(&r, err.Error(), "error")
				}
			}
			r.UpdatePortals = append(r.UpdatePortals, &web.PortalUpdate{Name: portal, Body: body(ctx)})
			return r, nil
		}
	}
	wb := pb.GetWebBuilder()
	wb.RegisterEventFunc(evAccept, answer(true))
	wb.RegisterEventFunc(evDecline, answer(false))
	pb.AddMenuTopItemFunc("shared_access_invites", func(ctx *web.EventContext) h.HTMLComponent {
		if currentUser(ctx.R) == "" {
			return nil
		}
		return web.Portal(body(ctx)).Name(portal)
	})
}

// ListUsers / ListGroups: the users and the groups of the admin to choose.
func (a adminUsers) ListUsers(ctx context.Context) (out []Choice) {
	var rows []struct{ ID, Name string }
	model := a.b.NewUser()
	stmt := &gorm.Statement{DB: a.db}
	if err := stmt.Parse(model); err != nil {
		return nil
	}
	a.db.WithContext(ctx).Table(stmt.Schema.Table).Select("id, name").Order("name").Limit(500).Scan(&rows)
	for _, r := range rows {
		out = append(out, Choice{Value: r.ID, Title: r.Name})
	}
	return
}

func (a adminUsers) ListGroups(ctx context.Context) (out []Choice) {
	var gs []user.Group
	a.db.WithContext(ctx).Order("name").Find(&gs)
	for _, g := range gs {
		out = append(out, Choice{Value: g.ID.String(), Title: g.Name})
	}
	return
}
