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
              <NTimePicker
                :default-value="defaultTime"
                v-model:value="model.timeToUpdate"
              />
            </NFormItem>
          </NForm>
          <NButton type="primary">提交</NButton>
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
import useToggle from '@/composables/useToggle'

export default defineComponent({
  components: { ImageItem, NImage, NCol, NRow, NForm, NFormItem, NSwitch, NInput, NTimePicker, NButton },
  setup() {
    const { settings, setLoading } = inject(GlobalData)
    const message = useMessage()
    const { getTodayImage } = useApi()

    const todayImage = ref()
    const formModel = ref({})

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
    const defaultTime = new Date('2000/01/01 08:00')
    watch(settings, (value) => {
      formModel.value = {
        ...value,
        timeToUpdate: value.timeToUpdate ? value.timeToUpdate : null,
      }
    })
    watch(() => unref(settings).currentSource, () => {
      fetchImage()
    })
    return {
      todayImage,
      defaultTime,
      model: formModel,
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