<template>
  <q-card flat="" square="">
    <q-card-section>
      <div class="text-h5">{{print "{{ $t(title) }}"}}</div>
      <div class="text-subtitle2 text-grey">
        - {{print "{{ $t(subtitle) }}"}}
      </div>
    </q-card-section>
    <q-card-section class="q-pt-none">
      {{print "{{ weatherString }}"}}
    </q-card-section>
  </q-card>
</template>

<script>
// {{range $i,$e:=.Config.Components}}{{if eq $e.Key "blocks"}}{{range $ni,$ne:=$e.Values}}
// import |{{$ne.ComponentHash}}| from '../components/{{$ne.ProjectFeaturesInstallName}}/Index.vue';{{end}}{{end}}{{end}}
export default {
  name: 'Index',
  // components: { {{range $i,$e:=.Config.Components}}{{if eq $e.Key "blocks"}}{{range $ni,$ne:=$e.Values}}|{{$ne.ComponentHash}}|,{{end}}{{end}}{{end}}Swiper, SwiperSlide },
  data() {
    return {
      title: '{{.DataValues.title}}',
      subtitle: '{{.DataValues.subtitle}}',
      weatherString: '',
    };
  },
  computed: {
    locale() {
      return this.$i18n ? this.$i18n.locale : 'zh_CN';
    },
  },
  mounted() {
    this.loadWeather();
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
              location: '{{.DataValues.location}}',
              lang: this.locale.substring(0, 2),
            },
          })
          .then(function(resp) {
            if (resp.status === 200) {
              const timestamp = nowTime + 1800 * 1000;
              const body = JSON.parse(resp.body);
              self.$storage.set(
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
  },
};
</script>
