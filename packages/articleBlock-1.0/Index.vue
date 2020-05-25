<template>
  <q-card flat="" square="">
    <q-card-section>
      <div class="text-h5">{{print "{{ $t('title') }}"}}</div>
      <div class="text-subtitle2 q-ml-md text-grey">
        - {{print "{{ $t('subtitle') }}"}}
      </div>
    </q-card-section>
    <q-card-section class="q-pt-none">
      <q-list separator>
        <q-item
          v-for="(item, index) in articles"
          :key="`articles-${index}`"
          v-ripple
          clickable
          @click="featureClick(item)"
        >
          <q-item-section avatar="" class="text-grey"
            >{{print "{{ item.date }}"}}
          </q-item-section>
          <q-item-section>{{print "{{ item.article_title }}"}}</q-item-section>
        </q-item>
      </q-list>
    </q-card-section>
  </q-card>
</template>
<script>
export default {
  data() {
    return {
      title:
        '{{range $i,$e:=.Config.Data.Values}}{{if eq $e.Key "title"}}{{ $e.Value }}{{ end }}{{ end }}',
      subtitle:
        '{{range $i,$e:=.Config.Data.Values}}{{if eq $e.Key "subtitle"}}{{ $e.Value }}{{ end }}{{ end }}',
      articles: [],
    };
  },
  computed: {
    locale() {
      return this.$i18n ? this.$i18n.locale : 'zh_CN';
    },
  },
  mounted() {
    this.loadArticles();
  },
  methods: {
    featureClick(item) {
      let path =
        '{{range $i,$e:=.Config.Components}}{{if eq $e.Key "articleDetails"}}{{range $ia,$ie:=$e.Values}}{{range $pi,$pe:=$ie.ProjectFeaturesConfig.Data.Values}}{{if eq $pe.Key "routePath"}}{{$pe.Value}}{{end}}{{end}}{{end}}{{end}}{{end}}';
      const idIndex = path.lastIndexOf(':');
      path = path.replace(path.substring(idIndex), item.article_id);
      this.$router.push({ path: path, replace: true });
    },
    loadArticles() {
      const self = this;
      const articleCategoryId =
        '{{range $i,$e:=.Config.Data.Values}}{{if eq $e.Key "articleCategoryId"}}{{$e.Value}}{{end}}{{end}}';
      let articleLimit = parseInt(
        '{{range $i,$e:=.Config.Data.Values}}{{if eq $e.Key "articleShowCount"}}{{$e.Value}}{{end}}{{end}}',
      );
      if (isNaN(articleLimit) || articleLimit <= 0) {
        articleLimit = 3;
      }
      const articleDbTag = articleCategoryId + '_' + self.locale;
      self.articles = [];
      const readStore = self.$store.state.global.indexedDB
        .transaction('articles')
        .objectStore('articles')
        .get(articleDbTag);
      readStore.onsuccess = function(e) {
        const r = e.target.result;
        const nowTime = new Date().getTime();
        if (r && r.timestamp > nowTime && r.val.length > 0) {
          for (let i = 0; i < r.val.length; i++) {
            if (i === articleLimit) break;
            const _d = new Date(parseInt(r.val[i].article_create_at) * 1000);
            r.val[i].date = _d.getMonth() + 1 + '/' + _d.getDate();
            self.articles.push(r.val[i]);
          }
        } else {
          self.$http
            .get(self.apiHost + '/articles', {
              params: {
                categoryID: articleCategoryId,
                locale: self.$locale,
                limit: articleLimit,
                projectID: self.$store.state.global.startPointInfo.project_id,
                timestamp: parseInt(nowTime / 1000),
              },
            })
            .then(function(resp) {
              if (resp.status === 200) {
                const timestamp = nowTime + 300 * 1000;
                const writeStore = self.$store.state.global.indexedDB
                  .transaction('articles', 'readwrite')
                  .objectStore('articles');
                writeStore.put({
                  article_category_id: articleDbTag,
                  timestamp: timestamp,
                  val: resp.body,
                });
                for (let i = 0; i < resp.body.length; i++) {
                  if (i === articleLimit) {
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
