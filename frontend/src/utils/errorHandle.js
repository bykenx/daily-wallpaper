/**
 *
 * @param res {{code?:Number, data?:any, msg?:string}}
 * @returns {String}
 */
export function getErrorIfExist(res) {
  if (!res) {
    return '网络错误'
  }
  if (res.code >= 400 && res.code < 500) {
    return res.msg ?? '请求参数错误'
  } else if (res.code >= 500) {
    return res.msg ?? '服务器错误'
  }
}