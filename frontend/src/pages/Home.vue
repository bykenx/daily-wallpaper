<template>
  <div class="Home">
    <div class="Home-content">
      <NRow gutter="12">
        <NCol span="16">
          <ImageItem v-if="todayImage" :value="todayImage" />
        </NCol>
        <NCol span="8">
          <NForm :model="model">
            <NFormItem path="currentImage" label="当前图片">
              <NInput
                :disabled="model.autoUpdate"
                :placeholder="model.autoUpdate ? '开启每天更新桌面时此项无效' : '贴入链接或选择图片'"
                v-model:value="model.currentImage"
              />
            </NFormItem>
            <NFormItem path="autoRunAtSystemBoot" label="开机自启">
              <NSwitch v-model:value="model.autoRunAtSystemBoot" />
            </NFormItem>
            <NFormItem path="autoUpdate" label="每天自动更新桌面">
              <NSwitch v-model:value="model.autoUpdate" />
            </NFormItem>
            <NFormItem v-if="model.autoUpdate" path="timeToUpdate" label="更新时间">
              <NTimePicker v-model:value="model.timeToUpdate" format="hh:mm" />
            </NFormItem>
            <NFormItem path="qualityFirst" label="更高品质">
              <NSwitch v-model:value="model.qualityFirst" />
            </NFormItem>
          </NForm>
          <NButton type="primary" @click="handleSubmit">提交</NButton>
        </NCol>
      </NRow>
    </div>
  </div>
</template>

<script>
import { defineComponent, inject, onMounted, ref, unref, watch } from 'vue'
import { NButton, NCol, NForm, NFormItem, NImage, NInput, NRow, NSwitch, NTimePicker, useMessage } from 'naive-ui'
import ImageItem from '@/components/ImageItem.vue'
import GlobalData from '@/injections/GlobalData'
import useApi from '@/composables/useApi'

export default defineComponent({
  components: { ImageItem, NImage, NCol, NRow, NForm, NFormItem, NSwitch, NInput, NTimePicker, NButton },
  setup() {
    const { settings, setLoading } = inject(GlobalData)
    const message = useMessage()
    const { getTodayImage, updateSettings } = useApi()

    const todayImage = ref()
    const formModel = ref({
      autoUpdate: false,
      currentImage: '',
      autoRunAtSystemBoot: false,
      timeToUpdate: null,
      qualityFirst: false,
    })

    const fetchImage = () => {
      setLoading(true)
      getTodayImage()
        .then(data => {
          todayImage.value = data
        })
        .catch(err => {
          message.error(err)
        })
        .finally(() => {
          setLoading(false)
        })
    }
    onMounted(() => {
      fetchImage()
    })
    const setFormModel = (value) => {
      value = { ...value }
      if (value.timeToUpdate) {
        value.timeToUpdate = new Date(`2000/01/01 ${value.timeToUpdate}`).getTime()
      } else {
        value.timeToUpdate = null
      }
      formModel.value = value
    }
    setFormModel(unref(settings))
    watch(settings, (value) => {
      setFormModel(value)
    })
    watch(() => unref(settings).currentSource, () => {
      fetchImage()
    })
    const handleSubmit = () => {
      const value = { ...unref(formModel) }
      if (value.timeToUpdate) {
        const time = new Date(value.timeToUpdate)
        const hourStr = time.getHours().toString().padStart(2, '0')
        const minuteStr = time.getMinutes().toString().padStart(2, '0')
        const secondStr = time.getSeconds().toString().padStart(2, '0')
        value.timeToUpdate = `${hourStr}:${minuteStr}:${secondStr}`
      }
      console.log('settings: ', value)
      updateSettings(value)
        .then(() => {
          message.success('设置更新成功')
        })
    }
    return {
      todayImage,
      model: formModel,
      handleSubmit,
    }
  },
})
</script>

<style lang="less">
.Home {
  height: 100%;
  display: flex;
  justify-content: center;

  &-content {
    width: 70vw;
  }
}
</style>