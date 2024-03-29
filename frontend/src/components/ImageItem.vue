<template>
  <div class="ImageItem">
    <img :src="`/api/image/get?link=${encodeURIComponent(value.url)}`" alt="">
    <div class="ImageItem-mask" />
    <div class="ImageItem-toolbar">
      <div class="ImageItem-toolbar-desc">
        {{ value.copyright }}
      </div>
      <NButton text @click="handleSetWallpaper">
        <template #icon>
          <NIcon size="20">
            <IconWallpaper />
          </NIcon>
        </template>
        设为壁纸
      </NButton>
    </div>
  </div>
</template>

<script>
import IconDownload from '@/assets/icon-download.svg'
import IconWallpaper from '@/assets/icon-wallpaper.svg'
import useApi from '@/composables/useApi'
import GlobalData from '@/injections/GlobalData'
import { NButton, NIcon, useMessage } from 'naive-ui'
import { defineComponent, inject, unref } from 'vue'

export default defineComponent({
  components: { NIcon, IconWallpaper, NButton, IconDownload },
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { settings, refreshSettings } = inject(GlobalData)
    const message = useMessage()
    const { updateSettings } = useApi()

    const handleSetWallpaper = () => {
      const { url, urlHS } = props.value
      const { qualityFirst } = unref(settings)
      const data = {
        currentImage: qualityFirst && urlHS ? urlHS : url,
      }
      updateSettings(data)
        .then(data => {
          refreshSettings()
          message.success('壁纸设置成功')
        })
        .catch(err => {
          message.error(err)
        })
    }
    return {
      handleSetWallpaper,
    }
  },
})
</script>

<style lang="less">
.ImageItem {
  position: relative;
  overflow: hidden;

  &:hover {
    .ImageItem {
      &-mask {
        background: rgba(0, 0, 0, 0.3);
      }

      &-toolbar {
        transform: translateY(0);
      }
    }
  }

  &-img {
    display: block;
    width: 100%;
  }

  img {
    width: 100%;
    display: block;
  }

  &-mask {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    transition: background 0.3s cubic-bezier(0.4, 0, 0.2, 1) 0s;
  }

  &-toolbar {
    transform: translateY(100%);
    transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1) 0s;
    display: flex;
    justify-content: flex-end;
    align-items: center;
    position: absolute;
    bottom: 0;
    width: 100%;
    padding: 0 1rem;
    box-sizing: border-box;
    height: 4rem;
    background: rgba(255, 255, 255, 0.7);
    gap: 1rem;

    &-desc {
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
}
</style>