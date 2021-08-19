<template>
  <div class="More">
    <div class="More-gallery">
      <NGrid :cols="2">
        <NGi :key="index" v-for="(image, index) in data.items">
          <ImageItem :value="image" />
        </NGi>
      </NGrid>
    </div>
    <div class="More-gasket" />
  </div>
</template>

<script>
import { defineComponent, inject, onMounted, onUnmounted, reactive, unref, watch } from 'vue'
import { NGi, NGrid, NImage, useMessage } from 'naive-ui'
import ImageItem from '@/components/ImageItem.vue'
import GlobalData from '@/injections/GlobalData'
import useApi from '@/composables/useApi'
import useToggle from '@/composables/useToggle'

export default defineComponent({
  components: { ImageItem, NImage, NGrid, NGi },
  setup() {
    /** @type {IntersectionObserver} */
    let observer

    const { settings, setLoading } = inject(GlobalData)
    const { getArchiveImages } = useApi()
    const message = useMessage()
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
          message.error(err)
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