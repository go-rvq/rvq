package presets

import (
	"context"
	"sort"
	"strings"

	"github.com/go-rvq/rvq/x/perm"
	"github.com/iancoleman/strcase"
	"github.com/mpvl/unique"
)

func FieldPerm(name string) string {
	return "#" + name
}

type ResourcePermActions []*ResourcePermAction

func (a *ResourcePermActions) Add(action ...*ResourcePermAction) {
	for _, action := range action {
		if action != nil {
			*a = append(*a, action)
		}
	}
}

type CustomPerm struct {
	Name     string                           `yaml:",omitempty" json:",omitempty"`
	Title    func(ctx context.Context) string `yaml:"-" json:"-"`
	Children []*CustomPerm                    `yaml:",omitempty" json:",omitempty"`
}

type ModelPermObject struct {
	Custom  []*CustomPerm       `yaml:",omitempty" json:",omitempty"`
	Actions ResourcePermActions `yaml:",omitempty" json:",omitempty"`
}
type ModelPerm struct {
	Parent      *ModelPerm                       `yaml:"-" json:"-"`
	Model       *ModelBuilder                    `yaml:"-" json:"-"`
	Name        string                           `yaml:",omitempty" json:",omitempty"`
	Title       func(ctx context.Context) string `yaml:"-" json:"-"`
	Custom      []*CustomPerm                    `yaml:",omitempty" json:",omitempty"`
	Object      *ModelPermObject                 `yaml:",omitempty" json:",omitempty"`
	Actions     ResourcePermActions              `yaml:",omitempty" json:",omitempty"`
	ListActions ResourcePermActions              `yaml:"list_actions,omitempty" json:",omitempty"`
	Children    []*ModelPerm                     `yaml:",omitempty" json:",omitempty"`
}

type ResourcePermActionAction struct {
	Name  string                           `yaml:",omitempty" json:",omitempty"`
	Title func(ctx context.Context) string `yaml:"-" json:"-"`
}

type ResourcePermAction struct {
	Title    func(ctx context.Context) string `yaml:"-" json:"-"`
	Resource string                           `yaml:",omitempty" json:",omitempty"`
	Name     string                           `yaml:",omitempty" json:",omitempty"`
	Fields   []string                         `yaml:",omitempty" json:",omitempty"`
}

type PermMenu struct {
	Parent    *PermMenu                        `yaml:"-" json:"-"`
	Name      string                           `yaml:",omitempty" json:",omitempty"`
	Title     func(ctx context.Context) string `yaml:"-" json:"-"`
	Resources []*ModelPerm                     `yaml:",omitempty" json:",omitempty"`
	Children  []*PermMenu                      `yaml:",omitempty" json:",omitempty"`
}

func (m *PermMenu) AddChildren(children ...*PermMenu) {
	m.Children = append(m.Children, children...)
	for _, child := range children {
		child.Parent = m
	}
}

func (m *PermMenu) Tree() (n *PermNode) {
	return m.tree()
}

func (m *PermMenu) tree() (n *PermNode) {
	n = &PermNode{
		Name:  m.Name,
		Title: m.Title,
	}

	for _, res := range m.Resources {
		n.AddChildren(res.Tree())
	}
	for _, c := range m.Children {
		n.AddChildren(c.tree())
	}
	sort.Slice(n.Children, func(i, j int) bool {
		return n.Children[i].Name < n.Children[j].Name
	})
	return
}

func (m *ModelPerm) AddChildren(children ...*ModelPerm) {
	for _, child := range children {
		child.Parent = m
	}
	m.Children = append(m.Children, children...)
}

func (m *ModelPerm) Tree() (node *PermNode) {
	node = &PermNode{
		Name:  m.Name,
		Title: m.Title,
	}

	for _, s := range m.Custom {
		node.AddChildren(&PermNode{
			Name:  s.Name,
			Title: s.Title,
		})
	}

	for _, action := range m.Actions {
		node.Actions = append(node.Actions, &PermNodeAction{
			Name:   strings.TrimPrefix(action.Name, node.Name),
			Title:  action.Title,
			Fields: action.Fields,
		})
	}

	for _, action := range m.ListActions {
		node.Actions = append(node.Actions, &PermNodeAction{
			Name:   strings.TrimPrefix(action.Name, node.Name),
			Title:  action.Title,
			Fields: action.Fields,
		})
	}

	if m.Object != nil {
		objName := m.Name + "*:"
		obj := &PermNode{
			Name: objName,
			Title: func(ctx context.Context) string {
				return m.Model.TTitle(ctx)
			},
		}

		for _, s := range m.Object.Custom {
			obj.AddChildren(&PermNode{
				Name:  s.Name,
				Title: s.Title,
			})
		}

		for _, action := range m.Object.Actions {
			obj.Actions = append(obj.Actions, &PermNodeAction{
				Name:   strings.TrimPrefix(action.Name, objName),
				Title:  action.Title,
				Fields: action.Fields,
			})
		}
		for _, child := range m.Children {
			obj.AddChildren(child.Tree())
		}
		node.AddChildren(obj)
	} else {
		for _, child := range m.Children {
			node.AddChildren(child.Tree())
		}
	}

	return
}

