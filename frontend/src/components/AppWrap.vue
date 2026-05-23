<template>
  <div class="relative min-h-screen overflow-hidden bg-sky-50 text-slate-900">
    <Transition name="wallpaper-fade">
      <div
        v-if="appBackgroundSrc"
        :key="appBackgroundSrc"
        class="absolute -inset-8 scale-105 bg-cover bg-center opacity-55 blur-2xl"
        :style="{ backgroundImage: `url(${appBackgroundSrc})` }"
      />
    </Transition>
    <div class="absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(255,255,255,0.42),rgba(219,234,254,0.22)_42%,rgba(240,249,255,0.1)_72%,rgba(255,255,255,0.04))]" />
    <div class="absolute inset-0 bg-white/8 backdrop-blur-[1px]" />

    <aside class="fixed bottom-4 left-4 top-4 z-30 flex w-20 flex-col rounded-[2.8rem] border border-white/80 bg-white/78 px-3 py-4 shadow-fluent backdrop-blur-md sm:left-6 sm:w-64 sm:px-5">
      <RouterLink :to="{ name: 'Home' }" class="mb-8 flex items-center gap-3 rounded-3xl bg-white/88 p-3 text-slate-900 shadow-subtle ring-1 ring-white/90">
        <span class="flex h-11 w-11 shrink-0 items-center justify-center rounded-2xl bg-accent-500 text-white shadow-sm">
          <svg class="h-6 w-6" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <path d="M4 7.5A3.5 3.5 0 0 1 7.5 4h9A3.5 3.5 0 0 1 20 7.5v9a3.5 3.5 0 0 1-3.5 3.5h-9A3.5 3.5 0 0 1 4 16.5v-9Z" stroke="currentColor" stroke-width="1.8" />
            <path d="m6.8 16.8 3.2-3.2a1.4 1.4 0 0 1 2 0l1.1 1.1 1.9-1.9a1.4 1.4 0 0 1 2 0l2.2 2.2" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
            <path d="M15.6 8.2h.01" stroke="currentColor" stroke-width="3" stroke-linecap="round" />
          </svg>
        </span>
        <span class="hidden min-w-0 sm:block">
          <span class="block truncate text-[1.5rem] font-semibold">每日壁纸</span>
          <span class="block truncate text-[1.1rem] text-slate-500">Daily Wallpaper</span>
        </span>
      </RouterLink>

      <nav class="flex flex-1 flex-col gap-2">
        <RouterLink
          v-for="item in menuItems"
          :key="item.key"
          :to="{ name: item.key }"
          class="flex items-center gap-3 rounded-3xl px-3 py-3 text-[1.4rem] font-semibold"
          :class="currentMenuItem === item.key ? 'bg-white/92 text-accent-700 shadow-subtle ring-1 ring-white' : 'text-slate-600'"
        >
          <span
            class="flex h-11 w-11 shrink-0 items-center justify-center rounded-2xl"
            :class="currentMenuItem === item.key ? 'bg-accent-500 text-white shadow-sm' : 'bg-white/80 text-slate-500'"
          >
            <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path
                v-for="path in item.icon"
                :key="path"
                :d="path"
                stroke="currentColor"
                stroke-width="1.8"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </span>
          <span class="hidden min-w-0 sm:block">
            <span class="block truncate">{{ item.label }}</span>
            <span class="block truncate text-[1.1rem] font-normal opacity-70">{{ item.description }}</span>
          </span>
        </RouterLink>
      </nav>

    </aside>

    <div class="relative z-10 ml-28 h-screen overflow-y-auto sm:ml-80">
      <main class="mx-auto flex min-h-full w-full max-w-[1440px] flex-col px-4 py-5 sm:px-6 lg:px-8">
        <slot v-if="!firstLoad" />
      </main>
    </div>

    <div v-if="loading" class="fixed inset-0 z-40 flex items-center justify-center bg-sky-100/30 backdrop-blur-sm">
      <div class="flex items-center gap-4 rounded-3xl border border-white/80 bg-white/82 px-6 py-4 text-[1.4rem] font-semibold text-slate-700 shadow-fluent backdrop-blur-2xl">
        <span class="h-8 w-8 animate-spin rounded-full border-4 border-accent-100 border-t-accent-600" />
        加载中
      </div>
    </div>

    <div class="fixed right-4 top-4 z-50 flex w-[min(34rem,calc(100vw-2rem))] flex-col gap-3">
      <div
        v-for="item in notifications"
        :key="item.id"
        class="rounded-3xl border px-5 py-4 text-[1.35rem] shadow-fluent backdrop-blur-xl"
        :class="item.type === 'success' ? 'border-emerald-200 bg-emerald-50/92 text-emerald-900' : 'border-rose-200 bg-rose-50/92 text-rose-900'"
      >
        {{ item.content }}
      </div>
    </div>
  </div>
</template>

<script>
import { computed, defineComponent, onMounted, provide, ref, unref, watch } from 'vue'
import GlobalData from '@/injections/GlobalData'
import useApi from '@/composables/useApi'
import { RouterLink, useRoute } from 'vue-router'

