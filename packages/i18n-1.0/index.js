import Vue from 'vue';
import VueI18n from 'vue-i18n';
import i18nComponent from './i18n.vue';
const readI18nMessagesFromRemote = function(params) {
  return new Promise(function(resolve) {
    const readStore = params.store.state.global.indexedDB
      .transaction('language')
      .objectStore('language')
      .getAll();
    readStore.onsuccess = function(e) {
      const r = e.target.result;
      if (r && r.length > 0) {
        resolve(r);
      } else {
        const res = [];
        params.Vue.http
          .get(params.store.state.global.apiHost + 'language', {
            params: {
              project_id: 145,
              timestamp: parseInt(new Date().getTime() / 1000),
            },
          })
          .then(function(resp) {
            if (resp.status === 200 && resp.body) {
              const writeStore = params.store.state.global.indexedDB
                .transaction('language', 'readwrite')
                .objectStore('language');
              const lang = {};
              for (let i = 0; i < resp.body.length; i++) {
                if (!lang[resp.body[i].language_locale]) {
                  lang[resp.body[i].language_locale] = {};
                }
                lang[resp.body[i].language_locale][resp.body[i].language_key] =
                  resp.body[i].language_value;
              }
              const hasCall = Object.prototype.hasOwnProperty;
              for (const locale in lang) {
                if (hasCall.call(lang, locale)) {
                  const l = {
                    locale: locale,
                    val: lang[locale],
                  };
                  writeStore.put(l);
                  res.push(l);
                }
              }
            } else {
              console.error(new Error('读取语言包失败'));
            }
          })
          .catch(function() {
            console.error(new Error('读取语言包失败'));
          })
          .finally(function() {
            resolve(res);
          });
      }
    };
    readStore.onerror = function() {
      console.error(new Error('读取语言包失败'));
      resolve(res);
    };
  }).then(function(res) {
    for (let i = 0; i < res.length; i++) {
      params.app.i18n.setLocaleMessage(res[i].locale, res[i].val);
    }
    const localLanguage =
      params.Vue.prototype.$storage.get('locale') || 'zh_CN';
    params.app.i18n.locale = localLanguage;
    return Promise.resolve();
  });
};
Vue.use(VueI18n);

const i18n = new VueI18n({
  locale: 'zh_CN',
  fallbackLocale: 'zh_CN',
});

Vue.component('i18n', i18nComponent);

export default ({ app, store }) => {
  // Set i18n instance on app
  app.i18n = i18n;
  store.commit('global/REGISTERED_INIT_FUNCTION', {
    dependencies: ['indexedDB'],
    concurrency: true,
    fn: readI18nMessagesFromRemote,
  });
};

export { i18n };