type PermNodeAction struct {
	Name   string                           `yaml:",omitempty" json:",omitempty"`
	Fields []string                         `yaml:",omitempty" json:",omitempty"`
	Title  func(ctx context.Context) string `yaml:"-" json:"-"`
}

type PermNode struct {
	Parent   *PermNode                        `yaml:"-" json:"-"`
	Name     string                           `yaml:",omitempty" json:",omitempty"`
	Title    func(ctx context.Context) string `yaml:"-" json:"-"`
	Actions  []*PermNodeAction                `yaml:",omitempty" json:",omitempty"`
	Children []*PermNode                      `yaml:",omitempty" json:",omitempty"`
}

func (n *PermNode) AddChildren(children ...*PermNode) {
	n.Children = append(n.Children, children...)
	for _, child := range children {
		child.Parent = n
	}
}

func (n *PermNode) Walk(f func(parents []*PermNode, node *PermNode)) {
	n.walk(nil, f)
}

func (n *PermNode) walk(parents []*PermNode, f func(parents []*PermNode, node *PermNode)) {
	for _, child := range n.Children {
		f(parents, child)
		child.walk(append(parents, n), f)
	}
}

type PermZipEntry struct {
	Resource string            `yaml:",omitempty" json:",omitempty"`
	Actions  []*PermNodeAction `yaml:",omitempty" json:",omitempty"`
}

func (n *PermNode) Zip() (enties []*PermZipEntry) {
	n.Walk(func(parents []*PermNode, node *PermNode) {
		e := &PermZipEntry{
			Resource: node.Name,
		}
		e.Actions = node.Actions
		enties = append(enties, e)
	})
	return
}

