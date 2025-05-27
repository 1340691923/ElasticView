import request from '@/utils/request'

var api = '/api/ai/'

interface SearchBigModeRequest {
  sys_type: number;
  content: string;
}

interface AIConfigRequest {
  qwen_enabled: boolean;
  big_mode_key: string;
  openai_enabled: boolean;
  open_ai_key: string;
  deepseek_enabled: boolean;
  deep_seek_key: string;
}

export function SearchBigMode(data: SearchBigModeRequest) {
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

export function SaveAIConfig(data: AIConfigRequest) {
  return request({
    url: api + 'SaveAIConfig',
    method: 'post',
    data
  })
}

export function TestAIConnection(data: AIConfigRequest) {
  return request({
    url: api + 'TestAIConnection',
    method: 'post',
    data
  })
}
