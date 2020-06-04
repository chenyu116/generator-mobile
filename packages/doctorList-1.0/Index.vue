<template>
  <q-layout view="lHh Lpr lFf" class="bg-white">
    <q-page-container>
      <HeaderWithBack title="{{.DataValues.title}}" icon="{{.DataValues.icon}}" />
      <q-page class="q-pa-xs">
        <q-spinner-ios
          v-if="loading"
          color="grey"
          class="absolute-center"
          size="2em"
        />
        <div v-if="!loading">
          <q-tabs
            v-model="departmentSelect"
            class="bg-cyan text-white shadow-1"
            inline-label
          >
            <q-tab
              v-for="(item, index) in department"
              :key="index"
              :name="item.department_id"
              :label="item.department_name"
            />
          </q-tabs>
          <q-list separator>
            <q-item
              v-for="(item, index) in list"
              :key="index"
              v-ripple
              class="q-px-none"
              clickable=""
              @click="showDialog(item)"
            >
              <q-item-section top avatar="">
                <q-avatar v-if="item.doctor_photo">
                  <img :src="ossHost + item.doctor_photo" />
                </q-avatar>

                <q-avatar v-if="!item.doctor_photo" square color="grey-2" />
              </q-item-section>
              <q-item-section
                ><q-item-label class="text-subtitle1">{{print "{{item.doctor_name}}"}}</q-item-label>
                <q-item-label caption>{{print "{{item.doctor_job_title}}"}}</q-item-label></q-item-section
              >
              <q-item-section side>
                <q-btn
                  size="12px"
                  flat
                  dense
                  icon="details"
                  :label="$t('description')"
                  stack=""
                />
              </q-item-section>
            </q-item>
          </q-list>
        </div>
      </q-page> </q-page-container
  ></q-layout>
</template>

<script>
export default {
  data() {
    return {
      list: [],
      departmentSelect: '',
      department: [],
      loading: false,
    };
  },
  computed: {
    locale() {
      return this.$i18n ? this.$i18n.locale : 'zh_CN';
    },
  },
  watch: {
    departmentSelect(val) {
      if (!val) return;
      this.loadDoctorList();
    },
  },
  beforeCreate() {
    this.$q.loading.hide();
  },
  mounted() {
    this.initData();
  },
  methods: {
    showDialog(doctor) {
      this.$q
        .dialog({
          // persistent: true,
          ok: this.$t('close'),
          title: doctor.doctor_name,
          message: doctor.doctor_intro,
          html: true,
          focus: 'none',
        });
    },
    initData() {
      const self = this;
      this.loading = true;
      const initFunc = [this.loadDepartmentList];
      BP.mapSeries(initFunc, function(f) {
        return f();
      }).finally(function() {
        self.loading = false;
      });
    },
    sortByOrder(a, b) {
      const sortOrderA = parseInt(a.map_polygon_sort_order);
      const sortOrderB = parseInt(b.map_polygon_sort_order);
      return sortOrderB - sortOrderA;
    },
    loadDepartmentList() {
      this.list = [];
      this.department = [];
      const self = this;
      return new Promise(function(resolve, reject) {
        const readStore = self.$store.state.global.indexedDB
          .transaction('doctor_department')
          .objectStore('doctor_department')
          .getAll();
        readStore.onsuccess = function(e) {
          const r = e.target.result;
          if (r && r.length > 0) {
            self.department = r;
            self.departmentSelect = r[0].department_id;
            resolve();
          } else {
            self.$http
              .get(self.apiHost + '/doctor/department', {
                params: {
                  projectID: 146,
                  departmentID: self.departmentSelect,
                  timestamp: parseInt(new Date().getTime() / 1000),
                },
              })
              .then(function(resp) {
                if (resp.status === 200 && resp.body) {
                  const writeStore = self.$store.state.global.indexedDB
                    .transaction('doctor_department', 'readwrite')
                    .objectStore('doctor_department');
                  for (let i = 0; i < resp.body.length; i++) {
                    writeStore.put(resp.body[i]);
                  }
                  self.department = resp.body;
                  self.departmentSelect = resp.body[0].department_id;
                }
              })
              .finally(function() {
                resolve();
              });
          }
        };
      });
    },
    loadDoctorList() {
      this.list = [];
      const self = this;
      return new Promise(function(resolve, reject) {
        const readStore = self.$store.state.global.indexedDB
          .transaction('doctor')
          .objectStore('doctor')
          .get(self.departmentSelect);
        readStore.onsuccess = function(e) {
          const r = e.target.result;
          if (r && r.data.length > 0) {
            self.list = r.data;
            resolve();
          } else {
            self.$http
              .get(self.apiHost + '/doctor', {
                params: {
                  projectID: 146,
                  departmentID: self.departmentSelect,
                  timestamp: parseInt(new Date().getTime() / 1000),
                },
              })
              .then(function(resp) {
                if (resp.status === 200 && resp.body) {
                  const writeStore = self.$store.state.global.indexedDB
                    .transaction('doctor', 'readwrite')
                    .objectStore('doctor');
                  writeStore.put({
                    department_id: self.departmentSelect,
                    data: resp.body,
                  });
                  self.list = resp.body;
                }
              })
              .finally(function() {
                resolve();
              });
          }
        };
      });
    },
    viewDetails(item) {
      this.$router.replace({ path: '/polygon-details/' + item.map_gid });
    },
    route(item) {
      this.$store.commit('updateCurrentRoute', item);
      this.$router.replace('/route');
    },
  },
};
</script>
<style>
.amap-geolocation-con {
  z-index: 2 !important;
}
</style>
<style scoped>
.toMapBtn {
  position: absolute;
  bottom: 70px;
  left: 10px;
  z-index: 2;
}
.cate-item-image {
  background-size: 100%;
  float: left;
  height: 4rem;
  width: 4rem;
  display: flex;
}
.amap {
  width: 100% !important;
  height: 220px;
}
</style>
