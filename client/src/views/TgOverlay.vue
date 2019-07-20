<template>
	<div class="tg-overlay" v-if="image">
    <md-icon @click.native="prevImage" class="md-size-2x">arrow_back</md-icon>
		<img :src="image" alt="image" />
    <md-icon @click.native="nextImage" class="md-size-2x">arrow_forward</md-icon>
    <md-icon @click.native="close" class="md-size-2x tg-overlay__close">close</md-icon>

	</div>
</template>

<script lang="ts">
import { Component, Watch, Vue } from 'vue-property-decorator';

import { Image } from '../models';

@Component
export default class TgOverlay extends Vue {
  get image(): Image {
    return this.$store.state.selectedImage !== undefined &&
      this.$store.state.images.data[this.$store.state.selectedImage] ?
      this.$store.state.images.data[this.$store.state.selectedImage].file :
      undefined;
  }
  public prevImage() {
    if (this.$store.state.selectedImage) {
      this.$store.commit('selectImage', { index: this.$store.state.selectedImage - 1 });
    }
  }
  public nextImage() {
    if (this.$store.state.selectedImage !== undefined) {
      this.$store.commit('selectImage', { index: this.$store.state.selectedImage + 1 });
    }
  }
  public close() {
    this.$store.commit('selectImage', { index: undefined });
  }
}
</script>

<style lang="scss">
	.tg-overlay {
    background: #424242;
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    right: 0;
    z-index: 100;
    display: flex;
    justify-content: center;
    align-items: center;

    &__close {
      position: absolute;
      right: 0;
      top: 0;
    }
	}
</style>
