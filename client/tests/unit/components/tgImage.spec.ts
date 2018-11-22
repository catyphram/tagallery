import { shallowMount, createLocalVue } from '@vue/test-utils';
import TgImage from '@/components/TgImage.vue';

const localVue = createLocalVue();

describe('TgImage.vue', () => {
  it('should be instantiated', () => {
    const wrapper = shallowMount(TgImage, { localVue, propsData: {
      image: { file: 'test.jpg' },
    }});
    const view = wrapper.vm as any;
    expect(view).toBeTruthy();
  });
});
