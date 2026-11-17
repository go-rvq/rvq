// vx-files: the explorer of files (fs-tools; the files of a user of the OS
// in senir, the file system of the admin in HermonCMS), as a
// file manager (Nautilus): the folders by a path, a list or a grid, the
// files sent (a button, or dropped — folders too), downloaded (a folder, or
// several items, as a tar.gz — or each by itself), from the internet, extracted, renamed, copied, moved, deleted, a
// text edited; its folder in the IDE (ide-url). Its API is fs-tools/explorer
// at base (base + 'api/fs/…'); its terminal (F4), <vx-terminal>.
(function () {
  'use strict';

  const ARCHIVE = /\.(zip|tar|tgz|tar\.gz|gz)$/i;
  const IMAGE = /\.(png|jpe?g|gif|webp|svg|bmp|ico)$/i;
  const CODE = /\.(go|js|ts|vue|py|sh|rb|php|java|c|h|cpp|rs|gad|gadx|css|scss|html?|json|ya?ml|toml|xml|sql)$/i;
  const TEXT = /\.(txt|md|log|csv|env|conf|ini|cfg)$/i;

  function iconOf(e) {
    if (e.dir) return e.link ? 'mdi-folder-arrow-right-outline' : 'mdi-folder';
    const n = e.name;
    if (ARCHIVE.test(n)) return 'mdi-folder-zip-outline';
    if (IMAGE.test(n)) return 'mdi-file-image-outline';
    if (/\.pdf$/i.test(n)) return 'mdi-file-pdf-box';
    if (CODE.test(n)) return 'mdi-file-code-outline';
    if (TEXT.test(n) || n.startsWith('.')) return 'mdi-file-document-outline';
    return 'mdi-file-outline';
  }

  // the kind of an entry, for the color of its icon (vx-files.css: a soft tone
  // a kind, of the light theme and of the dark)
  function kindOf(e) {
    if (e.dir) return e.link ? 'link' : 'folder';
    const n = e.name;
    if (ARCHIVE.test(n)) return 'archive';
    if (IMAGE.test(n)) return 'image';
    if (/\.pdf$/i.test(n)) return 'pdf';
    if (CODE.test(n)) return 'code';
    if (TEXT.test(n) || n.startsWith('.')) return 'text';
    return 'file';
  }

  // the language of a file for <vx-codemirror> by its extension ("" plain text)
  const LANGS = ['json', 'js', 'jsx', 'ts', 'tsx', 'html', 'css', 'go', 'gad', 'gadt', 'gadx', 'md', 'markdown', 'yaml', 'yml'];
  const LANG_ALIAS = { htm: 'html', svg: 'html', xml: 'html', vue: 'html', mjs: 'js', cjs: 'js', tmpl: 'gadt', javascript: 'js', typescript: 'ts' };
  function languageOf(name) {
    const n = name.toLowerCase();
    const ext = n.includes('.') ? n.slice(n.lastIndexOf('.') + 1) : '';
    const l = LANG_ALIAS[ext] || ext;
    return LANGS.includes(l) ? l : '';
  }

  function size(n, loc) {
    if (n < 1024) return n + ' B';
    const u = ['kB', 'MB', 'GB', 'TB'];
    let i = -1;
    do { n /= 1000; i++; } while (n >= 1000 && i < u.length - 1);
    const d = n < 10 ? 1 : 0;
    let t;
    try { t = n.toLocaleString(loc || undefined, { minimumFractionDigits: d, maximumFractionDigits: d }); } catch (e) { t = n.toFixed(d); }
    return t + ' ' + u[i];
  }

  // The sending of a large file (sendParts): its parts' size, how many at
// once, the tries of each.
const UPLOAD_PART = 8 << 20, UPLOAD_PARALLEL = 3, UPLOAD_TRIES = 4;

// words is the word of key of the namespace vxFiles of i18n (the page's
// vue-i18n), its placeholders of params; the key when there is none
function words(i18n, key, params) {
  const k = 'vxFiles.' + key;
  try { if (i18n && i18n.te(k)) return i18n.t(k, params || {}); } catch (e) { }
  return key;
}
// locOf is the locale of i18n as the browser names it ("pt-BR")
function locOf(i18n) {
  const l = i18n && i18n.locale ? String(i18n.locale.value || i18n.locale) : '';
  return (l || (navigator.language || 'en')).replace('_', '-');
}
function join(a, b) { return (a.replace(/\/+$/, '') + '/' + b).replace(/\/+/g, '/'); }
  function parent(p) { const s = p.replace(/\/+$/, '').split('/'); s.pop(); return s.join('/') || '/'; }

  // the files of a drop: each with its path (folders walked)
  function droppedFiles(dt) {
    const items = dt.items ? Array.from(dt.items) : [];
    const entries = items.map(function (it) { return it.webkitGetAsEntry ? it.webkitGetAsEntry() : null; }).filter(Boolean);
    if (!entries.length) {
      return Promise.resolve(Array.from(dt.files || []).map(function (f) { return { file: f, path: f.name }; }));
    }
    const out = [];
    function walk(entry, prefix) {
      if (entry.isFile) {
        return new Promise(function (res) {
          entry.file(function (f) { out.push({ file: f, path: prefix + f.name }); res(); }, function () { res(); });
        });
      }
      if (entry.isDirectory) {
        const reader = entry.createReader();
        return new Promise(function (res) {
          const all = [];
          // readEntries gives at most 100 a call: read until none
          (function more() {
            reader.readEntries(function (batch) {
              if (!batch.length) {
                Promise.all(all.map(function (e) { return walk(e, prefix + entry.name + '/'); })).then(res);
                return;
              }
              all.push.apply(all, batch);
              more();
            }, function () { res(); });
          })();
        });
      }
      return Promise.resolve();
    }
    return Promise.all(entries.map(function (e) { return walk(e, ''); })).then(function () { return out; });
  }

  const VxFiles = {
    // the words: the namespace vxFiles of the page's vue-i18n (fs-tools'
    // FilesMessages, in the user's language)
    setup: function () {
      let i18n = null;
      try { i18n = window.VueI18n.useI18n({ useScope: 'global' }); } catch (e) { }
      return { i18n: i18n };
    },
    props: { base: String, ideUrl: String, start: String, title: String, gitBase: String, terminal: Boolean, vendor: String,
      // the terminal: its address (by default base + 'terminal'), the variable of the root of the files in its shell, whether it is recorded
      terminalUrl: String, termRootVar: { type: String, default: 'FILES_ROOT' }, termRecorded: Boolean,
      // fetch (from the internet) and extract (an archive): offered when the API serves them
      fetch: { type: Boolean, default: true }, extract: { type: Boolean, default: true } },
    data: function () {
      let view = 'list', hidden = false, termOpen = false, termHeight = 260, thumbs = false, cols = { dims: false };
      try {
        view = localStorage.getItem('vx-files-view') || 'list';
        hidden = localStorage.getItem('vx-files-hidden') === '1';
        thumbs = localStorage.getItem('vx-files-thumbs') === '1';
        cols = Object.assign(cols, JSON.parse(localStorage.getItem('vx-files-cols') || '{}'));
        termOpen = localStorage.getItem('vx-files-term') === '1';
        termHeight = Number(localStorage.getItem('vx-files-term-h')) || 260;
      } catch (e) { }
      return {
        path: '/', repo: false, entries: [], loading: false,
        // what may be done in the folder (list's can): reading — the
        // downloads, offered whenever it may —, writing — the rest
        can: { read: true, write: true },
        // the images drawn (their thumbnails): on, and those loaded
        thumbs: thumbs, thumbOk: {},
        // the columns that may be shown (list, tree): the images' dimensions
        cols: cols, error: '', view: view, hidden: hidden, filter: '',
        selected: [], anchor: null, back: [], fwd: [], sortBy: 'name', sortAsc: true,
        editingPath: false, pathText: '',
        dialog: '', form: { name: '', to: '', url: '', confirm: [] },
        editor: { open: false, path: '', content: '', orig: '', saving: false, truncated: false, writable: false, editing: false, language: '', size: 0, mode: '', confirm: false },
        menu: { open: false, x: 0, y: 0 },
        image: { open: false, name: '', actual: false },
        dropping: 0, upload: { busy: false, done: 0, total: 0, bytes: 0, sent: 0, failed: [] , pending: null },
        snack: { show: false, text: '', color: 'success' },
        // the tree: the folders open (by their path of here) and what they have
        opened: {}, kids: {},
        // the terminal below (as Dolphin's, F4): open, its height
        term: { open: termOpen && this.terminal, height: termHeight, key: 0, dragging: false },
        // the height of the files' area: to the window's bottom (the
        // terminal's room below it taken)
        areaH: 0,
      };
    },
    computed: {
      tb: function () { return { size: 'small', variant: 'text', density: 'comfortable' }; },
      views: function () {
        return [
          { value: 'list', icon: 'mdi-view-list', title: this.w('viewList') },
          { value: 'tree', icon: 'mdi-file-tree-outline', title: this.w('viewTree') },
          { value: 'grid', icon: 'mdi-view-grid-outline', title: this.w('viewIcons') },
        ];
      },
      // what is shown: the folder's entries, filtered and sorted — in the
      // tree, with those of its folders open under each (their name the path
      // from here, base the name, depth)
      shown: function () {
        const self = this, f = this.filter.toLowerCase();
        const k = this.sortBy, asc = this.sortAsc ? 1 : -1;
        const sorted = function (list) {
          return list.filter(function (e) {
            return (self.hidden || !e.name.startsWith('.')) && (!f || e.dir || e.name.toLowerCase().includes(f));
          }).sort(function (a, b) {
            if (a.dir !== b.dir) return a.dir ? -1 : 1;
            let x = a[k], y = b[k];
            if (k === 'name') { x = x.toLowerCase(); y = y.toLowerCase(); }
            return x < y ? -asc : x > y ? asc : 0;
          });
        };
        if (this.view !== 'tree') {
          return sorted(this.entries.filter(function (e) { return !f || e.name.toLowerCase().includes(f); }));
        }
        const out = [];
        const walk = function (list, prefix, depth) {
          sorted(list).forEach(function (e) {
            const name = prefix ? prefix + '/' + e.name : e.name;
            out.push(Object.assign({}, e, { name: name, base: e.name, depth: depth }));
            if (e.dir && self.opened[name] && self.kids[name]) walk(self.kids[name], name, depth + 1);
          });
        };
        walk(this.entries, '', 0);
        return out;
      },
      crumbs: function () {
        const parts = this.path.split('/').filter(Boolean);
        const out = [{ name: this.title || this.w('home'), path: '/' }];
        let acc = '';
        parts.forEach(function (p) { acc += '/' + p; out.push({ name: p, path: acc }); });
        return out;
      },
      sel: function () {
        const s = this.selected;
        return this.shown.filter(function (e) { return s.includes(e.name); });
      },
      selSize: function () { return this.sel.reduce(function (t, e) { return t + (e.dir ? 0 : e.size); }, 0); },
      one: function () { return this.sel.length === 1 ? this.sel[0] : null; },
      dirty: function () { return this.editor.content !== this.editor.orig; },
    },
    mounted: function () {
      this.go(this.start || '/', true);
      window.addEventListener('popstate', this.onPop);
      window.addEventListener('resize', this.fitArea);
      this.$nextTick(this.fitArea);
    },
    beforeUnmount: function () {
      window.removeEventListener('popstate', this.onPop);
      window.removeEventListener('resize', this.fitArea);
    },
    watch: {
      'term.open': function () { this.$nextTick(this.fitArea); },
      'term.height': function () { this.fitArea(); },
      'image.open': function (v) { const self = this; if (v) setTimeout(function () { self.$refs.imageBox && self.$refs.imageBox.focus(); }, 50); },
    },
    methods: {
      api: function (op, params, init) {
        const self = this, q = new URLSearchParams();
        Object.keys(params || {}).forEach(function (k) {
          const v = params[k];
          (Array.isArray(v) ? v : [v]).forEach(function (x) { if (x !== undefined && x !== null) q.append(k, x); });
        });
        return fetch(this.base + 'api/fs/' + op + '?' + q.toString(), init || {}).then(function (r) {
          if (r.status === 403) throw new Error(self.w('expired'));
          const ct = r.headers.get('Content-Type') || '';
          if (!ct.includes('json')) {
            if (!r.ok) throw new Error(r.status + ' ' + r.statusText);
            return r;
          }
          return r.json().then(function (j) { if (!r.ok) throw new Error(j.error || r.statusText); return j; });
        });
      },
      post: function (op, params, body) { return this.api(op, params, { method: 'POST', body: body }); },
      toast: function (text, color) { this.snack = { show: true, text: text, color: color || 'success' }; },
      failed: function (e) { this.toast(e.message || String(e), 'error'); },
      go: function (p, replace) {
        const self = this;
        p = p || '/';
        this.loading = true;
        this.error = '';
        return this.api('list', { path: p, dims: this.cols.dims ? 1 : undefined }).then(function (j) {
          if (!replace && self.path !== j.path) { self.back.push(self.path); self.fwd = []; }
          const moved = self.path !== j.path;
          if (moved) { self.opened = {}; self.kids = {}; self.thumbOk = {}; }
          self.path = j.path;
          self.repo = !!j.repo;
          self.can = Object.assign({ read: true, write: true }, j.can || {});
          self.entries = j.entries;
          // the terminal follows the folder: only when it changed
          if (moved) self.termCd();
          self.selected = [];
          self.anchor = null;
          // the address keeps the folder (a reload, back and forward)
          const u = new URL(window.location.href);
          if (j.path === '/') u.searchParams.delete('dir'); else u.searchParams.set('dir', j.path);
          window.history[replace ? 'replaceState' : 'pushState']({ vxFilesDir: j.path }, '', u.toString());
        }).catch(function (e) { self.error = e.message; }).finally(function () { self.loading = false; });
      },
      onPop: function (ev) {
        if (ev.state && ev.state.vxFilesDir !== undefined) this.go(ev.state.vxFilesDir, true);
      },
      goBack: function () { if (this.back.length) { this.fwd.push(this.path); this.go(this.back.pop(), true); } },
      goFwd: function () { if (this.fwd.length) { this.back.push(this.path); this.go(this.fwd.pop(), true); } },
      up: function () { if (this.path !== '/') this.go(parent(this.path)); },
      refresh: function () {
        const self = this;
        // the folders open in the tree: read again (those gone, closed)
        Object.keys(this.opened).forEach(function (n) { self.loadKids(n); });
        return this.go(this.path, true);
      },
      // the tree: a folder (its path from here) open or closed; its entries
      // read when opened
      toggleNode: function (e) {
        const n = e.name;
        if (this.opened[n]) {
          const o = Object.assign({}, this.opened);
          Object.keys(o).forEach(function (k) { if (k === n || k.startsWith(n + '/')) delete o[k]; });
          this.opened = o;
          return;
        }
        const self = this;
        this.loadKids(n).then(function (ok) { if (ok) self.opened = Object.assign({}, self.opened, { [n]: true }); });
      },
      loadKids: function (n) {
        const self = this;
        return this.api('list', { path: join(this.path, n), dims: this.cols.dims ? 1 : undefined }).then(function (j) {
          self.kids = Object.assign({}, self.kids, { [n]: j.entries });
          return true;
        }).catch(function (e) {
          const o = Object.assign({}, self.opened);
          delete o[n];
          self.opened = o;
          self.failed(e);
          return false;
        });
      },
      // the terminal below: shown or hidden (F4), its height by its bar
      // dragged; it follows the folder (cd)
      fitArea: function () {
        const el = this.$refs.area;
        if (!el) return;
        const top = el.getBoundingClientRect().top + window.scrollY;
        // the terminal and the bar above it (vx-files-sash: 10px, 4 of margin)
        const below = 28 + (this.term.open ? this.term.height + 14 : 0);
        this.areaH = Math.max(200, window.innerHeight - top - below - 8);
      },
      toggleTerm: function () {
        if (!this.terminal) return;
        this.term.open = !this.term.open;
        try { localStorage.setItem('vx-files-term', this.term.open ? '1' : '0'); } catch (e) { }
      },
      termURL: function () { return (this.terminalUrl || this.base + 'terminal') + '?dir=' + encodeURIComponent(this.path); },
      termCd: function () {
        const t = this.$refs.term;
        if (t && t.cd) t.cd(this.path, true);
      },
      // the bar between the files and the terminal, dragged (a mouse, a
      // finger): the terminal's height, kept in the browser
      termDrag: function (ev) {
        if (ev.button !== undefined && ev.button !== 0) return;
        const self = this, y0 = ev.clientY, h0 = this.term.height, bar = ev.currentTarget;
        ev.preventDefault();
        try { bar.setPointerCapture(ev.pointerId); } catch (e) { }
        this.term.dragging = true;
        document.body.style.userSelect = 'none';
        const move = function (e) { self.term.height = self.termClamp(h0 + y0 - e.clientY); };
        const up = function () {
          bar.removeEventListener('pointermove', move);
          bar.removeEventListener('pointerup', up);
          bar.removeEventListener('pointercancel', up);
          self.term.dragging = false;
          document.body.style.userSelect = '';
          self.termResize(self.term.height);
        };
        bar.addEventListener('pointermove', move);
        bar.addEventListener('pointerup', up);
        bar.addEventListener('pointercancel', up);
      },
      termClamp: function (h) { return Math.round(Math.max(120, Math.min(window.innerHeight - 160, h))); },
      termResize: function (h) {
        this.term.height = this.termClamp(h);
        try { localStorage.setItem('vx-files-term-h', String(this.term.height)); } catch (e) { }
      },
      editPath: function () { this.pathText = this.path; this.editingPath = true; },
      submitPath: function () { this.editingPath = false; this.go(this.pathText); },
      setSort: function (k) { if (this.sortBy === k) this.sortAsc = !this.sortAsc; else { this.sortBy = k; this.sortAsc = true; } },
      // the buttons of the bar: icons all of one size; on (a mode, a toggle),
      // in the primary color
      on: function (active) { return Object.assign({}, this.tb, active ? { color: 'primary', variant: 'tonal' } : {}); },
      setView: function (v) { this.view = v; try { localStorage.setItem('vx-files-view', v); } catch (e) { } },
      // a column shown or hidden (kept); the folder read again for it
      toggleCol: function (k) {
        this.cols = Object.assign({}, this.cols, { [k]: !this.cols[k] });
        try { localStorage.setItem('vx-files-cols', JSON.stringify(this.cols)); } catch (e) { }
        this.refresh();
      },
      dims: function (e) { return e.w ? e.w + ' × ' + e.h : ''; },
      toggleThumbs: function () {
        this.thumbs = !this.thumbs;
        try { localStorage.setItem('vx-files-thumbs', this.thumbs ? '1' : '0'); } catch (e) { }
      },
      toggleHidden: function () { this.hidden = !this.hidden; try { localStorage.setItem('vx-files-hidden', this.hidden ? '1' : '0'); } catch (e) { } },
      isSel: function (e) { return this.selected.includes(e.name); },
      click: function (e, ev) {
        const names = this.shown.map(function (x) { return x.name; });
        if (ev.shiftKey && this.anchor !== null) {
          const a = names.indexOf(this.anchor), b = names.indexOf(e.name);
          this.selected = names.slice(Math.min(a, b), Math.max(a, b) + 1);
        } else if (ev.ctrlKey || ev.metaKey) {
          this.selected = this.isSel(e) ? this.selected.filter(function (n) { return n !== e.name; }) : this.selected.concat([e.name]);
          this.anchor = e.name;
        } else {
          this.selected = [e.name];
          this.anchor = e.name;
        }
      },
      selectAll: function () { this.selected = this.shown.map(function (x) { return x.name; }); },
      open: function (e) {
        if (e.dir) { this.go(join(this.path, e.name)); return; }
        if (this.extract && this.can.write && ARCHIVE.test(e.name)) { this.selected = [e.name]; this.ask('extract'); return; }
        if (IMAGE.test(e.name)) { this.image = { open: true, name: e.name, actual: false }; return; }
        this.edit(e);
      },
      // the viewer of the images of the folder: one shown, the others by the
      // arrows
      imageURL: function (name) { return this.base + 'api/fs/download?inline=1&path=' + encodeURIComponent(join(this.path, name)); },
      images: function () { return this.shown.filter(function (e) { return !e.dir && IMAGE.test(e.name); }); },
      imageStep: function (d) {
        const list = this.images(), i = list.findIndex(function (e) { return e.name === this.image.name; }, this);
        if (!list.length) return;
        const next = list[(i + d + list.length) % list.length];
        this.image = { open: true, name: next.name, actual: false };
        this.selected = [next.name];
      },
      imageKey: function (ev) {
        if (ev.key === 'ArrowRight') { ev.preventDefault(); this.imageStep(1); }
        else if (ev.key === 'ArrowLeft') { ev.preventDefault(); this.imageStep(-1); }
        else if (ev.key === 'Escape') { ev.preventDefault(); this.image.open = false; }
      },
      imageEntry: function () { const n = this.image.name; return this.entries.find(function (e) { return e.name === n; }) || { name: n, size: 0 }; },
      ctx: function (e, ev) {
        if (e && !this.isSel(e)) { this.selected = [e.name]; this.anchor = e.name; }
        if (!e) this.selected = [];
        this.menu = { open: true, x: ev.clientX, y: ev.clientY };
      },
      key: function (ev) {
        if (this.dialog || this.editor.open) return;
        if (ev.target && /INPUT|TEXTAREA/.test(ev.target.tagName)) return;
        if (ev.key === 'Delete' && this.sel.length && this.can.write) { ev.preventDefault(); this.ask('delete'); }
        else if (ev.key === 'F2' && this.one && this.can.write) { ev.preventDefault(); this.ask('rename'); }
        else if (ev.key === 'Enter' && this.one) { ev.preventDefault(); this.open(this.one); }
        else if (ev.key === 'Backspace' || (ev.altKey && ev.key === 'ArrowUp')) { ev.preventDefault(); this.up(); }
        else if ((ev.ctrlKey || ev.metaKey) && ev.key === 'a') { ev.preventDefault(); this.selectAll(); }
        else if ((ev.ctrlKey || ev.metaKey) && ev.key === 'l') { ev.preventDefault(); this.editPath(); }
        else if (ev.key === 'F4' && this.terminal) { ev.preventDefault(); this.toggleTerm(); }
      },
      // the dialogs: name (new folder, new file, rename), to (copy, move),
      // url (from the internet), delete, extract
      ask: function (kind) {
        const one = this.one;
        this.form = { name: kind === 'rename' && one ? (one.base || one.name) : '', to: this.path, url: '', confirm: this.sel.map(function (e) { return e.name; }) };
        if (kind === 'extract' && one) this.form.to = join(this.path, one.name.replace(/(\.tar\.gz|\.tgz|\.tar|\.zip|\.gz)$/i, ''));
        this.dialog = kind;
      },
      run: function () {
        const self = this, k = this.dialog, f = this.form, here = this.path;
        let p;
        const names = this.sel.map(function (e) { return join(here, e.name); });
        if (k === 'mkdir') p = this.post('mkdir', { path: join(here, f.name) });
        else if (k === 'newfile') p = this.post('write', { path: join(here, f.name) }, '');
        else if (k === 'rename') {
          // in the tree: in its own folder
          const n = this.one.name, d = n.includes('/') ? n.slice(0, n.lastIndexOf('/')) : '';
          p = this.post('rename', { path: join(here, n), to: join(join(here, d), f.name) });
        }
        else if (k === 'copy' || k === 'move') {
          const op = k === 'copy' ? 'copy' : 'rename';
          p = names.reduce(function (acc, src) {
            return acc.then(function () { return self.post(op, { path: src, to: join(f.to, src.split('/').pop()) }); });
          }, Promise.resolve());
        } else if (k === 'delete') p = this.post('delete', { path: names });
        else if (k === 'extract') p = this.post('extract', { path: names[0], to: f.to });
        else if (k === 'fetch') p = this.post('fetch', { path: here, url: f.url, name: f.name || undefined });
        if (!p) return;
        p.then(function (j) {
          self.dialog = '';
          self.toast({ mkdir: self.w('doneMkdir'), newfile: self.w('doneNewfile'), rename: self.w('doneRename'), copy: self.w('doneCopy'), move: self.w('doneMove'),
            delete: self.w('doneDelete'), extract: j && j.files !== undefined ? self.w('doneExtractN', { n: j.files }) : self.w('doneExtract'),
            fetch: j && j.path ? self.w('doneFetchPath', { path: j.path }) : self.w('doneFetch') }[k]);
          return self.refresh();
        }).catch(this.failed);
      },
      // fetched: the address of the download of paths (of the file system),
      // saved as name
      fetchDownload: function (paths, name) {
        const a = document.createElement('a');
        a.href = this.base + 'api/fs/download?' + paths.map(function (p) { return 'path=' + encodeURIComponent(p); }).join('&');
        a.download = name;
        document.body.appendChild(a);
        a.click();
        a.remove();
      },
      // download: the item e — a file as it is, a folder as a tar.gz ({dir:
      // true, name: ''}: the folder open)
      download: function (e) {
        const base = (e.base || e.name || (this.path === '/' ? this.w('archiveName') : this.path.split('/').pop())).split('/').pop();
        this.fetchDownload([join(this.path, e.name)], e.dir ? base + '.tar.gz' : base);
      },
      // downloadSel: the items selected — one by itself; several, in one
      // tar.gz (each: false) or each by itself (each: true)
      downloadSel: function (each) {
        const self = this;
        if (this.sel.length === 1 || each) { this.sel.forEach(function (e, i) { setTimeout(function () { self.download(e); }, i * 300); }); return; }
        const name = (this.path === '/' ? this.w('archiveName') : this.path.split('/').pop()) + '.tar.gz';
        this.fetchDownload(this.sel.map(function (e) { return join(self.path, e.name); }), name);
      },
      // the title of the download of the selection (or of the folder open)
      downloadTitle: function () {
        if (!this.sel.length) return this.w('downloadFolder');
        if (this.sel.length > 1) return this.w('downloadN', { n: this.sel.length });
        return this.one.dir ? this.w('downloadTgz') : this.w('download');
      },
      ideHere: function (sub) { return this.ideUrl + '?dir=' + encodeURIComponent(sub || this.path); },
      // the repository at sub (the folder open) in Git: the IDE with its
      // Changes and Git panels alone
      gitHere: function (sub) { return this.ideHere(sub) + '&git=1'; },
      // the URL that clones the repository at the path p of here (git, by
      // HTTPS: the login of SSH)
      cloneURL: function (p) {
        p = p.replace(/^\/+/, '');
        return window.location.origin + (this.gitBase || '/') + (/\.git$/.test(p) ? p : p + '.git');
      },
      // copyClone copies the clone URL of e (none: the folder open)
      copyClone: function (e) {
        const self = this, url = this.cloneURL(e ? join(this.path, e.name) : this.path);
        const done = function () { self.toast(self.w('cloneCopied', { url: url })); };
        if (navigator.clipboard && navigator.clipboard.writeText) {
          navigator.clipboard.writeText(url).then(done).catch(function () { self.toast(url, 'info'); });
        } else {
          self.toast(url, 'info');
        }
      },
      // the viewer of a text: CodeMirror read only; edited (Editar) when the
      // user of the OS may write it and it was read whole
      edit: function (e, editing) {
        const self = this, p = join(this.path, e.name);
        this.api('read', { path: p }).then(function (j) {
          if (j.binary) { self.download(e); return; }
          self.editor = { open: true, path: p, content: j.content, orig: j.content, saving: false, truncated: j.truncated,
            writable: !!j.writable && !j.truncated, editing: false, language: languageOf(e.name), size: j.size, mode: j.mode, confirm: false };
          if (editing && self.editor.writable) self.editor.editing = true;
        }).catch(this.failed);
      },
      startEdit: function () { if (this.editor.writable) this.editor.editing = true; },
      cancelEdit: function () { this.editor.content = this.editor.orig; this.editor.editing = false; this.editor.confirm = false; },
      save: function () {
        const self = this;
        if (!this.editor.editing || !this.dirty) return;
        this.editor.saving = true;
        this.post('write', { path: this.editor.path }, this.editor.content).then(function () {
          self.editor.orig = self.editor.content;
          self.editor.editing = false;
          self.toast(self.w('saved'));
          self.refresh();
        }).catch(this.failed).finally(function () { self.editor.saving = false; });
      },
      editorKey: function (ev) {
        if ((ev.ctrlKey || ev.metaKey) && ev.key === 's') { ev.preventDefault(); this.save(); }
        else if ((ev.ctrlKey || ev.metaKey) && ev.key === 'e' && !this.editor.editing) { ev.preventDefault(); this.startEdit(); }
        else if (ev.key === 'Escape' && !this.editor.editing) { ev.preventDefault(); this.closeEditor(); }
      },
      // closing with changes not saved: asked first
      closeEditor: function (force) {
        if (this.dirty && force !== true) { this.editor.confirm = true; return; }
        this.editor.open = false;
        this.editor.confirm = false;
      },
      // sending files: picked (files, a folder) or dropped
      pick: function (folder) { (folder ? this.$refs.folder : this.$refs.files).click(); },
      picked: function (ev) {
        const list = Array.from(ev.target.files || []).map(function (f) { return { file: f, path: f.webkitRelativePath || f.name }; });
        ev.target.value = '';
        this.send(list, false);
      },
      dragOver: function (ev) {
        if (!this.can.write || !ev.dataTransfer || !Array.from(ev.dataTransfer.types || []).includes('Files')) return;
        ev.preventDefault();
        ev.dataTransfer.dropEffect = 'copy';
      },
      dragEnter: function (ev) { if (this.can.write && ev.dataTransfer && Array.from(ev.dataTransfer.types || []).includes('Files')) this.dropping++; },
      dragLeave: function () { if (this.dropping > 0) this.dropping--; },
      drop: function (ev, into) {
        if (!this.can.write || !ev.dataTransfer || !Array.from(ev.dataTransfer.types || []).includes('Files')) return;
        ev.preventDefault();
        this.dropping = 0;
        const self = this, dir = into ? join(this.path, into.name) : this.path;
        droppedFiles(ev.dataTransfer).then(function (list) { self.send(list, false, dir); });
      },
      // sendParts sends a large file in parts (upload-part: UPLOAD_PART
      // bytes each, UPLOAD_PARALLEL at once, each tried again up to
      // UPLOAD_TRIES times), then ends it (upload-done): no request large
      // enough for a proxy's timeout, a part failed sent again alone.
      // onProgress gets the bytes sent of it. It resolves 'ok', 'conflict'
      // (the file there: its parts kept, x.partsId — replaced by ending it
      // again with replace) or 'failed' (its parts removed).
      sendParts: function (x, dir, replace, onProgress) {
        const self = this, file = x.file, total = file.size;
        const rnd = function () { return Math.random().toString(36).slice(2); };
        const id = x.partsId || ((window.crypto && crypto.randomUUID) ? crypto.randomUUID().replace(/-/g, '') : (Date.now().toString(36) + rnd() + rnd()).slice(0, 32));
        x.partsId = id;
        const url = function (op, extra) {
          const q = new URLSearchParams({ path: dir, name: x.path, id: id, total: String(total) });
          Object.keys(extra || {}).forEach(function (k) { q.set(k, extra[k]); });
          return self.base + 'api/fs/' + op + '?' + q.toString();
        };
        const done = function () {
          return fetch(url('upload-done', replace ? { replace: '1' } : {}), { method: 'POST' }).then(function (r) {
            if (r.status === 409) return 'conflict';
            if (r.ok) { delete x.partsId; return 'ok'; }
            throw new Error('done ' + r.status);
          });
        };
        const fail = function () {
          fetch(url('upload-abort'), { method: 'POST' }).catch(function () {});
          delete x.partsId;
          return 'failed';
        };
        // the parts sent already (a conflict replaced): only ended again
        if (x.partsSent) return done().catch(fail);
        const loaded = {};
        const report = function () { onProgress(Object.keys(loaded).reduce(function (t, k) { return t + loaded[k]; }, 0)); };
        const offsets = [];
        for (let o = 0; o < total; o += UPLOAD_PART) offsets.push(o);
        const part = function (o) {
          const end = Math.min(o + UPLOAD_PART, total);
          return new Promise(function (resolve, reject) {
            let tries = 0;
            const go = function () {
              const xhr = new XMLHttpRequest();
              xhr.open('POST', url('upload-part', { offset: String(o) }));
              xhr.upload.onprogress = function (e) { loaded[o] = e.loaded; report(); };
              xhr.onload = function () { if (xhr.status < 300) { loaded[o] = end - o; report(); resolve(); } else again(); };
              xhr.onerror = again;
              xhr.send(file.slice(o, end));
            };
            const again = function () {
              loaded[o] = 0; report();
              if (++tries >= UPLOAD_TRIES) reject(new Error('part ' + o)); else setTimeout(go, 1000 * tries);
            };
            go();
          });
        };
        let next = 0;
        const worker = function () { return next < offsets.length ? part(offsets[next++]).then(worker) : Promise.resolve(); };
        const workers = [];
        for (let i = 0; i < Math.min(UPLOAD_PARALLEL, offsets.length); i++) workers.push(worker());
        return Promise.all(workers).then(function () { x.partsSent = true; return done(); }).catch(fail);
      },
      // send posts the files (with their paths) into dir: the small ones a
      // few at a time, a large one in parts (sendParts); those that are
      // there, asked to be replaced
      send: function (list, replace, dir) {
        const self = this;
        dir = dir || this.path;
        if (!list.length) return;
        const total = list.reduce(function (t, x) { return t + x.file.size; }, 0);
        this.upload = { busy: true, done: 0, total: list.length, bytes: total, sent: 0, failed: [], pending: null };
        let sentBefore = 0;
        const batches = [], large = [];
        let cur = [], curSize = 0;
        list.forEach(function (x) {
          if (x.file.size > UPLOAD_PART || x.partsId) { large.push(x); return; }
          cur.push(x); curSize += x.file.size;
          if (cur.length >= 20 || curSize > 32 << 20) { batches.push(cur); cur = []; curSize = 0; }
        });
        if (cur.length) batches.push(cur);
        const conflicts = [];
        const one = function (batch) {
          return new Promise(function (resolve) {
            const fd = new FormData();
            batch.forEach(function (x) { fd.append('file', x.file, x.path); });
            const xhr = new XMLHttpRequest();
            xhr.open('POST', self.base + 'api/fs/upload?path=' + encodeURIComponent(dir) + (replace ? '&replace=1' : ''));
            xhr.upload.onprogress = function (e) { self.upload.sent = sentBefore + e.loaded; };
            xhr.onload = function () {
              sentBefore += batch.reduce(function (t, x) { return t + x.file.size; }, 0);
              if (xhr.status === 409) conflicts.push.apply(conflicts, batch);
              else if (xhr.status >= 300) self.upload.failed.push.apply(self.upload.failed, batch.map(function (x) { return x.path; }));
              self.upload.done += batch.length;
              resolve();
            };
            xhr.onerror = function () { self.upload.failed.push.apply(self.upload.failed, batch.map(function (x) { return x.path; })); resolve(); };
            xhr.send(fd);
          });
        };
        // a large file: in parts, its progress over what was sent before it
        const oneLarge = function (x) {
          const base = sentBefore;
          return self.sendParts(x, dir, replace, function (n) { self.upload.sent = base + n; }).then(function (res) {
            sentBefore = base + x.file.size;
            if (res === 'conflict') conflicts.push(x);
            else if (res !== 'ok') self.upload.failed.push(x.path);
            self.upload.done++;
          });
        };
        const steps = batches.map(function (b) { return function () { return one(b); }; })
          .concat(large.map(function (x) { return function () { return oneLarge(x); }; }));
        steps.reduce(function (acc, step) { return acc.then(step); }, Promise.resolve()).then(function () {
          self.upload.busy = false;
          if (conflicts.length && !replace) {
            // which, and how each differs from the one there
            self.conflictRows(conflicts, dir).then(function (rows) {
              self.upload.pending = { list: conflicts, dir: dir, rows: rows };
              self.dialog = 'replace';
            });
          } else if (self.upload.failed.length) {
            self.toast(self.w('notSent', { list: self.upload.failed.slice(0, 5).join(', ') + (self.upload.failed.length > 5 ? '…' : '') }), 'error');
          } else {
            self.toast(list.length === 1 ? self.w('sentOne') : self.w('sentMany', { n: list.length }));
          }
          self.refresh();
        });
      },
      replaceAll: function () { const p = this.upload.pending; this.dialog = ''; if (p) this.send(p.list, true, p.dir); },
      // keepTheirs keeps the files there: the parts sent of the large ones
      // removed
      keepTheirs: function () {
        const self = this, p = this.upload.pending;
        this.dialog = '';
        (p ? p.list : []).forEach(function (x) {
          if (!x.partsId) return;
          self.api('upload-abort', { path: p.dir, name: x.path, id: x.partsId }, { method: 'POST' }).catch(function () {});
          delete x.partsId; delete x.partsSent;
        });
      },
      // conflictRows are the files sent that are there (list, into dir),
      // each with its size and its time there and of the one sent
      conflictRows: function (list, dir) {
        const self = this, folders = {};
        list.forEach(function (x) {
          const i = x.path.lastIndexOf('/');
          const folder = join(dir, i < 0 ? '' : x.path.slice(0, i));
          (folders[folder] = folders[folder] || []).push(x);
        });
        return Promise.all(Object.keys(folders).map(function (folder) {
          return self.api('list', { path: folder }).then(function (r) {
            const there = {};
            ((r && r.entries) || []).forEach(function (e) { there[e.name] = e; });
            return folders[folder].map(function (x) {
              const e = there[x.path.slice(x.path.lastIndexOf('/') + 1)] || {};
              return { path: x.path, newSize: x.file.size, newTime: x.file.lastModified, oldSize: e.size, oldTime: e.mtime ? Date.parse(e.mtime) : null };
            });
          }).catch(function () {
            return folders[folder].map(function (x) { return { path: x.path, newSize: x.file.size, newTime: x.file.lastModified }; });
          });
        })).then(function (groups) { return [].concat.apply([], groups); });
      },
      // timeDiff says how the one sent is of the one there in time: newer or
      // older by how long
      timeDiff: function (r) {
        if (!r.oldTime || !r.newTime) return '';
        const d = r.newTime - r.oldTime, a = Math.abs(d) / 1000;
        if (a < 60) return this.w('sameTime');
        const n = a < 3600 ? Math.round(a / 60) + ' min' : a < 86400 ? Math.round(a / 3600) + ' h' : Math.round(a / 86400) + ' ' + (Math.round(a / 86400) === 1 ? this.w('day') : this.w('days'));
        return this.w(d > 0 ? 'newer' : 'older', { n: n });
      },
      // sizeDiff says how the one sent is of the one there in size
      sizeDiff: function (r) {
        if (r.oldSize === undefined || r.oldSize === null) return '';
        const d = r.newSize - r.oldSize;
        return d === 0 ? this.w('sameSize') : (d > 0 ? '+' : '−') + size(Math.abs(d), this.loc());
      },
      // moving by dragging: the files selected onto a folder
      dragStart: function (e, ev) {
        if (!this.isSel(e)) this.selected = [e.name];
        ev.dataTransfer.setData('application/x-vx-files', JSON.stringify(this.selected));
        ev.dataTransfer.effectAllowed = 'move';
      },
      dropOn: function (target, ev) {
        if (!this.can.write) return;
        const data = ev.dataTransfer && ev.dataTransfer.getData('application/x-vx-files');
        if (!data) { this.drop(ev, target); return; }
        ev.preventDefault();
        ev.stopPropagation();
        this.dropping = 0;
        const self = this, here = this.path, names = JSON.parse(data).filter(function (n) { return n !== target.name; });
        names.reduce(function (acc, n) {
          return acc.then(function () { return self.post('rename', { path: join(here, n), to: join(join(here, target.name), n.split('/').pop()) }); });
        }, Promise.resolve()).then(function () { self.toast(self.w('movedTo', { name: target.name })); self.refresh(); }).catch(this.failed);
      },
      dropOver: function (target, ev) {
        if (!target.dir || !this.can.write) return;
        const types = Array.from((ev.dataTransfer && ev.dataTransfer.types) || []);
        if (types.includes('application/x-vx-files') || types.includes('Files')) { ev.preventDefault(); }
      },
      date: function (s) { const d = new Date(s), l = this.loc(); return d.toLocaleDateString(l) + ' ' + d.toLocaleTimeString(l, { hour: '2-digit', minute: '2-digit' }); },
      fmtSize: function (n) { return size(n, this.loc()); },
      icon: iconOf,
      iconClass: function (e) { return 'vx-files-ic vx-files-ic-' + kindOf(e); },
      isArchive: function (e) { return e && !e.dir && ARCHIVE.test(e.name); },
      isImage: function (e) { return e && !e.dir && IMAGE.test(e.name); },
      // paths, for the template (no regular expression there: a "\\/" of
      // this file's template string would be "//", a comment, in its code)
      parentOf: function (p) { return parent(p); },
      // w is the word of key (vxFiles), its placeholders ({n}) of params
      w: function (key, params) { return words(this.i18n, key, params); },
      // loc is the locale of the page ("pt-BR"): its dates and numbers
      loc: function () { return locOf(this.i18n); },
      childOf: function (name) { return join(this.path, name); },
    },
    template: `
<div class="vx-files" tabindex="0" @keydown="key" @dragenter="dragEnter" @dragleave="dragLeave" @dragover="dragOver" @drop="drop($event)">
  <div class="d-flex flex-wrap align-center ga-1 mb-2">
    <v-btn icon="mdi-arrow-left" v-bind="tb" :disabled="!back.length" :title="w('back')" @click="goBack"></v-btn>
    <v-btn icon="mdi-arrow-right" v-bind="tb" :disabled="!fwd.length" :title="w('forward')" @click="goFwd"></v-btn>
    <v-btn icon="mdi-arrow-up" v-bind="tb" :disabled="path === '/'" :title="w('up')" @click="up"></v-btn>
    <div class="vx-files-path flex-grow-1" @dblclick="editPath">
      <v-text-field v-if="editingPath" v-model="pathText" density="compact" hide-details variant="outlined" autofocus
        @keydown.enter="submitPath" @keydown.esc="editingPath = false" @blur="editingPath = false"></v-text-field>
      <div v-else class="d-flex align-center flex-wrap">
        <template v-for="(c, i) in crumbs" :key="c.path">
          <v-icon v-if="i" size="small" class="text-medium-emphasis">mdi-chevron-right</v-icon>
          <v-btn size="small" variant="text" class="text-none px-2" :prepend-icon="i ? undefined : 'mdi-home-outline'"
            @click="go(c.path)" @dragover="dropOver({dir: true, name: ''}, $event)">{{ c.name }}</v-btn>
        </template>
        <v-btn icon="mdi-pencil-outline" size="x-small" variant="text" :title="w('typePath')" @click="editPath"></v-btn>
      </div>
    </div>
    <v-text-field v-model="filter" density="compact" hide-details variant="outlined" prepend-inner-icon="mdi-magnify"
      :placeholder="w('filter')" clearable style="max-width: 200px"></v-text-field>
    <v-btn icon="mdi-refresh" v-bind="tb" :title="w('refresh')" :loading="loading" @click="refresh"></v-btn>
    <template v-if="can.read">
      <v-menu v-if="sel.length > 1">
        <template #activator="{ props }"><v-btn v-bind="Object.assign({}, tb, props)" icon="mdi-download" :title="downloadTitle()"></v-btn></template>
        <v-list density="compact">
          <v-list-item prepend-icon="mdi-folder-zip-outline" :title="w('downloadNTgz', { n: sel.length })" @click="downloadSel(false)"></v-list-item>
          <v-list-item prepend-icon="mdi-download-multiple" :title="w('downloadNEach', { n: sel.length })" @click="downloadSel(true)"></v-list-item>
        </v-list>
      </v-menu>
      <v-btn v-else icon="mdi-download" v-bind="tb" :title="downloadTitle()" @click="sel.length ? downloadSel() : download({ dir: true, name: '' })"></v-btn>
    </template>
    <v-menu v-if="can.write">
      <template #activator="{ props }"><v-btn v-bind="Object.assign({}, tb, props)" icon="mdi-plus" :title="w('new')"></v-btn></template>
      <v-list density="compact">
        <v-list-item prepend-icon="mdi-folder-plus-outline" :title="w('folder')" @click="ask('mkdir')"></v-list-item>
        <v-list-item prepend-icon="mdi-file-plus-outline" :title="w('textFile')" @click="ask('newfile')"></v-list-item>
        <v-divider></v-divider>
        <v-list-item prepend-icon="mdi-upload" :title="w('uploadFiles')" @click="pick(false)"></v-list-item>
        <v-list-item prepend-icon="mdi-folder-upload-outline" :title="w('uploadFolder')" @click="pick(true)"></v-list-item>
        <v-list-item v-if="fetch" prepend-icon="mdi-web" :title="w('fetch')" @click="ask('fetch')"></v-list-item>
      </v-list>
    </v-menu>
    <div class="vx-files-views" role="group" :aria-label="w('view')">
      <v-btn v-for="v in views" :key="v.value" :icon="v.icon" v-bind="on(view === v.value)" :title="v.title"
        :aria-pressed="view === v.value" @click="setView(v.value)"></v-btn>
    </div>
    <v-btn v-if="terminal" icon="mdi-console" v-bind="on(term.open)"
      :title="w(term.open ? 'termClose' : 'termOpen')" @click="toggleTerm"></v-btn>
    <v-btn icon="mdi-image-outline" v-bind="on(thumbs)"
      :title="w(thumbs ? 'thumbsHide' : 'thumbsShow')" @click="toggleThumbs"></v-btn>
    <v-btn icon="mdi-eye-outline" v-bind="on(hidden)"
      :title="w(hidden ? 'hiddenHide' : 'hiddenShow')" @click="toggleHidden"></v-btn>
    <template v-if="repo">
      <v-btn :href="gitHere()" icon="mdi-source-branch" v-bind="tb" :title="w('gitOpen')"></v-btn>
      <v-btn icon="mdi-git" v-bind="tb" :title="w('cloneCopy')" @click="copyClone()"></v-btn>
    </template>
    <v-btn v-if="ideUrl" :href="ideHere()" icon="mdi-code-braces" v-bind="tb" :title="w('ideOpen')"></v-btn>
    <slot name="toolbar"></slot>
    <input ref="files" type="file" multiple hidden @change="picked">
    <input ref="folder" type="file" webkitdirectory hidden @change="picked">
  </div>

  <v-alert v-if="error" type="error" density="compact" class="mb-2">{{ error }}</v-alert>
  <v-progress-linear v-if="upload.busy" :model-value="upload.bytes ? upload.sent / upload.bytes * 100 : 0" color="primary" height="18" rounded class="mb-2">
    <span class="text-caption">{{ w('uploading', { done: upload.done, total: upload.total, sent: fmtSize(upload.sent), bytes: fmtSize(upload.bytes) }) }}</span>
  </v-progress-linear>

  <div ref="area" class="vx-files-area" :class="{ 'vx-files-dropping': dropping > 0 }" :style="areaH ? { height: areaH + 'px', overflow: 'auto', minHeight: 0 } : {}"
    @contextmenu.prevent="ctx(null, $event)" @click.self="selected = []">
    <div v-if="dropping > 0" class="vx-files-drophint"><v-icon size="48">mdi-tray-arrow-down</v-icon><div>{{ w('dropHere', { path: path }) }}</div></div>
    <v-table v-if="view === 'list' || view === 'tree'" density="compact" hover class="vx-files-table" fixed-header>
      <thead><tr>
        <th @click="setSort('name')">{{ w('name') }} <v-icon v-if="sortBy === 'name'" size="x-small">{{ sortAsc ? 'mdi-arrow-up' : 'mdi-arrow-down' }}</v-icon></th>
        <th class="text-end" @click="setSort('size')" style="width: 110px">{{ w('size') }} <v-icon v-if="sortBy === 'size'" size="x-small">{{ sortAsc ? 'mdi-arrow-up' : 'mdi-arrow-down' }}</v-icon></th>
        <th v-if="cols.dims" class="text-end" style="width: 120px">{{ w('dimensions') }}</th>
        <th @click="setSort('mtime')" style="width: 170px">{{ w('modified') }} <v-icon v-if="sortBy === 'mtime'" size="x-small">{{ sortAsc ? 'mdi-arrow-up' : 'mdi-arrow-down' }}</v-icon></th>
        <th style="width: 150px" @contextmenu.prevent.stop>{{ w('permissions') }}
          <v-menu location="bottom end">
            <template #activator="{ props }"><v-btn v-bind="props" icon="mdi-table-cog" size="x-small" variant="text" density="comfortable"
              class="float-end" :title="w('columns')" @click.stop></v-btn></template>
            <v-list density="compact">
              <v-list-subheader>{{ w('columns') }}</v-list-subheader>
              <v-list-item :prepend-icon="cols.dims ? 'mdi-checkbox-marked-outline' : 'mdi-checkbox-blank-outline'"
                :title="w('imageDims')" @click="toggleCol('dims')"></v-list-item>
            </v-list>
          </v-menu>
        </th>
      </tr></thead>
      <tbody>
        <tr v-for="e in shown" :key="e.name" :class="{ 'vx-files-sel': isSel(e) }" :draggable="can.write"
          @click="click(e, $event)" @dblclick="open(e)" @contextmenu.prevent.stop="ctx(e, $event)"
          @dragstart="dragStart(e, $event)" @dragover="dropOver(e, $event)" @drop="e.dir && dropOn(e, $event)">
          <td :style="view === 'tree' ? { paddingLeft: (8 + e.depth * 18) + 'px' } : undefined"><template v-if="view === 'tree'"><v-icon v-if="e.dir" size="small" class="vx-files-twisty"
              @click.stop="toggleNode(e)" @dblclick.stop>{{ opened[e.name] ? 'mdi-chevron-down' : 'mdi-chevron-right' }}</v-icon><span v-else class="vx-files-twisty"></span></template><span v-if="thumbs && isImage(e)" class="vx-files-thumb vx-files-thumb-sm me-2"><v-icon size="small" :class="[iconClass(e), { 'vx-files-thumb-off': thumbOk[e.name] }]">{{ icon(e) }}</v-icon><img
              :class="{ 'vx-files-thumb-on': thumbOk[e.name] }" :src="imageURL(e.name)" loading="lazy" alt="" draggable="false" @load="thumbOk[e.name] = true"></span><v-icon v-else size="small" class="me-2" :class="iconClass(e)">{{ icon(e) }}</v-icon>{{ e.base || e.name }}
            <v-icon v-if="e.repo" size="x-small" color="deep-orange" class="ms-1" :title="w('gitRepo')">mdi-git</v-icon>
            <span v-if="e.link" class="text-caption text-medium-emphasis"> → {{ e.link }}</span></td>
          <td class="text-end text-medium-emphasis">{{ e.dir ? '' : fmtSize(e.size) }}</td>
          <td v-if="cols.dims" class="text-end text-medium-emphasis">{{ dims(e) }}</td>
          <td class="text-medium-emphasis">{{ date(e.mtime) }}</td>
          <td class="text-medium-emphasis vx-files-mode">{{ e.mode }}</td>
        </tr>
        <tr v-if="!shown.length && !loading"><td :colspan="cols.dims ? 5 : 4" class="text-center text-medium-emphasis py-6">{{ w(filter ? 'noMatch' : 'emptyFolder') }}</td></tr>
      </tbody>
    </v-table>
    <div v-else class="vx-files-grid">
      <div v-for="e in shown" :key="e.name" class="vx-files-cell" :class="{ 'vx-files-sel': isSel(e) }" :draggable="can.write" :title="e.name"
        @click="click(e, $event)" @dblclick="open(e)" @contextmenu.prevent.stop="ctx(e, $event)"
        @dragstart="dragStart(e, $event)" @dragover="dropOver(e, $event)" @drop="e.dir && dropOn(e, $event)">
        <div v-if="thumbs && isImage(e)" class="vx-files-thumb vx-files-thumb-lg">
          <v-icon size="48" :class="[iconClass(e), { 'vx-files-thumb-off': thumbOk[e.name] }]">{{ icon(e) }}</v-icon>
          <img :class="{ 'vx-files-thumb-on': thumbOk[e.name] }" :src="imageURL(e.name)" loading="lazy" alt="" draggable="false" @load="thumbOk[e.name] = true">
        </div>
        <v-badge v-else-if="e.repo" icon="mdi-git" color="deep-orange" location="bottom end" :title="w('gitRepo')">
          <v-icon size="48" :class="iconClass(e)">{{ icon(e) }}</v-icon>
        </v-badge>
        <v-icon v-else size="48" :class="iconClass(e)">{{ icon(e) }}</v-icon>
        <div class="vx-files-name">{{ e.name }}</div>
      </div>
      <div v-if="!shown.length && !loading" class="text-center text-medium-emphasis py-6 w-100">{{ w(filter ? 'noMatch' : 'emptyFolder') }}</div>
    </div>
  </div>

  <div class="d-flex text-caption text-medium-emphasis mt-1 vx-files-status">
    <span>{{ shown.length === 1 ? w('itemsOne') : w('itemsMany', { n: shown.length }) }}</span>
    <span v-if="sel.length" class="ms-2">· {{ w('selected', { n: sel.length }) }}{{ selSize ? ' (' + fmtSize(selSize) + ')' : '' }}</span>
    <v-spacer></v-spacer>
    <span>{{ w('hint') }}</span>
  </div>

  <div v-if="terminal && term.open" class="vx-files-sash" :class="{ active: term.dragging }" role="separator"
    aria-orientation="horizontal" tabindex="0" :aria-valuenow="term.height" :aria-label="w('termHeight')"
    :title="w('termSash')"
    @pointerdown="termDrag" @dblclick="termResize(260)"
    @keydown.up.prevent="termResize(term.height + 20)" @keydown.down.prevent="termResize(term.height - 20)"><span></span></div>
  <div v-if="terminal && term.open" class="vx-files-term" :style="{ height: term.height + 'px' }">
    <vx-terminal ref="term" :ws-url="termURL()" :vendor="vendor" :root-var="termRootVar" :recorded="termRecorded" fill></vx-terminal>
  </div>

  <v-menu v-model="menu.open" :target="[menu.x, menu.y]" location="bottom start">
    <v-list density="compact" min-width="220">
      <template v-if="sel.length">
        <v-list-item v-if="one" prepend-icon="mdi-open-in-app" :title="one.dir ? w('open') : (extract && isArchive(one) ? w('extractDots') : w('open'))" @click="open(one)"></v-list-item>
        <v-list-item v-if="can.write && one && !one.dir && !isArchive(one) && !isImage(one)" prepend-icon="mdi-pencil-outline" :title="w('edit')" @click="edit(one, true)"></v-list-item>
        <v-list-item v-if="one && one.dir" prepend-icon="mdi-code-braces" :title="w('openIde')" :href="ideHere(childOf(one.name))"></v-list-item>
        <v-list-item v-if="one && one.repo" prepend-icon="mdi-source-branch" :title="w('openGit')" :href="gitHere(childOf(one.name))"></v-list-item>
        <v-list-item v-if="one && one.repo" prepend-icon="mdi-git" :title="w('cloneCopy')" @click="copyClone(one)"></v-list-item>
        <template v-if="can.read">
          <template v-if="sel.length > 1">
            <v-list-item prepend-icon="mdi-folder-zip-outline" :title="w('downloadNTgz', { n: sel.length })" @click="downloadSel(false)"></v-list-item>
            <v-list-item prepend-icon="mdi-download-multiple" :title="w('downloadNEach', { n: sel.length })" @click="downloadSel(true)"></v-list-item>
          </template>
          <v-list-item v-else prepend-icon="mdi-download" :title="downloadTitle()" @click="downloadSel()"></v-list-item>
        </template>
        <template v-if="can.write">
          <v-list-item v-if="extract && isArchive(one)" prepend-icon="mdi-package-variant" :title="w('extractDots')" @click="ask('extract')"></v-list-item>
          <v-divider></v-divider>
          <v-list-item v-if="one" prepend-icon="mdi-form-textbox" :title="w('renameF2')" @click="ask('rename')"></v-list-item>
          <v-list-item prepend-icon="mdi-content-copy" :title="w('copyTo')" @click="ask('copy')"></v-list-item>
          <v-list-item prepend-icon="mdi-file-move-outline" :title="w('moveTo')" @click="ask('move')"></v-list-item>
          <v-divider></v-divider>
          <v-list-item prepend-icon="mdi-delete-outline" base-color="error" :title="w('deleteKey')" @click="ask('delete')"></v-list-item>
        </template>
      </template>
      <template v-else>
        <template v-if="can.write">
          <v-list-item prepend-icon="mdi-folder-plus-outline" :title="w('newFolder')" @click="ask('mkdir')"></v-list-item>
          <v-list-item prepend-icon="mdi-file-plus-outline" :title="w('newTextFile')" @click="ask('newfile')"></v-list-item>
          <v-list-item prepend-icon="mdi-upload" :title="w('uploadFiles')" @click="pick(false)"></v-list-item>
          <v-list-item prepend-icon="mdi-folder-upload-outline" :title="w('uploadFolder')" @click="pick(true)"></v-list-item>
          <v-list-item v-if="fetch" prepend-icon="mdi-web" :title="w('fetch')" @click="ask('fetch')"></v-list-item>
        </template>
        <v-list-item v-if="can.read" prepend-icon="mdi-download" :title="w('downloadFolder')" @click="download({ dir: true, name: '' })"></v-list-item>
        <v-list-item prepend-icon="mdi-select-all" :title="w('selectAll')" @click="selectAll"></v-list-item>
      </template>
    </v-list>
  </v-menu>

  <v-dialog :model-value="!!dialog && dialog !== 'replace'" max-width="520" @update:model-value="v => { if (!v) dialog = '' }">
    <v-card>
      <v-card-title class="text-subtitle-1">{{ w({ mkdir: 'newFolder', newfile: 'newFile', rename: 'rename', copy: 'copyToTitle', move: 'moveToTitle', delete: 'delete', extract: 'extract', fetch: 'fetch' }[dialog] || 'new') }}</v-card-title>
      <v-card-text>
        <v-text-field v-if="['mkdir', 'newfile', 'rename'].includes(dialog)" v-model="form.name" :label="w('name')" autofocus density="compact" @keydown.enter="form.name && run()"></v-text-field>
        <template v-if="dialog === 'copy' || dialog === 'move' || dialog === 'extract'">
          <div class="text-caption mb-2" v-if="dialog !== 'extract'">{{ form.confirm.join(', ') }}</div>
          <v-text-field v-model="form.to" :label="w('toFolder')" :hint="w('toFolderHint')" persistent-hint autofocus density="compact" @keydown.enter="run()"></v-text-field>
        </template>
        <template v-if="dialog === 'fetch'">
          <v-text-field v-model="form.url" :label="w('address')" autofocus density="compact" placeholder="https://…"></v-text-field>
          <v-text-field v-model="form.name" :label="w('saveAs')" density="compact" :hint="w('inPath', { path: path })" persistent-hint @keydown.enter="form.url && run()"></v-text-field>
        </template>
        <div v-if="dialog === 'delete'">{{ w('deleteConfirm', { what: form.confirm.length === 1 ? form.confirm[0] : w('itemsMany', { n: form.confirm.length }) }) }}</div>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn variant="text" @click="dialog = ''">{{ w('cancel') }}</v-btn>
        <v-btn :color="dialog === 'delete' ? 'error' : 'primary'" variant="flat" @click="run">{{ w({ delete: 'delete', fetch: 'download', extract: 'extract', copy: 'copy', move: 'move', rename: 'rename' }[dialog] || 'create') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog :model-value="dialog === 'replace'" max-width="860" scrollable @update:model-value="v => { if (!v) keepTheirs() }">
    <v-card>
      <v-card-title class="text-subtitle-1">{{ w('existTitle') }}</v-card-title>
      <v-card-text>
        <div class="mb-2">{{ w('existText', { n: upload.pending ? upload.pending.list.length : 0, dir: upload.pending && upload.pending.dir }) }}</div>
        <v-table density="compact" class="vx-files-conflicts">
          <thead><tr><th>{{ w('colFile') }}</th><th>{{ w('colThere') }}</th><th>{{ w('colSent') }}</th><th>{{ w('colDiff') }}</th></tr></thead>
          <tbody>
            <tr v-for="r in (upload.pending && upload.pending.rows) || []" :key="r.path">
              <td class="text-no-wrap"><code>{{ r.path }}</code></td>
              <td class="text-no-wrap text-medium-emphasis">{{ r.oldSize !== undefined ? fmtSize(r.oldSize) : '?' }} · {{ r.oldTime ? date(r.oldTime) : '?' }}</td>
              <td class="text-no-wrap">{{ fmtSize(r.newSize) }} · {{ r.newTime ? date(r.newTime) : '?' }}</td>
              <td class="text-no-wrap">
                <span :class="{'text-success': r.newTime > r.oldTime, 'text-warning': r.newTime < r.oldTime}">{{ timeDiff(r) }}</span>
                <span v-if="sizeDiff(r)" class="text-medium-emphasis"> · {{ sizeDiff(r) }}</span>
              </td>
            </tr>
          </tbody>
        </v-table>
      </v-card-text>
      <v-card-actions><v-spacer></v-spacer>
        <v-btn variant="text" @click="keepTheirs">{{ w('keepTheirs') }}</v-btn>
        <v-btn color="warning" variant="flat" @click="replaceAll">{{ w('replace') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog v-model="editor.open" fullscreen persistent>
    <v-card class="d-flex flex-column" @keydown="editorKey">
      <v-toolbar density="compact">
        <v-btn icon="mdi-close" :title="w('closeEsc')" @click="closeEditor"></v-btn>
        <v-toolbar-title class="text-subtitle-1">{{ editor.path }}<span v-if="dirty" :title="w('changed')"> •</span></v-toolbar-title>
        <v-chip size="small" variant="tonal" class="me-2">{{ editor.language || w('text') }}</v-chip>
        <v-chip size="small" variant="tonal" class="me-2">{{ fmtSize(editor.size) }}</v-chip>
        <v-chip v-if="editor.truncated" size="small" color="warning" class="me-2">{{ w('truncated') }}</v-chip>
        <v-chip v-else-if="!editor.writable" size="small" color="grey" prepend-icon="mdi-lock-outline" class="me-2">{{ w('noWrite') }}</v-chip>
        <template v-if="editor.editing">
          <v-btn variant="text" :disabled="editor.saving" @click="cancelEdit">{{ w('cancel') }}</v-btn>
          <v-btn color="primary" variant="flat" prepend-icon="mdi-content-save-outline" class="me-2" :loading="editor.saving" :disabled="!dirty" @click="save">{{ w('saveCtrlS') }}</v-btn>
        </template>
        <v-btn v-else-if="editor.writable" color="primary" variant="tonal" prepend-icon="mdi-pencil-outline" class="me-2" :title="w('editCtrlE')" @click="startEdit">{{ w('edit') }}</v-btn>
        <v-btn :href="ideHere(parentOf(editor.path))" prepend-icon="mdi-code-braces" variant="text" class="me-2" :title="w('ideFolder')">IDE</v-btn>
      </v-toolbar>
      <div class="flex-grow-1 vx-files-viewer" :class="{ 'vx-files-viewer-edit': editor.editing }">
        <vx-codemirror v-if="editor.open" v-model="editor.content" :language="editor.language" :readonly="!editor.editing"
          :tab-size="4" min-height="calc(100vh - 64px)" max-height="calc(100vh - 64px)"></vx-codemirror>
      </div>
    </v-card>
  </v-dialog>

  <v-dialog v-model="image.open" fullscreen>
    <div class="vx-files-image" tabindex="0" @keydown="imageKey" ref="imageBox">
      <div class="vx-files-image-bar">
        <v-btn icon="mdi-close" variant="text" color="white" :title="w('closeEsc')" @click="image.open = false"></v-btn>
        <div class="text-subtitle-1 text-truncate flex-grow-1">{{ image.name }}</div>
        <span class="text-caption me-3">{{ fmtSize(imageEntry().size) }} · {{ images().findIndex(e => e.name === image.name) + 1 }}/{{ images().length }}</span>
        <v-btn :icon="image.actual ? 'mdi-fit-to-screen-outline' : 'mdi-magnify-plus-outline'" variant="text" color="white"
          :title="w(image.actual ? 'fitScreen' : 'actualSize')" @click="image.actual = !image.actual"></v-btn>
        <v-btn icon="mdi-download" variant="text" color="white" :title="w('download')" @click="download(imageEntry())"></v-btn>
      </div>
      <v-btn v-if="images().length > 1" class="vx-files-image-prev" icon="mdi-chevron-left" variant="tonal" color="white" :title="w('prev')" @click="imageStep(-1)"></v-btn>
      <div class="vx-files-image-stage" :class="{ 'vx-files-image-actual': image.actual }" @click.self="image.open = false">
        <img :src="imageURL(image.name)" :alt="image.name" @click="image.actual = !image.actual">
      </div>
      <v-btn v-if="images().length > 1" class="vx-files-image-next" icon="mdi-chevron-right" variant="tonal" color="white" :title="w('next')" @click="imageStep(1)"></v-btn>
    </div>
  </v-dialog>

  <v-dialog v-model="editor.confirm" max-width="460">
    <v-card>
      <v-card-title class="text-subtitle-1">{{ w('discardTitle') }}</v-card-title>
      <v-card-text>{{ w('discardText', { path: editor.path }) }}</v-card-text>
      <v-card-actions><v-spacer></v-spacer>
        <v-btn variant="text" @click="editor.confirm = false">{{ w('keepEditing') }}</v-btn>
        <v-btn color="error" variant="flat" @click="cancelEdit(); closeEditor(true)">{{ w('discard') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-snackbar v-model="snack.show" :color="snack.color" timeout="3500">{{ snack.text }}</v-snackbar>
</div>`,
  };

  const registers = window.__goplaidVueComponentRegisters = window.__goplaidVueComponentRegisters || [];
  registers.push(function (app) {
    app.component('vx-files', VxFiles);
  });
})();
