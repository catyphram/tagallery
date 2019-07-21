import Vue from 'vue';
import Vuex from 'vuex';

import * as api from '@/api/api';
import { Category, Image, LIST_MODE } from './models';

Vue.use(Vuex);

export interface State  {
  categories: {
    loading: boolean;
    data: Category[];
    error?: string;
  };
  images: {
    loading: boolean;
    completed: boolean;
    data: Image[];
  };
  selectedImage?: number;
  listMode: LIST_MODE;
  selectedCategories: Category[];
  listUncategorized: boolean;
}

export const state: State = {
  categories: {
    loading: false,
    data: [] as Category[],
  },
  images: {
    loading: false,
    completed: false,
    data: [] as Image[],
  },
  selectedImage: undefined,
  listMode: LIST_MODE.MODE_VIEW,
  selectedCategories: [] as Category[],
  listUncategorized: false,
};

export const getters = {
  isCategorySelected: (state: State) => (category: Category) => state.selectedCategories.includes(category),
};

export const mutations = {
  updateCategories(state: State, {
    categories, loading = false, error,
  }: { categories: Category[], loading?: boolean, error?: any }) {
    state.categories.data = categories;
    state.categories.loading = loading;
    state.categories.error = error ? `${error}` : undefined;
  },
  updateImages(state: State, {
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
  updateImage(state: State, { image }: { image: Image }) {
    const imageIndex = state.images.data.findIndex((storeImage) => {
      return storeImage.file === image.file;
    });
    if (imageIndex >= 0) {
      state.images.data[imageIndex].assignedCategories = image.assignedCategories;
      state.images.data[imageIndex].proposedCategories = image.proposedCategories;
      state.images.data[imageIndex].starredCategory = image.starredCategory;
    }
  },
  selectImage(state: State, { index }: { index?: number }) {
    state.selectedImage = index;
  },
  setMode(state: State, { mode }: { mode: LIST_MODE }) {
    state.listMode = mode;
  },
  selectCategory(state: State, { category }: { category: Category }) {
    if (!state.selectedCategories.includes(category)) {
      state.selectedCategories.push(category);
      state.listUncategorized = false;
    }
  },
  unselectCategory(state: State, { category }: { category: Category }) {
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
  async toggleImageCategory(context, { image, category }: {
    image: Image, category: Category,
  }) {
    const newImage: Image = { ...image };

    if (image.proposedCategories) {
      newImage.proposedCategories = image.proposedCategories.filter((proposedCategory) => {
        return proposedCategory !== category.key;
      });
    }

    let assignCategory = true;

    if (image.assignedCategories) {
      newImage.assignedCategories = image.assignedCategories.filter((assignedCategory) => {
        if (assignedCategory === category.key) {
          assignCategory = false;
          return false;
        }
        return true;
      });
    } else {
      newImage.assignedCategories = [];
    }

    if (assignCategory) {
      newImage.assignedCategories.push(category.key);
    }

    context.commit('updateImage', { image: newImage });
    try {
      await api.updateImage(newImage);
    } catch (error) {
      context.commit('updateImage', { image });
    }
  },
  async loadCategories(context) {
    context.commit('updateCategories', {
      categories: [],
      loading: true,
    });

    try {
      context.commit('updateCategories', {
        categories: await api.loadCategories(),
        loading: false,
      });
    } catch (error) {
      context.commit('updateCategories', {
        categories: [],
        loading: false,
        error,
      });
    }

    context.dispatch('loadImages');
  },
  async loadImages(context, append = false) {
    context.commit('updateImages', {
      images: [],
      loading: true,
      append,
    });

    // @TODO: Add error handling, set loading to false
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
