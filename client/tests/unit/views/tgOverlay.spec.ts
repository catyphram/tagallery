import { shallowMount, createLocalVue } from '@vue/test-utils';
import Vuex from 'vuex';
import { MdIcon } from 'vue-material/dist/components';
import TgOverlay from '@/views/TgOverlay.vue';
import { Image } from '@/models';

const localVue = createLocalVue();
localVue.use(Vuex);
localVue.use(MdIcon);

const state: {
  images: {
    data: Image[],
  },
  selectedImage?: number,
} = {
  images: {
    data: [
      { file: 'test1.jpg' }, { file: 'test2.jpg' }, { file: 'test3.jpg' },
    ] as Image[],
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
    expect(view.image).toEqual(state.images.data[0].file);
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
});
