import request from '@/utils/request'


export function GetI18nCfg(data) {
  return request({
    url: '/api/GetI18nCfg',
    method: 'post',
     data
  })
}
