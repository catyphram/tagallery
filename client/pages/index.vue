<template>
  <v-row dense>
    <v-col
      v-for="(image, index) in images"
      :key="image.file"
      v-intersect.once="
        (e, o, intersecting) =>
          intersecting && isIntersectImage(index) && getImages(true)
      "
      :cols="12"
      :sm="4"
      :md="3"
      :xl="2"
    >
      <v-img :src="image.file" :lazy-src="image.file" aspect-ratio="1">
        <template v-slot:placeholder>
          <v-row class="fill-height ma-0" align="center" justify="center">
            <v-progress-circular indeterminate></v-progress-circular>
          </v-row>
        </template>
      </v-img>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { Component, namespace, Vue } from 'nuxt-property-decorator'
import type { ActionMethod } from 'vuex'

import { Image } from '~/store/image/state'

const imageModule = namespace('image')

@Component
export default class extends Vue {
  @imageModule.Action('load') loadImages!: ActionMethod
  @imageModule.State('data') images!: Image[]

  fetch() {
    this.getImages()
  }

  getImages(append = false) {
    this.loadImages({
      append,
      filter: this.$store.state.filter,
      lastImage: this.images.length
        ? this.images[this.images.length - 1].file
        : undefined,
    })
  }

  isIntersectImage(index: number) {
    if (this.images.length <= 12) {
      return index + 1 === this.images.length
    }
    return index + 1 === this.images.length - 12
  }
}
</script>
