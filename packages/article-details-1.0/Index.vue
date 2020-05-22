<template>
  <q-layout view="lHh Lpr lFf" class="bg-white">
    <HeaderWithBack />
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
  mounted() {},
};
</script>
