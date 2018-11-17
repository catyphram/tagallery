<template>
  <div class="tg-menu">
    <form>
      <md-list>
        <md-list-item class="mode-toggle">
          <md-switch v-model="verifyMode" v-on:change="updateMode">Verify Mode</md-switch>
        </md-list-item>
        <md-list-item v-for="category in categories" :key="category.key"
          v-if="!categoriesLoading && categories.length > 0">
          <md-button
            v-on:click="$store.dispatch('toggleCategory', category)"
            v-bind:class="{'md-raised': $store.getters.isCategorySelected(category)}"
          >{{category.name}}</md-button>
        </md-list-item>
        <md-list-item>
          <md-button
            v-on:click="$store.dispatch('toggleListUncategorized')"
            v-bind:class="{'md-raised': $store.state.listUncategorized}"
          >Uncategorized</md-button>
        </md-list-item>
      </md-list>
    </form>
  </div>
</template>

<script lang="ts">
import { Component, Watch, Vue } from 'vue-property-decorator';

import { LIST_MODE } from '@/models';

@Component
export default class TgMenu extends Vue {
  public verifyMode = false;
  get categories() {
    return this.$store.state.categories.data;
  }
  get categoriesLoading() {
    return this.$store.state.categories.loading;
  }
  public updateMode(value: boolean) {
    this.$store.dispatch('setMode', value ? LIST_MODE.MODE_VERIFY : LIST_MODE.MODE_VIEW);
  }
  @Watch('$store.state.listMode', { immediate: true })
  public watchMode(newMode: LIST_MODE, prevMode: LIST_MODE) {
    this.verifyMode = newMode === LIST_MODE.MODE_VERIFY;
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
