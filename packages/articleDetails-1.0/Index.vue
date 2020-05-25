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
      return {};
    },
    computed: {
      // tabWidth() {
      //   return parseInt(
      //     this.$store.state.global.windowSize.width / this.countsPerRow
      //   );
      // },
    },
    mounted() {
    },
  };
</script>
