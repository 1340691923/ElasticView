<template>
  <div class="custom-select-wrapper">
    <el-select
      v-model="data.linkID"
      filterable
      default-first-option
      :placeholder="$t('请选择数据源')"
      @change="change"
      class="custom-select"
      popper-class="custom-select-dropdown"
    >
      <template #prefix>
        <el-icon><DataLine /></el-icon>
      </template>
      <el-option :value="Number(0)" :label="$t('请选择数据源')" class="placeholder-option" />
      <el-option 
        v-for="item in data.opt" 
        :key="item.id" 
        :value="Number(item.id)" 
        :label="`${item.remark} [${item.version}]`"
      >
        <div class="custom-option">
          <div class="option-content">
            <span class="option-name">{{ item.remark }}</span>
            <el-tag size="small" effect="plain" class="version-tag">{{ item.version }}</el-tag>
          </div>
          <div class="option-indicator" :class="{ active: Number(data.linkID) === Number(item.id) }">
            <el-icon><Check /></el-icon>
          </div>
        </div>
      </el-option>
    </el-select>
  </div>
</template>

<script lang="ts" setup>
import { OptAction } from '@/api/es-link'
import {GetEsConnect, SaveEsConnect, SaveEsConnectVer} from "@/utils/es_link";
import { DataLine, Check } from '@element-plus/icons-vue'

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
    type: 'success',offset:60,
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

<style lang="scss" scoped>
.custom-select-wrapper {
  display: inline-block;
  
  :deep(.custom-select) {
    width: 280px;
    
    .el-input {
      &__wrapper {
        background: var(--el-bg-color);
        border-radius: 0.75rem;
        padding: 0.25rem 0.5rem;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
        border: 1px solid var(--el-border-color-light);
        transition: all 0.3s ease;

        &:hover {
          border-color: var(--el-color-primary-light-5);
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
        }

        &.is-focus {
          border-color: var(--el-color-primary);
          box-shadow: 0 0 0 2px var(--el-color-primary-light-8);
        }
      }

      &__prefix {
        margin-right: 8px;
        color: var(--el-text-color-secondary);
      }
    }

    .el-select__tags {
      margin-left: 6px;
    }
  }
}

:deep(.custom-select-dropdown) {
  border: none;
  border-radius: 0.75rem;
  padding: 0.5rem;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  background: var(--el-bg-color);

  .el-select-dropdown__item {
    height: auto;
    padding: 0.5rem;
    border-radius: 0.5rem;
    margin: 0.25rem 0;

    &.hover {
      background: var(--el-color-primary-light-9);
    }

    &.selected {
      background: var(--el-color-primary-light-9);
      font-weight: normal;
      
      .option-name {
        color: var(--el-color-primary);
        font-weight: 500;
      }
    }
  }

  .placeholder-option {
    color: var(--el-text-color-secondary);
    font-style: italic;
  }

  .custom-option {
    display: flex;
    align-items: center;
    justify-content: space-between;
    
    .option-content {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      
      .option-name {
        font-size: 0.875rem;
      }
      
      .version-tag {
        font-size: 0.75rem;
        background: var(--el-color-info-light-9);
        border: 1px solid var(--el-border-color-lighter);
        padding: 0 0.5rem;
      }
    }

    .option-indicator {
      opacity: 0;
      color: var(--el-color-primary);
      transition: all 0.3s ease;
      
      &.active {
        opacity: 1;
      }
    }
  }

  .el-scrollbar__wrap {
    margin-bottom: 0 !important;
  }
}

// 深色模式适配
:deep(.dark) {
  .custom-select {
    .el-input__wrapper {
      background: var(--el-bg-color-overlay);
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);

      &:hover {
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
      }
    }
  }

  .custom-select-dropdown {
    background: var(--el-bg-color-overlay);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);

    .el-select-dropdown__item {
      &.hover {
        background: var(--el-color-primary-dark-2);
      }

      &.selected {
        background: var(--el-color-primary-dark-2);
      }
    }

    .version-tag {
      background: var(--el-bg-color);
      border-color: var(--el-border-color);
    }
  }
}
</style>
