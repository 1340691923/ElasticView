import request from "@/utils/request";

export function setIndexCfg(data) {
  return request({
    url: "/api/search/setIndexCfg",
    method: "post",
    data
  });
}

export function getIndexCfg(data) {
  return request({
    url: "/api/search/getIndexCfg",
    method: "post",
    data
  });
}

export function SearchLog(data) {
  return request({
    url: "/api/search/SearchLog",
    method: "post",
    data
  });
}

export function SetMappingAlias(data) {
  return request({
    url: "/api/search/SetMappingAlias",
    method: "post",
    data
  });
}
export function GetMappingAlias(data) {
  return request({
    url: "/api/search/GetMappingAlias",
    method: "post",
    data
  });
}

