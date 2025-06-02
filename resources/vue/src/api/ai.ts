import request from '@/utils/request'

var api = '/api/ai/'

export function SearchBigMode(data: any) {
  return request({
    url: api + 'SearchBigMode',
    method: 'post',
    data
  })
}
