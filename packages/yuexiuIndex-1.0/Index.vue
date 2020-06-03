<template>
  <q-layout view="lHh Lpr lFf" class="bg-white">
    <q-header
      style="background: rgba(0, 0, 0, 0.7);"
      :class="`${$store.state.global.config.theme.textColor} shadow-3`"
    >
      <q-toolbar>
        <q-toolbar-title class="text-subtitle1">
          <q-icon
            name="location_on"
            size="sm"
            left=""
          />{{print "{{$t($store.state.global.startPointInfo.name)}}"}}</q-toolbar-title
        >
      </q-toolbar>
      <i18n v-if="hasI18n" />
    </q-header>
    <q-page-container class="bg">
      <div class="row featureRow">
        <div
          class="roundBorder"
          :style="`width:100%;background: rgba(0, 0, 0, 0.4)`"
        >
          <div
            class="roundBorder1"
            style="width:100%;border-top: 2px solid #ccc;"
          >
            <div :style="`width: ${$store.state.global.windowSize.width}px`">
              <q-tabs
                dense
                class="text-grey-1 transparent"
                align="justify"
                :breakpoint="0"
                no-caps=""
              >
                <template v-for="(item, itemIndex) in splitFeatures[0]">
                  <q-separator
                    v-if="itemIndex > 0"
                    :key="`separator-0-${itemIndex}`"
                    vertical
                    color="grey"
                  />

                  <q-tab
                    :key="`tab-0-${itemIndex}`"
                    class="q-px-xs"
                    :style="`width:${tabWidth}px;`"
                    @click="featureClick(item)"
                    ><q-img
                      :key="`tab-0-image-${itemIndex}`"
                      :src="item.image"
                      style="width:60px;height:40px"
                    />{{print "{{ $t(item.name) }}"}}</q-tab
                  >
                </template>
              </q-tabs>
            </div>
          </div>
        </div>
        <div
          v-if="typeof splitFeatures[1] !== 'undefined'"
          class="transparentBlack"
          :style="`width: ${$store.state.global.windowSize.width}px`"
        >
          <q-separator color="grey" />
          <q-tabs
            dense
            class="text-grey-1 transparent"
            align="justify"
            :breakpoint="0"
            no-caps=""
          >
            <template v-for="(item, itemIndex) in splitFeatures[1]">
              <q-separator
                v-if="itemIndex > 0"
                :key="`separator-0-${itemIndex}`"
                vertical
                color="grey"
              />

              <q-tab
                :key="`tab-0-${itemIndex}`"
                class="q-px-xs"
                :style="`width:${tabWidth}px;`"
                @click="featureClick(item)"
                ><q-img
                  :key="`tab-0-image-${itemIndex}`"
                  :src="item.image"
                  style="width:60px;height:40px"
                />
                {{print "{{ $t(item.name) }}"}}</q-tab
              >
            </template>
          </q-tabs>
        </div>
      </div>
      <q-page>
        {{range $i,$e:=.Config.Components}}{{if eq $e.Key "blocks"}}{{range $ni,$ne:=$e.Values}}
        <{{$ne.ComponentHash}} />{{end}}{{end}}{{end}}
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script>
{{range $i,$e:=.Config.Components}}{{if eq $e.Key "blocks"}}{{range $ni,$ne:=$e.Values}}
import |{{$ne.ComponentHash}}| from '../{{$ne.ProjectFeaturesInstallName}}/Index.vue';{{end}}{{end}}{{end}}

