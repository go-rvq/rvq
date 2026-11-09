package shared

import (
	"context"
	"fmt"
	"github.com/go-rvq/rvq/admin/model"
	"net/http"
	"slices"
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
	// fieldParentRead: of a record under another, its parent read too
	fieldParentRead = "shared_parent_read"
	fieldKeep       = "shared_keep"
	fieldNewOwner   = "shared_new_owner"
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
	// a model under another: its records reached by its parents'
	// permissions, a role on one of them its own, its parent read
	// (nestedPolicies); a model on top: its records' own
	policies := s.nestedPolicies(mb, o.ReadScope)
	if mb.Parent() == nil {
		policies = func(tx *gorm.DB, a *Access, role *ShareTemplate, removed bool) error {
			ref := "sharedaccess:" + a.ID.String()
			// the listing of the model, scoped to the records reached
			// (ReadScope): who reaches one lists them — and the columns of
			// the listing (#<field>, of the model)
			listRef, colsRef := ref+":list", ref+":list_columns"
			if removed || role == nil || role.None {
				for _, r := range []string{listRef, colsRef} {
					if err := perms.Revoke(tx, r); err != nil {
						return err
					}
				}
				return perms.Revoke(tx, ref)
			}
			rid, err := mb.ParseRecordID(a.RecordID)
			if err != nil {
				return err
			}
			_ = perms.Revoke(tx, ref)
			_, err = perms.GrantShared(tx, ref, a.Subject, perms.RecordResource(mb, rid)+"*", []string(role.Actions), a.ID)
			if err == nil && o.ReadScope && role.Allows(presets.PermList) {
				list := perms.ListResource(mb)
				_, err = perms.GrantShared(tx, listRef, a.Subject, list+presets.PermList, []string{presets.PermList}, a.ID)
				if err == nil {
					_, err = perms.GrantShared(tx, colsRef, a.Subject, list+"#*", []string{presets.PermList}, a.ID)
				}
			}
			return err
		}
	}
	s.OnAccess(res.Name, func(tx *gorm.DB, a *Access, role *ShareTemplate, removed bool) error {
		if err := policies(tx, a, role, removed); err != nil {
			return err
		}
		s.loadPolicies(mb, tx)
		return nil
	})
	s.mu.Lock()
	s.installed = append(s.installed, installedModel{mb: mb, resource: res.Name, policies: policies})
	s.mu.Unlock()
	// the policies of the accesses of before, as they are made now — of a
	// model under another, by SyncPolicies, once every model is in the admin
	if mb.Parent() == nil && mb.Builder().GetPermission() != nil {
		if err := s.syncPolicies(installedModel{mb: mb, resource: res.Name, policies: policies}); err != nil {
			return err
		}
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
			body := h.Div(
				v.VAutocomplete().Label(m.Subject).Items(s.choices(ctx.R.Context(), m)).ItemTitle("title").ItemValue("value").
					Attr(web.VField(fieldSubject, "")...).Attr("data-shared-subject", true),
				v.VSelect().Label(m.Role).Items(s.roleChoices(ctx.R.Context(), res.Name)).ItemTitle("title").ItemValue("value").
					Attr(web.VField(fieldRole, "")...).Attr("data-shared-role", true),
			)
			if mb.Parent() != nil {
				// a record under another: its parent read too, by default
				body.AppendChildren(v.VSwitch().Label(m.ParentRead).Color(v.ColorPrimary).HideDetails(true).
					Attr(web.VField(fieldParentRead, true)...).Attr("data-shared-parent-read", true))
			}
			return body, nil
		}).
		UpdateFunc(func(id string, ctx *web.EventContext) error {
			roleID, err := uuid.Parse(ctx.R.FormValue(fieldRole))
			if err != nil {
				return ErrRole
			}
			var opts []ShareOption
			if mb.Parent() != nil && ctx.R.FormValue(fieldParentRead) == "false" {
				opts = append(opts, NoParentRead())
			}
			req, err := s.Share(ctx.R.Context(), res.Name, recordOf(id), ctx.R.FormValue(fieldSubject), roleID,
				currentUser(ctx.R), originOf(ctx.R), opts...)
			if err != nil {
				return err
			}
			if req == nil {
				// None: at once
				ctx.Flash = msgs(ctx.Context()).Granted
			} else {
				ctx.Flash = msgs(ctx.Context()).RequestSent
			}
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
		rows.AppendChildren(h.Tr(h.Td(h.Text(e.At.Local().Format("2006-01-02 15:04"))), h.Td(h.Text(m.EventName(e.Event))),
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
				h.Tr(h.Th(m.Event), h.Td(h.Text(m.KindName(req.Kind)+" — "+m.StateName(req.State)))),
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

// NotificationCount is the number of the user's invites pending.
func (s *Service) NotificationCount(ctx *web.EventContext) int {
	uid := currentUser(ctx.R)
	if uid == "" {
		return 0
	}
	invs, _ := s.Invites(ctx.R.Context(), uid)
	return len(invs)
}

// InstallInvites puts the user's invites — accept, decline — in the
// notification center.
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
			return h.Div()
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
		return v.VList(items...).Density(v.DensityCompact)
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
	pb.AddNotificationSource(presets.NotificationSource{
		Name:  "shared_access_invites",
		Count: s.NotificationCount,
		Content: func(ctx *web.EventContext) h.HTMLComponent {
			if s.NotificationCount(ctx) == 0 {
				return nil
			}
			return web.Portal(body(ctx)).Name(portal)
		},
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

// nestedPolicies are the policies of the accesses of a model under another
// (mb):
//   - None takes the record out: a policy denying the subject everything of
//     it and under it;
//   - a role: its actions on the record and under it, the listing of the
//     model under the parent (scoped: readScope); its parent read (not
//     NoParentRead) — the parent's detail and listing, and its supporting
//     models, the ones under it not shared by themselves: listed and read
//     (its other shared records stay out).
func (s *Service) nestedPolicies(mb *presets.ModelBuilder, readScope bool) AccessHook {
	return func(tx *gorm.DB, a *Access, role *ShareTemplate, removed bool) error {
		ref := "sharedaccess:" + a.ID.String()
		if err := revokeRefs(tx, ref); err != nil {
			return err
		}
		if removed || role == nil {
			return nil
		}
		levels, err := s.chainIn(tx, a.Resource, a.RecordID)
		if err != nil {
			return err
		}
		parents := mb.Parents()
		if len(levels)-1 != len(parents) {
			return fmt.Errorf("shared: the parents of %s %s: %d records, %d models", a.Resource, a.RecordID,
				len(levels)-1, len(parents))
		}
		// the ids of the parents, the outermost first (presets' order)
		pids := make([]model.ID, len(parents))
		for i, p := range parents {
			if pids[i], err = p.ParseRecordID(levels[len(levels)-1-i].id); err != nil {
				return err
			}
		}
		rid, err := mb.ParseRecordID(a.RecordID)
		if err != nil {
			return err
		}
		record := perms.RecordResource(mb, rid, pids...)
		if role.None {
			_, err = perms.DenyShared(tx, ref, a.Subject, record+"*", []string{"*"}, a.ID)
			return err
		}
		read := []string{presets.PermList, presets.PermGet}
		grant := func(suffix, resource string, actions []string) error {
			_, err := perms.GrantShared(tx, ref+suffix, a.Subject, resource, actions, a.ID)
			return err
		}
		if err := grant("", record+"*", role.Actions); err != nil {
			return err
		}
		if readScope && role.Allows(presets.PermList) {
			list := perms.ListResource(mb, pids...)
			if err := grant(":list", list+presets.PermList, []string{presets.PermList}); err != nil {
				return err
			}
			if err := grant(":list_columns", list+"#*", []string{presets.PermList}); err != nil {
				return err
			}
		}
		if a.NoParentRead {
			return nil
		}
		parent, ppids := parents[len(parents)-1], pids[:len(pids)-1]
		pid := pids[len(pids)-1]
		// the parent: its detail and listing, its fields declared
		pfields := s.supportOf(parent).fieldsPattern()
		precord := perms.RecordResource(parent, pid, ppids...)
		plist := perms.ListResource(parent, ppids...)
		gs := []struct{ suffix, resource string }{
			{":parent", precord + presets.PermGet},
			{":parent_list", plist + presets.PermList},
		}
		if pfields != "" {
			gs = append(gs, struct{ suffix, resource string }{":parent_fields", precord + pfields},
				struct{ suffix, resource string }{":parent_list_columns", plist + pfields})
		}
		// its supporting models: the ones declared (Support), not shared by
		// themselves — their listing and records, the fields and children
		// declared
		for _, c := range parent.Children() {
			o, ok := s.supportDecl(c)
			if c == mb || !ok || s.isShared(c) {
				continue
			}
			l := perms.ListResource(c, pids...)
			k := ":support:" + c.UriName()
			gs = append(gs, struct{ suffix, resource string }{k, l + presets.PermList},
				struct{ suffix, resource string }{k + ":get", l + presets.PermGet},
				struct{ suffix, resource string }{k + ":record", l + "<*>:" + presets.PermGet})
			if f := o.fieldsPattern(); f != "" {
				gs = append(gs, struct{ suffix, resource string }{k + ":columns", l + f},
					struct{ suffix, resource string }{k + ":fields", l + "<*>:" + f})
			}
			for _, ch := range o.Children {
				gs = append(gs, struct{ suffix, resource string }{k + ":child:" + ch, l + "<*>:" + ch + ":*"})
			}
		}
		for _, g := range gs {
			if err := grant(g.suffix, g.resource, read); err != nil {
				return err
			}
		}
		return nil
	}
}

// SupportOptions say what of a supporting model — under the parent of a
// shared model, not shared by itself — is read by whom reaches the parent
// only as support (a record under it shared with them): none declared,
// the model is not read.
type SupportOptions struct {
	// Fields are the fields read (in its listing and its records); none,
	// none.
	Fields []string
	// Children are the URI names of its models read under its records;
	// none, none.
	Children []string
	// Related narrows its records to the ones related to the records
	// shared with the user (the ids of sharedIDs, of the shared models
	// under the same parent); nil, all of the parent's.
	Related func(db *gorm.DB, sharedIDs []string) *gorm.DB
}

// supportModels are the supporting models declared (Support).
type supportModels map[*presets.ModelBuilder]SupportOptions

// fieldsPattern is the resource of the fields of o: "#{A,B}:*" ("" none).
func (o SupportOptions) fieldsPattern() string {
	switch len(o.Fields) {
	case 0:
		return ""
	case 1:
		return "#" + o.Fields[0] + ":*"
	}
	return "#{" + strings.Join(o.Fields, ",") + "}:*"
}

// Support declares mb a supporting model, read as o says — mb under the
// parent of a shared model (or that parent itself: its fields read). Before
// SyncPolicies. A model under the parent not declared is not read.
func (s *Service) Support(mb *presets.ModelBuilder, o SupportOptions) {
	s.mu.Lock()
	if s.support == nil {
		s.support = supportModels{}
	}
	s.support[mb] = o
	s.mu.Unlock()
	if o.Related != nil {
		s.scopeSupport(mb, o.Related)
	}
}

func (s *Service) supportDecl(mb *presets.ModelBuilder) (SupportOptions, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.support[mb]
	return o, ok
}

func (s *Service) supportOf(mb *presets.ModelBuilder) SupportOptions {
	o, _ := s.supportDecl(mb)
	return o
}

// installedOf is the model shared of mb (nil, not shared).
func (s *Service) installedOf(mb *presets.ModelBuilder) *installedModel {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.installed {
		if s.installed[i].mb == mb {
			m := s.installed[i]
			return &m
		}
	}
	return nil
}

// scopeSupport narrows the reads of the supporting model mb, for whom
// reaches its parent only as support, to the records related to the ones
// shared with them.
func (s *Service) scopeSupport(mb *presets.ModelBuilder, related func(db *gorm.DB, sharedIDs []string) *gorm.DB) {
	mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		return do.(*gorm2op.DataOperatorBuilder).
			WithReadCallbacks(func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
				cb.Pre(func(state *gorm2op.CallbackState) error {
					if state.Ctx == nil {
						return nil
					}
					ids, ok, err := s.supportedIDs(state.Ctx, mb)
					if err != nil || !ok {
						return err
					}
					state.DB = related(state.DB, ids)
					return nil
				})
			})
	})
}

// supportedIDs are the records shared with the user of ctx under the parent
// of the path of mb — ok when the user reaches that parent only as support.
func (s *Service) supportedIDs(ctx *web.EventContext, mb *presets.ModelBuilder) (ids []string, ok bool, err error) {
	uid := currentUser(ctx.R)
	parentMB := mb.Parent()
	if uid == "" || parentMB == nil {
		return nil, false, nil
	}
	pids, err := mb.ParseParentsID(ctx.R)
	if err != nil || len(pids) == 0 {
		return nil, false, nil
	}
	return s.supportedFor(ctx.R.Context(), uid, parentMB, strings.Join(pids[len(pids)-1].StringValues(), "_"))
}

// supportedFor are the records shared with user under the record pid of
// parentMB — ok when user reaches that record only as support.
func (s *Service) supportedFor(c context.Context, uid string, parentMB *presets.ModelBuilder, pid string) (ids []string, ok bool, err error) {
	pm := s.installedOf(parentMB)
	if pm == nil {
		return nil, false, nil
	}
	// a real access to the parent (owned, shared, inherited): all
	reached, err := s.accessible(c, uid, pm.resource, false)
	if err != nil || slices.Contains(reached, pid) {
		return nil, false, err
	}
	s.mu.RLock()
	var shared []installedModel
	for _, m := range s.installed {
		if m.mb.Parent() == parentMB {
			shared = append(shared, m)
		}
	}
	s.mu.RUnlock()
	for _, m := range shared {
		r, err := s.resource(m.resource)
		if err != nil {
			return nil, false, err
		}
		got, err := s.accessible(c, uid, m.resource, false)
		if err != nil {
			return nil, false, err
		}
		if len(got) == 0 {
			continue
		}
		var under []string
		if err := s.o.DB.WithContext(c).Table(r.Table).Where(r.IDColumn+" IN ? AND "+r.ParentColumn+" = ?", got, pid).
			Pluck(r.IDColumn, &under).Error; err != nil {
			return nil, false, err
		}
		ids = append(ids, under...)
	}
	return ids, true, nil
}

// installedModel is a model shared in the admin (Install): its resource,
// the policies of its accesses.
type installedModel struct {
	mb       *presets.ModelBuilder
	resource string
	policies AccessHook
}

// SyncPolicies makes again the policies of the accesses of every model
// shared in the admin (Install), as they are made now: to call once the
// admin has all its models — the policies of a model under another name its
// parent's (supporting) models.
func (s *Service) SyncPolicies() error {
	s.mu.RLock()
	ms := slices.Clone(s.installed)
	s.mu.RUnlock()
	for _, m := range ms {
		if err := s.syncPolicies(m); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) syncPolicies(m installedModel) error {
	if m.mb.Builder().GetPermission() == nil {
		return nil
	}
	var accesses []Access
	if err := s.o.DB.Preload("Role").Where("resource = ?", m.resource).Find(&accesses).Error; err != nil {
		return err
	}
	for i := range accesses {
		a := &accesses[i]
		if err := m.policies(s.o.DB, a, a.Role, false); err != nil {
			return fmt.Errorf("shared: the policies of %s %s: %w", m.resource, a.RecordID, err)
		}
	}
	if len(accesses) > 0 {
		s.loadPolicies(m.mb, s.o.DB)
	}
	return nil
}

// loadPolicies puts the policies just written in the permissions of the
// admin of mb.
func (s *Service) loadPolicies(mb *presets.ModelBuilder, tx *gorm.DB) {
	if pb := mb.Builder().GetPermission(); pb != nil {
		from := s.o.Now().Add(-time.Minute)
		pb.LoadDBPoliciesToMemory(tx, &from)
	}
}

// revokeRefs takes the policies of the reference ref away, and its
// derived ones (ref:…).
func revokeRefs(tx *gorm.DB, ref string) error {
	return tx.Where("refer_id = ? OR refer_id LIKE ?", ref, ref+":%").Delete(&perm.DefaultDBPolicy{}).Error
}

// shared says the model mb is shared by s (Install).
func (s *Service) isShared(mb *presets.ModelBuilder) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, m := range s.installed {
		if m.mb == mb {
			return true
		}
	}
	return false
}
