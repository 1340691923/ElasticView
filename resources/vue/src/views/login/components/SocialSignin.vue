<template>
  <div class="social-signup-container">
    <div class="sign-btn" @click="handleOAuthClick('wechat')">
      <span class="wx-svg-container"><svg-icon icon-class="wechat" class="icon"/></span>
      企业微信
    </div>
    <div class="sign-btn" @click="handleOAuthClick('dingtalk')">
      <span class="dingtalk-svg-container"><svg-icon icon-class="dingtalk" class="icon"/></span>
      钉钉
    </div>
    <div class="sign-btn" @click="handleOAuthClick('feishu')">
      <span class="feishu-svg-container"><svg-icon icon-class="feishu" class="icon"/></span>
      飞书
    </div>
  </div>
</template>

<script>
import openWindow from '@/utils/open-window'

export default {
  name: 'SocialSignin',
  data() {
    return {
      oauthConfig: {
        wechat: {
          appId: 'wx39f7f96daaec1f96',
          authUrl: 'https://open.weixin.qq.com/connect/qrconnect',
          scope: 'snsapi_login',
          responseType: 'code',
          hash: '#wechat_redirect'
        },
        dingtalk: {
          appId: 'your_dingtalk_app_id',
          authUrl: 'https://login.dingtalk.com/oauth2/auth',
          scope: 'openid',
          responseType: 'code',
          hash: ''
        },
        feishu: {
          appId: 'your_feishu_app_id',
          authUrl: 'https://open.feishu.cn/open-apis/authen/v1/index',
          scope: 'user_info',
          responseType: 'code',
          hash: ''
        }
      }
    }
  },
  methods: {
    handleOAuthClick(provider) {
      const config = this.oauthConfig[provider]
      if (!config) return
      
      const state = provider + '_' + Date.now()
      const redirect_uri = encodeURIComponent('http://localhost:9528/#/auth-redirect?provider=' + provider)
      
      let url = `${config.authUrl}?`
      
      // Different providers use different parameter names
      if (provider === 'feishu') {
        url += `app_id=${config.appId}&redirect_uri=${redirect_uri}&state=${state}`
      } else {
        url += `client_id=${config.appId}&redirect_uri=${redirect_uri}&response_type=${config.responseType}&scope=${config.scope}&state=${state}`
      }
      
      if (config.hash) {
        url += config.hash
      }
      
      openWindow(url, provider, 540, 540)
    }
  }
}
</script>

<style lang="scss" scoped>
.social-signup-container {
  margin: 20px 0;

.sign-btn {
  display: inline-block;
  cursor: pointer;
}

.icon {
  color: #fff;
  font-size: 24px;
  margin-top: 8px;
}

.wx-svg-container,
.qq-svg-container {
  display: inline-block;
  width: 40px;
  height: 40px;
  line-height: 40px;
  text-align: center;
  padding-top: 1px;
  border-radius: 4px;
  margin-bottom: 20px;
  margin-right: 5px;
}

.wx-svg-container {
  background-color: #24da70;
}

.qq-svg-container {
  background-color: #6BA2D6;
  margin-left: 50px;
}

}
</style>
