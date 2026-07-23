<template>
  <div class="group relative overflow-hidden rounded-fluent bg-slate-200 shadow-subtle">
    <img class="image-item__image block aspect-[16/10] w-full object-cover transition duration-500" :src="imageSrc" alt="">
    <div class="absolute inset-0 bg-gradient-to-t from-slate-950/65 via-slate-950/10 to-transparent opacity-70 transition group-hover:opacity-100" />
    <div class="absolute inset-x-0 bottom-0 flex translate-y-3 items-end gap-4 p-4 opacity-0 transition duration-300 group-hover:translate-y-0 group-hover:opacity-100 sm:p-5">
      <div class="min-w-0 flex-1 text-[1.3rem] leading-6 text-white drop-shadow">
        {{ value.copyright }}
      </div>
      <button class="inline-flex shrink-0 items-center gap-2 rounded-2xl bg-white/90 px-4 py-2.5 text-[1.3rem] font-semibold text-slate-900 shadow-sm backdrop-blur transition hover:bg-white focus:outline-none focus:ring-4 focus:ring-white/40" @click="handleSetWallpaper">
        <IconWallpaper class="h-5 w-5" />
        设为壁纸
      </button>
    </div>
  </div>
</template>

<script>
import IconWallpaper from '@/assets/icon-wallpaper.svg'
import useApi from '@/composables/useApi'
import GlobalData from '@/injections/GlobalData'
import { computed, defineComponent, inject, unref } from 'vue'

export default defineComponent({
  components: { IconWallpaper },
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { settings, refreshSettings, message } = inject(GlobalData)
    const { updateSettings } = useApi()

    const imageSrc = computed(() => {
      const { url, urlHS } = props.value
      const { qualityFirst } = unref(settings)
      const link = qualityFirst && urlHS ? urlHS : url
      return `/api/image/get?link=${encodeURIComponent(link)}`
    })
    const handleSetWallpaper = () => {
      const { url, urlHS } = props.value
      const { qualityFirst } = unref(settings)
      const currentImage = qualityFirst && urlHS ? urlHS : url
      const data = {
        currentImage,
      }
      updateSettings(data)
        .then(data => {
          refreshSettings()
          message.success('壁纸设置成功')
        })
        .catch(err => {
          message.error(err?.message || '设置壁纸失败')
        })
    }
    return {
      imageSrc,
      handleSetWallpaper,
    }
  },
})
</script>

<style scoped>
.image-item__image {
  backface-visibility: hidden;
  transform: translateZ(0);
}

.group:hover .image-item__image {
  transform: translateZ(0) scale(1.03);
}
</style>