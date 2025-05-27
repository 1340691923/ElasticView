import request from '@/utils/request'

var api = '/api/ai/'

export function SearchBigMode(data) {
  return request({
    url: api + 'SearchBigMode',
    method: 'post',
    data
  })
}

export function GetAIConfig() {
  return request({
    url: api + 'GetAIConfig',
    method: 'post'
  })
}

export function SaveAIConfig(data) {
  return request({
    url: api + 'SaveAIConfig',
    method: 'post',
    data
  })
}

export function TestAIConnection(data) {
  return request({
    url: api + 'TestAIConnection',
    method: 'post',
    data
  })
}
