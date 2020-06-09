const mapList = function(params) {
  return new Promise(function(resolve, reject) {
    const readStore = params.store.state.global.indexedDB
      .transaction('mapList')
      .objectStore('mapList')
      .getAll();
    readStore.onsuccess = function(e) {
      const r = e.target.result;
      if (r && r.length > 0) {
        resolve();
      } else {
        params.Vue.http
          .get(params.Vue.apiHost + '/map/list', {
            params: {
              projectID: params.store.state.global.startPointInfo.project_id,
              timestamp: parseInt(new Date().getTime() / 1000),
            },
          })
          .then(function(resp) {
            if (resp.status === 200 && resp.body) {
              const writeStore = params.store.state.global.indexedDB
                .transaction('mapList', 'readwrite')
                .objectStore('mapList');
              const hasCall = Object.prototype.hasOwnProperty;
              for (const i in resp.body) {
                if (hasCall.call(resp.body, i)) {
                  writeStore.put(resp.body[i]);
                }
              }
              resolve();
            } else {
              reject(new Error('读取地图数据失败'));
            }
          })
          .catch(function() {
            reject(new Error('读取地图数据失败'));
          });
      }
    };
    readStore.onerror = function() {
      reject(new Error('读取地图数据失败'));
    };
  });
};

export default ({ store }) => {
  store.commit('global/REGISTERED_INIT_FUNCTION', {
    dependencies: ['indexedDB'],
    concurrency: true,
    fn: mapList,
    name: 'mapList',
  });
};
