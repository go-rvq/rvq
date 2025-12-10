package db_tools

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/worker"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/js"
	"github.com/go-rvq/rvq/web/vue"
	"github.com/go-rvq/rvq/x/i18n"
	db_tools "github.com/go-rvq/rvq/x/packages/db-tools"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"gorm.io/gorm"
)

const (
	actionCreateBackup    = "createBackup"
	actionConfigureBackup = "configureBackup"

	eventRemoveBackup        = "removeBackup"
	eventRemoveBackupConfirm = "removeBackupConfirm"

	JobCreateBackup = "DBTools:CreateBackup"
	JobAutoBackup   = "DBTools:AutoBackup"
)

type Builder struct {
	db                       *gorm.DB
	verifier                 *perm.PermVerifierBuilder
	backupController         db_tools.BackupController
	p                        *presets.Builder
	createAction             *presets.ActionBuilder
	configuePesistenceAction *presets.ActionBuilder
	wb                       *worker.Builder
	backupCronSpec           string
}

type MessageForm struct {
	Message string `admin:"required"`
}

func New(db *gorm.DB, i18nB *i18n.Builder, wb *worker.Builder) *Builder {
	if err := db.AutoMigrate(&DbBackupConfig{}); err != nil {
		panic(err)
	}

	ConfigureMessages(i18nB)

	return &Builder{
		db: db,
		wb: wb,
	}
}

func (b *Builder) DailyBackup() *Builder {
	b.backupCronSpec = "@daily"
	return b
}

func (b *Builder) SetBackupCronSpec(cronSpec string) *Builder {
	b.backupCronSpec = cronSpec
	return b
}

func (b *Builder) BackupCronSpec() string {
	return b.backupCronSpec
}

func (b *Builder) GetConfig() (config *DbBackupConfig, err error) {
	config = &DbBackupConfig{}
	if err = b.db.Session(&gorm.Session{}).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			config.Persistence.Data = &db_tools.Persistence{}
		} else {
			return
		}
	} else if config.Persistence.Data == nil {
		config.Persistence.Data = &db_tools.Persistence{}
	}
	return
}

func (b *Builder) GetPersistence() (per *db_tools.Persistence, err error) {
	var cfg *DbBackupConfig
	if cfg, err = b.GetConfig(); err != nil {
		return
	}
	per = cfg.Persistence.Data
	return
}

func (b *Builder) setupCreateAction(p *presets.Builder, page *presets.PageBuilder) {
	b.createAction = page.PrivateAction(actionCreateBackup, func(v *perm.Verifier) *perm.Verifier {
		return v.Do(actionCreateBackup)
	}).
		SetEnabled(func(id string, ctx *web.EventContext) (ok bool, err error) {
			return b.backupController != nil, nil
		}).
		SetI18nLabel(func(ctx context.Context) string {
			return GetMessages(ctx).CreateBackup
		})

	ed := presets.NewModelBuilder(p, &MessageForm{}, presets.ModelConfig().SetModuleKey(MessagesKey))
	presets.ActionForm[*MessageForm](b.createAction, ed.Editing(), func(ctx *presets.ActionFormContext[*MessageForm]) (err error) {
		// 	msg := strings.TrimSpace(ctx.Form.Message)
		_, err = b.wb.CreateJob(ctx.Context, JobCreateBackup, false, &MessageForm{
			Message: strings.TrimSpace(ctx.Form.Message),
		})
		if err == nil {
			ctx.Context.Data().(*presets.PageDoActionOptions).AutoReloadDisabled = true
			ctx.Context.Flash = &presets.FlashMessage{
				HtmlText: string(GetMessages(ctx.Context.Context()).BackupStartedInBackgroundJob(b.wb, ctx.Context.Context())),
			}
		}
		return
	}).Build()
}

