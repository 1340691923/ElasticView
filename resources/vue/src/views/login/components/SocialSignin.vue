<template>
  <div class="social-container">
    <div v-if="oauthProviders.length > 0" class="social-title">第三方登录</div>
    <div class="social-icons">
      <a v-for="provider in oauthProviders" :key="provider.name" @click="handleLogin(provider)">
        <svg-icon :icon-class="provider.img" />
      </a>
    </div>
  </div>
</template>

<script>
import { getOAuthConfig } from '@/api/user'

export default {
  name: 'SocialSignin',
  data() {
    return {
      oauthProviders: []
    }
  },
  created() {
    this.getOAuthProviders()
  },
  methods: {
    async getOAuthProviders() {
      try {
        const { data } = await getOAuthConfig()
        if (data && data.code === 0) {
          this.oauthProviders = data.data.filter(provider => provider.enable)
        }
      } catch (error) {
        console.error('Failed to load OAuth providers:', error)
      }
    },
    handleLogin(provider) {
      const redirectUri = encodeURIComponent(`${window.location.origin}/#/auth-redirect`)
      const state = encodeURIComponent(JSON.stringify({
        redirect: this.$route.query.redirect || '/'
      }))
      
      window.location.href = `/api/oauth/${provider.name}/login?redirect_uri=${redirectUri}&state=${state}`
    }
  }
}
</script>

<style lang="scss" scoped>
.social-container {
  margin: 20px 0;
  
  .social-title {
    font-size: 14px;
    color: #606266;
    margin-bottom: 10px;
    text-align: center;
  }

  .social-icons {
    display: flex;
    justify-content: center;
    
    a {
      margin: 0 10px;
      color: #606266;
      font-size: 24px;
      cursor: pointer;
      
      &:hover {
        color: #409EFF;
      }
      
      .svg-icon {
        width: 30px;
        height: 30px;
      }
    }
  }
}
</style>
