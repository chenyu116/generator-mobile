<template>
  <q-layout view="lHh Lpr lFf" class="bg-white">
    <HeaderWithBack title="{{.DataValues.title}}" icon="{{.DataValues.icon}}" />
    <q-page-container>
      <q-page>
        <template v-for="(item, index) in list">
          <q-separator v-if="index > 0" :key="`separator-${index}`" inset="" />
          <q-item
            :key="`item-${index}`"
            v-ripple
            v-viewer
            @click.native="href(item.href)"
          >
            <q-item-section>{{print "{{ item.text }}"}}</q-item-section>
            <q-space />
            <q-item-section thumbnail="" style="padding-right:16px">
              <img v-if="item.src" :src="item.src" class="item-qrcode" />
              <q-item-label v-if="item.sub"
                >{{print "{{ item.sub }}"}}</q-item-label
              >
            </q-item-section>
          </q-item>
        </template>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script>
export default {
  components: {},
  data() {
    return {
      list: [
        {
          text: '广州越秀发布',
          icon: 'fa-comments',
          color: 'indigo',
          href: '',
          sub: '',
          src: 'https://o.signp.cn/project/145/yuexiufabu.png',
        },
        {
          text: 'wechatNumber',
          icon: 'fa-comments',
          color: 'indigo',
          href: '',
          sub: '广州越秀文体',
          src: 'https://o.signp.cn/project/145/wechatQrcode.png',
        },
        {
          text: 'travelAdvisory',
          icon: 'fa-question-circle',
          color: 'indigo',
          href: 'tel:8602081078225',
          sub: '+86-020-81078225',
        },
        {
          text: 'police',
          icon: 'fa-user-shield',
          color: 'red lighten-3',
          href: 'tel:86110',
          sub: '+86-110',
        },
        {
          text: 'firstAid',
          icon: 'fa-user-shield',
          color: 'red lighten-3',
          href: 'tel:86120',
          sub: '+86-120',
        },
        {
          text: 'fireAlarm',
          icon: 'fa-user-shield',
          color: 'red lighten-3',
          href: 'tel:86119',
          sub: '+86-119',
        },
        {
          text: 'tourismSupervision',
          icon: 'fa-star',
          color: 'indigo',
          href: 'tel:8612301',
          sub: '+86-12301',
        },
        {
          text: 'travelReport',
          icon: 'fa-star',
          color: 'indigo',
          href: 'tel:8612345',
          sub: '+86-12345',
        },
      ],
    };
  },
  beforeCreate() {
      this.$q.loading.hide();
    },
  methods: {
    href(href) {
      if (href) {
        window.location.href = href;
      }
    },
  },
};
</script>
<style scoped>
.item-qrcode {
  max-width: 40px;
  max-height: 40px;
}
</style>