func (b *Builder) setupConfigureBackupAction(p *presets.Builder, page *presets.PageBuilder) {
	type Form = db_tools.Persistence

	action := page.PrivateAction(actionConfigureBackup, func(v *perm.Verifier) *perm.Verifier {
		return v.Do(actionConfigureBackup)
	}).
		SetEnabled(func(id string, ctx *web.EventContext) (ok bool, err error) {
			return b.backupController != nil, nil
		}).
		SetI18nLabel(func(ctx context.Context) string {
			return GetMessages(ctx).ConfigureBackupPersistence
		})

	mb := presets.NewModelBuilder(p, &Form{}, presets.ModelConfig().SetModuleKey(MessagesKey))
	ed := mb.Editing()
	presets.ConfigureSelectField(mb, "Other", presets.EDIT, &presets.SelectConfig{
		KeyLabelsFunc: func(ctx *presets.FieldContext, key []string) (v []string) {
			m := GetMessages(ctx.EventContext.Context())
			v = make([]string, len(key))
			for i, k := range key {
				switch k {
				case "":
					v[i] = ""
				case "days":
					v[i] = m.Persistence.OtherDays
				case "weeks":
					v[i] = m.Persistence.OtherWeeks
				case "months":
					v[i] = m.Persistence.OtherMonths
				case "years":
					v[i] = m.Persistence.OtherYears
				}
			}
			return
		},
		AvailableKeysFunc: func(ctx *presets.FieldContext) []string {
			return []string{"", "days", "weeks", "months", "years"}
		},
		SetSelectedKeyFunc: func(ctx *presets.FieldContext, key string) (err error) {
			o := ctx.Obj.(*Form)
			err = o.Other.Parse(key)
			return
		},
	})

	presets.ActionForm[*Form](action, ed, func(ctx *presets.ActionFormContext[*Form]) (err error) {
		var cfg *DbBackupConfig
		if cfg, err = b.GetConfig(); err != nil {
			return
		}
		cfg.Persistence.Data = ctx.Form
		if err = b.db.Session(&gorm.Session{}).Save(&cfg).Error; err != nil {
			return
		}
		return
	}).InitForm(func(fctx *presets.ActionFormContext[*Form]) (err error) {
		var cfg *DbBackupConfig
		if cfg, err = b.GetConfig(); err != nil {
			return
		}
		fctx.Form = cfg.Persistence.Data
		return nil
	}).
		Build()

	b.configuePesistenceAction = action
}

func (b *Builder) Install(p *presets.Builder) (err error) {
	b.p = p
	page := p.PagesRegistrator().New(
		presets.HttpPage("/db-tools").
			MenuIcon("mdi-database-cog").
			TitleFunc(func(ctx context.Context) string {
				return GetMessages(ctx).Database
			})).
		Private().
		Layout(b.pageFunc)

	defer page.Build()

	page.
		PrivateEventFunc(eventRemoveBackupConfirm, func(ctx *web.EventContext) (r web.EventResponse, err error) {
			var (
				id  db_tools.BackupID
				bkp db_tools.Backuper
			)
			if err = id.Parse(ctx.Param(presets.ParamID)); err != nil {
				return
			}
			if bkp, err = b.backupController.Get(id); err != nil {
				return
			}
			r.Body = GetMessages(ctx.Context()).BackupRemoveConfirm(bkp.DetailString(ctx.Context()))
			return
		}).
		PrivateEventFunc(eventRemoveBackup, func(ctx *web.EventContext) (r web.EventResponse, err error) {
			var (
				id  db_tools.BackupID
				bkp db_tools.Backuper
			)
			if id.Parse(ctx.Param(presets.ParamID)); err != nil {
				return
			}
			if bkp, err = b.backupController.Get(id); err != nil {
				return
			}
			if err = b.backupController.Remove(id); err != nil {
				return
			}
			ctx.Flash = GetMessages(ctx.Context()).BackupRemoved(bkp.DetailString(ctx.Context()))
			return
		})

	b.verifier = page.Page().GetVerifier()

	b.setupCreateAction(p, page)
	b.setupConfigureBackupAction(p, page)

	b.wb.NewJob(JobCreateBackup).
		Title(func(ctx *web.EventContext) string {
			return GetMessages(ctx.Context()).CreateDatabaseBackup
		}).
		Resource(&MessageForm{}, func(mb *presets.ModelBuilder) {
			mb.SetModuleKey(MessagesKey)
		}).
		Handler(func(ctx context.Context, job worker.QorJobInterface) (err error) {
			var info *worker.JobInfo
			if info, err = job.GetJobInfo(); err != nil {
				return
			}
			job.AddLog("starting")
			defer job.AddLog("done")
			var bkp db_tools.Backuper
			if bkp, err = b.backupController.Create(false, info.Argument.(*MessageForm).Message); err != nil {
				return
			}
			job.AddLog(fmt.Sprintf(GetMessages(ctx).BackupDetailTemplate, bkp.DetailString(ctx)))
			return
		})

	if b.backupController != nil && b.backupCronSpec != "" {
		b.wb.NewJob(JobAutoBackup).
			Title(func(ctx *web.EventContext) string {
				return GetMessages(ctx.Context()).DatabaseAutoBackup
			}).
			System(true).
			CronConfig(worker.JobCronConfig{Spec: b.backupCronSpec, Once: true}).
			Handler(func(ctx context.Context, job worker.QorJobInterface) (err error) {
				job.AddLog("starting")
				defer job.AddLog("done")
				var bkp db_tools.Backuper
				if bkp, err = b.backupController.Create(true, "auto backup"); err != nil {
					return
				}
				job.AddLog(fmt.Sprintf(GetMessages(ctx).BackupDetailTemplate, bkp.DetailString(ctx)))

				var per *db_tools.Persistence
				if per, err = b.GetPersistence(); err != nil {
					return
				}

				var bkps []db_tools.Backuper
				if bkps, err = b.backupController.RemoveOlder(true, per); err != nil {
					return
				}
				if len(bkps) > 0 {
					var lines = make([]string, len(bkps))
					for i, b := range bkps {
						lines[i] = fmt.Sprintf("\t%d. %s", i+1, b.DetailString(ctx))
					}
					job.AddLog("Older removed:\n" + strings.Join(lines, "\n"))
				} else {
					job.AddLog("No older to remove")
				}

				return
			})
	}

	return
}

