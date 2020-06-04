<template>
  <q-layout view="lHh Lpr lFf" class="bg-white">
    <q-page-container>
      <HeaderWithBack title="{{.DataValues.title}}" icon="{{.DataValues.icon}}" />
      <q-page class="q-pa-md">
        <q-spinner-ios
          v-if="loading"
          color="grey"
          class="absolute-center"
          size="2em"
        />
        <div v-if="!loading">
          <q-list separator>
            <q-item
              v-for="(item, index) in list"
              :key="index"
              v-ripple
              class="q-px-none"
              clickable=""
              @click="showDialog(item)"
            >
              <q-item-section
                ><q-item-label class="text-subtitle1">{{print "{{item.department_name}}"}}</q-item-label>
                <q-item-label caption>{{print "{{item.doctor_job_title"}}
                }}</q-item-label></q-item-section
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
      //   this.loadDoctorList();
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
          //   persistent: true,
          ok: this.$t('close'),
          title: doctor.department_name,
          message: doctor.department_intro,
          html: true,
          focus: 'none',
        })
        .onOk(() => {
          // console.log('OK')
        })
        .onCancel(() => {
          // console.log('Cancel')
        })
        .onDismiss(() => {
          // console.log('I am triggered on both OK and Cancel')
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
        console.log('self.loading ', self.loading);
      });
    },
    sortByOrder(a, b) {
      const sortOrderA = parseInt(a.map_polygon_sort_order);
      const sortOrderB = parseInt(b.map_polygon_sort_order);
      return sortOrderB - sortOrderA;
    },
    loadDepartmentList() {
      this.list = [];
      const self = this;
      return new Promise(function(resolve, reject) {
        const readStore = self.$store.state.global.indexedDB
          .transaction('doctor_department')
          .objectStore('doctor_department')
          .getAll();
        readStore.onsuccess = function(e) {
          const r = e.target.result;
          console.log('r', r);
          if (r && r.length > 0) {
            self.list = r;
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
