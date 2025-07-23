import type { EventFuncID, EventResponse, Location, PortalUpdate, Queries, QueryValue } from './types'
import { buildPushState, objectToFormData } from '@/utils'
import querystring from 'query-string'

declare let window: any

export interface PreFetchData {
  builder: Builder
  options: any
  url: string
}

export class Builder {
  _eventFuncID: EventFuncID = { id: '__reload__' }
  _url?: string
  _method?: string
  _vars?: any
  _locals?: any
  _closer?: any
  _scope?: any = {}
  _loadPortalBody: boolean = false
  _form?: any = {}
  _popstate?: boolean
  _pushState?: boolean
  _location?: Location
  _updateRootTemplate?: any
  _buildPushStateResult?: any
  _parents?: any
  _preFetch: ((data: PreFetchData) => void)[] = []
  _noCache?: boolean = false

  readonly ignoreErrors = [
    'Failed to fetch', // Chrome
    'NetworkError when attempting to fetch resource.', // Firefox
    'The Internet connection appears to be offline.', // Safari
    'Network request failed' // `cross-fetch`
  ]

  isIgnoreError = (err: any) => {
    if (err instanceof Error) {
      return this.ignoreErrors?.includes(err.message)
    }
    return false
  }

  runScript = (r: EventResponse) => {
    if (r.runScript) {
      let script = r.runScript
      if (this._scope) {
        script = 'with(this._scope) {\n' + script + '\n}'
      }
      new Function('vars', 'locals', 'form', 'closer', 'scope', 'plaid', script).apply(this, [
        this._vars,
        this._locals,
        this._form,
        this._closer,
        this._scope,
        (): Builder => {
          return plaid()
            .vars(this._vars)
            .locals(this._locals)
            .form(this._form)
            .closer(this._closer)
            .updateRootTemplate(this._updateRootTemplate)
        }
      ])
    }
    return r
  }

  public parseUrl(url: string): Builder {
    const data = querystring.parseUrl(url, {
      arrayFormat: 'comma',
      parseFragmentIdentifier: true
    })

    this._url = data.url
    this._location = {}

    if (data.query.__execute_event__) {
      this.eventFunc(data.query.__execute_event__ as string)
      delete data.query.__execute_event__
    }

    Object.keys(data.query).forEach((key) => this.query(key, data.query[key] as QueryValue))

    return this
  }

  public eventFunc(id: string): Builder {
    this._eventFuncID.id = id
    return this
  }

  public updateRootTemplate(v: any): Builder {
    this._updateRootTemplate = v
    return this
  }

  public eventFuncID(v: EventFuncID): Builder {
    this._eventFuncID = v
    return this
  }

  public reload(): Builder {
    this._eventFuncID.id = '__reload__'
    return this
  }

  // if you call url(), it will post event func to this url, or else it will post to current window url
  public url(v: string): Builder {
    this._url = v
    return this
  }

  public vars(v: any): Builder {
    this._vars = v
    return this
  }

  public closer(v: any): Builder {
    this._closer = v
    return this
  }

  public loadPortalBody(v: boolean): Builder {
    this._loadPortalBody = v
    return this
  }

  public locals(v: any): Builder {
    this._locals = v
    return this
  }

  public query(key: string, val: QueryValue): Builder {
    if (!this._location) {
      this._location = {}
    }
    if (!this._location.query) {
      this._location.query = {}
    }
    this._location.query[key] = val
    return this
  }

  public mergeQuery(v: boolean): Builder {
    if (!this._location) {
      this._location = {}
    }
    this._location.mergeQuery = v
    return this
  }

  public clearMergeQuery(clearKeys: string[]): Builder {
    if (!this._location) {
      this._location = {}
    }
    this._location.mergeQuery = true
    this._location.clearMergeQueryKeys = clearKeys
    return this
  }

  public location(v: Location): Builder {
    this._location = v
    return this
  }

  public stringQuery(v: string): Builder {
    if (!this._location) {
      this._location = {}
    }

    this._location.stringQuery = v
    return this
  }

  public pushState(v: boolean): Builder {
    this._pushState = v
    return this
  }

  public queries(v: Queries): Builder {
    if (!this._location) {
      this._location = {}
    }
    this._location.query = v
    return this
  }

  public pushStateURL(v: string): Builder {
    if (!this._location) {
      this._location = {}
    }
    this._location.url = v
    this.pushState(true)
    return this
  }

  public form(v: any): Builder {
    this._form = v
    return this
  }

  public scope(v: any): Builder {
    if (v) {
      this._scope = { ...(this._scope ?? {}), ...v }
    }
    return this
  }

  public fieldValue(name: string, v: any): Builder {
    if (!this._form) {
      throw new Error('form not exist')
    }
    this._form[name] = v
    return this
  }

  public popstate(v: boolean): Builder {
    this._popstate = v
    return this
  }

  public preFetch(f: (d: PreFetchData) => void): Builder {
    this._preFetch.push(f)
    return this
  }

  public run(script: string): Builder {
    const f = new Function(script)
    f.apply(this)
    return this
  }

  public method(m: string): Builder {
    this._method = m
    return this
  }

  public parent(index: string, value: any): Builder {
    this._parents = this._parents || {}
    this._parents[index] = value
    return this
  }

