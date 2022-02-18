import request from '@/utils/request'

const api = '/api/es_crud/'

export function GetList(data) {
  return request({
    url: api + 'GetList',
    method: 'post',
    data
  })
}
