<template>
  <div class="app-container">
    <el-card shadow="never" class="table-container">
      <el-tabs v-model="data.authType" >
      <el-tab-pane label="企业微信认证(内部应用)" name="企业微信认证(内部应用)">
        <el-form label-width="100px" label-position="left">
          <el-form-item label="开启认证:" >
            <el-switch
              v-model="data.wecomeCfg.enable"
              active-text="开启"
              inactive-text="不开启">
            </el-switch>
          </el-form-item>
          <el-form-item label="回调域名:" >
            <el-input :placeholder="cbPlaceholder" v-model="data.wecomeCfg.rootUrl"></el-input>
          </el-form-item>
          <el-form-item label="corpid:" >
            <el-input placeholder="企业 ID，在企业微信管理后台可查，用于识别调用接口的企业身份。类似于企业的唯一编号。" v-model="data.wecomeCfg.corpid"></el-input>
          </el-form-item>
          <el-form-item label="agentId:" >
            <el-input placeholder="应用 ID，用于区分具体是哪个企业内部应用。" v-model="data.wecomeCfg.agentId"></el-input>
          </el-form-item>
          <el-form-item label="secert:" >
            <el-input placeholder="自建应用凭证密钥，用于获取该应用的 access_token。每个应用有不同的 secret，在企业微信管理后台可以查看。" type="password" show-password  v-model="data.wecomeCfg.secert"></el-input>
          </el-form-item>
          <el-form-item label="" >
            <el-button @click="saveOAuthConfig({
            application_name:data.authType,
            config:data.wecomeCfg,
            })" type="primary" >提交</el-button>
          </el-form-item>
        </el-form>

      </el-tab-pane>

    </el-tabs>
    </el-card>
  </div>
</template>

<script lang="ts" setup>

import {GetOAuthConfigs, SaveOAuthConfigs} from "@/api/user";

const data = reactive({
  authType:'企业微信认证(内部应用)',
  wecomeCfg:{
    rootUrl:"",
    agentId: "",
    corpid: "",
    enable: false,
    secert: ""
  }
})

const cbPlaceholder = computed(()=>{
  let href = window.location.href;
  if(!import.meta.env.PROD){
    href = import.meta.env.VITE_APP_API_URL
  }
  const protocol = href.split('//')[0] === 'http:' ? 'http' : 'https';
  const host = href.split('//')[1].split('/')[0];
  return `例如:${protocol}://${host}/`
})

const saveOAuthConfig = async (data)=>{
  const res= await SaveOAuthConfigs(data)
  if(res.code != 0 ){
    ElMessage.error({
      message: res.msg,
      type: 'error'
    })
    return
  }
  ElMessage.success({
    message: res.msg,
    type: 'success'
  })
}

const Init = async ()=>{
  const res = await GetOAuthConfigs({})
  if(res.code != 0 ){
    ElMessage.error({
      message: res.msg,
      type: 'error'
    })
    return
  }
  data.wecomeCfg = res.data["企业微信认证(内部应用)"]
}

onMounted(()=>{
  Init()
})

</script>

<style scoped>

</style>
