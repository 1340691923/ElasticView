<template>
  <div class="app-container">
    <el-card shadow="never" class="table-container">
      <el-tabs v-model="data.activeTab">
        <el-tab-pane label="LLM配置" name="llm">
          <el-form label-width="120px" label-position="left">
            <el-form-item label="Qwen模型:">
              <el-switch
                v-model="data.aiConfig.qwenEnabled"
                active-text="启用"
                inactive-text="禁用">
              </el-switch>
            </el-form-item>
            <el-form-item label="Qwen API Key:" v-if="data.aiConfig.qwenEnabled">
              <el-input 
                placeholder="请输入通义千问API Key" 
                v-model="data.aiConfig.bigModeKey"
                type="password"
                show-password>
              </el-input>
            </el-form-item>
            
            <el-form-item label="OpenAI模型:">
              <el-switch
                v-model="data.aiConfig.openaiEnabled"
                active-text="启用"
                inactive-text="禁用">
              </el-switch>
            </el-form-item>
            <el-form-item label="OpenAI API Key:" v-if="data.aiConfig.openaiEnabled">
              <el-input 
                placeholder="请输入OpenAI API Key" 
                v-model="data.aiConfig.openAIKey"
                type="password"
                show-password>
              </el-input>
            </el-form-item>
            
            <el-form-item label="DeepSeek模型:">
              <el-switch
                v-model="data.aiConfig.deepseekEnabled"
                active-text="启用"
                inactive-text="禁用">
              </el-switch>
            </el-form-item>
            <el-form-item label="DeepSeek API Key:" v-if="data.aiConfig.deepseekEnabled">
              <el-input 
                placeholder="请输入DeepSeek API Key" 
                v-model="data.aiConfig.deepSeekKey"
                type="password"
                show-password>
              </el-input>
            </el-form-item>
            
            <el-form-item>
              <el-button @click="saveAIConfig" type="primary">保存配置</el-button>
              <el-button @click="testConnection" type="success">测试连接</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script>
import { GetAIConfig, SaveAIConfig, TestAIConnection } from '@/api/ai'

export default {
  name: 'AIConfig',
  data() {
    return {
      data: {
        activeTab: 'llm',
        aiConfig: {
          qwenEnabled: false,
          bigModeKey: '',
          openaiEnabled: false,
          openAIKey: '',
          deepseekEnabled: false,
          deepSeekKey: ''
        }
      }
    }
  },
  created() {
    this.loadAIConfig()
  },
  methods: {
    async loadAIConfig() {
      try {
        const res = await GetAIConfig()
        if (res.code !== 0) {
          this.$message.error(res.msg)
          return
        }
        this.data.aiConfig = res.data
      } catch (error) {
        console.error(error)
        this.$message.error('加载AI配置失败')
      }
    },
    async saveAIConfig() {
      try {
        const res = await SaveAIConfig(this.data.aiConfig)
        if (res.code !== 0) {
          this.$message.error(res.msg)
          return
        }
        this.$message.success('AI配置保存成功')
      } catch (error) {
        console.error(error)
        this.$message.error('保存AI配置失败')
      }
    },
    async testConnection() {
      try {
        const res = await TestAIConnection(this.data.aiConfig)
        if (res.code !== 0) {
          this.$message.error(res.msg)
          return
        }
        this.$message.success('AI服务连接测试成功')
      } catch (error) {
        console.error(error)
        this.$message.error('AI服务连接测试失败')
      }
    }
  }
}
</script>

<style scoped>
.table-container {
  margin: 30px;
}
</style>
