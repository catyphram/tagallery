import { shallowMount, createLocalVue } from '@vue/test-utils';
import Vuex from 'vuex';
import TgGallery from '@/views/TgGallery.vue';
import InfiniteLoading from 'vue-infinite-loading';
import { Image } from '@/models';

const localVue = createLocalVue();
localVue.use(Vuex);
localVue.use(InfiniteLoading as any);

const state: {
  images: {
    data: Image[],
    completed: boolean,
  },
  selectedImage?: number,
} = {
  images: {
    data: [
      { file: 'test1.jpg' }, { file: 'test2.jpg' }, { file: 'test3.jpg' },
    ] as Image[],
    completed: false,
  },
};

const actions = {
  loadImages: jest.fn(),
};

const mutations = {
  selectImage(state, { index }: { index?: number }) {
    state.selectedImage = index;
  },
};

const store = new Vuex.Store({
  mutations,
  state,
  actions,
});

describe('TgGallery.vue', () => {
  it('images should select the images from the store', () => {
    const wrapper = shallowMount(TgGallery, { store, localVue });
    const view = wrapper.vm as any;
    expect(view.images).toEqual(state.images.data);
  });

  it('scroll should load another batch of images', async () => {
    const scrollState = {
      complete: jest.fn(),
      loaded: jest.fn(),
    };
    const wrapper = shallowMount(TgGallery, { store, localVue });
    const view = wrapper.vm as any;

    await view.scroll(scrollState);
    expect(actions.loadImages).toHaveBeenCalledWith(expect.anything(), true, undefined);
    expect(scrollState.loaded).toHaveBeenCalled();
  });

  it('scroll should not load images in parallel', async () => {
    const scrollState = {
      complete: jest.fn(),
      loaded: jest.fn(),
    };
    const store = new Vuex.Store({
      state: {
        images: {
          data: [] as Image[],
          loading: true,
        },
      },
      actions,
    });

    const wrapper = shallowMount(TgGallery, { store, localVue });
    const view = wrapper.vm as any;

    await view.scroll(scrollState);
    expect(scrollState.complete).not.toHaveBeenCalled();
    expect(scrollState.loaded).not.toHaveBeenCalled();
  });

  it('scroll should complete the infinite scroll once all images have been loaded', async () => {
    const scrollState = {
      complete: jest.fn(),
      loaded: jest.fn(),
    };
    const store = new Vuex.Store({
      state: {
        images: {
          data: [] as Image[],
          completed: true,
        },
      },
      actions,
    });

    const wrapper = shallowMount(TgGallery, { store, localVue });
    const view = wrapper.vm as any;

    await view.scroll(scrollState);
    expect(scrollState.complete).toHaveBeenCalled();
  });

  it('watchSelectedImage should scroll to the selected image', async (done) => {
    const store = new Vuex.Store({
      state: {
        images: {
          data: [{file: ''}] as Image[],
        },
        selectedImage: 0,
      },
      mutations,
    });

    const wrapper = shallowMount(TgGallery, { store, localVue });
    const view = wrapper.vm as any;

    view.$nextTick(() => {
      view.$refs.images[view.$store.state.selectedImage].$el.scrollIntoView = jest.fn();
      view.watchSelectedImage();

      expect(view.$refs.images[view.$store.state.selectedImage].$el.scrollIntoView).toHaveBeenCalled();
      done();
    });
  });

  it('watchSelectedImage should not scroll if there is no image selected', async (done) => {
    const store = new Vuex.Store({
      state: {
        images: {
          data: [] as Image[],
        },
        selectedImage: 0,
      },
      mutations,
    });

    const wrapper = shallowMount(TgGallery, { store, localVue });
    const view = wrapper.vm as any;

    view.$nextTick(() => {
      jest.spyOn(view, 'watchSelectedImage');
      view.watchSelectedImage();

      expect(view.$refs.images).toBeUndefined();
      done();
    });
  });

});
