<template>
  <section class="flex flex-1 flex-col gap-6">
    <header class="fluent-card sticky top-5 z-20 px-5 py-5 sm:px-7">
      <p class="text-[1.3rem] font-semibold uppercase tracking-[0.2em] text-accent-600">Today</p>
      <div class="mt-2">
        <h1 class="text-[3rem] font-bold tracking-tight text-slate-950 sm:text-[4rem]">今日一图</h1>
        <p class="mt-2 max-w-3xl text-[1.45rem] leading-7 text-slate-600">每天挑选一张适合桌面的图片，保持轻盈、沉浸的浏览体验。</p>
      </div>
    </header>

    <div class="grid flex-1 gap-6 xl:grid-cols-[minmax(0,1fr)_32rem]">
      <div class="fluent-card overflow-hidden p-3 sm:p-4">
        <ImageItem v-if="todayImage" :value="todayImage" />
        <div v-else class="flex min-h-[42rem] items-center justify-center rounded-3xl bg-white/60 text-[1.5rem] text-slate-500">
          暂无今日图片
        </div>
      </div>

      <aside class="fluent-card h-fit px-5 py-5 sm:px-6">
        <h2 class="text-[2rem] font-bold text-slate-950">壁纸状态</h2>
        <dl class="mt-5 space-y-4 text-[1.35rem]">
          <div class="rounded-2xl border border-white/90 bg-white/78 p-4 shadow-subtle ring-1 ring-sky-100/70">
            <dt class="font-semibold text-slate-500">当前壁纸源</dt>
            <dd class="mt-1 truncate font-semibold text-slate-900">{{ settings.currentSource || '未设置' }}</dd>
          </div>
          <div class="rounded-2xl border border-white/90 bg-white/78 p-4 shadow-subtle ring-1 ring-sky-100/70">
            <dt class="font-semibold text-slate-500">自动更新</dt>
            <dd class="mt-1 font-semibold" :class="settings.autoUpdate ? 'text-emerald-700' : 'text-slate-900'">
              {{ settings.autoUpdate ? `已开启${settings.timeToUpdate ? ` · ${settings.timeToUpdate}` : ''}` : '未开启' }}
            </dd>
          </div>
        </dl>
      </aside>
    </div>
  </section>
</template>

<script>
// @ts-check

import { defineComponent, inject, onMounted, ref, unref, watch } from 'vue'
import ImageItem from '@/components/ImageItem.vue'
import GlobalData from '@/injections/GlobalData'
import useApi from '@/composables/useApi'

export default defineComponent({
  components: { ImageItem },
  setup() {
    const { settings, setLoading, message } = inject(GlobalData)
    const { getTodayImage } = useApi()

    const todayImage = ref()

    const fetchImage = () => {
      setLoading(true)
      getTodayImage()
        .then(data => {
          todayImage.value = data
        })
        .catch(err => {
          message.error(err?.message || '获取今日图片失败')
        })
        .finally(() => {
          setLoading(false)
        })
    }
    onMounted(() => {
      fetchImage()
    })
    watch(() => unref(settings).currentSource, () => {
      fetchImage()
    })
    return {
      todayImage,
      settings,
    }
  },
})
</script>