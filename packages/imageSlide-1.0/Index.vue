<template>
  <q-card flat="" square="">
    <q-card-section>
      <div class="text-h5">{{print "{{ $t(title) }}"}}</div>
      <div class="text-subtitle2 text-grey">
        - {{print "{{ $t(subtitle) }}"}}
      </div>
    </q-card-section>
    <q-card-section class="q-pt-none">
      <swiper :options="swiperOption">
        <swiper-slide v-for="(b, i) in images" :key="i"
          ><q-img
            width="{{.DataValues.imageWidth}}"
            height="{{.DataValues.imageHeight}}"
            :src="b.img"
            @click="imageClick(b)"
          ></q-img
        ></swiper-slide>
      </swiper>
    </q-card-section>
  </q-card>
</template>

<script>
export default {
  data() {
    return {
      images: [],
      title: '{{.DataValues.title}}',
      subtitle: '{{.DataValues.subtitle}}',
      imageType: '{{.DataValues.imageType}}',
      swiperOption: {
        slidesPerView: {{.DataValues.slidesPerView}},
        spaceBetween: {{.DataValues.spaceBetween}},
        freeMode: true,
        pagination: {
          el: '.swiper-pagination',
          clickable: true,
        },
      },
    };
  },
  computed: {
    locale() {
      return this.$i18n ? this.$i18n.locale : 'zh_CN';
    },
  },
  mounted() {
    this.loadImage();
  },
  methods: {
    loadImage() {
      const self = this;
      self.images = [];
      const readStore = self.$store.state.global.indexedDB
        .transaction('images')
        .objectStore('images')
        .get(self.imageType);
      readStore.onsuccess = function(e) {
        const images = [];
        const r = e.target.result;
        const nowTime = new Date().getTime();
        if (r && r.timestamp > nowTime && r.val.length > 0) {
          for (let i = 0; i < r.val.length; i++) {
            banner.push({
              img: self.ossHost + r.val[i].images_path,
              title: '',
              bindType: r.val[i].images_bind_type,
              link: r.val[i].images_link,
              mapGid: r.val[i].images_map_gid,
            });
          }
          self.images = images;
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
                const writeStore = self.$store.state.global.indexedDB
                  .transaction('images', 'readwrite')
                  .objectStore('images');
                writeStore.put({
                  image_type: self.bannerType,
                  timestamp: nowTime + 300 * 1000,
                  val: resp.body,
                });
                for (let i = 0; i < resp.body.length; i++) {
                  images.push({
                    img: self.ossHost + resp.body[i].images_path,
                    title: '',
                    bindType: resp.body[i].images_bind_type,
                    link: resp.body[i].images_link,
                    mapGid: resp.body[i].images_map_gid,
                  });
                }
                self.images = images;
              }
            });
        }
      };
    },
    imageClick(item) {
      console.log(item);
      switch (item.bindType) {
        case 'link':
          location.href = item.link;
          break;
        // case 'mapGid':
        //   this.navPage('details/' + item.mapGid);
        //   break;
        // case 'internal':
        //   this.navPage(item.link);
        //   break;
      }
    },
  },
};
</script>
