<template>
  <div class="Local">
    <div class="Local-list">
      <template v-for="item, index in data.items" :key="index">
        <ImageItem class="Local-list-item" :value="{ url: item }" />
      </template>
    </div>
    <div class="Local-gasket" />
  </div>
</template>

<script>

// @ts-check
import { defineComponent, inject, onMounted, onUnmounted, reactive, ref, unref, watch } from 'vue'
import { NH2, useMessage } from 'naive-ui'
import ImageItem from '@/components/ImageItem.vue'
import useApi from '@/composables/useApi'
import useToggle from '@/composables/useToggle'
import GlobalData from '@/injections/GlobalData'

export default defineComponent({
  components: {
    NH2,
    ImageItem,
  },
  setup() {
    /** @type {IntersectionObserver} */
    let observer

    const { setLoading } = inject(GlobalData)
    const message = useMessage()
    const [loadMore, setLoadMore] = useToggle()
    const { getImageList } = useApi()
    const data = reactive({
      items: [],
      pagination: { start: 0, limit: 8 },
      end: false,
    })

    const fetchImages = () => {
      setLoading(true)
      return getImageList(data.pagination)
        .then(items => {
          data.items.push(...items)
          data.end = items.length < data.pagination.limit
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
    onMounted(() => {
      fetchImages()
        .then(() => {
          observer = new IntersectionObserver((entries) => {
            console.log('load')
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

<style lang="less" scoped>
.Local {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  overflow-y: auto;

  &-list {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
  }

  &-gasket {
    width: 100%;
    flex: 0 0 1rem;
  }
}
</style>