<template>
  <section class="flex flex-1 flex-col gap-6">
    <header class="fluent-card sticky top-5 z-20 px-5 py-5 sm:px-7">
      <p class="text-[1.3rem] font-semibold uppercase tracking-[0.2em] text-accent-600">Gallery</p>
      <h1 class="mt-2 text-[3rem] font-bold tracking-tight text-slate-950 sm:text-[4rem]">更多图片</h1>
      <p class="mt-2 max-w-3xl text-[1.45rem] leading-7 text-slate-600">继续浏览历史推荐，找到更适合桌面的那一张。</p>
    </header>

    <div class="grid gap-5 md:grid-cols-2 2xl:grid-cols-3">
      <ImageItem
        v-for="(image, index) in data.items"
        :key="index"
        :value="image"
      />
    </div>

    <div v-if="data.items.length === 0" class="fluent-card flex min-h-[24rem] items-center justify-center text-[1.5rem] text-slate-500">
      暂无图片
    </div>

    <div class="h-8 shrink-0 More-gasket" />
  </section>
</template>

<script>
// @ts-check

import { defineComponent, inject, onMounted, onUnmounted, reactive, unref, watch } from 'vue'
import ImageItem from '@/components/ImageItem.vue'
import GlobalData from '@/injections/GlobalData'
import useApi from '@/composables/useApi'
import useToggle from '@/composables/useToggle'

export default defineComponent({
  components: { ImageItem },
  setup() {
    /** @type {IntersectionObserver} */
    let observer

    const { settings, setLoading, message } = inject(GlobalData)
    const { getArchiveImages } = useApi()
    const [loadMore, setLoadMore] = useToggle()

    const data = reactive({
      items: [],
      pagination: { current: 0, pageSize: 8 },
      end: true,
    })
    const fetchImages = () => {
      setLoading(true)
      return getArchiveImages(data.pagination)
        .then(r => {
          const { current, pageSize, end, items } = r
          data.pagination = { current, pageSize }
          data.end = end
          data.items.push(...items)
        })
        .catch(err => {
          message.error(err?.message || '获取更多图片失败')
        })
        .finally(() => {
          setLoading(false)
        })
    }
    watch(loadMore, () => {
      fetchImages()
    })
    watch(() => unref(settings).currentSource, () => {
      data.items = []
      data.pagination.current = 0
      data.end = true
      setLoadMore()
    })
    onMounted(() => {
      fetchImages()
        .then(() => {
          observer = new IntersectionObserver((entries) => {
            if (entries[0].isIntersecting && !data.end) {
              data.pagination.current += 1
              setLoadMore()
            }
          }, {})
          observer.observe(document.querySelector('.More-gasket'))
        })
    })
    onUnmounted(() => {
      observer?.disconnect()
    })
    return {
      data
    }
  },
})
</script>