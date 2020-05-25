<template>
  <q-layout view="lHh Lpr lFf" class="bg-white">
    <HeaderWithBack icon="{{range $i,$e:=.Config.Data.Values}}{{if eq $e.Key "icon"}}{{$e.Value}}{{end}}{{end}}" title="{{range $i,$e:=.Config.Data.Values}}{{if eq $e.Key "title"}}{{$e.Value}}{{end}}{{end}}"/>
    <q-page-container>
      <q-page></q-page>
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
          sub: '',
          qrcode: 'project/145/yuexiufabu.png',
        },
        {
          text: 'wechatNumber',
          icon: 'fa-comments',
          color: 'indigo',
          sub: '广州越秀文体',
          qrcode: 'project/145/wechatQrcode.png',
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
  mounted() {},
  methods: {},
};
</script>
