<template>
  <q-layout view="lHh Lpr lFf" class="bg-white">
    <HeaderWithBack
      title="{{.DataValues.title}}"
      icon="{{.DataValues.icon}}"
      :loading="loading.header"
      class="z-max"
    ></HeaderWithBack>
    <q-page-container>
      <q-page class="row">
        <div
          id="map"
          :style="
            `width:${$store.state.global.windowSize.width}px;height:${$store
              .state.global.windowSize.height - 188}px;`
          "
        ></div>
        <q-footer
          v-if="!loading.header"
          style="height:138px"
          :class="`bg-grey-2 text-grey-7 shadow-up-1`"
        >
          <q-list separator>
            <q-item class="q-px-md">
              <q-item-section
                ><q-carousel
                  v-model="slide"
                  transition-prev="scale"
                  transition-next="scale"
                  animated
                  padding=""
                  :arrows="!simulate.simulating"
                  height="60px"
                  control-color="grey-7"
                  class="bg-grey-2"
                >
                  <q-carousel-slide
                    v-for="(item, index) in stepText"
                    :key="`step-${index}`"
                    :name="index"
                    class="column flex-center"
                  >
                    <div class="text-center">
                      {{print "{{ item }}"}}
                    </div>
                  </q-carousel-slide>
                </q-carousel></q-item-section
              >
              <q-item-section side>
                <q-btn
                  :class="`text-grey-7`"
                  size="12px"
                  :icon="simulate.simulating ? 'stop' : 'directions_walk'"
                  :label="simulate.simulating ? $t('cancel') : $t('simulate')"
                  stack=""
                  @click="simulate.simulating ? cancelSimulate() : goSimulate()"
                />
              </q-item-section>
            </q-item>
            <q-item class="q-px-md">
              <q-item-section
                ><q-item-label
                  v-if="routeDest.name"
                  class="text-h6 text-black"
                  >{{print "{{ $t(routeDest.name) }}"}}</q-item-label
                ><q-item-label
                  v-if="routeDest.map_name"
                  class="ellipsis"
                  caption=""
                  >{{print "{{ routeDest.map_name }}"}}</q-item-label
                ></q-item-section
              >
              <q-item-section side top>
                <q-icon name="timer" />
                {{print "{{ parseTime($mapUtil.routing.cost) }}"}}
              </q-item-section>
              <q-item-section side top>
                <q-icon name="timeline" />
                {{print "{{ parseDistance($mapUtil.routing.cost) }}"}}
              </q-item-section>
            </q-item>
          </q-list>
        </q-footer>
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
      slide: '',
      loading: { select: true, header: false },
      mapList: {},
      mapPolygons: {},
      mapCategoryOptions: [],
      mapListOptions: [],
      mapCategory: {},
      showDetails: false,
      routeDest: self.$store.state.global.routeDest || {},
      clickedMapPolygon: '',
      markerOffset: [0, 0],
      stepText: [],
      startMarker: null,
      endMarker: null,
      simulate: {
        simulating: false,
        step: 0,
      },
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
      currentMapId: null,
    };
  },
  watch: {
    slide(val) {
      console.log('watch slide', val);
      if (
        !this.simulate.simulating &&
        this.$mapUtil.routing.routingPathArr[val]
      ) {
        const self = this;
        this.changeMap(
          this.$mapUtil.routing.routingPathArr[val][0].map_id,
        ).then(function() {
          self.$mapUtil.renderRoute(val);
        });
      }
    },
    showDetails(val) {
      if (val === false) {
        this.selectMapGID = '';
      }
    },
  },
  beforeDestroy() {
    if (typeof this.dismiss === 'function') {
      this.dismiss();
    }
    this.$store.commit('global/SET_STATE_PROPERTY', {
      name: 'routeDest',
      value: null,
    });
    this.$mapUtil.reset();
  },
  mounted() {
    this.$q.loading.hide();
    if (!this.$store.state.global.routeDest && !this.$route.params.categoryId) {
      return this.$router.replace('/');
    }
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
        layers: [],
      },
      zoom: 18,
      minZoom: 18,
      maxZoom: 21,
      center: [0, 0],
      pitch: 0,
      attributionControl: false,
      localIdeographFontFamily: '"Noto Sans CJK SC",sans-serif',
    });
    this.$mapUtil.setMap(map);
    this.initData();
  },
  methods: {
    cancelSimulate() {
      this.simulate.simulating = false;
      this.$mapUtil.cancelSimulate();
    },
    goSimulate() {
      const self = this;
      this.simulate.simulating = true;
      if (
        this.$mapUtil.routing.routingPathArr &&
        typeof this.$mapUtil.routing.routingPathArr.forEach === 'function'
      ) {
        console.log('simulate');
        let simulatingInveral;
        BP.mapSeries(self.$mapUtil.routing.routingPathArr, function(
          arr,
          index,
        ) {
          return new Promise(function(resolve, reject) {
            if (!self.simulate.simulating) {
              self.$mapUtil.cancelSimulate();
              return;
            }
            console.log(arr[0].map_id, self.currentMapId);
            if (arr[0].map_id !== self.currentMapId) {
              self.changeMap(arr[0].map_id);
            }
            self.slide = index;
            self.$mapUtil.renderRoute(index);
            simulatingInveral = setInterval(() => {
              if (!self.$mapUtil.routing.simulating) {
                clearInterval(simulatingInveral);
                resolve();
              }
            }, 2000);
            self.$mapUtil.simulate(arr);
          });
        })
          .then(function() {
            self.cancelSimulate();
          })
          .catch(function(err) {});
      }
    },
    loadCategoryRouteRemote() {
      const self = this;
      return new Promise(function(resolve, reject) {
        self.$http
          .post(
            'https://apis.signp.cn/category/navi',
            {
              projectID: self.$store.state.global.startPointInfo.project_id,
              pointID: self.$store.state.global.startPointInfo.id,
              categoryID: self.$route.params.categoryId,
              refresh: 1,
            },
            { emulateJSON: true },
          )
          .then(function(resp) {
            console.log(resp);
            if (
              resp.status === 200 &&
              resp.body &&
              resp.body.code === 0 &&
              resp.body.data.length > 0
            ) {
              const routeDest = resp.body.data[resp.body.data.length - 1];
              console.log('routeDest', routeDest);
              routeDest.name = routeDest.point_name;
              self.$store.commit('global/SET_STATE_PROPERTY', {
                name: 'routeDest',
                value: routeDest,
              });
              self.routeDest = routeDest;
              resolve(resp.body.data);
            } else {
              reject(new Error('读取路径信息失败'));
            }
          })
          .catch(function() {
            reject(new Error('读取路径信息失败'));
          });
      });
    },
    loadRouteRemote() {
      const self = this;
      return new Promise(function(resolve, reject) {
        self.$http
          .get(self.apiHost + '/route', {
            params: {
              projectID: self.$store.state.global.startPointInfo.project_id,
              dest: self.routeDest.map_gid,
              startPointID: self.$store.state.global.startPointInfo.id,
              timestamp: parseInt(new Date().getTime() / 1000),
            },
          })
          .then(function(resp) {
            if (resp.status === 200) {
              resolve(resp.body);
            } else {
              reject(new Error('读取路径信息失败'));
            }
          });
      });
    },
    changeMap(mapId) {
      const self = this;
      return new Promise(function(resolve) {
        if (!self.mapList[mapId]) {
          return resolve();
        }
        const mapDetails = self.mapList[mapId];
        self.$mapUtil.removeMapPolygonClickEvent(self.clickPolygonListener);
        self.$mapUtil.clearMarker();
        self.$mapUtil.parseMapSetting(mapDetails.map_setting);
        self.$mapUtil.render(mapId);
        BP.delay(1000).then(resolve);
      });
    },
    initData() {
      const self = this;
      this.loading.header = true;
      this.loadMapList()
        .then(function() {
          self.currentMapId = self.$store.state.global.startPointInfo.map_id;
          if (
            self.$route.params.categoryId &&
            self.$route.params.categoryId > 0
          ) {
            return self.loadCategoryRouteRemote();
          }
          return self.loadRouteRemote();
        })
        .then(function(routeList) {
          console.log('routeList', routeList);
          self.$mapUtil.genRoutingPath(routeList);
          self.stepText = self.$mapUtil.routing.stepText;
          self.slide = 0;
          console.log(self.stepText);
          return BP.delay(1000);
        })
        .then(function() {
          self.loading.header = false;
        })
        .catch(function() {
          self.dismiss = self.$q.notify({
            message: '数据加载失败',
            timeout: 0,
            type: 'negative',
            actions: [
              {
                label: self.$t('retry'),
                color: 'yellow',
                handler: () => {
                  self.initData();
                },
              },
            ],
          });
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
    parseTime(time) {
      const hour = parseInt(time / 3600);
      let minutes = parseInt((time - hour * 3600) / 60);
      if (minutes < 1) minutes = 1;
      return (
        (hour ? hour + ' ' + this.$t('hour') : '') +
        (minutes ? minutes + ' ' + this.$t('minutes') : '')
      );
    },
    parseDistance(d) {
      const kilometer = d / 1000;
      if (kilometer < 1) {
        return d + ' ' + this.$t('meter');
      }
      return kilometer.toFixed(1) + ' ' + this.$t('kilometer');
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
