const mapPolygons = function(params) {
  return new Promise(function(resolve, reject) {
    const readStore = params.store.state.global.indexedDB
      .transaction('mapPolygonCategory')
      .objectStore('mapPolygonCategory')
      .getAll();
    readStore.onsuccess = function(e) {
      const r = e.target.result;
      if (r && r.length > 0) {
        const mapPolygonsStore = params.store.state.global.indexedDB
          .transaction('mapPolygons')
          .objectStore('mapPolygons')
          .getAll();
        mapPolygonsStore.onsuccess = function(m) {
          const mps = m.target.result;
          if (mps && mps.length > 0) {
            resolve();
          } else {
            const ids = [];
            for (let i = 0; i < r.length; i++) {
              ids.push(r[i].map_polygon_category_id);
            }
            params.Vue.http
              .get(params.Vue.apiHost + '/category/ids', {
                params: {
                  categoryIDs: ids.join(','),
                  projectID:
                    params.store.state.global.startPointInfo.project_id,
                  timestamp: parseInt(new Date().getTime() / 1000),
                },
              })
              .then(function(resp) {
                console.log(resp);
                if (resp.status === 200 && resp.body) {
                  const writeStore = params.store.state.global.indexedDB
                    .transaction('mapPolygons', 'readwrite')
                    .objectStore('mapPolygons');
                  const hasCall = Object.prototype.hasOwnProperty;
                  for (const i in resp.body) {
                    if (hasCall.call(resp.body, i)) {
                      for (let x = 0; x < resp.body[i].length; x++) {
                        resp.body[i][x].changedPoint = params.Vue.$changePoint(
                          resp.body[i][x].point,
                        );
                        writeStore.put(resp.body[i][x]);
                      }
                    }
                  }
                  resolve();
                } else {
                  reject(new Error('读取点位信息失败'));
                }
              })
              .catch(function(err) {
                reject(new Error('读取点位信息失败'));
              });
          }
        };
      } else {
        reject(new Error('读取地图分类信息失败'));
      }
    };
    readStore.onerror = function() {
      reject(new Error('读取地图分类信息失败'));
    };
  });
};

export default ({ store }) => {
  store.commit('global/REGISTERED_INIT_FUNCTION', {
    dependencies: ['indexedDB', 'mapPolygonCategory'],
    concurrency: true,
    fn: mapPolygons,
  });
};
