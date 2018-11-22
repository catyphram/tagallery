import { config, shallowMount, createLocalVue } from '@vue/test-utils';
import Vuex from 'vuex';
import App from '@/App.vue';
import { LIST_MODE, Category } from '@/models';
import { MdApp, MdButton, MdIcon, MdToolbar } from 'vue-material/dist/components';

// https://github.com/vuejs/vue-test-utils/issues/973
config.logModifiedComponents = false;

const localVue = createLocalVue();
localVue.use(Vuex);
localVue.use(MdApp);
localVue.use(MdButton);
localVue.use(MdIcon);
localVue.use(MdToolbar);

const store = new Vuex.Store({
  state: {
    selectedCategories: [
      { name: 'Category 1' }, { name: 'Category 2' }, { name: 'Category 3' },
    ] as Category[],
    listMode: LIST_MODE.MODE_VERIFY,
    listUncategorized: false,
  },
  actions: {
    loadCategories: jest.fn(),
  },
});

describe('App.vue', () => {
  it('menuVisible should be visible by default', () => {
    const wrapper = shallowMount(App, { store, localVue });
    const view = wrapper.vm as any;
    expect(view.menuVisible).toEqual(true);
  });

  it('toggleMenu should toggle #menuVisible', () => {
    const wrapper = shallowMount(App, { store, localVue });
    const view = wrapper.vm as any;
    const prevState: boolean = view.menuVisible;
    view.toggleMenu();
    expect(view.menuVisible).toEqual(!prevState);
  });

  it('title should be set correctly when verifying images from multiple categories', () => {
    const wrapper = shallowMount(App, { store, localVue });
    const view = wrapper.vm as any;
    expect(view.title).toEqual('Verifying images from categories Category 1, Category 2, Category 3');
  });

  it('title should be set correctly when viewing images from one category', () => {
    const store = new Vuex.Store({
      state: {
        selectedCategories: [
          { name: 'Category 2' },
        ] as Category[],
        listMode: LIST_MODE.MODE_VIEW,
        listUncategorized: false,
      },
      actions: {
        loadCategories: jest.fn(),
      },
    });
    const wrapper = shallowMount(App, { store, localVue });
    const view = wrapper.vm as any;
    expect(view.title).toEqual('Viewing images from category Category 2');
  });

  it('title should be set correctly when viewing uncategorized images', () => {
    const store = new Vuex.Store({
      state: {
        selectedCategories: [] as Category[],
        listMode: LIST_MODE.MODE_VIEW,
        listUncategorized: true,
      },
      actions: {
        loadCategories: jest.fn(),
      },
    });
    const wrapper = shallowMount(App, { store, localVue });
    const view = wrapper.vm as any;
    expect(view.title).toEqual('Viewing uncategorized images');
  });
});
