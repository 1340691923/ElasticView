import request from '@/utils/request'

var api = '/api/ai/'

export function SearchBigMode(data: any) {
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

export function SaveAIConfig(data: any) {
  return request({
    url: api + 'SaveAIConfig',
    method: 'post',
    data
  })
}

export function TestAIConnection(data: any) {
  return request({
    url: api + 'TestAIConnection',
    method: 'post',
    data
  })
}
