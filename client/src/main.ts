import Vue from 'vue';
import InfiniteLoading from 'vue-infinite-loading';
import App from './App.vue';
import store from './store';

import './material';

Vue.use(InfiniteLoading as any);

Vue.config.productionTip = false;

new Vue({
  store,
  render: (h) => h(App),
}).$mount('#app');
