import Vue from 'vue';
import Vuex, { Store } from 'vuex';

import * as api from '@/api/api';
import { Category, LIST_MODE } from './models';

Vue.use(Vuex);

export const state = {
  categories: {
    loading: false,
    data: [] as Category[],
  },
  listMode: LIST_MODE.MODE_VIEW,
  selectedCategories: [] as Category[],
  listUncategorized: false,
};

export const getters = {
  isCategorySelected: (state) => (category: Category) => state.selectedCategories.includes(category),
};

export const mutations = {
  updateCategories(state, { categories, loading }) {
    state.categories.data = categories;
    state.categories.loading = loading;
  },
  setMode(state, { mode }: {mode: LIST_MODE}) {
    state.listMode = mode;
  },
  selectCategory(state, { category }: { category: Category }) {
    if (!state.selectedCategories.includes(category)) {
      state.selectedCategories.push(category);
      state.listUncategorized = false;
    }
  },
  unselectCategory(state, { category }: { category: Category }) {
    for (let i = 0; i < state.selectedCategories.length; i++) {
      if (state.selectedCategories[i] === category) {
        state.selectedCategories.splice(i, 1);
        break;
      }
    }
  },
  setListUncategorized(state, { flag }: { flag: boolean }) {
    state.listUncategorized = flag;
    if (flag) {
      state.selectedCategories = [];
    }
  },
};

export const actions = {
  async loadCategories(context) {
    context.commit('updateCategories', {
      categories: [],
      loading: true,
    });
    context.commit('updateCategories', {
      categories: await api.loadCategories(),
      loading: false,
    });
  },
  setMode(context, mode: LIST_MODE) {
    return context.commit('setMode', { mode });
  },
  toggleCategory(context, category: Category) {
    if (context.state.selectedCategories.includes(category)) {
      return context.commit('unselectCategory', { category });
    } else {
      return context.commit('selectCategory', { category });
    }
  },
  toggleListUncategorized(context) {
    return context.commit('setListUncategorized', { flag: !context.state.listUncategorized });
  },
};

export default new Vuex.Store({
  state,
  getters,
  mutations,
  actions,
});
