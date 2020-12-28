<template>
  <v-app>
    <v-navigation-drawer v-model="drawer" app>
      <v-list dense>
        <v-subheader>Status</v-subheader>
        <v-list-item>
          <v-list-item-content>
            <v-btn-toggle v-model="status" mandatory>
              <v-btn small :value="Status.Processed">Processed</v-btn>
              <v-btn small :value="Status.Unprocessed">Unprocessed</v-btn>
            </v-btn-toggle>
          </v-list-item-content>
        </v-list-item>
        <template v-if="status === Status.Processed">
          <v-subheader>Categorization</v-subheader>
          <v-list-item>
            <v-list-item-content>
              <v-btn-toggle v-model="categorization" mandatory>
                <v-btn small :value="Categorization.Done">Done</v-btn>
                <v-btn small :value="Categorization.Open">Open</v-btn>
              </v-btn-toggle>
            </v-list-item-content>
          </v-list-item>
        </template>
        <template
          v-if="
            status === Status.Processed &&
            categorization === Categorization.Done
          "
        >
          <v-subheader>Mode</v-subheader>
          <v-list-item>
            <v-list-item-content>
              <v-btn-toggle v-model="mode" mandatory>
                <v-btn small :value="Mode.View">View</v-btn>
                <v-btn small :value="Mode.Verify">Verify</v-btn>
              </v-btn-toggle>
            </v-list-item-content>
          </v-list-item>
          <v-subheader>Categories</v-subheader>
          <v-list-item-group v-model="selectedCategories" multiple>
            <v-list-item
              v-for="(category, i) in categories"
              :key="i"
              :value="category"
            >
              <v-list-item-content>
                <v-list-item-title v-text="category.name"></v-list-item-title>
                <v-list-item-subtitle
                  v-text="category.description"
                ></v-list-item-subtitle>
              </v-list-item-content>
            </v-list-item>
          </v-list-item-group>
        </template>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app>
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Tagallery</v-toolbar-title>
      <v-spacer></v-spacer>

      <NuxtLink v-if="viewingSettings" to="/">
        <v-btn icon>
          <v-icon>mdi-image-multiple</v-icon>
        </v-btn>
      </NuxtLink>
      <NuxtLink v-else to="/settings">
        <v-btn icon>
          <v-icon>mdi-cog</v-icon>
        </v-btn>
      </NuxtLink>
    </v-app-bar>

    <v-main>
      <v-container fluid>
        <Nuxt />
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import type { ActionMethod } from 'vuex'
import { Component, namespace, Vue, Watch } from 'nuxt-property-decorator'
import { mapState } from 'vuex'

import { Category } from '~/store/category/state'
import { Status, Mode, Categorization } from '~/store/filter/state'
import { Image } from '~/store/image/state'

const filterModule = namespace('filter')
const imageModule = namespace('image')

@Component({
  computed: mapState('category', {
    categories: 'data',
  }),
})
export default class extends Vue {
  @filterModule.Action('set') setFilter!: ActionMethod
  @imageModule.Action('load') loadImages!: ActionMethod
  @imageModule.State('data') images!: Image[]

  Status = Status
  Mode = Mode
  Categorization = Categorization

  drawer = true

  get status(): Status {
    return this.$store.state.filter.status
  }

  set status(value) {
    this.setFilter({ status: value })
  }

  get mode(): Mode {
    return this.$store.state.filter.mode
  }

  set mode(value) {
    this.setFilter({ mode: value })
  }

  get categorization(): Categorization {
    return this.$store.state.filter.categorization
  }

  set categorization(value) {
    this.setFilter({ categorization: value })
  }

  get selectedCategories(): Category[] {
    return this.$store.state.filter.categories
  }

  set selectedCategories(value) {
    this.setFilter({ categories: value })
  }

  @Watch('$store.state.filter', {
    deep: true,
  })
  reloadImages() {
    this.loadImages({
      append: false,
      filter: this.$store.state.filter,
      lastImage: this.images.length
        ? this.images[this.images.length - 1].file
        : undefined,
    })
  }

  fetch() {
    this.$store.dispatch('category/load')
  }

  get viewingSettings() {
    return this.$route.path === '/settings'
  }
}
</script>

<style lang="scss" scoped>
.v-btn-toggle {
  width: 100%;
  .v-btn {
    width: 50%;
  }
}
</style>
