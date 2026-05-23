<template>
  <section class="flex flex-1 flex-col gap-6">
    <header class="fluent-card sticky top-5 z-20 px-5 py-5 sm:px-7">
      <p class="text-[1.3rem] font-semibold uppercase tracking-[0.2em] text-accent-600">Local</p>
      <h1 class="mt-2 text-[3rem] font-bold tracking-tight text-slate-950 sm:text-[4rem]">本地图片</h1>
      <p class="mt-2 max-w-3xl text-[1.45rem] leading-7 text-slate-600">浏览本地可用图片，并快速设置为桌面壁纸。</p>
    </header>

    <div class="grid gap-5 md:grid-cols-2 2xl:grid-cols-3">
      <ImageItem
        v-for="(item, index) in data.items"
        :key="index"
        :value="{ url: item }"
      />
    </div>

    <div v-if="data.items.length === 0" class="fluent-card flex min-h-[24rem] items-center justify-center text-[1.5rem] text-slate-500">
      暂无本地图片
    </div>

    <div class="h-8 shrink-0 Local-gasket" />
  </section>
</template>

<script>

// @ts-check
import { defineComponent, inject, onMounted, onUnmounted, reactive, watch } from 'vue'
import ImageItem from '@/components/ImageItem.vue'
import useApi from '@/composables/useApi'
import useToggle from '@/composables/useToggle'
import GlobalData from '@/injections/GlobalData'

export default defineComponent({
  components: {
    ImageItem,
  },
  setup() {
    /** @type {IntersectionObserver} */
    let observer

    const { setLoading, message } = inject(GlobalData)
    const [loadMore, setLoadMore] = useToggle()
    const { getImageList } = useApi()
    const data = reactive({
      items: [],
      pagination: { start: 0, limit: 8 },
      end: false,
    })

    const fetchImages = async () => {
      setLoading(true)
      try {
        try {
          const items = await getImageList(data.pagination)
          data.items.push(...items)
          data.end = items.length < data.pagination.limit
        } catch (err) {
          message.error(err?.message || '获取图片列表失败')
        }
      } finally {
        setLoading(false)
      }
    }

    watch(loadMore, () => {
      fetchImages()
    })

    onMounted(() => {
      fetchImages()
        .then(() => {
          observer = new IntersectionObserver((entries) => {
            if (entries[0].isIntersecting && !data.end) {
              data.pagination.start += 1
              setLoadMore()
            }
          }, {})
          observer.observe(document.querySelector('.Local-gasket'))
        })
    })
    onUnmounted(() => {
      observer?.disconnect()
    })

    return {
      data,
    }
  },
})
</script>