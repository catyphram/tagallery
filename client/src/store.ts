import Vue from 'vue';
import Vuex from 'vuex';

import * as api from '@/api/api';
import { Category, Image, LIST_MODE } from './models';

Vue.use(Vuex);

export const state = {
  categories: {
    loading: false,
    data: [] as Category[],
  },
  images: {
    loading: false,
    completed: false,
    data: [] as Image[],
  },
  listMode: LIST_MODE.MODE_VIEW,
  selectedCategories: [] as Category[],
  listUncategorized: false,
};

export const getters = {
  isCategorySelected: (state) => (category: Category) => state.selectedCategories.includes(category),
};

export const mutations = {
  updateCategories(state, { categories, loading = false }: { categories?: Category[], loading?: boolean }) {
    state.categories.data = categories;
    state.categories.loading = loading;
  },
  updateImages(state, {
    images, loading = false, completed = false, append = false,
  }: { images: Image[], loading?: boolean, completed?: boolean, append?: boolean }) {
    if (append) {
      state.images.data.push(...images);
    } else {
      state.images.data = images;
    }
    state.images.loading = loading;
    state.images.completed = completed;
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
    context.dispatch('loadImages');
  },
  async loadImages(context, append = false) {
    context.commit('updateImages', {
      images: [],
      loading: true,
      append,
    });

    const images = await api.loadImages();
    context.commit('updateImages', {
      images,
      loading: false,
      completed: !images.length,
      append,
    });
  },
  setMode(context, mode: LIST_MODE) {
    context.commit('setMode', { mode });
    context.dispatch('loadImages');
  },
  toggleCategory(context, category: Category) {
    if (context.state.selectedCategories.includes(category)) {
      context.commit('unselectCategory', { category });
    } else {
      context.commit('selectCategory', { category });
    }
    context.dispatch('loadImages');
  },
  toggleListUncategorized(context) {
    context.commit('setListUncategorized', { flag: !context.state.listUncategorized });
    context.dispatch('loadImages');
  },
};

export default new Vuex.Store({
  state,
  getters,
  mutations,
  actions,
});
