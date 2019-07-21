import { shallowMount, createLocalVue } from '@vue/test-utils';
import Vuex from 'vuex';
import { MdIcon, MdButton } from 'vue-material/dist/components';
import TgOverlay from '@/views/TgOverlay.vue';
import { Image, Category } from '@/models';

const localVue = createLocalVue();
localVue.use(Vuex);
localVue.use(MdIcon);
localVue.use(MdButton);

const state: {
  images: {
    data: Image[],
  },
  categories: {
    data: Category[],
  },
  selectedImage?: number,
} = {
  images: {
    data: [
      { file: 'test1.jpg' }, { file: 'test2.jpg' }, { file: 'test3.jpg' },
    ] as Image[],
  },
  categories: {
    data: [
      { name: 'Category 1', key: 'category1' },
      { name: 'Category 2', key: 'category2' },
    ],
  },
};

const mutations = {
  selectImage(state, { index }: { index?: number }) {
    state.selectedImage = index;
  },
};

const store = new Vuex.Store({
  mutations,
  state,
});

describe('TgOverlay.vue', () => {
  it('image should return undefined when no image is selected', () => {
    const wrapper = shallowMount(TgOverlay, { store, localVue });
    const view = wrapper.vm as any;
    expect(view.image).toBeUndefined();
  });

  it('image should select the selected image from the store', () => {
    const wrapper = shallowMount(TgOverlay, { store: new Vuex.Store({
      state: {
        ...state,
        selectedImage: 0,
      },
    }), localVue });
    const view = wrapper.vm as any;
    expect(view.image).toEqual(state.images.data[0]);
  });

  it('categories should return undefined when no image is selected', () => {
    const wrapper = shallowMount(TgOverlay, { store, localVue });
    const view = wrapper.vm as any;
    expect(view.categories).toBeUndefined();
  });

  it('categories should select and map the categories from the store', () => {
    const wrapper = shallowMount(TgOverlay, { store: new Vuex.Store({
      state: {
        ...state,
        selectedImage: 0,
        images: {
          data: [{
            file: 'test1.jpg',
            proposedCategories: ['category1'],
            assignedCategories: ['category2'],
          }] as Image[],
        },
      },
    }), localVue });
    const view = wrapper.vm as any;
    expect(view.categories).toMatchSnapshot();
  });

  it('prevImage should not select the previous image if no image is selected', () => {
    const wrapper = shallowMount(TgOverlay, { store, localVue });
    const view = wrapper.vm as any;

    store.commit = jest.fn();
    view.prevImage();

    expect(store.commit).not.toHaveBeenCalled();
  });

  it('prevImage should select the next image', () => {
    const localStore = new Vuex.Store({
      state: {
        ...state,
        selectedImage: 1,
      },
      mutations,
    });
    const wrapper = shallowMount(TgOverlay, { store: localStore, localVue });
    const view = wrapper.vm as any;

    localStore.commit = jest.fn();
    view.prevImage();

    expect(localStore.commit).toHaveBeenCalledWith('selectImage', { index: 0 });
  });

  it('nextImage should not select the next image if no image is selected', () => {
    const wrapper = shallowMount(TgOverlay, { store, localVue });
    const view = wrapper.vm as any;

    store.commit = jest.fn();
    view.nextImage();

    expect(store.commit).not.toHaveBeenCalled();
  });

  it('nextImage should select the next image', () => {
    const localStore = new Vuex.Store({
      state: {
        ...state,
        selectedImage: 1,
      },
      mutations,
    });
    const wrapper = shallowMount(TgOverlay, { store: localStore, localVue });
    const view = wrapper.vm as any;

    localStore.commit = jest.fn();
    view.nextImage();

    expect(localStore.commit).toHaveBeenCalledWith('selectImage', { index: 2 });
  });

  it('close should unselect the image', () => {
    const localStore = new Vuex.Store({
      state: {
        ...state,
        selectedImage: 0,
      },
      mutations,
    });
    const wrapper = shallowMount(TgOverlay, { store: localStore, localVue });
    const view = wrapper.vm as any;

    localStore.commit = jest.fn();
    view.close();

    expect(localStore.commit).toHaveBeenCalledWith('selectImage', { index: undefined });
  });

  it('toggleImageCategory should toggle a assigned category of an image', () => {
    const localStore = new Vuex.Store({
      state: {
        ...state,
        selectedImage: 0,
      },
      mutations,
    });
    const wrapper = shallowMount(TgOverlay, { store: localStore, localVue });
    const view = wrapper.vm as any;

    localStore.dispatch = jest.fn();
    view.toggleImageCategory({ name: 'Category 1', key: 'category1' } as Category);

    expect((localStore.dispatch as jest.Mock).mock.calls).toMatchSnapshot();
  });
});
