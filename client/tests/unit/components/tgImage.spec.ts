import { shallowMount, createLocalVue } from '@vue/test-utils';
import Vuex from 'vuex';
import TgImage from '@/components/TgImage.vue';

const localVue = createLocalVue();
localVue.use(Vuex);

const state = {
  selectedImage: undefined,
};

const mutations = {
  selectImage: jest.fn(),
};

const store = new Vuex.Store({
  state,
  mutations,
});

describe('TgImage.vue', () => {
  it('should be instantiated', () => {
    const wrapper = shallowMount(TgImage, { localVue, propsData: {
      image: { file: 'test.jpg' },
    }});
    const view = wrapper.vm as any;
    expect(view).toBeTruthy();
  });

  it('openOverlay should open the overlay', () => {
    const wrapper = shallowMount(TgImage, { store, localVue, propsData: {
      image: { file: 'test.jpg' },
      index: 1,
    }});
    const view = wrapper.vm as any;

    view.openOverlay(true);

    expect(mutations.selectImage).toHaveBeenCalledWith(
      expect.anything(),
      expect.objectContaining({ index: 1 }),
    );
  });
});
