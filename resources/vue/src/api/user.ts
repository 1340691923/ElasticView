import request from "@/utils/request";


class UserAPI {

  /**
   * 获取当前登录用户信息
   *
   * @returns 登录用户昵称、头像信息，包括角色和权限
   */
  static getInfo() {
    return request<any, UserInfo>({
      url: `/api/gm_user/infoV2`,
      method: "post",
    });
  }

}

export default UserAPI;

/** 登录用户信息 */
export interface UserInfo {
  /** 用户ID */
  userId?: number;

  /** 用户名 */
  username?: string;

  /** 昵称 */
  nickname?: string;

  /** 头像URL */
  avatar?: string;

}


export function login(data) {
  return request({
    url: '/api/gm_user/login',
    method: 'post',
    data
  })
}

export function NoAuthRoute(data) {
  return request({
    url: '/api/NoAuthRoute',
    method: 'post',
    data
  })
}

export function getInfo() {
  return request({
    url: '/api/gm_user/info',
    method: 'post',
  })
}

export function GetRoutesConfig(data) {
  return request({
    url: '/api/gm_user/GetRoutesConfig',
    method: 'post',
    data
  })
}

export function logout() {
  return request({
    url: '/api/gm_user/logout',
    method: 'post'
  })
}

export function userList(data) {
  return request({
    url: '/api/gm_user/userlist',
    method: 'post',
    data
  })
}

export function roleOption() {
  return request({
    url: '/api/gm_user/roleOption',
    method: 'post'
  })
}

export function getUserById(data) {
  return request({
    url: '/api/gm_user/getUserById',
    method: 'post',
    data
  })
}

export function UpdateUser(data) {
  return request({
    url: '/api/gm_user/UpdateUser',
    method: 'post',
    data
  })
}

export function InsertUser(data) {
  return request({
    url: '/api/gm_user/InsertUser',
    method: 'post',
    data
  })
}

export function DelUser(data) {
  return request({
    url: '/api/gm_user/DelUser',
    method: 'post',
    data
  })
}

export function ModifyPwdByUserId(data) {
  return request({
    url: '/api/gm_user/ModifyPwdByUserId',
    method: 'post',
    data
  })
}

export function GetOAuthList(data){
  return request({
    url: '/api/GetOAuthList',
    method: 'post',
    data
  })
}

export function GetOAuthConfigs(data){
  return request({
    url: '/api/gm_user/GetOAuthConfigs',
    method: 'post',
    data
  })
}

export function SaveOAuthConfigs(data){
  return request({
    url: '/api/gm_user/SaveOAuthConfigs',
    method: 'post',
    data
  })
}

export function SealUserAction(data){
  return request({
    url: '/api/gm_user/SealUserAction',
    method: 'post',
    data
  })
}

export function UnSealUserAction(data){
  return request({
    url: '/api/gm_user/UnSealUserAction',
    method: 'post',
    data
  })
}