export default {
  name: 'Index',
  components: { {{range $i,$e:=.Config.Components}}{{if eq $e.Key "blocks"}}{{range $ni,$ne:=$e.Values}}|{{$ne.ComponentHash}}|,{{end}}{{end}}{{end}} },
  data() {
    const self = this;
    return {
      splitFeatures: [],
      features: [
        {{range $i,$e:=.Config.Components}}{{if eq $e.Key "nav"}}{{range $ni,$ie:=$e.Values}}
        {
          name: '{{range $pi,$ne:=$ie.ProjectFeaturesConfig.Data.Values}}{{if eq $ne.Key "title"}}{{$ne.Value}}{{end}}{{end}}',
          image: require('../{{$ie.ProjectFeaturesInstallName}}/{{range $pi,$ne:=$ie.ProjectFeaturesConfig.Data.Values}}{{if eq $ne.Key "iconImage"}}{{$ne.Value}}{{end}}{{end}}'),
          onClick: function() {
            self.$router.replace('{{range $pi,$ne:=$ie.ProjectFeaturesConfig.Data.Values}}{{if eq $ne.Key "routePath"}}{{$ne.Value}}{{end}}{{end}}');
          },
        },
        {{end}}{{end}}{{end}}
      ],
      countsPerRow: 4,
      hasI18n: Object.keys(self.$i18n.messages).length > 1,
    };
  },
  computed: {
    locale() {
      return this.$i18n ? this.$i18n.locale : 'zh_CN';
    },
    tabWidth() {
      return parseInt(
        this.$store.state.global.windowSize.width / this.countsPerRow,
      );
    },
  },
  mounted() {
    const splitFeatures = [];
    let splitCounts = 0;
    for (let i = 0; i < this.features.length; i++) {
      if (!splitFeatures[splitCounts]) splitFeatures.push([]);
      splitFeatures[splitCounts].push(this.features[i]);
      if (this.features[i + 1] && (i + 1) % this.countsPerRow === 0) {
        splitCounts++;
      }
    }
    this.splitFeatures = splitFeatures;
  },
  methods: {
    loadWeather() {
      const self = this;
      const weatherTag = 'weather-' + this.locale;
      const weather = JSON.parse(this.$storage.get(weatherTag));
      const nowTime = new Date().getTime();
      if (!weather || weather.timestamp < nowTime) {
        this.$http
          .get(this.apiHost + '/weather', {
            params: {
              location: 'haizhu',
              lang: this.locale.substring(0, 2),
            },
          })
          .then(function(resp) {
            if (resp.status === 200) {
              const timestamp = nowTime + 1800 * 1000;
              const body = JSON.parse(resp.body);
              localStorage.setItem(
                weatherTag,
                JSON.stringify({
                  timestamp: timestamp,
                  data: body,
                }),
              );
              self.parseWeather(body);
            }
          });
      } else {
        self.parseWeather(weather.data);
      }
    },
    parseWeather(w) {
      if (w.HeWeather6 && w.HeWeather6.length > 0 && w.HeWeather6[0].now) {
        const now = w.HeWeather6[0].now;
        const wa = [];
        wa.push(now.cond_txt);
        wa.push(now.tmp + '℃');
        wa.push(this.$t('relativeHumidity') + ' ' + now.hum + '%');
        wa.push(now.wind_dir + ' ' + now.wind_sc + this.$t('level'));
        this.weatherString = wa.join(' ');
      }
    },
    loadBanner() {
      const self = this;
      self.articles = [];
      const readStore = self.$store.state.global.indexedDB
        .transaction('banner')
        .objectStore('banner')
        .get(self.bannerType);
      readStore.onsuccess = function(e) {
        const banner = [];
        const r = e.target.result;
        const nowTime = new Date().getTime();
        if (r && r.timestamp > nowTime && r.val.length > 0) {
          for (let i = 0; i < r.val.length; i++) {
            banner.push({
              img: 'https://o.signp.cn/' + r.val[i].images_path,
              title: '',
              bindType: r.val[i].images_bind_type,
              link: r.val[i].images_link,
              mapGid: r.val[i].images_map_gid,
            });
          }
          self.banner = banner;
        } else {
          self.$http
            .get(self.apiHost + '/images', {
              params: {
                type: self.bannerType,
                projectID: self.$store.state.global.startPointInfo.project_id,
                timestamp: parseInt(nowTime / 1000),
              },
            })
            .then(function(resp) {
              if (resp.status === 200) {
                // const writeStore = self.$store.state.db
                //   .transaction('banner', 'readwrite')
                //   .objectStore('banner');
                // writeStore.put({
                //   image_type: self.bannerType,
                //   timestamp: nowTime + 300 * 1000,
                //   val: resp.body,
                // });
                for (let i = 0; i < resp.body.length; i++) {
                  banner.push({
                    img: 'https://o.signp.cn/' + resp.body[i].images_path,
                    title: '',
                    bindType: resp.body[i].images_bind_type,
                    link: resp.body[i].images_link,
                    mapGid: resp.body[i].images_map_gid,
                  });
                }
                self.banner = banner;
              }
            });
        }
      };
    },
    featureClick(item) {
      if (typeof item.onClick === 'function') {
        item.onClick();
      }
    },
    loadArticles() {
      const self = this;
      self.articles = [];
      const readStore = self.$store.state.global.indexedDB
        .transaction('articles')
        .objectStore('articles')
        .get('122_' + self.locale);
      readStore.onsuccess = function(e) {
        const r = e.target.result;
        const nowTime = new Date().getTime();

        const __datalimit__ = 3;
        if (r && r.timestamp > nowTime && r.val.length > 0) {
          for (let i = 0; i < r.val.length; i++) {
            if (i === __datalimit__) break;
            const _d = new Date(parseInt(r.val[i].article_create_at) * 1000);
            r.val[i].date = _d.getMonth() + 1 + '/' + _d.getDate();
            self.articles.push(r.val[i]);
          }
        } else {
          self.$http
            .get(self.apiHost + '/articles', {
              params: {
                categoryID: 122,
                locale: self.$locale,
                limit: __datalimit__,
                projectID: self.$store.state.global.startPointInfo.project_id,
                timestamp: parseInt(new Date().getTime() / 1000),
              },
            })
            .then(function(resp) {
              if (resp.status === 200) {
                const timestamp = new Date().getTime() + 300 * 1000;
                const writeStore = self.$store.state.global.indexedDB
                  .transaction('articles', 'readwrite')
                  .objectStore('articles');
                writeStore.put({
                  article_category_id: '122_' + self.locale,
                  timestamp: timestamp,
                  val: resp.body,
                });
                for (let i = 0; i < resp.body.length; i++) {
                  if (i === __datalimit__) {
                    break;
                  }
                  const _d = new Date(
                    parseInt(resp.body[i].article_create_at) * 1000,
                  );
                  resp.body[i].date = _d.getMonth() + 1 + '/' + _d.getDate();
                  self.articles.push(resp.body[i]);
                }
              }
            });
        }
      };
    },
  },
};
</script>
<style scoped>
.bg {
  width: 100%;
  height: 38em;
  top: 0;
  {{if .DataValues.homeBg}}
  background: url('../{{.InstallDir}}/{{.DataValues.homeBg}}') no-repeat;{{end}}
  background-size: cover;
}
.transparentBlack {
  background-color: rgba(0, 0, 0, 0.4);
}
.featureRow {
  margin-top: 20em;
}
.roundBorder {
  border-radius: 50% 50% 0 0 !important;
}
.roundBorder1 {
  margin-top: 10px;
  padding-top: 20px;
  border-radius: 50% 50% 0 0 !important;
  background-color: rgba(0, 0, 0, 0);
}
</style>
