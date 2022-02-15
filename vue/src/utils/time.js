
export function timestampToTime(timestamp) {
  var date = new Date(timestamp * 1000) // 时间戳为10位需*1000，时间戳为13位的话不需乘1000
  const Y = date.getFullYear() + '-'
  const M = (date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1) + '-'
  const D = change(date.getDate()) + ' '
  const h = change(date.getHours()) + ':'
  const m = change(date.getMinutes()) + ':'
  const s = change(date.getSeconds())
  return Y + M + D + h + m + s
}

function change(t) {
  if (t < 10) {
    return '0' + t
  } else {
    return t
  }
}
