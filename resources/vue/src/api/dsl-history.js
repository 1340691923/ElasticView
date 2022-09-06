import request from '@/utils/request'

const api = '/api/dslHistory/'

export function CleanAction(data) {
  return request({
    url: api + `CleanAction`,
    method: 'post',
    data
  })
}

export function ListAction(data) {
  return request({
    url: api + `ListAction`,
    method: 'post',
    data
  })
}
