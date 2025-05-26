import request from '@/utils/request'

const api = '/api/notice/'

export function GetList(data) {
  return request({
    url: api + `GetList`,
    method: 'post',
    data
  })
}

export function MarkReadNotice(data) {
  return request({
    url: api + `MarkReadNotice`,
    method: 'post',
    data
  })
}

export function Truncate(data) {
  return request({
    url: api + `Truncate`,
    method: 'post',
    data
  })
}


