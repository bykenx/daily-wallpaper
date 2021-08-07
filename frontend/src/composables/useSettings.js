import request, { wrappedErrHandleGet } from '@/utils/request'
import { getErrorIfExist } from '@/utils/errorHandle'

function getAllSettings() {
  return wrappedErrHandleGet('/settings')
}

function updateSettings(settings) {
  return request.put(`/settings`, { data: settings })
    .then(res => {
      const err  =getErrorIfExist(res)
      if (err) {
        return Promise.reject(err)
      }
      return res.data
    })
}

export default function useSettings() {
  return {
    getAllSettings,
    updateSettings,
  }
}