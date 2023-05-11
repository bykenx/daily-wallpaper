import request, { wrappedErrHandleGet } from '@/utils/request'
import { getErrorIfExist } from '@/utils/errorHandle'

function getAllSettings() {
  return wrappedErrHandleGet('/settings')
}

function updateSettings(settings) {
  return request.put(`/settings`, { data: settings })
    .then(res => {
      const err = getErrorIfExist(res)
      if (err) {
        return Promise.reject(err)
      }
      return res.data
    })
}

function getAllSources() {
  return wrappedErrHandleGet(`/image/sources`)
}

function getArchiveImages(params) {
  return wrappedErrHandleGet(`/image/archive`, { params })
}

function getTodayImage() {
  return wrappedErrHandleGet(`/image/today`)
}

function getImageList(params) {
  return wrappedErrHandleGet(`/image/list`, { params })
}

export default function useApi() {
  return {
    getAllSettings,
    updateSettings,
    getAllSources,
    getArchiveImages,
    getTodayImage,
    getImageList,
  }
}
