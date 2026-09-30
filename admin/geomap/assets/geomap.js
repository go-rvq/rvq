// rvq-geo-map: a map of points — where the likes and the comments of a
// record came from, where a contact message was sent from, a user's signup
// and sessions —, on Leaflet with OpenStreetMap's tiles. Each point is
// {lat, lng, kind: "like"|"comment"|"place"|"session", count, label}: a
// circle, as large as its count, the colour of its kind, and its label on
// hover. zoom: how close a map of a single point is.
(function () {
  window.__goplaidVueComponentRegisters = window.__goplaidVueComponentRegisters || [];
  window.__goplaidVueComponentRegisters.push(function (app) {
    app.component('rvq-geo-map', {
      props: {
        points: { type: Array, default: function () { return []; } },
        height: { type: String, default: '420px' },
        zoom: { type: Number, default: 6 },
      },
      template: '<div ref="el" :style="{height: height, width: \'100%\', borderRadius: \'8px\'}"></div>',
      mounted: function () {
        if (!window.L) return;
        var map = window.L.map(this.$refs.el, { scrollWheelZoom: false }).setView([20, 0], 2);
        window.L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
          maxZoom: 18,
          attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
        }).addTo(map);
        var bounds = [];
        (this.points || []).forEach(function (p) {
          if (p.lat == null || p.lng == null) return;
          var colour = { like: '#e91e63', comment: '#1976d2', place: '#2e7d32', session: '#f57c00' }[p.kind] || '#1976d2';
          window.L.circleMarker([p.lat, p.lng], {
            radius: 5 + Math.min(20, Math.sqrt(p.count || 1) * 3),
            color: colour,
            fillColor: colour,
            fillOpacity: 0.45,
            weight: 1,
          }).bindTooltip(p.label || '').addTo(map);
          bounds.push([p.lat, p.lng]);
        });
        if (bounds.length > 1) map.fitBounds(bounds, { padding: [24, 24], maxZoom: 8 });
        else if (bounds.length === 1) map.setView(bounds[0], this.zoom);
        this.map = map;
        // drawn inside a card that sized itself after: measure again
        setTimeout(function () { map.invalidateSize(); }, 200);
      },
      beforeUnmount: function () {
        if (this.map) this.map.remove();
      },
    });
  });
})();
