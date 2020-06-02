<template>
  <q-btn icon="language" :label="localeName" dense="" unelevated>
    <q-menu anchor="bottom left" self="top left" auto-close="">
      <template v-for="(item, index) in $i18n.messages">
        <q-item :key="index" clickable="" @click="changeLanguage(index)">
          <q-item-section>{{ $t(index) }}</q-item-section> </q-item
        ><q-separator :key="`${index}-separator`"
      /></template>
    </q-menu>
  </q-btn>
</template>
<script>
import { QSpinnerGears } from 'quasar';
export default {
  computed: {
    localeName() {
      const name = this.$i18n ? this.$i18n.locale : 'zh_CN';
      return this.$t(name);
    },
  },
  methods: {
    changeLanguage(lang) {
      console.log(lang, this.$i18n.locale);
      if (lang === this.$i18n.locale) {
        return;
      }
      this.$q.loading.show({
        QSpinnerGears,
        spinnerColor: 'info',
        messageColor: 'white',
        spinnerSize: 'md',
        // message: 'Updated message',
      });
      this.$storage.set('locale', lang);
      this.$i18n.locale = lang;
      setTimeout(() => {
        this.$q.loading.hide();
      }, 1000);
    },
  },
};
</script>
