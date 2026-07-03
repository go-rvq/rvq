import { $ } from 'bun'
import { mkdir, readdir, readFile, rm, writeFile } from 'node:fs/promises'
import { existsSync } from 'node:fs'
import path from 'node:path'

type ApiProp = {
  name: string
  type?: string | string[]
}

type ApiSlot = {
  name: string
}

type ApiComponent = {
  displayName?: string
  pathName?: string
  props?: ApiProp[] | Record<string, unknown>
  slots?: ApiSlot[] | Record<string, unknown>
}

type Api = Record<string, ApiComponent>

const jsRoot = path.resolve(import.meta.dirname, '../..')
const repoRoot = path.resolve(jsRoot, '..')

const args = parseArgs(process.argv.slice(2))
const outDir = path.resolve(jsRoot, args.out ?? 'api/vuetify')
const packageName = args.package ?? 'vuetify'
const version = normalizeVersion(args.version ?? await vuetifyVersion())
const apiPath = args.api ? path.resolve(jsRoot, args.api) : await officialApiPath(version)

const api = await loadApi(apiPath)
await rm(outDir, { recursive: true, force: true })
await mkdir(outDir, { recursive: true })

await writeFile(path.join(outDir, 'vuetify.go'), baseFile(packageName))

for (const [tagName, component] of Object.entries(api).sort(([a], [b]) => a.localeCompare(b))) {
  if (!tagName.startsWith('v-')) continue
  await writeFile(path.join(outDir, `${tagName.slice(2)}.go`), componentFile(packageName, tagName, component))
}

async function vuetifyVersion(): Promise<string> {
  const pkg = JSON.parse(await readFile(path.join(jsRoot, 'vuetifyx/package.json'), 'utf8'))
  const raw = pkg.dependencies?.vuetify ?? pkg.devDependencies?.vuetify
  if (!raw) throw new Error('vuetify dependency not found in js/vuetifyx/package.json')
  return raw
}

function normalizeVersion(raw: string): string {
  const match = raw.match(/\d+\.\d+\.\d+(?:[-+][0-9A-Za-z.-]+)?/)
  if (!match) throw new Error(`cannot parse vuetify version from ${JSON.stringify(raw)}`)
  return match[0]
}

async function officialApiPath(version: string): Promise<string> {
	const workDir = path.join(repoRoot, '.tmp', `vuetify-api-generator-${version}`)
	const apiGeneratorDir = path.join(workDir, 'packages/api-generator')
	const webTypesPath = path.join(apiGeneratorDir, 'node_modules/vuetify/dist/json/web-types.json')

	if (!existsSync(webTypesPath)) {
		await rm(workDir, { recursive: true, force: true })
		await mkdir(workDir, { recursive: true })

    const tag = `v${version}`
    await $`curl -L ${`https://codeload.github.com/vuetifyjs/vuetify/tar.gz/refs/tags/${tag}`} | tar -xz --strip-components=1 -C ${workDir}`

    const pkgPath = path.join(apiGeneratorDir, 'package.json')
    const pkg = JSON.parse(await readFile(pkgPath, 'utf8'))
    pkg.dependencies.vuetify = version
    pkg.dependencies.mkdirp ??= '^3.0.1'
    pkg.dependencies.rimraf ??= '^6.0.1'
    pkg.dependencies.upath ??= '^2.0.1'
    pkg.dependencies.yargs ??= '^17.7.2'
    await writeFile(pkgPath, `${JSON.stringify(pkg, null, 2)}\n`)

		await $`bun install`.cwd(apiGeneratorDir)
	}

	return webTypesPath
}

function parseArgs(values: string[]): Record<string, string> {
  const out: Record<string, string> = {}
  for (let i = 0; i < values.length; i++) {
    const arg = values[i]
    if (!arg.startsWith('--')) continue
    const key = arg.slice(2)
    const next = values[i + 1]
    if (next && !next.startsWith('--')) {
      out[key] = next
      i++
    } else {
      out[key] = 'true'
    }
  }
  return out
}

async function loadApi(input: string): Promise<Api> {
	if (input.endsWith('web-types.json')) {
		const webTypes = JSON.parse(await readFile(input, 'utf8'))
		const api: Api = {}
		for (const tag of webTypes.contributions?.html?.tags ?? []) {
			const tagName = toTagName(tag.name)
			api[tagName] = {
				displayName: tag.name,
				pathName: tagName,
				props: (tag.attributes ?? []).map((attr: any) => ({
					name: attr.name,
					type: attr.value?.type,
				})),
				slots: (tag.slots ?? []).map((slot: any) => ({
					name: slot.name,
				})),
			}
		}
		return api
	}

	const entries = await readdir(input, { withFileTypes: true })
	const api: Api = {}
	for (const entry of entries) {
		if (!entry.isFile() || !entry.name.startsWith('V') || !entry.name.endsWith('.json')) continue
		const component = JSON.parse(await readFile(path.join(input, entry.name), 'utf8')) as ApiComponent
    const tagName = component.pathName ?? kebab(component.displayName ?? entry.name.slice(0, -'.json'.length))
    api[tagName] = component
  }
  return api
}

