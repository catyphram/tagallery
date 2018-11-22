import { shallowMount, createLocalVue } from '@vue/test-utils';
import Vuex from 'vuex';
import Gallery from '@/views/Gallery.vue';
import InfiniteLoading from 'vue-infinite-loading';
import { Image } from '@/models';

const localVue = createLocalVue();
localVue.use(Vuex);
localVue.use(InfiniteLoading as any);

const state = {
  images: {
    data: [
      { file: 'test1.jpg' }, { file: 'test2.jpg' }, { file: 'test3.jpg' },
    ] as Image[],
    completed: false,
  },
};

const actions =  {
  loadImages: jest.fn(),
};

const store = new Vuex.Store({
  state,
  actions,
});

let scrollState;

describe('Gallery.vue', () => {
  beforeEach(() => {
    scrollState = {
      complete: jest.fn(),
      loaded: jest.fn(),
    };
  });

  it('images should select the images from the store', () => {
    const wrapper = shallowMount(Gallery, { store, localVue });
    const view = wrapper.vm as any;
    expect(view.images).toEqual(state.images.data);
  });

  it('scroll should load another batch of images', async () => {
    const wrapper = shallowMount(Gallery, { store, localVue });
    const view = wrapper.vm as any;

    await view.scroll(scrollState);
    expect(actions.loadImages).toHaveBeenCalledWith(expect.anything(), true, undefined);
    expect(scrollState.loaded).toHaveBeenCalled();
  });

  it('scroll should complete the infinite scroll once all images have been loaded', async () => {
    const store = new Vuex.Store({
      state: {
        images: {
          data: [] as Image[],
          completed: true,
        },
      },
      actions,
    });

    const wrapper = shallowMount(Gallery, { store, localVue });
    const view = wrapper.vm as any;

    await view.scroll(scrollState);
    expect(scrollState.complete).toHaveBeenCalled();
  });
});
