<template>
  <NConfigProvider :locale="zhCN">
    <NMessageProvider>
      <AppWrap>
        <RouterView />
      </AppWrap>
    </NMessageProvider>
  </NConfigProvider>
</template>

<script>
import './common.less'
import { defineComponent } from 'vue'
import AppWrap from './components/AppWrap.vue'
import { RouterView } from 'vue-router'
import { NConfigProvider, NMessageProvider, zhCN } from 'naive-ui'

export default defineComponent({
  components: { AppWrap, RouterView, NConfigProvider, NMessageProvider },
  setup() {
    return {
      zhCN,
    }
  },
})
</script>

<style>
</style>
