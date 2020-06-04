<template>
  <q-layout view="lHh Lpr lFf" class="bg-white">
    <HeaderWithBack
      title="{{.DataValues.title}}"
      icon="{{.DataValues.icon}}"
      :loading="loading.header"
    ></HeaderWithBack>
    <q-page-container>
      <q-page class="row">
        <div class="col-6">
          <q-select
            v-model="selectedMapCategory"
            :disable="loading.select"
            :loading="loading.select"
            bottom-slots=""
            filled
            :options="mapCategoryOptions"
            hide-hint=""
            hide-bottom-space=""
            dense=""
            class="z-max"
            behavior="menu"
            :label="$t('building')"
          >
            <template v-slot:prepend>
              <q-icon name="apartment" @click.stop />
            </template>
          </q-select>
        </div>
        <div class="col-6 z-top">
          <q-select
            v-model="selectedMap"
            :disable="loading.select"
            :loading="loading.select"
            bottom-slots=""
            filled
            :label="$t('floor')"
            option-label="map_name"
            option-value="map_id"
            :options="mapListOptions"
            hide-hint=""
            hide-bottom-space=""
            dense=""
            class="z-max"
            behavior="menu"
          >
            <template v-slot:prepend>
              <q-icon name="meeting_room" @click.stop />
            </template>
          </q-select>
        </div>
        <div
          id="map"
          :style="
            `width:${$store.state.global.windowSize.width}px;height:${$store
              .state.global.windowSize.height - 90}px`
          "
        ></div>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script>
