<template>
  <AppWrap>
    <RouterView />
  </AppWrap>
</template>

<script>
import './common.css'
import { defineComponent } from 'vue'
import AppWrap from './components/AppWrap.vue'
import { RouterView } from 'vue-router'

export default defineComponent({
  components: { AppWrap, RouterView },
})
</script>

<style>
</style>
