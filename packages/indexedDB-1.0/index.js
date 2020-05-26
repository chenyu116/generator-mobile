const indexedDB = function(params) {
  return new BP(function(resolve, reject) {
    if (params.store.state.global.db) {
      return resolve();
    }
    const indexedDB =
      window.indexedDB ||
      window.msIndexedDB ||
      window.mozIndexedDB ||
      window.webkitIndexedDB;
    if (!indexedDB) {
      return reject(new Error('不支持 IndexedDB'));
    }
    console.log(params.Vue.prototype.$storage);
    const localDBVersion = parseInt(
      params.Vue.prototype.$storage.get('dbVersion'),
    );
    let newDBVersion = parseInt(
      params.store.state.global.startPointInfo.dbVersion,
    );
    const dbName =
      'signp-' + params.store.state.global.startPointInfo.project_id;
    newDBVersion = isNaN(newDBVersion) ? 1 : newDBVersion;
    if (
      isNaN(localDBVersion) ||
      (localDBVersion !== newDBVersion && localDBVersion !== 0)
    ) {
      indexedDB.deleteDatabase(dbName);
      params.Vue.prototype.$storage.set('dbVersion', newDBVersion);
      location.reload();
    }

    const openRequest = indexedDB.open(dbName, newDBVersion);

    openRequest.onupgradeneeded = function(e) {
      const _db = e.target.result;
      {{range $key,$v:=.DataValues}}
      _db.createObjectStore("{{$key}}", {{$v}});{{end}}
    };
    openRequest.onsuccess = function(e) {
      params.store.commit('global/SET_STATE_PROPERTY', {
        name: 'indexedDB',
        value: e.target.result,
      });
      resolve();
    };
    openRequest.onerror = function() {
      reject(new Error('IndexedDB 初始化失败'));
    };
  });
};

export default ({ store }) => {
  store.commit('global/REGISTERED_INIT_FUNCTION', {
    concurrency: false,
    dependencies: [],
    fn: indexedDB,
  });
};