import mapboxgl from 'mapbox-gl';
export default {
  components: {},
  data() {
    const self = this;
    return {
      loading: { select: true, header: false },
      mapList: {},
      mapPolygons: {},
      mapCategoryOptions: [],
      mapListOptions: [],
      mapCategory: {},
      selectedMapCategory: null,
      selectedMap: null,
      selectMapGID: '',
      showDetails: false,
      currentDetails: {},
      clickedMapPolygon: '',
      markerOffset: [0, 0],
      markerList: [],
      startMarker: null,
      endMarker: null,
      addressUpdate: false,
      address: {},
      loadedInterval: null,
      zoom: 16,
      center: self.$store.state.global.startPointInfo.center.split(','),
      events: {
        init: function() {},
        click: function() {
          //   self.$refs.select.blur();
        },
      },
      amapPosition: false,
      plugin: [],
      dismiss: null,
    };
  },
  computed: {
    locale() {
      return this.$i18n ? this.$i18n.locale : 'zh_CN';
    },
  },
  watch: {
    selectedMapCategory(val) {
      if (typeof this.dismiss === 'function') {
        this.dismiss();
      }
      console.log('watch selectedMapCategoryId', val);
      if (!val || !this.mapCategory[val.value]) return;
      this.loading.select = true;
      this.mapListOptions = this.mapCategory[val.value].mapList;
      if (this.selectedMap) {
        this.selectedMap = this.mapCategory[val.value].mapList[0];
        console.log('selectedMapCategory set selectedMap');
      }

      setTimeout(() => {
        this.loading.select = false;
      }, 500);
    },
    selectedMap(val) {
      console.log('watch selectedMap', val);
      if (val && val.map_id) {
        if (typeof this.dismiss === 'function') {
          this.dismiss();
        }
        this.changeMap(val.map_id);
      }
    },
    showDetails(val) {
      if (val === false) {
        this.selectMapGID = '';
      }
    },
  },
  beforeDestroy() {
    clearInterval(this.loadedInterval);
    if (typeof this.dismiss === 'function') {
      this.dismiss();
    }
    this.$mapUtil.reset();
  },
  mounted() {
    this.$q.loading.hide();
    const map = new mapboxgl.Map({
      container: 'map',
      style: {
        version: 8,
        sources: {},
        sprite: this.ossHost + 'project/113/sprite/sprite',
        glyphs: this.ossHost + 'project/113/{fontstack}/{range}.pbf',
        light: {
          intensity: 0.1,
        },
        layers: [
          // {
          //   id: 'background',
          //   type: 'background',
          //   paint: {
          //     'background-color': 'rgba(0, 0, 0, 0.05)',
          //   },
          // },
        ],
      },
      zoom: 18,
      minZoom: 18,
      maxZoom: 21,
      center: [0, 0],
      pitch: 0,
      attributionControl: false,
      localIdeographFontFamily: '"Noto Sans CJK SC",sans-serif',
    });
    map.addControl(new mapboxgl.NavigationControl());
    this.$mapUtil.setMap(map);
    this.initData();
  },
  methods: {
    clickPolygonListener(e) {
      if (typeof this.dismiss === 'function') {
        this.dismiss();
      }
      if (!e.features || !e.features[0]) return;
      const feature = e.features[0];
      this.currentDetails = {};

      const mapGId = feature.properties.map_gid;
      if (
        !mapGId ||
        !this.mapPolygons[mapGId] ||
        this.clickedMapPolygon.map_gid === mapGId
      ) {
        return;
      }
      this.currentDetails = this.mapPolygons[mapGId];
      this.clickedMapPolygon = mapGId;
      const filter = ['==', 'map_gid', mapGId];
      this.$mapUtil.map.setFilter('polygonSelected', filter);
      if (
        this.$mapUtil.map.getLayoutProperty('polygonSelected', 'visibility') !==
        'visible'
      ) {
        this.$mapUtil.map.setLayoutProperty(
          'polygonSelected',
          'visibility',
          'visible',
        );
      }
      const self = this;
      this.dismiss = this.$q.notify({
        timeout: 0,
        html: true,
        message:
          '<div class="text-black"><span class="text-h6">' +
          this.$t(this.currentDetails.name) +
          '</span><br />' +
          this.currentDetails.map_name +
          '</div>',
        color: 'grey-3',
        classes: 'max-width',
        multiLine: true,
        actions: [
          {
            label: this.$t('mapIt'),
            color: 'black',
            icon: 'near_me',
            handler: () => {
              self.$store.commit('global/SET_STATE_PROPERTY', {
                name: 'routeDest',
                value: self.currentDetails,
              });
              const path =
                      '{{range $i,$e:=.Config.Components}}{{if eq $e.Key "mapRoute"}}{{range $ia,$ie:=$e.Values}}{{range $pi,$pe:=$ie.ProjectFeaturesConfig.Data.Values}}{{if eq $pe.Key "routePath"}}{{$pe.Value}}{{end}}{{end}}{{end}}{{end}}{{end}}';
              self.$router.replace(path);
            },
          },
        ],
      });
      if (this.currentDetails.changedPoint) {
        this.$mapUtil.map.easeTo({
          center: this.currentDetails.changedPoint,
        });
      }
    },
    changeMap(mapId) {
      if (!this.mapList[mapId]) {
        return;
      }
      const self = this;
      this.loading.header = true;
      this.$mapUtil.removeMapPolygonClickEvent(this.clickPolygonListener);
      this.$mapUtil.removeMarker('start');
      const mapDetails = this.mapList[mapId];
      this.$mapUtil.parseMapSetting(mapDetails.map_setting);
      this.$mapUtil.render(mapId);
      this.$mapUtil.addMapPolygonClickEvent(this.clickPolygonListener);
      if (this.$store.state.global.startPointInfo.map_id === mapId) {
        const point = this.$changePoint(
          this.$store.state.global.startPointInfo.point,
        );
        this.$mapUtil.createMarker({
          id: 'start',
          point: point,
          image: require('assets/start.png'),
        });
      }

      this.loadedInterval = setInterval(() => {
        if (this.$mapUtil.map.loaded()) {
          clearInterval(self.loadedInterval);
          if (!this.$mapUtil.map.getLayer('polygonSelected')) {
            this.$mapUtil.map.addLayer({
              id: 'polygonSelected',
              type: 'fill-extrusion',
              source: 'polygon',
              paint: {
                'fill-extrusion-color': '#567ffb',
                'fill-extrusion-height': [
                  'match',
                  ['get', 'type'],
                  'building',
                  ['get', 'height'],
                  'parking',
                  0.05,
                  'walking',
                  0.05,
                  'below',
                  0.05,
                  'aquare',
                  0.05,
                  'road',
                  0.05,
                  'square',
                  0.05,
                  'exit',
                  0.05,
                  'green',
                  0.05,
                  'bottom',
                  0.05,
                  2,
                ],
              },
            });
            this.$mapUtil.map.setLayoutProperty(
              'polygonSelected',
              'visibility',
              'none',
            );
            this.$mapUtil.map.moveLayer('polygonSelected', 'label');
          }
          this.loading.header = false;
        }
      }, 1000);
    },
    initData() {
      const self = this;
      const initFunc = [
        this.loadMapCategory,
        this.loadMapList,
        this.loadPolygons,
      ];
      BP.mapSeries(initFunc, function(f) {
        return f();
      })
        .then(function() {
          let selectMapCategory = null;
          let selectMap = null;
          const hasCall = Object.prototype.hasOwnProperty;
          for (const i in self.mapList) {
            if (hasCall.call(self.mapList, i)) {
              const m = self.mapCategory[self.mapList[i].map_category_id];
              if (m) {
                m.mapList.push(self.mapList[i]);
              }
              if (
                self.$store.state.global.startPointInfo.map_id ==
                self.mapList[i].map_id
              ) {
                selectMapCategory = {
                  label: m.map_category_name,
                  value: m.map_category_id,
                };
                if (self.mapList[self.mapList[i].map_id]) {
                  selectMap = self.mapList[self.mapList[i].map_id];
                }
              }
            }
          }
          self.selectedMapCategory = selectMapCategory;
          return Promise.resolve(selectMap);
        })
        .then(function(selectMap) {
          setTimeout(() => {
            self.selectedMap = selectMap;
          }, 500);
        })
        .finally(function(selectMap) {
          setTimeout(() => {
            self.loading.select = false;
          }, 500);
        });
    },
    loadMapCategory() {
      const self = this;
      return new Promise(function(resolve) {
        const readStore = self.$store.state.global.indexedDB
          .transaction('mapCategory')
          .objectStore('mapCategory')
          .getAll();
        readStore.onsuccess = function(e) {
          const r = e.target.result;
          if (r && r.length > 0) {
            for (let i = 0; i < r.length; i++) {
              self.mapCategory[r[i].map_category_id] = r[i];
              self.mapCategory[r[i].map_category_id].mapList = [];
              self.mapCategoryOptions.push({
                label: self.$t(r[i].map_category_name),
                value: r[i].map_category_id,
              });
            }
          }
          resolve();
        };
      });
    },
    loadMapList() {
      const self = this;
      return new Promise(function(resolve) {
        const readStore = self.$store.state.global.indexedDB
          .transaction('mapList')
          .objectStore('mapList')
          .getAll();
        readStore.onsuccess = function(e) {
          const r = e.target.result;
          if (r && r.length > 0) {
            for (let i = 0; i < r.length; i++) {
              self.mapList[r[i].map_id] = r[i];
            }
          }
          resolve();
        };
      });
    },
    loadPolygons() {
      const self = this;
      return new Promise(function(resolve) {
        const readStore = self.$store.state.global.indexedDB
          .transaction('mapPolygons')
          .objectStore('mapPolygons')
          .getAll();
        readStore.onsuccess = function(e) {
          const r = e.target.result;
          if (r && r.length > 0) {
            for (let i = 0; i < r.length; i++) {
              self.mapPolygons[r[i].map_gid] = r[i];
            }
          }
          resolve();
        };
      });
    },
    route(item) {
      this.$store.commit('updateCurrentRoute', item);
      this.$router.replace('/route');
    },
  },
};
</script>
<style>
@import url('https://api.tiles.mapbox.com/mapbox-gl-js/v1.10.1/mapbox-gl.css');
</style>
<style>
.max-width {
  width: 100% !important;
}
</style>
