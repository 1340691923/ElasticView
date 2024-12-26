import request from '@/utils/request'

var api = '/api/ai/'

export function SearchBigMode(data) {
  return request({
    url: api + 'SearchBigMode',
    method: 'post',
    data
  })
}
