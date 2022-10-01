import request from '@/utils/request'

let api = "/api/TimingController/"

export function ListAction(data) {
  return request({
    url: api + 'ListAction',
    method: 'post',
    data
  })
}

export function CancelAction(data) {
  return request({
    url: api + 'CancelAction',
    method: 'get',
    params: data
  })
}
