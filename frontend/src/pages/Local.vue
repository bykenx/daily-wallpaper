<template>
  <div class="Local">
    <NH2>正在施工</NH2>
  </div>
</template>

<script>
import { defineComponent } from 'vue'
import { NH2 } from 'naive-ui'

export default defineComponent({
  components: {
    NH2,
  },
  setup() {
    return {}
  },
})
</script>

<style scoped>
.Local {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>