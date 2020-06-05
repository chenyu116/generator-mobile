<template>
  <q-layout view="lHh Lpr lFf" class="bg-white">
    <q-page-container>
      <HeaderWithBack title="search" icon="search"></HeaderWithBack>
      <q-page class="q-pa-md">
        <div class="row">
          <div class="col-1 q-mr-md">
            <q-btn flat="" round="" @click="speechMode = !speechMode"
              ><q-icon v-if="!speechMode" name="record_voice_over"/><q-icon
                v-if="speechMode"
                name="voice_over_off"
            /></q-btn>
          </div>
          <div class="col-7">
            <q-input
              v-model="keywords"
              :loading="loading.searching"
              :disable="loading.searching"
              filled
              :label="$t('keywords')"
              dense=""
            />
          </div>
          <div class="q-ml-md">
            <q-btn @click="search()">{{print "{{ $t('search') }}"}}</q-btn>
          </div>
        </div>
        <div v-if="!speechMode">
          <!--<div v-if="list.length === 0 && polygons.length > 0">
            <q-card flat="">
              <q-card-section class="q-px-none">推荐</q-card-section>
              <q-chip v-for="(item, index) in polygons" :key="`chip-${index}`">
                {{print "{{ $t(item.point_name) }}"}}</q-chip
              >
            </q-card>
          </div>-->
          <div v-if="list.length > 0" class="q-mt-md">
            <q-list separator>
              <q-item-label header class="q-pt-none q-px-none"
                >{{print "{{ $t('searchResult') }}"}}</q-item-label
              >
              <q-item
                v-for="(item, index) in list"
                :key="index"
                class="q-px-none"
              >
                <q-item-section>{{print "{{ item.name }}"}}</q-item-section>
                <q-item-section side>
                  <div class="text-grey-8 q-gutter-md">
                    <q-btn
                      size="12px"
                      flat
                      dense
                      icon="view_list"
                      :label="$t('details')"
                      stack=""
                    />
                    <q-btn
                      size="12px"
                      flat
                      dense
                      icon="navigation"
                      :label="$t('mapIt')"
                      stack=""
                    />
                  </div>
                </q-item-section>
              </q-item>
            </q-list>
          </div>
        </div>
        <div v-if="speechMode">
          <q-card
            flat=""
            :style="`height:${$store.state.global.windowSize.height - 320}px`"
          >
            <q-card-section>{{print "{{ $t('speech') }}"}}</q-card-section>
            <q-card-section class="absolute-center row text-center">
              <q-icon name="settings_voice" size="md" />
            </q-card-section>
          </q-card>
          <q-card flat="" class="text-center">
            <q-card-section>
              <q-btn
                icon="fingerprint"
                round=""
                size="40px"
                dense=""
                color="pink"
                class="shadow-4"
                @touchstart="startRecord"
                @touchend="stopRecord"
              />
            </q-card-section>
            <q-card-section>{{print "{{ $t('pressToTalk') }}"}}</q-card-section>
          </q-card>
        </div>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script>