export default defineComponent({
  components: {
    RouterLink,
  },
  setup() {
    const menuItems = [
      {
        label: '今日一图',
        description: '每日精选推荐',
        key: 'Home',
        icon: [
          'M4 7.5A3.5 3.5 0 0 1 7.5 4h9A3.5 3.5 0 0 1 20 7.5v9a3.5 3.5 0 0 1-3.5 3.5h-9A3.5 3.5 0 0 1 4 16.5v-9Z',
          'm6.8 16.8 3.2-3.2a1.4 1.4 0 0 1 2 0l1.1 1.1 1.9-1.9a1.4 1.4 0 0 1 2 0l2.2 2.2',
        ],
      },
      {
        label: '更多图片',
        description: '浏览历史图库',
        key: 'More',
        icon: [
          'M4.5 6.5A2.5 2.5 0 0 1 7 4h10a2.5 2.5 0 0 1 2.5 2.5v11A2.5 2.5 0 0 1 17 20H7a2.5 2.5 0 0 1-2.5-2.5v-11Z',
          'M8 8h8',
          'M8 12h8',
          'M8 16h5',
        ],
      },
      {
        label: '本地图片',
        description: '查看本地收藏',
        key: 'Local',
        icon: [
          'M3.8 7.5A2.5 2.5 0 0 1 6.3 5h3l2 2h6.4a2.5 2.5 0 0 1 2.5 2.5v6.8a2.7 2.7 0 0 1-2.7 2.7H6.5a2.7 2.7 0 0 1-2.7-2.7V7.5Z',
          'm7 15 2.3-2.3a1.2 1.2 0 0 1 1.7 0l.8.8 1.6-1.6a1.2 1.2 0 0 1 1.7 0L18 14.8',
        ],
      },
      {
        label: '设置',
        description: '偏好和图片源',
        key: 'Settings',
        icon: [
          'M12 15.2a3.2 3.2 0 1 0 0-6.4 3.2 3.2 0 0 0 0 6.4Z',
          'M19.4 15a1.8 1.8 0 0 0 .36 1.98l.04.05a2.15 2.15 0 0 1-3.04 3.04l-.05-.04a1.8 1.8 0 0 0-1.98-.36 1.8 1.8 0 0 0-1.08 1.65V21.5a2.15 2.15 0 0 1-4.3 0v-.18a1.8 1.8 0 0 0-1.08-1.65 1.8 1.8 0 0 0-1.98.36l-.05.04a2.15 2.15 0 0 1-3.04-3.04l.04-.05A1.8 1.8 0 0 0 4.6 15a1.8 1.8 0 0 0-1.65-1.08H2.8a2.15 2.15 0 0 1 0-4.3h.15A1.8 1.8 0 0 0 4.6 8.54a1.8 1.8 0 0 0-.36-1.98l-.04-.05a2.15 2.15 0 0 1 3.04-3.04l.05.04a1.8 1.8 0 0 0 1.98.36 1.8 1.8 0 0 0 1.08-1.65V2.1a2.15 2.15 0 0 1 4.3 0v.12a1.8 1.8 0 0 0 1.08 1.65 1.8 1.8 0 0 0 1.98-.36l.05-.04a2.15 2.15 0 0 1 3.04 3.04l-.04.05a1.8 1.8 0 0 0-.36 1.98 1.8 1.8 0 0 0 1.65 1.08h.15a2.15 2.15 0 0 1 0 4.3h-.15A1.8 1.8 0 0 0 19.4 15Z',
        ],
      },
    ]
    const settings = ref({})
    const refresh = ref(false)
    const route = useRoute()
    const notifications = ref([])
    const refreshSettings = () => refresh.value = !refresh.value
    const setLoading = (value) => loading.value = value
    const pushMessage = (type, content) => {
      const id = `${Date.now()}-${Math.random()}`
      notifications.value.push({ id, type, content })
      window.setTimeout(() => {
        notifications.value = notifications.value.filter(item => item.id !== id)
      }, 3200)
    }
    const message = {
      success: content => pushMessage('success', content),
      error: content => pushMessage('error', content),
    }
    provide(GlobalData, {
      settings,
      refreshSettings,
      setLoading,
      message,
    })
    const firstLoad = ref(true)
    const loading = ref(false)
    const { getAllSettings } = useApi()
    const appBackgroundSrc = computed(() => {
      const link = unref(settings)?.currentImage
      if (link) {
        return `/api/image/get?link=${encodeURIComponent(link)}`
      }
      return null
    })
    const fetchSettings = async () => {
      setLoading(true)
      try {
        try {
          const data = await getAllSettings();
          loading.value = false;
          settings.value = data;
        } catch (err) {
          message.error(err?.message || '获取设置失败');
        }
      } finally {
        setLoading(false);
      }
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
      loading,
      settings,
      currentMenuItem,
      menuItems,
      firstLoad,
      appBackgroundSrc,
      notifications,
    }
  },
})
</script>

<style scoped>
.wallpaper-fade-enter-active,
.wallpaper-fade-leave-active {
  transition: opacity 0.8s ease, transform 0.8s ease, filter 0.8s ease;
}

.wallpaper-fade-enter-from,
.wallpaper-fade-leave-to {
  opacity: 0;
  filter: blur(32px);
  transform: scale(1.08);
}

.wallpaper-fade-enter-to,
.wallpaper-fade-leave-from {
  opacity: 0.55;
  filter: blur(24px);
  transform: scale(1.05);
}
</style>