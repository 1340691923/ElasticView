import request from '@/utils/request'

var api = '/api/gm_user/'

export function UrlConfig(data) {
  return request({
    url: api + 'UrlConfig',
    method: 'post',
     data
  })
}

export function SaveRbac(data) {
  return request({
    url: api + 'SaveRbac',
    method: 'post',
     data
  })
}

export function RbacList() {
  return request({
    url: api + 'RbacList',
    method: 'post'
  })
}