export default {
  components: {},
  data() {
    return {
      loading: {
        searching: false,
      },
      errMsg: '',
      keywords: '',
      list: [],
      polygons: [],
      currentPosition: '',
      listHeight: window.innerHeight - 200,
      canUseSpeech: false,
      speechMode: false,
      startRecordTime: 0,
      endRecordTime: 0,
      recordTime: 0,
      isRecording: false,
      recordTimeout: null,
      speechRowStyle: '',
      speechLocalId: '',
    };
  },

  computed: {
    windowHeight() {
      return this.$store.state.windowHeight;
    },
  },
  watch: {
    windowHeight: {
      handler(val) {
        this.speechRowStyle = 'height:' + (val - 270) + 'px;';
      },
      immediate: true,
    },
    speechLocalId(val) {
      if (!val) return;
      const self = this;
      this.$wechat.wx.translateVoice({
        localId: val, // 需要识别的音频的本地Id，由录音相关接口获得
        isShowProgressTips: 1, // 默认为1，显示进度提示
        success: function(res) {
          self.keywords = res.translateResult
            ? res.translateResult.replace('。', '')
            : ''; // 语音识别的结果
          self.speechLocalId = '';
          self.search();
        },
      });
    },
  },
  beforeCreate() {
    this.$q.loading.hide();
  },
  mounted() {
    const self = this;
    if (this.$store.state.global.isWx) {
      this.$wechat.wx.ready(function() {
        self.canUseSpeech = true;
      });
    }

    this.loadPolygons();
    if (this.$store.state.global.searchKeywords) {
      this.keywords = this.$store.state.global.searchKeywords;
      this.search();
    }
  },
  beforeDestroy() {
    if (this.$route.name !== 'route') {
      this.$store.commit('updateSearchKeywords', null);
    }
  },
  methods: {
    search() {
      const self = this;
      this.keywords = this.keywords ? this.keywords.trim() : '';
      if (this.keywords === '') {
        this.errMsg = this.$t('needKeywords');
        return;
      }
      //   this.$http.put(
      //     this.apiHost + '/project/spm',
      //     {
      //       type: 'search',
      //       content: JSON.stringify({
      //         keywords: self.keywords,
      //         locale: self.$i18n.locale,
      //         agent: self.$store.state.navigator.userAgent,
      //       }),
      //       contentType: 'json',
      //       tag: self.$store.state.userId,
      //       projectID: self.$store.state.startPointInfo.project_id,
      //       timestamp: parseInt(new Date().getTime() / 1000),
      //     },
      //     {
      //       emulateJSON: false,
      //     },
      //   );
      //   this.$store.commit('updateSearchKeywords', this.keywords);
      this.errMsg = '';
      this.loading.searching = true;
      if (window.AMap) {
        const searchLang = this.$i18n.locale === 'zh_CN' ? 'zh_CN' : 'en';
        const placeSearch = new window.AMap.PlaceSearch({
          pageSize: 20,
          lang: searchLang,
        });
        let currentPoint = this.changePoint(
          this.$store.state.startPointInfo.point,
        );
        try {
          window.AMap.plugin('AMap.Geolocation', function() {
            const geolocation = new window.AMap.Geolocation();
            geolocation.getCurrentPosition(function(status, result) {
              if (result.position) {
                currentPoint = [result.position.lng, result.position.lat];
              }
              placeSearch.searchNearBy(
                self.keywords,
                currentPoint,
                5000,
                function(status, result) {
                  if (
                    result.info === 'OK' &&
                    result.poiList &&
                    result.poiList.pois.length > 0
                  ) {
                    self.list = result.poiList.pois;
                    self.speechMode = false;
                  } else {
                    self.errMsg = self.$t('noSearchResult');
                  }
                  self.loading.searching = false;
                },
              );
            });
          });
        } catch (e) {
          self.errMsg = self.$t('anErrorOccurred');
          self.loading.searching = false;
        }
      } else {
        this.searchRemote();
      }
    },
    loadPolygons() {
      const self = this;
      const readStore = self.$store.state.global.indexedDB
        .transaction('mapPolygons')
        .objectStore('mapPolygons')
        .getAll();
      readStore.onsuccess = function(e) {
        const r = e.target.result;
        if (r && r.length > 0) {
          self.polygons = r;
          console.log('self.polygons', self.polygons);
        }
      };
    },
    startRecord(e) {
      const self = this;
      e.preventDefault();
      console.log('startRecord');
      this.startRecordTime = new Date().getTime();
      this.recordTimeout = setTimeout(function() {
        console.log('startRecord');
        clearTimeout(self.recordTimeout);
        self.isRecording = true;
        self.$wechat.wx.startRecord();
      }, 200);
    },
    stopRecord() {
      const self = this;
      this.isRecording = false;
      this.recordTime = new Date().getTime() - this.startRecordTime;
      console.log(this.recordTime);
      if (this.recordTime <= 200) {
        clearTimeout(this.recordTimeout);
        return;
      }
      this.$wechat.wx.stopRecord({
        success: function(res) {
          self.speechLocalId = res.localId;
        },
        fail: function() {},
      });
    },
    searchRemote() {
      const self = this;
      return new Promise(function(resolve, reject) {
        self.$http
          .post(
            'https://apis.signp.cn/search',
            {
              keywords: self.keywords,
              type: '',
              projectID: self.$store.state.global.startPointInfo.project_id,
              timestamp: parseInt(new Date().getTime() / 1000),
              access_token: 'test',
            },
            { emulateJSON: true },
          )
          .then(function(resp) {
            console.log(resp);
            // if (resp.status === 200 && resp.body) {
            //   const writeStore = params.store.state.global.indexedDB
            //     .transaction('mapCategory', 'readwrite')
            //     .objectStore('mapCategory');
            //   const hasCall = Object.prototype.hasOwnProperty;
            //   for (const i in resp.body) {
            //     if (hasCall.call(resp.body, i)) {
            //       writeStore.put(resp.body[i]);
            //     }
            //   }
            //   resolve();
            // } else {
            //   reject(new Error('读取地图楼层分类失败'));
            // }
          })
          .catch(function() {
            // reject(new Error('读取地图楼层分类失败'));
          });
      });
    },
    route(item) {
      const details = item;
      details.lnglat = [item.location.lng, item.location.lat];
      this.$store.commit('updateCurrentRoute', details);
      this.$router.replace('/route');
    },
  },
};
</script>
