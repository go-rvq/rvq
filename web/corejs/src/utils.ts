import 'formdata-polyfill'
import querystring from 'query-string'
import union from 'lodash/union'
import without from 'lodash/without'
import type { EventFuncID, ValueOp } from './types'

export function buildPushState(eventFuncId: EventFuncID, url: string): any {
  const loc = eventFuncId.location
  const orig = querystring.parseUrl(loc?.url || url, {
    arrayFormat: 'comma',
    parseFragmentIdentifier: true
  })

  const resultQuery: any = {}
  let locQuery
  // If pushState is string, then replace query string to it
  // If pushState it object, merge url query
  if (loc) {
    if (loc.stringQuery) {
      const strQuery = querystring.parse(loc.stringQuery, { arrayFormat: 'comma' })
      // @ts-ignore
      loc.query = { ...strQuery, ...loc.query }
    }

    if (loc.mergeQuery) {
      const clearKeys = loc.clearMergeQueryKeys || []
      for (const [key, value] of Object.entries(orig.query)) {
        // If clearMergeQueryKeys is present then skip current location queries which contained by clearMergeQueryKeys
        // If clearMergeQueryKeys is empty, all queries from current location will be kept
        if (clearKeys.indexOf(key.split('.')[0]) < 0) {
          resultQuery[key] = value
        }
      }
      if (!loc.query) {
        loc.query = {}
      }
    }
    locQuery = loc.query
  }

  const st = locQuery || orig.query
  let addressBarQuery = ''
  for (const [key, v] of Object.entries(st)) {
    if (Array.isArray(v)) {
      resultQuery[key] = v
    } else if (typeof v === 'object') {
      const valueOp = v as ValueOp
      queryUpdateByValueOp(resultQuery, key, valueOp)
    } else {
      resultQuery[key] = v
    }
  }

  const requestQuery = { ...resultQuery, ...{ __execute_event__: eventFuncId.id } }

  for (const [key, value] of Object.entries(requestQuery)) {
    if (value === null || value === undefined || value === '') {
      delete requestQuery[key]
    }
  }

  addressBarQuery = querystring.stringify(resultQuery, { arrayFormat: 'comma' })
  if (addressBarQuery.length > 0) {
    addressBarQuery = `?${addressBarQuery}`
  }

  let newUrl = orig.url + addressBarQuery
  if (orig.fragmentIdentifier) {
    newUrl = newUrl + '#' + orig.fragmentIdentifier
  }
  const pushedState = { query: resultQuery, url: newUrl }

  return {
    pushStateArgs: [pushedState, '', newUrl],
    eventURL: `${orig.url}?${querystring.stringify(requestQuery, { arrayFormat: 'comma' })}`
  }
}

function queryUpdateByValueOp(query: any, key: string, valueOp: ValueOp): void {
  if (!valueOp.value) {
    return
  }

  let opValues: any = valueOp.value
  if (!Array.isArray(valueOp.value)) {
    opValues = [valueOp.value]
  }

  let values = query[key]
  if (values && !Array.isArray(values)) {
    values = [values]
  }

  if (valueOp.add) {
    query[key] = union(values, opValues)
    return
  }

  if (valueOp.remove) {
    const newValues = without(values, ...opValues)
    if (newValues.length === 0) {
      delete query[key]
    } else {
      query[key] = newValues
    }
  }
  return
}

export function setFormValue(form: FormData, fieldName: string, val: any): boolean {
  // console.log("setFormValue", inspectFormData(form), fieldName, val)
  if (!fieldName || fieldName.length === 0) {
    return false
  }

  if (val instanceof Event) {
    return setFormValue(form, fieldName, val.target)
  }

  if (val instanceof HTMLInputElement) {
    // console.log("target.value = ", target.value, ", target.type = ", target.type, ", target.checked = ", target.checked)
    if (val.files) {
      return setFormValue(form, fieldName, val.files)
    }

    switch (val.type) {
      case 'checkbox':
        if (val.checked) {
          return formSet(form, fieldName, val.value)
        } else {
          if (form.has(fieldName)) {
            form.delete(fieldName)
            return true
          }
        }
        return false
      case 'radio':
        if (val.checked) {
          return formSet(form, fieldName, val.value)
        }
        return false
      default:
        return formSet(form, fieldName, val.value)
    }
  }

  if (val instanceof HTMLTextAreaElement) {
    return formSet(form, fieldName, val.value)
  }

  if (val instanceof HTMLSelectElement) {
    return formSet(form, fieldName, val.value)
  }

  if (val === null || val === undefined) {
    return formSet(form, fieldName, '')
  }

  let changed = false
  if (form.has(fieldName)) {
    changed = true
    form.delete(fieldName)
  }

  // console.log('val', val, 'Array.isArray(val)', Array.isArray(val));
  if (Array.isArray(val) || val instanceof FileList) {
    for (let i = 0; i < val.length; i++) {
      if (val[i] instanceof File) {
        changed = true
        form.append(fieldName, val[i], val[i].name)
      } else {
        changed = true
        form.append(fieldName, val[i])
      }
    }
    return changed
  }

  if (val instanceof File) {
    form.set(fieldName, val, val.name)
    return true
  } else {
    return formSet(form, fieldName, val)
  }
}

function formSet(form: FormData, fieldName: string, val: string): boolean {
  if (form.get(fieldName) === val) {
    return false
  }
  form.set(fieldName, val)
  return true
}

export function registerEvent(el: any, event: string, listener: any, options: any) {
  el.addEventListener(event, listener, options)
  return () => el.removeEventListener(event, listener, options)
}

export function objectToFormData(obj: any, form: FormData, parentKey = '') {
  if (obj === undefined || obj === null) {
    return
  }
  const isArr = Array.isArray(obj)
  if (isArr && obj.length > 0) {
    if (obj[0] instanceof File || obj[0] instanceof Blob) {
      setFormValue(form, parentKey, obj)
      return
    }
  }
  Object.keys(obj).forEach((key) => {
    if (key === '$parent') {
      return
    }

    const value = obj[key]
    // Construct the form key
    const formKey = parentKey ? (isArr ? `${parentKey}[${key}]` : `${parentKey}.${key}`) : key

    if (typeof value === 'object' && !(value instanceof File) && !(value instanceof Date)) {
      objectToFormData(value, form, formKey)
    } else {
      setFormValue(form, formKey, value)
    }
  })

  return form
}
