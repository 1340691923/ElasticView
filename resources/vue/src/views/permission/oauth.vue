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
          <el-form-item label="corpid:" >
            <el-input v-model="data.wecomeCfg.corpid"></el-input>
          </el-form-item>
          <el-form-item label="agentId:" >
            <el-input v-model="data.wecomeCfg.agentId"></el-input>
          </el-form-item>
          <el-form-item label="secert:" >
            <el-input  type="password" show-password  v-model="data.wecomeCfg.secert"></el-input>
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
    agentId: "",
    corpid: "",
    enable: false,
    secert: ""
  }
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
