<template>
  <div id="app">
    <md-app md-waterfall md-mode="fixed">
      <md-app-toolbar class="md-primary">
        <md-button class="md-icon-button" @click="toggleMenu" v-if="!menuVisible">
          <md-icon>menu</md-icon>
        </md-button>
        <span class="md-title">{{title}}</span>
      </md-app-toolbar>

      <md-app-drawer :md-active.sync="menuVisible" md-persistent="full">
        <md-toolbar class="md-transparent" md-elevation="0">
          <span class="app-name">Tagallery</span>
          <div class="md-toolbar-section-end">
            <md-button class="md-icon-button md-dense" @click="toggleMenu">
              <md-icon>keyboard_arrow_left</md-icon>
            </md-button>
          </div>
        </md-toolbar>
        <tg-menu />
      </md-app-drawer>

      <md-app-content>
        <gallery />
      </md-app-content>
    </md-app>
  </div>
</template>

<style lang="scss" scoped>
.md-drawer {
  width: 260px;
}
.md-app {
  height: 100vh;
}
.app-name {
  padding-left: 16px;
}
</style>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Gallery from '@/views/Gallery.vue';
import TgMenu from '@/components/TgMenu.vue';
import { LIST_MODE } from '@/models';

@Component({
  components: {
    Gallery,
    TgMenu,
  },
})
export default class App extends Vue {
  get title() {
    const selectedCategories = this.$store.state.selectedCategories;
    const modeLabel = (this.$store.state.listMode === LIST_MODE.MODE_VERIFY ? 'Verifying' : 'Viewing') as string;
    const uncategorizedLabel = (this.$store.state.listUncategorized ? 'uncategorized ' : '') + 'images' as string;
    let categoryLabel = '';

    if (selectedCategories.length === 1) {
      categoryLabel = `from category ${selectedCategories[0].name}`;
    } else if (selectedCategories.length > 1) {
      const categoryNames = selectedCategories.map((category) => category.name);
      categoryLabel = `from categories ${categoryNames.join(', ')}`;
    }

    return `${modeLabel} ${uncategorizedLabel}` + (categoryLabel ? ' ' + categoryLabel : '');
  }
  public menuVisible = false;
  constructor() {
    super();
    this.$store.dispatch('loadCategories');
  }
  public toggleMenu() {
    this.menuVisible = !this.menuVisible;
  }
}
</script>