  public noCache(): Builder {
    this._noCache = true
    return this
  }

  public buildFetchURL(): string {
    this.ensurePushStateResult()
    return this._buildPushStateResult.eventURL
  }

  public buildPushStateArgs(): [data: any, title: string, url?: string | null] {
    this.ensurePushStateResult()
    return this._buildPushStateResult.pushStateArgs
  }

  public onpopstate(event: any): Promise<void | EventResponse> {
    if (!event.state) {
      // hashtag changes will trigger popstate, when this happens, event.state is null.
      return Promise.reject('event state is undefined')
    }
    return this.popstate(true).location(event.state).reload().go()
  }

  public runPushState() {
    if (this._popstate !== true && this._pushState === true) {
      if (window.history.length <= 2) {
        window.history.pushState({ url: window.location.href }, '', window.location.href)
      }
      const args = this.buildPushStateArgs()
      if (args) {
        window.history.pushState(...args)
      }
    }
  }

  public async fetch(): Promise<EventResponse> {
    if (this._eventFuncID.id == '__reload__') {
      this._buildPushStateResult = null
    }

    const fetchOpts: RequestInit = {
      method: 'POST',
      redirect: 'follow'
    }

    if (this._method) {
      fetchOpts.method = this._method
    }

    if (fetchOpts.method === 'POST') {
      const formData = new FormData()
      objectToFormData(this._form, formData)
      fetchOpts.body = formData
    }

    window.dispatchEvent(new Event('fetchStart'))
    let fetchURL = this.buildFetchURL()

    if (this._preFetch.length) {
      const data = {
        builder: this,
        options: fetchOpts,
        url: fetchURL
      }
      this._preFetch.forEach((f) => {
        f(data)
      })
      fetchURL = data.url
    }

    let er: EventResponse = {}

    try {
      const r = await fetch(fetchURL, fetchOpts).catch((error) => {
        console.error(error)
        Promise.reject(error)
        if (!this.isIgnoreError(error)) {
          alert('Fetch Unknown Error: ' + error)
        }
        // document.location.reload();
      })

      if (!r) {
        return er
      }

      if (!r.ok) {
        this._vars.presetsMessage = {
          show: true,
          message: `Server Response Error: <code>${r.status}</code> (${r.statusText}). See window console for detail.`,
          color: 'error'
        }

        Promise.reject(new Error(`Response status: ${r.status} ${r.statusText}`))
        return er
      }

      if (r.redirected) {
        document.location.replace(r.url)
        return er
      }

      er = await r.json()
    } finally {
      window.dispatchEvent(new Event('fetchEnd'))
    }

    return er
  }

  public go(): Promise<EventResponse> {
    return this.goPre((b: Builder) => null)
  }

  public goPre(pre: (b: Builder) => void): Promise<EventResponse> {
    pre(this)

    this.runPushState()

    return this.fetch()
      .then(this.runScript)
      .then((r: EventResponse) => {
        if (r.pageTitle) {
          document.title = r.pageTitle
        }

        if (r.redirectURL) {
          document.location.replace(r.redirectURL)
        }

        if (r.reloadPortals && r.reloadPortals.length > 0) {
          for (const portalName of r.reloadPortals) {
            const portal = window.__goplaid.portals[portalName]
            if (portal) {
              portal.reload()
            }
          }
        }

        if (r.updatePortals && r.updatePortals.length > 0) {
          for (const pu of r.updatePortals) {
            if (window.__goplaid.portals[pu.name]) {
              const { updatePortalTemplate } = window.__goplaid.portals[pu.name]
              if (updatePortalTemplate) {
                updatePortalTemplate(pu.body)
              }
            } else if (pu.defer) {
              setTimeout(
                (pu: PortalUpdate) => {
                  if (window.__goplaid.portals[pu.name]) {
                    const { updatePortalTemplate } = window.__goplaid.portals[pu.name]
                    if (updatePortalTemplate) {
                      updatePortalTemplate(pu.body)
                    }
                  } else {
                    console.error("go-plaid-portal defered '" + pu.name + "' does not exists.")
                  }
                },
                10,
                pu
              )
            } else {
              console.error("go-plaid-portal '" + pu.name + "' does not exists.")
            }
          }
        }

        if (r.pushState) {
          return plaid()
            .updateRootTemplate(this._updateRootTemplate)
            .reload()
            .pushState(true)
            .location(r.pushState)
            .go()
        }

        if (this._loadPortalBody && r.body) {
          return r
        }

        if (r.body) {
          this._updateRootTemplate(r.body)
          return r
        }

        return r
      })
  }

  public json(): Promise<EventResponse> {
    return this.fetch().then((r: EventResponse) => {
      return r.data
    })
  }

  private ensurePushStateResult() {
    if (!this._noCache && this._buildPushStateResult) {
      return
    }

    const defaultURL = window.location.href

    let url = this._url

    if (url && this._parents) {
      const keys = Object.keys(this._parents)
      keys.forEach((k: any) => {
        url = (url as string).replace('{parent_' + k + '_id}', this._parents[k].toString())
      })
    }

    this._buildPushStateResult = buildPushState(
      {
        ...this._eventFuncID,
        ...{ location: this._location }
      },
      url || defaultURL
    )
  }
}

export function plaid(): Builder {
  return new Builder()
}
