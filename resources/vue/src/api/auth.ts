import request from "@/utils/request";

class AuthAPI {
  /** 登录 接口*/
  static login(data: LoginData) {
    return request({
      url: `/api/gm_user/login`,
      method: "post",
      data: data,
    });
  }

  /** 注销 接口*/
  static logout() {
    return request({
      url: `/api/gm_user/logout`,
      method: "delete",
    });
  }

}

export default AuthAPI;

/** 登录请求参数 */
export interface LoginData {
  /** 用户名 */
  username: string;
  /** 密码 */
  password: string;
}
