<template>
  <div>
    <v-form ref="editForm" class="pt-3">
      <v-row>
        <template v-for="(category, i) in categories">
          <v-col :key="`name-${category.id}`" cols="12" sm="6">
            <v-text-field
              v-model="category.name"
              dense
              :rules="nameRules"
              required
              label="Name"
              @focusout="() => syncCategory(i)"
            ></v-text-field>
          </v-col>
          <v-col :key="`desc-${category.id}`" cols="12" sm="6">
            <v-text-field
              v-model="category.description"
              append-outer-icon="mdi-delete"
              dense
              :rules="nameRules"
              label="Description"
              @focusout="() => syncCategory(i)"
              @click:append-outer="() => deleteCategory(i)"
            ></v-text-field>
          </v-col>
        </template>
      </v-row>
    </v-form>
    <v-form ref="createForm">
      <v-row>
        <v-col cols="12" sm="6">
          <v-text-field
            v-model.trim="newCategory.name"
            dense
            :rules="nameRules"
            required
            label="Name"
          ></v-text-field>
        </v-col>
        <v-col cols="12" sm="6">
          <v-text-field
            v-model.trim="newCategory.description"
            append-outer-icon="mdi-check"
            dense
            :rules="nameRules"
            label="Description"
            @click:append-outer="createCategory"
          ></v-text-field>
        </v-col>
      </v-row>
    </v-form>
  </div>
</template>

<script lang="ts">
import { Component, namespace, Vue, Watch } from 'nuxt-property-decorator'
import type { ActionMethod } from 'vuex'

import { Category } from '~/store/category/state'

const categoryModule = namespace('category')

@Component
export default class extends Vue {
  @categoryModule.Action('update') updateCategory!: ActionMethod
  @categoryModule.Action('add') addCategory!: ActionMethod

  categories = [] as Category[]
  newCategory = {
    name: '',
    description: '',
  }

  nameRules = [
    (v: string) => v.trim().length > 0 || 'Field is required',
    (v: string) => this.isCategoryUnique(v) || 'Name must be unique',
  ]

  isCategoryUnique(cn: string) {
    return this.categories.filter((c) => c.name === cn).length <= 1
  }

  @Watch('$store.state.category.data', {
    deep: true,
    immediate: true,
  })
  syncCategories(newValue: Category[]) {
    this.categories = JSON.parse(JSON.stringify(newValue))
  }

  // Categories are only synced when the entire form is valid and neiher the name,
  // nor the description of the edited category are currently in focus.
  syncCategory(index: number) {
    const formIsValid = this.validateForm(this.$refs.editForm as Vue)

    if (formIsValid) {
      const categoryNameInput = (this.$refs.editForm as Vue).$children[
        index * 2
      ].$refs.input
      const categoryDescInput = (this.$refs.editForm as Vue).$children[
        index * 2 + 1
      ].$refs.input
      setTimeout(() => {
        if (
          window.document.activeElement !== categoryNameInput &&
          window.document.activeElement !== categoryDescInput
        ) {
          const stateCategories = this.$store.state.category.data
          const updatedCategories = this.categories.filter(
            (category, index) =>
              category.name !== stateCategories[index].name ||
              category.description !== stateCategories[index].description
          )
          updatedCategories.forEach((category) =>
            this.updateCategory({ category })
          )
        }
      })
    }
  }

  validateForm(form: Vue) {
    return (form as Vue & { validate: () => boolean }).validate()
  }

  deleteCategory(index: number) {
    this.categories.splice(index, 1)
    this.validateForm(this.$refs.editForm as Vue)
  }

  createCategory() {
    if (this.validateForm(this.$refs.createForm as Vue)) {
      this.addCategory({
        category: {
          id: `${Math.random() * 10e15}`,
          name: this.newCategory.name,
          description: this.newCategory.description,
        },
      })
      this.newCategory.name = ''
      this.newCategory.description = ''
    }
  }
}
</script>
