<template>
  <q-card flat="" square="">
    <q-card-section>
      <div class="text-h5">__data.title__</div>
      <div class="text-subtitle2 q-ml-md text-grey">- __data.subtitle</div>
    </q-card-section>
    <q-card-section class="q-pt-none">
      <q-list separator>
        <q-item v-ripple clickable>
          <q-item-section avatar="" class="text-grey">5/22 </q-item-section>
          <q-item-section
            >大雨！暴雨！大暴雨！@越秀街坊注意：“龙舟水”即将到货大雨！暴雨！大暴雨！@越秀街坊注意：“龙舟水”即将到货</q-item-section
          >
        </q-item>
      </q-list>
    </q-card-section>
  </q-card>
</template>
<script>
export default {
  name: "Index",

  components: {},
  data() {
    const self = this;
    return {
      articles: [],
    };
  },
  computed: {},
  mounted() {
    this.loadArticles();
  },
  methods: {
    featureClick(onClick) {
      if (typeof onClick === "function") {
        onClick();
      }
    },
    loadArticles() {
      const self = this;
      self.articles = [];
      const readStore = self.$store.state.global.indexedDB
        .transaction("articles")
        .objectStore("articles")
        .get(
          self.$store.state.startPointInfo.index_article_category +
            "_" +
            self.$i18n.locale
        );
      readStore.onsuccess = function (e) {
        const r = e.target.result;
        const nowTime = new Date().getTime();
        if (r && r.timestamp > nowTime && r.val.length > 0) {
          for (let i = 0; i < r.val.length; i++) {
            if (i === self.$store.state.startPointInfo.showArticleLimit) break;
            const _d = new Date(parseInt(r.val[i].article_create_at) * 1000);
            r.val[i].date =
              _d.getFullYear() + "/" + (_d.getMonth() + 1) + "/" + _d.getDate();
            self.articles.push(r.val[i]);
          }
        } else {
          self.$http
            .get(self.apiHost + "/articles", {
              params: {
                categoryID:
                  self.$store.state.startPointInfo.index_article_category,
                locale: self.$i18n.locale,
                limit: self.$store.state.startPointInfo.articleLimit,
                projectID: self.$store.state.startPointInfo.project_id,
                timestamp: parseInt(new Date().getTime() / 1000),
              },
            })
            .then(function (resp) {
              if (resp.status === 200) {
                const timestamp = new Date().getTime() + 300 * 1000;
                const writeStore = self.$store.state.db
                  .transaction("articles", "readwrite")
                  .objectStore("articles");
                writeStore.put({
                  article_category_id:
                    self.$store.state.startPointInfo.index_article_category +
                    "_" +
                    self.$i18n.locale,
                  timestamp: timestamp,
                  val: resp.body,
                });
                for (let i = 0; i < resp.body.length; i++) {
                  if (i === self.$store.state.startPointInfo.showArticleLimit) {
                    break;
                  }
                  const _d = new Date(
                    parseInt(resp.body[i].article_create_at) * 1000
                  );
                  resp.body[i].date =
                    _d.getFullYear() +
                    "/" +
                    (_d.getMonth() + 1) +
                    "/" +
                    _d.getDate();
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
