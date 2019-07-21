<template>
	<div class="tg-overlay" v-if="image">
    <div class="tg-overlay__menu">
      <span v-if="categories" class="tg-overlay__categories">
        <md-button
          v-for="categoryItem in categories" :key="categoryItem.category.key"
          @click="() => toggleImageCategory(categoryItem.category)"
          :class="{
            'md-raised': categoryItem.assigned,
            'tg-overlay__category--proposed': categoryItem.proposed,
          }"
        >{{categoryItem.category.name}}</md-button>
      </span>
      <md-icon @click.native="close" class="md-size-2x tg-overlay__close">close</md-icon>
    </div>
    <md-icon @click.native="prevImage" class="md-size-2x">arrow_back</md-icon>
		<img :src="image.file" alt="image" />
    <md-icon @click.native="nextImage" class="md-size-2x">arrow_forward</md-icon>
	</div>
</template>

<script lang="ts">
import { Component, Watch, Vue } from 'vue-property-decorator';

import { Image, Category } from '../models';

@Component
export default class TgOverlay extends Vue {
  get image(): Image {
    return this.$store.state.selectedImage !== undefined &&
      this.$store.state.images.data[this.$store.state.selectedImage] ?
      this.$store.state.images.data[this.$store.state.selectedImage] :
      undefined;
  }
  get categories(): Array<{
    category: Category,
    assigned: boolean,
    proposed: boolean,
  }> | undefined {
    if (this.image !== undefined) {
      return this.$store.state.categories.data.map((category: Category) => {
        return {
          category,
          assigned: this.imageHasAssignedCategory(category),
          proposed: this.imageHasProposedCategory(category),
        };
      });
    }
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
  public toggleImageCategory(category) {
    this.$store.dispatch('toggleImageCategory', { image: this.image, category });
  }
  protected imageHasAssignedCategory(category: Category) {
    return !!(this.image.assignedCategories && this.image.assignedCategories.find((assignedCategory) => {
      return assignedCategory === category.key;
    }));
  }
  protected imageHasProposedCategory(category: Category) {
    return !!(this.image.proposedCategories && this.image.proposedCategories.find((proposedCategory) => {
      return proposedCategory === category.key;
    }));
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

    &__menu {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      display: inline-flex;
    }

    &__categories {
      flex: 1;
      display: inline-flex;
      justify-content: center;
      margin-left: 48px;
    }
    &__category {
      &--proposed {
        color: deepskyblue !important;
      }
    }

    &__close {
      flex: 1;
      display: inline-flex;
      justify-content: flex-end;
      flex-grow: 0;
    }
	}
</style>
