<template>
  <div class="app-container">
    <el-card shadow="hover" class="search-card">
      <el-form :inline="true" class="flex flex-wrap items-end gap-4">
        <el-form-item :label="$t('接口名')" class="!mb-0">
          <el-select
            v-model="input.operater_action"
            reserve-keyword
            collapse-tags
            :placeholder="$t('接口名')"
            class="!w-[250px]"
            clearable
            filterable
            @change="getList(1)"
          >
            <el-option
              v-for="item in urlConfig"
              :key="item.url"
              :label="item.remark"
              :value="item.url"
            />
          </el-select>
        </el-form-item>

        <el-form-item :label="$t('用户')" class="!mb-0">
          <el-select
            v-model="input.operater_id"
            :placeholder="$t('请选择用户')"
            clearable
            class="!w-[120px]"
            filterable
            @change="getList(1)"
          >
            <el-option :label="$t('请选择用户')" :value="Number(0)" />
            <el-option
              v-for="item in userConfig"
              :key="item.id"
              :label="item.username"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item :label="$t('时间')" class="!mb-0 date-picker-item">
          <date v-model="input.date" class="filter-item date-picker" />
        </el-form-item>

        <el-form-item class="!mb-0">
          <el-button type="primary" @click="getList(1)">
            <el-icon class="mr-1"><Search /></el-icon>
            {{ $t('查询') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="hover" class="log-table">
      <el-table
      v-loading="tableLoading"
      :data="logList"
      max-height="800"
      :element-loading-text="`${$t('请给我点时间')}！`"
        class="w-full"
      >
        <el-table-column align="center" :label="$t('操作用户')" min-width="100">
          <template #default="scope">
            <el-tag size="small" type="info">{{ scope.row.operater_name }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column align="center" :label="$t('操作接口')" min-width="220">
        <template #default="scope">
            <el-tooltip :content="scope.row.operater_action" placement="top">
              <span class="text-gray-600 dark:text-gray-300">{{ urlConfigMap[scope.row.operater_action] }}</span>
            </el-tooltip>
        </template>
      </el-table-column>

        <el-table-column align="center" :label="$t('所花时间')" min-width="120">
        <template #default="scope">
            <span :class="getTimeClass(scope.row.cost_time)">{{ scope.row.cost_time }}</span>
        </template>
      </el-table-column>

        <el-table-column align="center" :label="$t('状态')" min-width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === '正常' ? 'success' : 'danger'" size="small">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column align="center" :label="$t('创建时间')" min-width="180">
          <template #default="scope">
            <span class="text-gray-500">{{ scope.row.created }}</span>
          </template>
        </el-table-column>

        <el-table-column align="center" :label="$t('请求数据')" min-width="100" fixed="right">
        <template #default="scope">
            <el-button 
              type="primary" 
              link
              @click="showJsonDetail(scope.row)"
              v-if="scope.row.dataFormat"
            >
              查看
            </el-button>
        </template>
      </el-table-column>
    </el-table>
    </el-card>

    <div class="pagination-container">
      <el-pagination
        background
        :current-page="input.page"
        :page-size="input.limit"
        :total="count"
        :pager-count="isMobile ? 5 : 7"
        :small="isMobile"
        layout="prev, pager, next"
        @current-change="getList"
        @size-change="handleSizeChange"
        class="pagination"
      />
    </div>

    <el-drawer
      v-model="jsonDrawerVisible"
      :title="urlConfigMap[currentLog?.operater_action] || currentLog?.operater_action"
      direction="rtl"
      :size="isMobile ? '100%' : '50%'"
      :with-header="true"
      class="json-drawer"
    >
      <div class="json-viewer-container">
        <!-- 移动端信息概览 -->
        <div v-if="isMobile" class="request-info">
          <div class="info-grid">
            <div class="info-item">
              <span class="info-label">操作用户</span>
              <el-tag size="small" type="info">{{ currentLog?.operater_name }}</el-tag>
            </div>
            <div class="info-item">
              <span class="info-label">状态</span>
              <el-tag :type="currentLog?.status === '正常' ? 'success' : 'danger'" size="small">
                {{ currentLog?.status }}
              </el-tag>
            </div>
            <div class="info-item">
              <span class="info-label">耗时</span>
              <span :class="getTimeClass(currentLog?.cost_time)">{{ currentLog?.cost_time }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">时间</span>
              <span class="text-gray-500 text-sm">{{ currentLog?.created }}</span>
            </div>
          </div>
        </div>

        <!-- 工具栏 -->
        <div class="json-toolbar">
          <el-radio-group v-model="viewMode" size="small">
            <el-radio-button label="pretty">美化</el-radio-button>
            <el-radio-button label="raw">原始</el-radio-button>
          </el-radio-group>
          <el-button 
            type="primary" 
            size="small" 
            @click="copyJson"
          >
            <el-icon class="mr-1"><DocumentCopy /></el-icon>
            复制
          </el-button>
        </div>

        <!-- JSON内容 -->
        <div class="json-content">
          <json-viewer
            v-if="viewMode === 'pretty'"
            :value="currentLog?.dataFormat"
            :expand-depth="2"
            copyable
            sort
            boxed
            theme="custom"
          />
          <pre v-else class="raw-json">{{ currentLog?.dataFormat }}</pre>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script>
import { UrlConfig } from '@/api/api-rbac'
import { getList } from '@/api/operate'
import { userList } from '@/api/user'
import dayjs from "dayjs";
import { Search, DocumentCopy } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

export default {
  components: {
    Date: () => import('@/components/Date/index.vue'),
    Search,
    DocumentCopy
  },
  data() {
    return {
      input: {
        page: 1,
        limit: 10,
        operater_action: '',
        operater_id: 0,
        date: [
          dayjs().format('YYYY-MM-DD 00:00:00'), 
          dayjs().format('YYYY-MM-DD 23:59:59')
        ]
      },
      count: 0,
      logList: [],
      tableLoading: false,
      urlConfig: [],
      urlConfigMap: [],
      userConfig: [],
      jsonDrawerVisible: false,
      currentLog: null,
      viewMode: 'pretty',
      isMobile: false
    }
  },
  computed: {


  },
  created() {
    this.initList()
  },

  mounted() {
    this.checkDevice()
    window.addEventListener('resize', this.checkDevice)
  },

  beforeUnmount() {
    window.removeEventListener('resize', this.checkDevice)
  },

  methods: {
    lookData(index, typ) {
      for (const i in this.logList) {
        if (this.logList[i].index == index) {
          this.logList[i].isFormatData = typ
        }
      }
    },
    changeDate(v) {
      this.input.date = v
      this.getList(1)
    },
    async initList() {
      const urlConfigRes = await UrlConfig()
      if (urlConfigRes) {
        for (var v of urlConfigRes.data.cfg) {
          this.urlConfig.push(v)
          this.urlConfigMap[v['url']] = v['remark']
        }
      }
      const userListRes = await userList()
      if (userListRes) {
        for (var v of userListRes.data.list) {
          this.userConfig.push(v)
        }
      }
      this.getList(1)
    },
    handleSizeChange(v) {
      this.input.limit = v
      this.getList(1)
    },
    getUrlOpt() {
      UrlConfig().then(res => {
        if (res) {
          for (var v of res.data.cfg) {
            this.urlConfig.push(v)
            this.urlConfigMap[v['url']] = v['remark']
          }
        }
      })
    },

    getUserOpt() {
      userList().then(res => {
        if (res) {
          for (var v of res.data) {
            this.userConfig.push(v)
          }
        }
      })
    },

    getList(page) {
      !page ? this.input.page = 1 : this.input.page = page
      this.tableLoading = true

      if (this.input.operater_id == '') {
        this.input.operater_id = 0
      }


      this.input.date[0] = dayjs(this.input.date[0]).format('YYYY-MM-DD HH:mm:ss')
      this.input.date[1] = dayjs(this.input.date[1]).format('YYYY-MM-DD HH:mm:ss')

      getList(this.input).then(res => {
        if (res.code == 0) {
          if (res.data.list == null) {
            res.data.list = []
          }
          const list = []
          let index = 0
          for (const v of res.data.list) {
            if (v['body_str'] != '') {
              v['dataFormat'] = JSON.stringify(JSON.parse(v['body_str']), null, '\t')
            }
            v['isFormatData'] = false
            v['index'] = index
            list.push(v)
            index++
          }
          this.logList = list

          this.tableLoading = false
          this.count = Number(res.data.count)
        } else {
          ElMessage.error({
            offset: 60,
            type: 'error',
            message: res.msg
          })
        }
        this.tableLoading = false
      }).catch(err => {
        this.tableLoading = false
        console.log('err', err)
      })
    },
    showJsonDetail(row) {
      this.currentLog = row
      this.jsonDrawerVisible = true
    },
    copyJson() {
      if (!this.currentLog?.dataFormat) return
      navigator.clipboard.writeText(this.currentLog.dataFormat)
      ElMessage.success('复制成功')
    },
    getTimeClass(time) {
      const duration = parseFloat(time)
      return {
        'text-success': duration < 100,
        'text-warning': duration >= 100 && duration < 1000,
        'text-danger': duration >= 1000
      }
    },
    checkDevice() {
      this.isMobile = window.innerWidth <= 768
    }
  }
}
</script>

<style lang="scss" scoped>
.search-card {
  @apply mb-6;
  @apply bg-white/90 dark:bg-gray-800/90;
  @apply backdrop-blur-sm;
  @apply transition-all duration-300;
  @apply rounded-lg;
  @apply shadow-sm;

  :deep(.el-card__body) {
    @apply p-4;
  }

  :deep(.el-form-item__label) {
    @apply text-gray-600 dark:text-gray-300;
  }

  :deep(.el-input__wrapper),
  :deep(.el-select .el-input__wrapper) {
    @apply shadow-sm;
    @apply bg-white dark:bg-gray-700;
    @apply border-gray-200 dark:border-gray-600;
  }

  :deep(.el-button--primary) {
    @apply shadow-sm;
    @apply transition-all duration-300;
    
    &:hover {
      @apply transform scale-105;
      @apply shadow-md;
    }
  }

  .date-picker-item {
    @media (max-width: 768px) {
      @apply w-full;
      
      :deep(.el-form-item__content) {
        @apply w-full;
      }
    }
  }

  :deep(.date-picker) {
    @media (max-width: 768px) {
      @apply w-full;
      
      .el-range-editor {
        @apply w-full max-w-full;
        
        .el-range-input {
          @apply w-[42%];
        }
        
        .el-range-separator {
          @apply w-[8%];
        }
      }
    }
  }
}

.log-table {
  @apply rounded-lg overflow-hidden;
  @apply bg-white/90 dark:bg-gray-800/90;
  @apply backdrop-blur-sm;
  @apply transition-all duration-300;
  @apply shadow-sm;

  :deep(.el-table) {
    @apply bg-transparent;
    
    th {
      @apply bg-gray-50/90 dark:bg-gray-700/90;
      @apply text-gray-600 dark:text-gray-300;
      @apply font-medium;
      @apply border-b border-gray-200 dark:border-gray-600;
    }
    
    td {
      @apply border-b border-gray-100 dark:border-gray-700;
    }
  }
}

.json-drawer {
  :deep(.el-drawer__header) {
    @apply mb-0 p-4;
    @apply border-b border-gray-200 dark:border-gray-600;
    margin-bottom: 0 !important;
  }

  :deep(.el-drawer__body) {
    @apply p-0;
    height: calc(100% - 54px);
  }
}

.json-viewer-container {
  @apply h-full flex flex-col;

  // 添加请求信息样式
  .request-info {
    @apply p-4;
    @apply border-b border-gray-200 dark:border-gray-600;
    @apply bg-white/90 dark:bg-gray-800/90;
    @apply backdrop-blur-sm;

    .info-grid {
      @apply grid grid-cols-2 gap-4;

      .info-item {
        @apply flex flex-col gap-1;

        .info-label {
          @apply text-xs text-gray-500 dark:text-gray-400;
        }
      }
    }
  }

  .json-toolbar {
    @apply p-4 flex items-center justify-between;
    @apply border-b border-gray-200 dark:border-gray-700;
    @apply bg-white/90 dark:bg-gray-800/90;
    @apply backdrop-blur-sm;
    @apply sticky top-0 z-10;

    @media (max-width: 768px) {
      @apply p-2;
      
      :deep(.el-radio-button__inner) {
        @apply px-2 py-1 text-sm;
      }
    }
  }

  .json-content {
    @apply flex-1 overflow-auto p-4;
    
    @media (max-width: 768px) {
      @apply p-2;
    }

    :deep(.jv-container) {
      @apply bg-transparent;
      
      .jv-key {
        @apply text-blue-600 dark:text-blue-400;
      }
      
      .jv-item {
        @apply text-gray-700 dark:text-gray-300;
      }
      
      .jv-string {
        @apply text-green-600 dark:text-green-400;
      }
      
      .jv-number {
        @apply text-orange-600 dark:text-orange-400;
      }
      
      .jv-boolean {
        @apply text-purple-600 dark:text-purple-400;
      }
      
      .jv-null {
        @apply text-gray-500 dark:text-gray-500;
      }
    }
  }

  .raw-json {
    @apply p-4 rounded-lg;
    @apply bg-gray-50 dark:bg-gray-900;
    @apply text-gray-700 dark:text-gray-300;
    @apply font-mono text-sm;
    @apply whitespace-pre-wrap;
  }
}

.text-success {
  @apply text-green-600 dark:text-green-400;
}

.text-warning {
  @apply text-amber-600 dark:text-amber-400;
}

.text-danger {
  @apply text-red-600 dark:text-red-400;
}

.pagination-container {
  @apply mt-4 flex justify-end;
  @apply overflow-x-auto;
  
  :deep(.el-pagination) {
    @apply bg-white/90 dark:bg-gray-800/90;
    @apply rounded-lg p-2;
    @apply shadow-sm;
    @apply min-w-fit;
    
    @media (max-width: 768px) {
      @apply w-full;
      @apply flex justify-center;
      @apply text-sm;
      
      .el-pager {
        @apply flex-wrap;
      }
      
      .btn-prev,
      .btn-next {
        @apply min-w-[24px];
      }
      
      li {
        @apply min-w-[24px];
      }
    }
  }
}
</style>

