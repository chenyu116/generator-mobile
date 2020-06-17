<template>
  <q-layout view="lHh Lpr lFf" class="bg-white">
    <q-page-container class="bg">
      <q-page>
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
            <q-space />
            <i18n v-if="hasI18n" />
          </q-toolbar>
        </q-header>

        <div class="row featureRow">
          <div
            class="roundBorder"
            :style="`width:100%;background: rgba(0, 0, 0, 0.4)`"
          >
            <div
              class="roundBorder1"
              style="width: 100%; border-top: 2px solid #ccc;"
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
                        v-if="item.image"
                        :key="`tab-0-image-${itemIndex}`"
                        :src="item.image"
                        style="width: 60px; height: 40px;"
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
                  :key="`separator-1-${itemIndex}`"
                  vertical
                  color="grey"
                />

                <q-tab
                  :key="`tab-0-${itemIndex}`"
                  class="q-px-xs"
                  :style="`width:${tabWidth}px;`"
                  @click="featureClick(item)"
                  ><q-img
                    v-if="item.image"
                    :key="`tab-1-image-${itemIndex}`"
                    :src="item.image"
                    style="width: 60px; height: 40px;"
                  />
                  {{print "{{ $t(item.name) }}"}}</q-tab
                >
              </template>
            </q-tabs>
          </div>
        </div>
        {{range $i,$e:=.Config.Components}}{{if eq $e.Key "blocks"}}{{range $ni,$ne:=$e.Values}} <{{ $ne.ComponentHash }} />{{ end}}{{ end }}{{ end }}
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
          name: '{{range $pi,$ne:=$ie.ProjectFeaturesConfig.Data.Values}}{{if eq $ne.Key "title"}}{{$ne.Value}}{{end}}{{end}}',{{range $pi,$ne:=$ie.ProjectFeaturesConfig.Data.Values}}{{if eq $ne.Key "iconImage"}}{{if $ne.Value}}
          image: require('../{{$ie.ProjectFeaturesInstallName}}/{{$ne.Value}}'),{{end}}{{end}}{{if eq $ne.Key "icon"}}{{if $ne.Value}}
          icon: '{{$ne.Value}}',{{end}}{{end}}{{end}}
          onClick: function() {
            self.$router.replace('{{range $pi,$ne:=$ie.ProjectFeaturesConfig.Data.Values}}{{if eq $ne.Key "routePath"}}{{$ne.Value}}{{end}}{{end}}');
          },
        },{{end}}{{end}}{{end}}
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
    featureClick(item) {
      if (typeof item.onClick === 'function') {
        item.onClick();
      }
    },
  },
};
</script>
<style scoped>
.bg {
  width: 100%;
  height: 38em;
  top: 0;{{if .DataValues.homeBg}}
  background: url('./{{.DataValues.homeBg}}') no-repeat;{{end}}
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