func (b *Builder) SetBackupController(v db_tools.BackupController) *Builder {
	b.backupController = v
	return b
}

func (b *Builder) pageFunc(ctx *web.EventContext) (r web.PageResponse, err error) {
	if download := ctx.R.FormValue("download"); download != "" {
		if b.backupController == nil {
			return
		}

		var id db_tools.BackupID
		if err = id.Parse(download); err != nil {
			err = fmt.Errorf("bad download id: %s", err)
			return
		}
		err = b.backupController.Download(ctx.W, ctx.R, id)
		return
	}

	type entry struct {
		ID        string
		createdAt time.Time
		CreatedAt string
		Message   string
		DbName    string
		Size      string
		Auto      string
	}

	var (
		m       = GetMessages(ctx.Context())
		pm      = presets.MustGetMessages(ctx.Context())
		im      = i18n.GetMessages(ctx.Context())
		backups []*entry
		headers = DataTableHeaderBasicSlice{
			{Key: "CreatedAt", Title: m.CreatedAt},
			{Key: "DbName", Title: m.DbName},
			{Key: "Message", Title: m.Message},
			{Key: "Size", Title: m.Size},
			{Key: "Auto", Title: m.Auto},
			{Key: "Actions", Title: m.Actions, Width: "100px"},
		}
		currentName string
	)

	if b.backupController != nil {
		if currentName, err = b.backupController.CurrentName(ctx.Context()); err != nil {
			return
		}

		df := i18n.GetMessages(ctx.Context()).DateTimeFormatter().Full()(ctx.Context())
		cb := func(bkp db_tools.Backuper) error {
			backups = append(backups, &entry{
				createdAt: bkp.GetCreatedAt(),
				ID:        bkp.GetID().String(),
				CreatedAt: df(bkp.GetCreatedAt()),
				Message:   bkp.GetMessage(),
				DbName:    bkp.GetDbName(),
				Size:      humanize.Bytes(uint64(bkp.GetSize())),
				Auto:      im.YesOrNo(bkp.IsAuto()),
			})
			return nil
		}

		b.backupController.List(true, cb, db_tools.NewListFilter())
		b.backupController.List(false, cb, db_tools.NewListFilter())

		sort.Slice(backups, func(i, j int) bool {
			return backups[i].createdAt.After(backups[j].createdAt)
		})

	}

	const itemToDeleteVar = "itemToDelete"

	var (
		tempPortal = ctx.UID()
		doDelete   = presets.OpenConfirmDialog().
				Portal(tempPortal).
				Handler(js.Raw(web.DELETE().
					EventFunc(eventRemoveBackup).
					URL(ctx.R.URL.Path).
					Query("id", web.Var(itemToDeleteVar+".item.ID")).
					ThenScript(itemToDeleteVar + `.cb()`).Go())).
				PromptHandler(js.Raw(web.DELETE().
					EventFunc(eventRemoveBackupConfirm).
					URL(ctx.R.URL.Path).
					Query("id", web.Var(itemToDeleteVar+".item.ID")).String())).
				Build().
				Scope(web.Var("{ " + itemToDeleteVar + " }")).
				Go()

		config      DbBackupConfig
		persistence h.RawHTML
	)

	if err = b.db.Session(&gorm.Session{}).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			config.Persistence.Data = &db_tools.Persistence{}
		} else {
			return
		}
	} else if config.Persistence.Data == nil {
		config.Persistence.Data = &db_tools.Persistence{}
	}

	if persistence, err = m.Persistence.Format(config.Persistence.Data); err != nil {
		return
	}

	r.Body = VContainer(
		h.Iff(len(currentName) > 0, func() h.HTMLComponent {
			return presets.FieldComponentContainer(m.DbName, h.Text(currentName))
		}),
		VCard(
			VCardText(
				h.Iff(b.backupController != nil,
					func() h.HTMLComponent {
						uid := ctx.UID()
						uid2 := ctx.UID()

						return h.HTMLComponents{
							web.Portal().Name(uid),
							VList(
								VListItem(
									web.Slot(
										VBtn("").
											Icon("mdi-cog").
											Density(DensityCompact).
											Attr("@click", web.GET().URL(ctx.R.URL.Path).
												EventFunc(actions.Action).
												Query(presets.ParamAction, actionConfigureBackup).
												Query(presets.ParamTargetPortal, uid).Go()),
									).Name("append"),
									VListItemTitle(h.Text(m.Persistence.Title)),
									VListItemSubtitle(persistence),
								).Class("pa-0"),
							).Class("pa-0 mb-3"),
							h.Div(
								web.Portal().Name(uid2),
								VBtn(m.CreateBackup).
									Attr("@click", web.GET().URL(ctx.R.URL.Path).
										EventFunc(actions.Action).
										Query(presets.ParamAction, actionCreateBackup).
										Query(presets.ParamTargetPortal, uid2).Go()),
							).Class("mb-3"),
						}
					},
				),
				vue.UserComponent(
					web.Portal().Name(tempPortal).Scope(itemToDeleteVar, js.Raw(itemToDeleteVar)),
					VDataTableVirtual(
						h.Iff(b.backupController != nil, func() h.HTMLComponent {
							return web.Slot(
								VBtn("").
									Tag("a").
									Icon("mdi-download").
									Color(ColorSuccess).
									Variant(VariantText).
									Density(DensityCompact).
									Attr(":href", fmt.Sprintf(`"?download="+ item.ID`)).
									Class("me-2"),
								VMenu(
									web.Slot(
										VBtn("").
											Icon("mdi-dots-vertical").
											Variant(VariantText).
											Density(DensityCompact).
											Attr("v-bind", "props").
											Attr("v-on", "on"),
									).Name("activator").
										Scope("{ on, props }"),
									VList(
										VListItem().
											Title(pm.Delete).
											Class("text-error").
											PrependIcon("mdi-delete").
											Attr("@click", "deleteItem(item)"),
									).Density(DensityCompact),
								),
							).Name("item.Actions").Scope("{ item }")
						}),
					).
						Attr(":items", "items").
						Headers(headers).
						NoDataText(pm.ListingNoRecordToShow).
						Density(DensityCompact),
				).Scope("items", backups).
					Scope("deleteItem").
					Scope(itemToDeleteVar).
					Setup(`({scope, window, Vue}) => {
	scope.items = Vue.ref(scope.items)
	scope.`+itemToDeleteVar+` = {}
	scope.deleteItem = (item) => {
		scope.`+itemToDeleteVar+`.item = item
		scope.`+itemToDeleteVar+`.cb = () => {
			const newItems = scope.items.value.filter(e => e.ID != item.ID)
			scope.items.value = newItems
		}
		`+doDelete+`
	}
}`),
			),
		).Title(m.Backups),
	)

	return
}
