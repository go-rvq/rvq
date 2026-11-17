// vx-terminal: a terminal (xterm.js, loaded of vendor when first used) on a
// websocket — its binary messages the bytes of the session, a text one
// {"exit": …} its end; sent {"resize": [cols, rows]} —: the shell of a user
// of the OS in senir, the sandboxed one of the file system in HermonCMS.
(function () {
  'use strict';

  const loaded = {};
  // load is a script (or a style) of base, once.
  function load(base, file) {
    const url = base + '/' + file;
    if (!loaded[url]) {
      loaded[url] = new Promise(function (resolve, reject) {
        let el;
        if (file.endsWith('.css')) {
          el = document.createElement('link');
          el.rel = 'stylesheet';
          el.href = url;
        } else {
          el = document.createElement('script');
          el.src = url;
        }
        el.onload = resolve;
        el.onerror = function () { reject(new Error('could not load ' + url)); };
        document.head.appendChild(el);
      });
    }
    return loaded[url];
  }

  function wsURL(path) {
    const u = new URL(path, window.location.href);
    u.protocol = u.protocol === 'https:' ? 'wss:' : 'ws:';
    return u.toString();
  }

  const VxTerminal = {
    props: { wsUrl: String, vendor: String, keep: Boolean, fill: Boolean,
      // rootVar: the variable of the root of the files in the shell (cd rooted); recorded: the session is recorded (said)
      rootVar: { type: String, default: 'FILES_ROOT' }, recorded: Boolean },
    data: function () { return { state: 'connecting', exit: '', fontSize: 12, expanded: false }; },
    mounted: function () {
      const self = this;
      Promise.all([
        load(this.vendor, 'xterm.css'),
        load(this.vendor, 'xterm.js'),
      ]).then(function () {
        return Promise.all([load(self.vendor, 'addon-fit.js'), load(self.vendor, 'addon-web-links.js'), load(self.vendor, 'addon-search.js')]);
      }).then(function () { self.start(); }).catch(function (e) { self.state = 'error'; self.exit = e.message; });
    },
    beforeUnmount: function () {
      if (this.sizeBox) window.removeEventListener('resize', this.sizeBox);
      if (this.ro) this.ro.disconnect();
      if (this.ws) this.ws.close();
      if (this.term) this.term.dispose();
    },
    methods: {
      start: function () {
        const self = this;
        this.state = 'connecting';
        this.exit = '';
        const dark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
        const term = new window.Terminal({
          cursorBlink: true, fontSize: this.fontSize, scrollback: 5000, allowProposedApi: true,
          fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Consolas, monospace',
          theme: dark ? {} : { background: '#1e1e1e' },
        });
        const fit = new window.FitAddon.FitAddon();
        term.loadAddon(fit);
        term.loadAddon(new window.WebLinksAddon.WebLinksAddon());
        this.search = new window.SearchAddon.SearchAddon();
        term.loadAddon(this.search);
        // the box down to the bottom of the window, where it really starts
        // (the header above it varies); again when the window changes
        const box = this.$refs.el;
        this.sizeBox = function () {
          if (self.fill) return;
          if (self.expanded) {
            // the whole window, but for the bar of the buttons
            box.style.height = Math.max(200, window.innerHeight - box.getBoundingClientRect().top - 8) + 'px';
            return;
          }
          const top = box.getBoundingClientRect().top + window.scrollY;
          box.style.height = Math.max(320, window.innerHeight - top - 16) + 'px';
        };
        this.sizeBox();
        window.addEventListener('resize', this.sizeBox);
        term.open(box);
        fit.fit();
        this.term = term;
        this.fit = fit;
        const ws = new WebSocket(wsURL(this.wsUrl) + '&cols=' + term.cols + '&rows=' + term.rows);
        ws.binaryType = 'arraybuffer';
        this.ws = ws;
        const enc = new TextEncoder();
        ws.onopen = function () { self.state = 'open'; term.focus(); };
        ws.onmessage = function (ev) {
          if (typeof ev.data === 'string') {
            const m = JSON.parse(ev.data);
            if (m.exit !== undefined) { self.exit = m.exit; self.state = 'closed'; }
            return;
          }
          term.write(new Uint8Array(ev.data));
        };
        ws.onclose = function () { if (self.state !== 'closed') { self.state = 'closed'; self.exit = self.exit || 'conexão encerrada'; } };
        term.onData(function (d) { if (ws.readyState === 1) ws.send(enc.encode(d)); });
        term.onBinary(function (d) {
          if (ws.readyState !== 1) return;
          const b = new Uint8Array(d.length);
          for (let i = 0; i < d.length; i++) b[i] = d.charCodeAt(i) & 255;
          ws.send(b);
        });
        term.onResize(function (sz) { if (ws.readyState === 1) ws.send(JSON.stringify({ resize: [sz.cols, sz.rows] })); });
        this.ro = new ResizeObserver(function () { try { fit.fit(); } catch (e) { } });
        this.ro.observe(this.$refs.el);
      },
      // expand: the terminal over the whole window — in a page or in a
      // panel (fill) —, the same button gives it back (Esc stays the
      // terminal's); the session goes on
      toggleExpand: function () {
        const self = this;
        this.expanded = !this.expanded;
        this.$nextTick(function () {
          if (self.sizeBox) self.sizeBox();
          if (self.fit) { try { self.fit.fit(); } catch (e) { } }
          if (self.term) self.term.focus();
        });
      },
      zoom: function (d) {
        this.fontSize = Math.max(9, Math.min(28, this.fontSize + d));
        if (this.term) { this.term.options.fontSize = this.fontSize; this.fit.fit(); }
      },
      find: function () {
        const q = window.prompt('Buscar no terminal:');
        if (q && this.search) this.search.findNext(q);
      },
      reopen: function () {
        if (!this.fill) { window.location.reload(); return; }
        // in place: a new session in the same box
        if (this.ro) this.ro.disconnect();
        if (this.ws) { this.ws.onclose = null; this.ws.close(); }
        if (this.term) this.term.dispose();
        this.$refs.el.innerHTML = '';
        this.start();
      },
      // cd: the shell to the folder dir — of the root of the files (the
      // explorer's: $<rootVar>) when rooted —, typed as Dolphin does: the
      // line cleared, a space first (out of the history)
      cd: function (dir, rooted) {
        if (!this.ws || this.ws.readyState !== 1) return;
        const q = "'" + String(dir).replace(/'/g, "'\\''") + "'";
        this.ws.send(new TextEncoder().encode('\x15 cd -- ' + (rooted ? '"$' + this.rootVar + '"' : '') + q + '\r'));
      },
    },
    template: `
<div class="vx-terminal" :class="{ 'vx-term-expanded': expanded, 'vx-term-fill': fill }">
  <div class="d-flex align-center ga-2" :class="fill ? 'vx-term-bar' : 'mb-2'">
    <v-chip size="small" :color="state === 'open' ? 'success' : (state === 'closed' ? 'grey' : 'warning')" variant="tonal">
      {{ state === 'open' ? 'conectado' : (state === 'closed' ? 'encerrado' : 'conectando…') }}
    </v-chip>
    <v-chip v-if="keep" size="small" variant="tonal" prepend-icon="mdi-layers-outline">tmux (a sessão continua ao fechar)</v-chip>
    <v-chip v-if="recorded" size="small" variant="tonal" color="info" prepend-icon="mdi-record-rec">gravado</v-chip>
    <v-spacer></v-spacer>
    <v-btn size="small" variant="text" icon="mdi-magnify" title="Buscar" @click="find"></v-btn>
    <v-btn size="small" variant="text" class="text-none px-1 vx-term-font" min-width="0" title="Diminuir a fonte" @click="zoom(-1)">A-</v-btn>
    <v-btn size="small" variant="text" class="text-none px-1 vx-term-font" min-width="0" title="Aumentar a fonte" @click="zoom(1)">A+</v-btn>
    <v-btn v-if="fill && state === 'closed'" size="small" variant="text" icon="mdi-restart" title="Reabrir" @click="reopen"></v-btn>
    <v-btn size="small" variant="text" :icon="expanded ? 'mdi-arrow-collapse' : 'mdi-arrow-expand'"
      :title="expanded ? 'Voltar ao tamanho de antes' : 'Expandir o terminal (a janela toda)'" @click="toggleExpand"></v-btn>
  </div>
  <div ref="el" class="vx-term-box"></div>
  <v-alert v-if="state === 'closed' && !fill" type="info" density="compact" class="mt-2">
    Sessão encerrada: {{ exit }}
    <template v-slot:append><v-btn size="small" variant="tonal" @click="reopen">Reabrir</v-btn></template>
  </v-alert>
</div>`,
  };

  const registers = window.__goplaidVueComponentRegisters = window.__goplaidVueComponentRegisters || [];
  registers.push(function (app) {
    app.component('vx-terminal', VxTerminal);
  });
})();
