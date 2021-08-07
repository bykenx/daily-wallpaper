import { extend } from "umi-request"
import config from '@/constants/config'
import { getErrorIfExist } from '@/utils/errorHandle'

const request = extend({
  prefix: process.env.NODE_ENV === 'development' ? '/api' : config.apiPrefix,
})

request.interceptors.request.use((url, options) => {
  return { options }
})

export default request

export function wrappedErrHandleGet(url, options) {
  return request.get(url, options)
    .then(res => {
      const err = getErrorIfExist(res)
      if (err) {
        return Promise.reject(err)
      }
      return res.data
    })
}