func (b *Builder) BuildPermissions() (rootMenu *PermMenu) {
	var (
		roots     []*ModelPerm
		m         = map[*ModelBuilder]*ModelPerm{}
		walkField = func(mb *ModelBuilder, fb *FieldsBuilder, f func(name string)) {
			fb.WalkOptions(mb.Info(), nil, nil, nil, &FieldWalkHandleOptions{
				SkipPerssionCheck: true,
				SkipMode:          true,
				InitializeObjects: true,
				InitializeSlices:  true,

				Handler: func(field *FieldContext) (s FieldWalkState) {
					f(field.Name)
					return
				},
			})
		}
	)

	WalkModels(b.Models(), func(mb *ModelBuilder) (state ModelWalkState, err error) {
		if !mb.IsInMenu() {
			return
		}

		p := &ModelPerm{
			Model: mb,
			Name:  mb.Permissioner().Default().Resource(),
			Title: func(ctx context.Context) string {
				return mb.TTitleAuto(ctx)
			},
		}

		create := &ResourcePermAction{
			Resource: p.Name,
			Name:     PermCreate,
		}

		walkField(mb, &mb.editing.CreatingBuilder().FieldsBuilder, func(name string) {
			create.Fields = append(create.Fields, name)
		})

		unique.Strings(&create.Fields)

		read := &ResourcePermAction{
			Resource: p.Name,
			Name:     PermGet,
		}

		walkField(mb, &mb.detailing.FieldsBuilder, func(name string) {
			read.Fields = append(read.Fields, name)
		})

		unique.Strings(&read.Fields)

		var readActions []*ResourcePermAction

		for _, action := range mb.detailing.actions {
			readActions = append(readActions, &ResourcePermAction{
				Resource: p.Name,
				Name:     action.PermName(),
				Title: func(ctx context.Context) string {
					return action.RequestTitle(mb, ctx)
				},
			})
		}

		var update *ResourcePermAction

		if !mb.editingDisabled {
			update = &ResourcePermAction{
				Resource: p.Name,
				Name:     PermUpdate,
			}

			walkField(mb, &mb.editing.FieldsBuilder, func(name string) {
				update.Fields = append(update.Fields, name)
			})

			unique.Strings(&update.Fields)
		}

		for _, verifier := range mb.verifiers {
			p.Custom = append(p.Custom, &CustomPerm{
				Name:  verifier.Build(mb.permissioner.Default()).Resource(),
				Title: verifier.GetTitle(),
			})
		}

		mp := make(map[*perm.PermVerifierBuilderNode]*CustomPerm)

		for n := range perm.WalkPermVerififierBuilders(mb.AllVerifiers()) {
			e := &CustomPerm{
				Name:  n.Elem.Build(mb.permissioner.Default()).Resource(),
				Title: n.Elem.GetTitle(),
			}
			mp[n] = e
			if n.Parent == nil {
				p.Custom = append(p.Custom, e)
			} else {
				mp[n.Parent].Children = append(mp[n.Parent].Children, e)
			}
		}

		if mb.singleton {
			for _, verifier := range mb.detailing.verifiers {
				p.Custom = append(p.Custom, &CustomPerm{
					Name:  verifier.Build(mb.permissioner.Default()).Resource(),
					Title: verifier.GetTitle(),
				})
			}
			p.Actions.Add(read, update)
			p.Actions.Add(readActions...)
		} else {
			list := &ResourcePermAction{
				Resource: p.Name,
				Name:     PermList,
			}
			walkField(mb, &mb.listing.FieldsBuilder, func(name string) {
				list.Fields = append(list.Fields, name)
			})

			unique.Strings(&list.Fields)
			p.Actions.Add(list, create)

			for _, action := range mb.listing.bulkActions {
				p.Actions = append(p.Actions, &ResourcePermAction{
					Resource: p.Name,
					Name:     "bulk:" + strcase.ToSnake(action.name),
				})
			}

			for _, action := range mb.listing.itemActions {
				p.Actions = append(p.Actions, &ResourcePermAction{
					Resource: p.Name,
					Name:     action.PermName(),
				})
			}

			p.Object = &ModelPermObject{}
			p.Object.Actions.Add(read, update)

			if !mb.deletingDisabled {
				p.Object.Actions.Add(&ResourcePermAction{
					Resource: p.Name + "*:",
					Name:     PermDelete,
				})
			}

			for _, verifier := range mb.detailing.verifiers {
				p.Object.Custom = append(p.Object.Custom, &CustomPerm{
					Name:  verifier.Build(mb.permissioner.Default().On("*")).Resource(),
					Title: verifier.GetTitle(),
				})
			}

			read.Resource += "*:"
			if update != nil {
				update.Resource += "*:"
			}

			for _, action := range readActions {
				action.Resource += "*:"
			}

			p.Object.Actions.Add(readActions...)
		}

		if parent := mb.Parent(); parent == nil {
			roots = append(roots, p)
		} else {
			m[parent].Children = append(m[parent].Children, p)
		}

		m[mb] = p
		return
	})

	menus := map[string]*PermMenu{}

	v := b.verifier.Spawn()
	for _, r := range roots {
		if m := menus[r.Model.menuGroupName]; m != nil {
			m.Resources = append(m.Resources, r)
		} else {
			menus[r.Model.menuGroupName] = &PermMenu{
				Name:      v.Spawn().SnakeOn(r.Model.menuGroupName).Resource(),
				Resources: []*ModelPerm{r},
			}
		}
	}

	for _, group := range b.menuGroups.menuGroups {
		if g := menus[group.name]; g != nil {
			g.Title = g.Title
		}
	}

	if rootMenu = menus[""]; rootMenu == nil {
		rootMenu = &PermMenu{}
	} else {
		delete(menus, "")
	}

	for _, page := range b.pagesRegistrator.httpPages {
		if page.verififer != nil {
			b.verifiers.Add(page.GetVerifier())
		}
	}

	for _, verifier := range b.verifiers {
		rootMenu.Children = append(rootMenu.Children, &PermMenu{
			Name:  verifier.Build(b.verifier.Spawn()).Resource(),
			Title: verifier.GetTitle(),
		})
	}

	for _, menu := range menus {
		rootMenu.Children = append(rootMenu.Children, menu)
	}

	sort.Slice(rootMenu.Children, func(i, j int) bool {
		return rootMenu.Children[i].Name < rootMenu.Children[j].Name
	})

	var sortRes func(res []*ModelPerm)
	sortRes = func(res []*ModelPerm) {
		sort.Slice(res, func(i, j int) bool {
			return res[i].Name < res[j].Name
		})
		for _, re := range res {
			sortRes(re.Children)
			sort.Slice(re.Custom, func(i, j int) bool {
				return re.Custom[i].Name < re.Custom[j].Name
			})
			sort.Slice(re.Actions, func(i, j int) bool {
				return re.Actions[i].Name < re.Actions[j].Name
			})
			if re.Object != nil {
				sort.Slice(re.Object.Custom, func(i, j int) bool {
					return re.Object.Custom[i].Name < re.Object.Custom[j].Name
				})
				sort.Slice(re.Object.Actions, func(i, j int) bool {
					return re.Object.Actions[i].Name < re.Object.Actions[j].Name
				})
			}
			sort.Slice(re.ListActions, func(i, j int) bool {
				return re.ListActions[i].Name < re.ListActions[j].Name
			})
		}
	}

	for _, child := range rootMenu.Children {
		sort.Slice(child.Resources, func(i, j int) bool {
			return child.Resources[i].Name < child.Resources[j].Name
		})
		sortRes(child.Resources)
	}

	return
}
