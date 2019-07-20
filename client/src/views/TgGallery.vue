<template>
  <div class="tg-gallery">
    <tg-image :index="index" :image="image" v-for="(image, index) in images" :key="image.file" ref="images"/>
    <infinite-loading @infinite="scroll"></infinite-loading>
  </div>
</template>

<script lang="ts">
import { Component, Watch, Vue } from 'vue-property-decorator';

import TgImage from '@/components/TgImage.vue';
import { Image, Category } from '../models';

@Component({
  components: {
    TgImage,
  },
})
export default class TgGallery extends Vue {
  get images(): Image[] {
    return this.$store.state.images.data;
  }
  public async scroll($state) {
    if (!this.$store.state.images.loading) {
      await this.$store.dispatch('loadImages', true);
      if (this.$store.state.images.completed) {
        $state.complete();
      } else {
        $state.loaded();
      }
    }
  }
  @Watch('$store.state.selectedImage')
  @Watch('$store.state.images.data')
  public watchSelectedImage() {
    if (this.$store.state.selectedImage !== undefined && this.$refs.images) {
      this.$refs.images[this.$store.state.selectedImage].$el.scrollIntoView({behavior: 'smooth'});
    }
  }
}
</script>

<style scoped lang="scss">
  .tg-gallery {
		display: flex;
		flex-flow: row wrap;
		margin: -5px;
  }
</style>
