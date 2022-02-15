import request from '@/utils/request'

const api = '/api/backUp/'

export function SnapshotListAction(data) {
  return request({
    url: api + 'SnapshotListAction',
    method: 'post',
    data
  })
}

export function SnapshotRepositoryListAction(data) {
  return request({
    url: api + 'SnapshotRepositoryListAction',
    method: 'post',
    data
  })
}

export function SnapshotCreateRepositoryAction(data) {
  return request({
    url: api + 'SnapshotCreateRepositoryAction',
    method: 'post',
    data
  })
}
export function SnapshotDeleteRepositoryAction(data) {
  return request({
    url: api + 'SnapshotDeleteRepositoryAction',
    method: 'post',
    data
  })
}

export function CleanupeRepositoryAction(data) {
  return request({
    url: api + 'CleanupeRepositoryAction',
    method: 'post',
    data
  })
}

export function CreateSnapshotAction(data) {
  return request({
    url: api + 'CreateSnapshotAction',
    method: 'post',
    data
  })
}

export function SnapshotDeleteAction(data) {
  return request({
    url: api + 'SnapshotDeleteAction',
    method: 'post',
    data
  })
}
export function SnapshotDetailAction(data) {
  return request({
    url: api + 'SnapshotDetailAction',
    method: 'post',
    data
  })
}
export function SnapshotRestoreAction(data) {
  return request({
    url: api + 'SnapshotRestoreAction',
    method: 'post',
    data
  })
}

export function SnapshotStatusAction(data) {
  return request({
    url: api + 'SnapshotStatusAction',
    method: 'post',
    data
  })
}
