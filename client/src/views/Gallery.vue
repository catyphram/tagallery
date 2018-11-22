<template>
  <div class="gallery">
    <tg-image :image="image" v-for="image in images" :key="image.file" />
    <infinite-loading @infinite="scroll"></infinite-loading>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

import TgImage from '@/components/TgImage.vue';

@Component({
  components: {
    TgImage,
  },
})
export default class Gallery extends Vue {
  get images() {
    return this.$store.state.images.data;
  }
  public async scroll($state) {
    await this.$store.dispatch('loadImages', true);
    if (this.$store.state.images.completed) {
      $state.complete();
    } else {
      $state.loaded();
    }
  }
}
</script>

<style scoped lang="scss">
  .gallery {
		display: flex;
		flex-flow: row wrap;
		margin: -5px;
  }
</style>
