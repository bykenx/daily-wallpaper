<template>
  <div class="More">
    <div class="More-gallery">
      <NGrid :cols="2">
        <NGi :key="index" v-for="(image, index) in archiveImages">
          <ImageItem :value="image" />
        </NGi>
      </NGrid>
    </div>
    <div class="More-gasket" />
  </div>
</template>

<script>
import { defineComponent, inject, onMounted, onUnmounted, reactive, watch } from 'vue'
import { NGi, NGrid, NImage, useMessage } from 'naive-ui'
import ImageItem from '@/components/ImageItem.vue'
import useImage from '@/composables/useImage'
import GlobalData from '@/injections/GlobalData'

export default defineComponent({
  components: { ImageItem, NImage, NGrid, NGi },
  setup() {
    const { settings, setLoading } = inject(GlobalData)
    const { getArchiveImages } = useImage()
    const message = useMessage()
    /** @type {IntersectionObserver} */
    let observer
    const data = reactive({
      items: [],
      pagination: { current: 0, pageSize: 8 },
      end: false,
    })
    const pagination = reactive({ current: 0, pageSize: 8 })
    const fetchImages = () => {
      setLoading(true)
      return getArchiveImages('/image/archive')
        .then(r => {
          const { current, pageSize, end, items } = r
          data.pagination = { current, pageSize }
          data.end = end
          data.items.push(...items)
        })
        .catch(err => {
          message.error(err)
        })
        .finally(() => {
          setLoading(false)
        })
    }
    watch(() => pagination.current, () => {
      if (!data.end) {
        fetchImages()
      }
    })
    onMounted(() => {
      fetchImages(settings.currentSource)
        .then(() => {
          observer = new IntersectionObserver((entries) => {
            if (entries[0].isIntersecting) {
              pagination.current += 1
            }
          }, {})
          observer.observe(document.querySelector('.More-gasket'))
        })
    })
    onUnmounted(() => {
      observer?.disconnect()
    })
    return {
      archiveImages: data.items,
    }
  },
})
</script>

<style lang="less" scoped>
.More {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  overflow-y: auto;

  &-gallery {
    width: 70vw;
  }

  &-gasket {
    width: 100%;
    flex: 0 0 1rem;
  }
}
</style>