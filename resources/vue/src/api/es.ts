import request from '@/utils/request'

const api = '/api/es/'

export function PingAction(data) {
  return request({
    url: api + `PingAction`,
    method: 'post',
    data
  })
}

export function IndexsCountAction(data) {
  return request({
    url: api + `IndexsCountAction`,
    method: 'post',
    data
  })
}

export function CatAction(data) {
  return request({
    url: api + `CatAction`,
    method: 'post',
    data
  })
}
