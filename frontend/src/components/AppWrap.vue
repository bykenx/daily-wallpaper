<template>
  <div class="AppWrap" :style="appWrapStyleObject">
    <div class="AppWrap-header">
      <div>
        <div class="AppWrap-header-sources">
          <NSelect
            placeholder="请选择一个图片源"
            :value="settings.currentSource"
            :options="options"
            :on-update:value="handleSwitchSource"
          />
        </div>
        <NMenu
          mode="horizontal"
          :value="currentMenuItem"
          :options="menuItems"
          :render-label="renderMenuLabel"
        />
      </div>
    </div>
    <div class="AppWrap-content">
      <div v-if="loading" class="AppWrap-content-loading">
        <NSpin />
      </div>
      <slot v-if="!firstLoad" />
    </div>
  </div>
</template>

<script>
import { computed, defineComponent, h, onMounted, provide, ref, unref, watch } from 'vue'
import { NMenu, NSelect, NSpin, useMessage, zhCN } from 'naive-ui'
import GlobalData from '@/injections/GlobalData'
import useApi from '@/composables/useApi'
import { RouterLink, useRoute } from 'vue-router'

export default defineComponent({
  components: {
    NMenu,
    NSelect,
    NSpin,
  },
  setup() {
    const menuItems = [
      { label: '今日一图', key: 'Home' },
      { label: '更多图片', key: 'More' },
      { label: '本地图片', key: 'Local' },
    ]
    const settings = ref({})
    const refresh = ref(false)
    const sourceOptions = ref([])
    const route = useRoute()
    const refreshSettings = () => refresh.value = !refresh.value
    const setLoading = (value) => loading.value = value
    provide(GlobalData, {
      settings,
      refreshSettings,
      setLoading,
    })
    const firstLoad = ref(true)
    const loading = ref(false)
    const { getAllSettings, updateSettings, getAllSources } = useApi()
    const message = useMessage()
    const renderMenuLabel = (options) => {
      return h(
        RouterLink,
        { to: { name: options.key } },
        { default: () => options.label },
      )
    }
    const appWrapStyleObject = computed(() => {
      const link = unref(settings)?.currentImage
      if (link) {
        return {
          backgroundImage: `url(${link})`,
        }
      }
      return {}
    })
    const fetchSettings = () => {
      setLoading(true)
      return Promise.all([
        getAllSettings(),
        getAllSources(),
      ])
        .then(([data, data2]) => {
          loading.value = false
          settings.value = data
          sourceOptions.value = data2.map(item => {
            return {
              label: item.description,
              value: item.name,
              disabled: !item.enabled,
            }
          })
        })
        .finally(() => {
          setLoading(false)
        })
    }

    const handleSwitchSource = (value) => {
      const data = {
        currentSource: value,
      }
      updateSettings(data)
        .then(data => {
          refreshSettings()
          message.success('切换图片源成功')
        })
        .catch(err => {
          message.error(err)
        })
    }
    watch(refresh, () => {
      fetchSettings()
    })
    onMounted(() => {
      fetchSettings()
        .finally(() => {
          firstLoad.value = false
        })
    })
    const currentMenuItem = computed(() => route.name)
    return {
      zhCN,
      loading,
      settings,
      currentMenuItem,
      menuItems,
      firstLoad,
      renderMenuLabel,
      handleSwitchSource,
      appWrapStyleObject,
      options: sourceOptions,
    }
  },
})
</script>

<style lang="less" scoped>
.AppWrap {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-size: cover;

  &-header {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 6rem;
    background-color: rgba(255, 255, 255, 0.5);

    & > * {
      display: flex;
      align-items: center;
      width: 70vw;
    }

    &-sources {
      width: 20rem;
    }
  }

  &-content {
    position: relative;
    height: 0;
    flex: 1;
    background-color: rgba(255, 255, 255, 0.5);

    &-loading {
      display: flex;
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      z-index: 1;
      align-items: center;
      justify-content: center;
    }
  }
}
</style>