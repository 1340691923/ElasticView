import request from '@/utils/request'

var api = '/api/analysis/'

export function GetConfigs(data) {
  return request({
    url: api + 'GetConfigs',
    method: 'post',
    data
  })
}

export function LoadPropQuotas(data) {
  return request({
    url: api + 'LoadPropQuotas',
    method: 'post',
    data
  })
}

export function GetValues(data) {
  return request({
    url: api + 'GetValues',
    method: 'post',
    data
  })
}