function baseFile(pkg: string): string {
  return `package ${pkg}

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/tag"
)

type (
	VTagBuilderGetter[T any] interface {
		tag.TagBuilderGetter[T]
		GetVTagBuilder() *VTagBuilder[T]
	}

	VTagBuilder[T any] struct {
		tag.TagBuilder[T]
	}
)

func VTag[T VTagBuilderGetter[T]](dot T, name string, children ...h.HTMLComponent) T {
	vtb := dot.GetVTagBuilder()
	vtb.TagBuilder = *tag.NewTag(dot, name, children...).GetTagBuilder()
	return dot
}

func (b *VTagBuilder[T]) GetVTagBuilder() *VTagBuilder[T] {
	return b
}

func (t *VTagBuilder[T]) RawWidth(v interface{}) T {
	return t.Attr(":width", v)
}

func (t *VTagBuilder[T]) RawHeight(v interface{}) T {
	return t.Attr(":height", v)
}

func (t *VTagBuilder[T]) RawClass(v interface{}) T {
	return t.Attr(":class", v)
}

func (t *VTagBuilder[T]) FormField(formKey string, v interface{}) T {
	return t.Attr(web.VField(formKey, v)...)
}
`
}

function componentFile(pkg: string, tagName: string, component: ApiComponent): string {
  const typeName = `${pascal(tagName)}Builder`
  const props = propsList(component.props)
  const slots = slotsList(component.slots)
  const methods = props.map(prop => propMethod(typeName, prop)).join('\n')
  const slotMethods = slots.map(slot => slotMethod(typeName, slot)).join('\n')
  return `package ${pkg}

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type ${typeName} struct {
	VTagBuilder[*${typeName}]
}

func ${pascal(tagName)}(children ...h.HTMLComponent) *${typeName} {
	return VTag(&${typeName}{}, ${JSON.stringify(tagName)}, children...)
}

${methods}func (b *${typeName}) On(name string, value string) (r *${typeName}) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *${typeName}) Bind(name string, value string) (r *${typeName}) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *${typeName}) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *${typeName}) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *${typeName}) Slot(name string, child ...h.HTMLComponent) (r *${typeName}) {
	b.SetSlot(name, child...)
	return b
}

func (b *${typeName}) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *${typeName}) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

${slotMethods}
`
}

function propsList(props: ApiComponent['props']): ApiProp[] {
  const values = Array.isArray(props)
    ? props
    : Object.entries(props ?? {}).map(([name, value]) => ({ name, ...(typeof value === 'object' && value ? value : {}) }))
  const seen = new Set<string>()
  const out: ApiProp[] = []
  for (const prop of values) {
    if (!prop.name || seen.has(prop.name)) continue
    seen.add(prop.name)
    out.push(prop)
  }
  return out
}

function slotsList(slots: ApiComponent['slots']): ApiSlot[] {
  const values = Array.isArray(slots)
    ? slots
    : Object.keys(slots ?? {}).map(name => ({ name }))
  const seen = new Set<string>()
  const out: ApiSlot[] = []
  for (const slot of values) {
    if (!slot.name || !/^[A-Za-z][A-Za-z0-9_.:-]*$/.test(slot.name) || seen.has(slot.name)) continue
    seen.add(slot.name)
    out.push(slot)
  }
  return out
}

function propMethod(typeName: string, prop: ApiProp): string {
  const method = pascal(prop.name)
  const attr = kebab(prop.name)
  const t = propType(prop)
  if (t === 'bool') {
    return `func (b *${typeName}) ${method}(v bool) (r *${typeName}) {
	b.Attr(":${attr}", fmt.Sprint(v))
	return b
}

`
  }
  if (t === 'string') {
    return `func (b *${typeName}) ${method}(v string) (r *${typeName}) {
	b.Attr("${attr}", v)
	return b
}

`
  }
  return `func (b *${typeName}) ${method}(v interface{}) (r *${typeName}) {
	b.Attr(":${attr}", h.JSONString(v))
	return b
}

`
}

function slotMethod(typeName: string, slot: ApiSlot): string {
  const method = pascal(slot.name)
  const name = JSON.stringify(slot.name)
  return `func (b *${typeName}) SetSlot${method}(child ...h.HTMLComponent) {
	b.SetSlot(${name}, child...)
}

func (b *${typeName}) SetScopedSlot${method}(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot(${name}, scope, child...)
}

func (b *${typeName}) Slot${method}(child ...h.HTMLComponent) (r *${typeName}) {
	b.SetSlot${method}(child...)
	return b
}

func (b *${typeName}) ScopedSlot${method}(scope string, child ...h.HTMLComponent) (r *${typeName}) {
	b.SetScopedSlot${method}(scope, child...)
	return b
}

`
}

function propType(prop: ApiProp): 'bool' | 'string' | 'interface' {
  const types = Array.isArray(prop.type) ? prop.type : prop.type ? [prop.type] : []
  if (types.length === 1 && types[0] === 'boolean') return 'bool'
  if (types.length === 1 && types[0] === 'string') return 'string'
  return 'interface'
}

function pascal(value: string): string {
  return value
    .replace(/^v-/, 'v-')
    .replace(/([a-z0-9])([A-Z])/g, '$1-$2')
    .split(/[-_:.]/g)
    .filter(Boolean)
    .map(part => part.charAt(0).toUpperCase() + part.slice(1))
    .join('')
    .replace(/Id$/g, 'ID')
}

function kebab(value: string): string {
	return value.replace(/([a-z0-9])([A-Z])/g, '$1-$2').toLowerCase()
}

function toTagName(value: string): string {
	if (/^V[A-Z]/.test(value)) return `v-${kebab(value.slice(1))}`
	return kebab(value)
}
