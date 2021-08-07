import { wrappedErrHandleGet } from '@/utils/request'

function getAllSources() {
  return wrappedErrHandleGet(`/image/sources`)
}

function getArchiveImages(source) {
  return wrappedErrHandleGet(`/image/archive`, { query: { source } })
}

function getTodayImage(source) {
  return wrappedErrHandleGet(`/image/today`, { query: { source } })
}

export default function useImage() {
  return {
    getAllSources,
    getTodayImage,
    getArchiveImages,
  }
}