<template>
  <section class="flex flex-1 flex-col gap-6">
    <header class="fluent-card sticky top-5 z-20 px-5 py-5 sm:px-7">
      <p class="text-[1.3rem] font-semibold uppercase tracking-[0.2em] text-accent-600">Settings</p>
      <h1 class="mt-2 text-[3rem] font-bold tracking-tight text-slate-950 sm:text-[4rem]">设置</h1>
      <p class="mt-2 max-w-3xl text-[1.45rem] leading-7 text-slate-600">管理图片源、自动更新和桌面壁纸偏好。</p>
    </header>

    <form class="grid gap-6" @submit.prevent="handleSubmit">
      <div class="fluent-card space-y-6 px-5 py-6 sm:px-7">
        <div class="setting-panel grid gap-3">
          <label class="fluent-label" for="currentSource">图片源</label>
          <select id="currentSource" v-model="model.currentSource" class="fluent-field">
            <option value="" disabled>请选择一个图片源</option>
            <option
              v-for="item in sourceOptions"
              :key="item.name"
              :value="item.name"
              :disabled="!item.enabled"
            >
              {{ item.description }}{{ item.enabled ? '' : '（不可用）' }}
            </option>
          </select>
        </div>

        <div class="setting-panel grid gap-3">
          <label class="fluent-label" for="currentImage">当前图片</label>
          <input
            id="currentImage"
            v-model="model.currentImage"
            class="fluent-field"
            :disabled="model.autoUpdate"
            :placeholder="model.autoUpdate ? '开启每天更新桌面时此项无效' : '贴入链接或选择图片'"
          >
        </div>

        <div class="grid gap-4 sm:grid-cols-2">
          <label class="setting-panel flex cursor-pointer items-center justify-between gap-4">
            <span>
              <span class="block text-[1.4rem] font-semibold text-slate-900">开机自启</span>
              <span class="mt-1 block text-[1.2rem] text-slate-500">启动系统后自动运行</span>
            </span>
            <input v-model="model.autoRunAtSystemBoot" class="peer sr-only" type="checkbox">
            <span class="relative h-8 w-14 rounded-full bg-slate-200 transition peer-checked:bg-accent-600 after:absolute after:left-1 after:top-1 after:h-6 after:w-6 after:rounded-full after:bg-white after:shadow-sm after:transition peer-checked:after:translate-x-6" />
          </label>

          <label class="setting-panel flex cursor-pointer items-center justify-between gap-4">
            <span>
              <span class="block text-[1.4rem] font-semibold text-slate-900">每天自动更新</span>
              <span class="mt-1 block text-[1.2rem] text-slate-500">按时间切换桌面</span>
            </span>
            <input v-model="model.autoUpdate" class="peer sr-only" type="checkbox">
            <span class="relative h-8 w-14 rounded-full bg-slate-200 transition peer-checked:bg-accent-600 after:absolute after:left-1 after:top-1 after:h-6 after:w-6 after:rounded-full after:bg-white after:shadow-sm after:transition peer-checked:after:translate-x-6" />
          </label>
        </div>

        <div v-if="model.autoUpdate" class="setting-panel grid gap-3">
          <label class="fluent-label" for="timeToUpdate">更新时间</label>
          <input id="timeToUpdate" v-model="model.timeToUpdate" class="fluent-field" type="time">
        </div>

        <label class="setting-panel flex cursor-pointer items-center justify-between gap-4">
          <span>
            <span class="block text-[1.4rem] font-semibold text-slate-900">更高品质</span>
            <span class="mt-1 block text-[1.2rem] text-slate-500">优先使用高清图片链接</span>
          </span>
          <input v-model="model.qualityFirst" class="peer sr-only" type="checkbox">
          <span class="relative h-8 w-14 rounded-full bg-slate-200 transition peer-checked:bg-accent-600 after:absolute after:left-1 after:top-1 after:h-6 after:w-6 after:rounded-full after:bg-white after:shadow-sm after:transition peer-checked:after:translate-x-6" />
        </label>

        <button class="fluent-button w-full sm:w-auto" type="submit">保存设置</button>
      </div>
    </form>
  </section>
</template>

<script>
// @ts-check

import { defineComponent, inject, onMounted, ref, unref, watch } from 'vue'
import GlobalData from '@/injections/GlobalData'
import useApi from '@/composables/useApi'

const defaultModel = {
  autoUpdate: false,
  currentImage: '',
  currentSource: '',
  autoRunAtSystemBoot: false,
  timeToUpdate: '',
  qualityFirst: false,
}

const normalizeTimeForInput = (value) => {
  if (!value) {
    return ''
  }
  if (typeof value === 'number') {
    const time = new Date(value)
    return `${time.getHours().toString().padStart(2, '0')}:${time.getMinutes().toString().padStart(2, '0')}`
  }
  return String(value).slice(0, 5)
}

export default defineComponent({
  setup() {
    const { settings, refreshSettings, setLoading, message } = inject(GlobalData)
    const { getAllSources, updateSettings } = useApi()
    const model = ref({ ...defaultModel })
    const sourceOptions = ref([])

    const setFormModel = (value) => {
      const nextValue = { ...defaultModel, ...value }
      nextValue.timeToUpdate = normalizeTimeForInput(nextValue.timeToUpdate)
      model.value = nextValue
    }

    const fetchSources = () => {
      setLoading(true)
      return getAllSources()
        .then(data => {
          sourceOptions.value = data
        })
        .catch(err => {
          message.error(err?.message || '获取图片源失败')
        })
        .finally(() => {
          setLoading(false)
        })
    }

    const handleSubmit = () => {
      const value = { ...unref(model) }
      value.timeToUpdate = value.autoUpdate && value.timeToUpdate ? `${value.timeToUpdate}:00` : null
      updateSettings(value)
        .then(() => {
          refreshSettings()
          message.success('设置更新成功')
        })
        .catch(err => {
          message.error(err?.message || '设置更新失败')
        })
    }

    setFormModel(unref(settings))
    watch(settings, value => {
      setFormModel(value)
    })
    onMounted(() => {
      fetchSources()
    })

    return {
      model,
      sourceOptions,
      handleSubmit,
    }
  },
})
</script>

<style scoped>
.setting-panel {
  padding: 1rem;
  border: 1px solid rgb(255 255 255 / 0.9);
  border-radius: 1.5rem;
  background: rgb(255 255 255 / 0.78);
  box-shadow: 0 0 0 1px rgb(224 242 254 / 0.7), 0 12px 32px rgb(14 116 144 / 0.12);
  backdrop-filter: blur(24px);
}

@media (min-width: 640px) {
  .setting-panel {
    padding: 1.25rem;
  }
}
</style>
