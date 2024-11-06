<!--  线 + 柱混合图 -->
<template>
  <el-row :gutter="10" class="mt-5">

    <el-col :xs="24" :sm="12" :lg="6">
      <el-card>
        <template #header>
          <div class="flex-x-between">
              <span class="text-[var(--el-text-color-secondary)]"
              >总分片数</span
              >
            <el-tag type="primary">{{articleData.count.allShared }}</el-tag>
          </div>
        </template>

        <div
          class="flex-x-between mt-2 text-sm text-[var(--el-text-color-secondary)]"
        >
          <span> 主分片数 </span>
          <el-tag type="success">{{articleData.count.pri }}</el-tag>
        </div>
        <div
          class="flex-x-between mt-2 text-sm text-[var(--el-text-color-secondary)]"
        >
          <span> 未分配的分片数 </span>
          <el-tag type="danger">{{articleData.count.unassigned }}</el-tag>
        </div>
      </el-card>
    </el-col>
    <el-col :xs="24" :sm="12" :lg="6">
      <el-card>
        <template #header>
          <div class="flex-x-between">
              <span class="text-[var(--el-text-color-secondary)]"
              >索引数</span
              >
            <el-tag type="primary">{{articleData.count.index}}</el-tag>
          </div>
        </template>

        <div
          class="flex-x-between mt-2 text-sm text-[var(--el-text-color-secondary)]"
        >
          <span> 文档数 </span>
          <el-tag type="success">{{articleData.count.document}}</el-tag>
        </div>
        <div
          class="flex-x-between mt-2 text-sm text-[var(--el-text-color-secondary)]"
        >
          <span> 索引所占空间大小 </span>
          <el-tag type="danger">{{articleData.count.size}}</el-tag>
        </div>
      </el-card>
    </el-col>

  </el-row>
</template>

<script setup lang="ts">
import {GetLocalPluginList, GetWxArticleList} from "@/api/plugins";
import {GetEsConnect} from "@/utils/es_link";
import {CatAction, IndexsCountAction} from "@/api/es";

const articleData = reactive({
  count: {
    allShared: 0,
    successfulShared: 0,
    index: 0,
    pri:0,
    unassigned:0,
    document: 0,
    size: 0
  },
})

const getIndexCount = async () =>{
  const form = {
    es_connect: GetEsConnect()
  }

  const { data, code, msg } = await IndexsCountAction(form)

  if (code == 0) {
    articleData.count.index = data
  }

}

const catAllocation = async () => {
  const form = {
    cat: 'CatAllocation',
    es_connect: GetEsConnect()
  }

  const { data, code, msg } = await CatAction(form)

  if (code == 0) {
    articleData.count.size = data[0]['disk.indices']
  }
}

const catStats = async ()=> {
  const form = {
    cat: 'CatStats',
    es_connect: GetEsConnect()
  }

  const { data, code, msg } = await CatAction(form)

  if (code == 0) {
    articleData.count.document = data.indices.docs.count
  }
}

const getSegments = async () =>{

  const form = {
    cat: 'CatHealth',
    es_connect: GetEsConnect()
  }

  const { data, code, msg } = await CatAction(form)

  if (code == 0) {
    if(data.length >0){
      let allShared = 0
      let unassigned = 0
      let pri = 0
      for(let i in data){

        let v = data[i]
        pri = pri + Number(v.pri)
        allShared = allShared + Number(v.shards)
        unassigned = unassigned + Number(v.unassign)
      }

      articleData.count.unassigned = unassigned
      articleData.count.allShared = allShared
      articleData.count.pri = pri
    }

  }
}


onMounted(()=>{
  getIndexCount()
  catAllocation()
  getSegments()
  catStats()
})
</script>
<style lang="scss" scoped></style>
