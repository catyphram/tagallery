<template>
  <div class="tg-menu">
    <form>
      <md-list>
        <md-list-item class="mode-toggle">
          <md-switch v-model="verifyMode" v-on:change="updateMode">Verify Mode</md-switch>
        </md-list-item>
        <template v-if="!categoriesLoading && categories.length > 0">
          <md-list-item v-for="category in categories" :key="category.key">
            <md-button
              v-on:click="$store.dispatch('toggleCategory', category)"
              :class="{'md-raised': $store.getters.isCategorySelected(category)}"
            >{{category.name}}</md-button>
          </md-list-item>
        </template>
        <md-list-item>
          <md-button
            v-on:click="$store.dispatch('toggleListUncategorized')"
            :class="{'md-raised': $store.state.listUncategorized}"
          >Uncategorized</md-button>
        </md-list-item>
      </md-list>
    </form>
  </div>
</template>

<script lang="ts">
import { Component, Watch, Vue } from 'vue-property-decorator';

import { LIST_MODE, Image, Category } from '../models';

@Component
export default class TgMenu extends Vue {
  get verifyMode(): boolean {
    return this.$store.state.listMode === LIST_MODE.MODE_VERIFY;
  }
  get categories(): Category[] {
    return this.$store.state.categories.data;
  }
  get categoriesLoading(): boolean {
    return this.$store.state.categories.loading;
  }
  public updateMode(value: boolean) {
    this.$store.dispatch('setMode', value ? LIST_MODE.MODE_VERIFY : LIST_MODE.MODE_VIEW);
  }
}
</script>

<style scoped lang="scss">
  .md-button {
    display: block;
  }
  .mode-toggle {
    margin-left: 16px;
  }
</style>
