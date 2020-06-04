const mapPolygonCategory = function(params) {
  return new Promise(function(resolve, reject) {
    const readStore = params.store.state.global.indexedDB
      .transaction('mapPolygonCategory')
      .objectStore('mapPolygonCategory')
      .getAll();
    readStore.onsuccess = function(e) {
      const r = e.target.result;
      if (r && r.length > 0) {
        resolve();
      } else {
        params.Vue.http
          .get(params.Vue.apiHost + '/category/all', {
            params: {
              projectID: params.store.state.global.startPointInfo.project_id,
              timestamp: parseInt(new Date().getTime() / 1000),
            },
          })
          .then(function(resp) {
            if (resp.status === 200) {
              const writeStore = params.store.state.global.indexedDB
                .transaction('mapPolygonCategory', 'readwrite')
                .objectStore('mapPolygonCategory');
              for (let i = 0; i < resp.body.length; i++) {
                writeStore.put(resp.body[i]);
              }
              resolve();
            } else {
              reject(new Error('读取地图分类信息失败'));
            }
          })
          .catch(function() {
            reject(new Error('读取地图分类信息失败'));
          });
      }
    };
    readStore.onerror = function() {
      reject(new Error('读取地图分类信息失败'));
    };
  });
};

export default ({ store }) => {
  store.commit('global/REGISTERED_INIT_FUNCTION', {
    dependencies: ['indexedDB'],
    concurrency: false,
    fn: mapPolygonCategory,
  });
};
