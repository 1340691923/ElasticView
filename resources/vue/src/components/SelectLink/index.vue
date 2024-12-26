<template>
  <div >

    <el-select
      v-model="data.linkID"
      style="width:10rem"
      filterable
      default-first-option
      :placeholder="$t('请选择数据源')"
      @change="change"
    >
      <el-option :value="Number(0)" label="请选择数据源" />
      <el-option v-for="item in data.opt" :key="item.id" :value="Number(item.id)" :label="`${item.remark} [${item.version}]`" />
    </el-select>
  </div>
</template>

<script lang="ts" setup>
import { OptAction } from '@/api/es-link'
import {GetEsConnect, SaveEsConnect, SaveEsConnectVer} from "@/utils/es_link";

const data = reactive({
  opt: [],
  linkID: '',
})

onMounted(()=>{
  let esConnect = GetEsConnect()
  if(esConnect){
    data.linkID = esConnect
  }
  getEsOpt()
})

const getEsOpt = async () => {
  const res = await OptAction({ 'getByLocal': 1 })
  if (res.data == null) res.data = []
  data.opt = res.data
}

const refresh = () => {
  ElMessage.success({
    type: 'success',
    message: '刷新数据源成功'
  })
}

const change = (link) => {
  let version = ''
  for(let v of data.opt){
    if(v.id == link){
      version = v.version
    }
  }

  SaveEsConnect(link)
  SaveEsConnectVer(version)

  window.location.reload()

}


</script>

<style>
.el-select-dropdown {
  z-index: 9999 !important; /* 设置一个较高的值 */
}
</style>
