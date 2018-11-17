import { shallowMount, createLocalVue } from '@vue/test-utils';
import Vuex from 'vuex';
import TgMenu from '@/components/TgMenu.vue';
import { LIST_MODE, Category } from '@/models';
import { MdButton, MdList, MdSwitch } from 'vue-material/dist/components';

const localVue = createLocalVue();
localVue.use(Vuex);
localVue.use(MdButton);
localVue.use(MdList);
localVue.use(MdSwitch);

const state = {
  selectedCategories: [
    { name: 'Category 1' },
  ] as Category[],
  categories: {
    data: [
      { name: 'Category 1' }, { name: 'Category 2' }, { name: 'Category 3' },
    ] as Category[],
    loading: false,
  },
  listMode: LIST_MODE.MODE_VIEW,
  listUncategorized: false,
};

const actions =  {
  setMode: jest.fn(),
};

const store = new Vuex.Store({
  state,
  actions,
  getters: {
    isCategorySelected: (state) => (category: Category) => state.selectedCategories.includes(category),
  },
});

describe('TgMenu.vue', () => {
  it('categories should select the categories from the store', () => {
    const wrapper = shallowMount(TgMenu, { store, localVue });
    const view = wrapper.vm as any;
    expect(view.categories).toEqual(state.categories.data);
  });

  it('categoriesLoading should select the loading status of the categories from the store', () => {
    const wrapper = shallowMount(TgMenu, { store, localVue });
    const view = wrapper.vm as any;
    expect(view.categoriesLoading).toEqual(state.categories.loading);
  });

  it('updateMode should toggle the list mode', () => {
    const wrapper = shallowMount(TgMenu, { store, localVue });
    const view = wrapper.vm as any;

    view.updateMode(true);
    expect(actions.setMode).toHaveBeenCalledWith(expect.anything(), LIST_MODE.MODE_VERIFY, undefined);
    view.updateMode(false);
    expect(actions.setMode).toHaveBeenCalledWith(expect.anything(), LIST_MODE.MODE_VIEW, undefined);
  });
});
