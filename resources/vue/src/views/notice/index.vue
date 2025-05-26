<template>
  <div class="app-container">
    <!-- 搜索区域 -->
    <div class="search-container">
      <div class="search-header" @click="toggleSearchForm">
        <span>筛选条件</span>
        <el-icon class="toggle-icon" :class="{ 'is-active': !isSearchCollapsed }">
          <arrow-down />
        </el-icon>
      </div>
      <el-form
        ref="queryFormRef"
        :model="queryParams"
        :inline="true"
        label-suffix=":"
        label-width="auto"
        v-show="!isSearchCollapsed"
        class="search-form"
      >
        <el-form-item label="标题" prop="title">
          <el-input
            v-model="queryParams.title"
            placeholder="标题"
            clearable
            @keyup.enter="handleQuery"
          />
        </el-form-item>

        <el-form-item label="阅读状态" prop="read_type">
          <el-select style="width: 10rem" v-model="queryParams.read_type" placeholder="请选择阅读状态" clearable>
            <el-option label="全部" :value="0" />
            <el-option label="未读" :value="2" />
            <el-option label="已读" :value="1" />
          </el-select>
        </el-form-item>

        <el-form-item class="search-buttons">
          <div class="button-group" :class="{ 'mobile-buttons': isMobile }">
            <el-button type="primary" icon="search" @click="handleQuery">搜索</el-button>
            <el-button icon="refresh" @click="handleResetQuery">重置</el-button>
            <el-button type="danger" @click="truncate">清空消息</el-button>
            <el-button v-loading="readLoading" :disabled="selectIds.length == 0" @click="markReadNotice(selectIds)">标为已读</el-button>
          </div>
        </el-form-item>
      </el-form>
    </div>

    <el-card shadow="hover" class="data-table">
      <el-table
        ref="dataTableRef"
        :row-key="row => row.id"
        v-loading="loading"
        :data="pageData"
        class="data-table__content"
        @selection-change="handleSelectionChange"
        max-height="500"
      >

        <el-table-column type="selection" width="55" align="center" />

        <el-table-column label="通知标题" prop="title" min-width="200" />
        <el-table-column align="center" label="通知类型" width="150" v-if="!isMobile">
          <template #default="scope">
            <el-tag :type="scope.row.level">{{ scope.row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" label="来源" prop="source" width="150" v-if="!isMobile" />
        <el-table-column align="center" label="阅读状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.is_read ? 'success' : 'warning'">
              {{ scope.row.is_read ? '已读' : '未读' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="发布时间" width="180" v-if="!isMobile">
          <template #default="scope">
            <span>{{ formatISODateTime(scope.row.created) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" fixed="right" label="操作" width="150">
          <template #default="scope">
            <el-button type="primary" size="small" link @click="openDetailDialog(scope.row)">
              详情
            </el-button>
            <el-button
              v-if="scope.row.btn_jump_url"
              size="small"
              link
              @click="handleJump(scope.row.btn_jump_type,scope.row.btn_jump_url)"
            >
              <template v-if="scope.row.btn_desc != ''">
                {{scope.row.btn_desc}}
              </template>
              <template v-else>
                跳转
              </template>
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页组件 -->
      <div class="pagination-container" v-if="total > 0">
        <el-pagination
          background
          :current-page="queryParams.page"
          :page-size="queryParams.page_size"
          :page-sizes="[10, 20, 30, 50,100,200,500]"
          :total="total"
          :layout="isMobile?'pager':'total, sizes, prev, pager, next, jumper'"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 通知公告详情 -->
    <el-drawer
      append-to-body
      v-model="detailDialog.visible"
      :title="currentNotice.title"
      :size="isMobile ? '90%' : '50%'"
      direction="rtl"
      class="notice-drawer"
      :before-close="closeDetailDialog"
    >
      <div class="drawer-content">
        <el-descriptions  :column="1" border>
          <el-descriptions-item width="10px"  label="标题" label-class-name="description-label">
            <span class="description-value">{{ currentNotice.title }}</span>
          </el-descriptions-item>

          <el-descriptions-item  label="类型" label-class-name="description-label">
            <el-tag
              :type="currentNotice.level"
              class="notice-tag"
            >
              {{ currentNotice.type }}
            </el-tag>
          </el-descriptions-item>

          <el-descriptions-item label="来源" label-class-name="description-label">
            <span class="description-value">{{ currentNotice.source || '无' }}</span>
          </el-descriptions-item>

          <el-descriptions-item label="通知影响范围" label-class-name="description-label">
            <span class="description-value">
              <template v-if="currentNotice.target_type == 'all'">全部用户</template>
               <template v-if="currentNotice.target_type == 'roles'">指定权限组</template>
               <template v-if="currentNotice.target_type == 'users'">指定用户</template>
            </span>
          </el-descriptions-item>

          <el-descriptions-item label="是否定时任务" label-class-name="description-label">
            <span  class="description-value">{{currentNotice.is_task==1?'是':'否'}}</span>
          </el-descriptions-item>


          <el-descriptions-item label="发布时间" label-class-name="description-label">
            <span class="description-value">{{ formatISODateTime(currentNotice.created) }}</span>
          </el-descriptions-item>

          <el-descriptions-item label="内容" label-class-name="description-label">
            <div class="notice-content" v-html="currentNotice.content" />
          </el-descriptions-item>
        </el-descriptions>

        <div class="drawer-footer">
          <el-button
            type="primary"
            @click="closeDetailDialog"
            class="close-button"
          >
            关闭
          </el-button>


        </div>
      </div>
    </el-drawer>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted, nextTick, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { GetList, MarkReadNotice,Truncate } from '@/api/notice'
import { useNoticeStore} from "@/store";
import { ArrowDown } from '@element-plus/icons-vue'

const noticeStore = useNoticeStore();

import router from "@/router";

defineOptions({
  name: "Notice",
  inheritAttrs: false,
});

const queryFormRef = ref()
const loading = ref(false)
const total = ref(0)
const dataTableRef = ref()

// 移动端适配
const isMobile = computed(() => {
  return window.innerWidth < 768
})

// 筛选条件收起状态
const isSearchCollapsed = ref(isMobile.value)

// 切换筛选表单的显示/隐藏
function toggleSearchForm() {
  isSearchCollapsed.value = !isSearchCollapsed.value
}

const queryParams = reactive({
  title: '',
  read_type: 2, // 0:全部 2:未读 1:已读
  page: 1,
  page_size: 10
})

// 通知公告表格数据
const pageData = ref([])

// 跨页选中数据：id => row
const selectedMap = reactive(new Map<number, any>())
const selectIds = computed(() => Array.from(selectedMap.keys()))

// 详情弹窗
const detailDialog = reactive({
  visible: false
})
const currentNotice = ref({})

async function truncate(){
  const res = await Truncate({})

  if(res.code == 0){
    ElMessage.success(res.msg)
    handleQuery()
    return
  }
  ElMessage.error(res.msg)
}

// 查询通知公告
async function handleQuery() {
  try {
    loading.value = true
    const res = await GetList(queryParams)
    if (res.code === 0) {
      pageData.value = res.data.list
      total.value = res.data.count

      // 回显当前页选中项
      nextTick(() => {
        if (dataTableRef.value) {
          dataTableRef.value.clearSelection() // 清空旧的选择
          pageData.value.forEach(row => {
            if (selectedMap.has(row.id)) {
              dataTableRef.value.toggleRowSelection(row, true)
            }
          })
        }
      })
    } else {
      ElMessage.error(res.msg || '获取通知列表失败')
    }
  } catch (error) {
    console.error('获取通知列表失败:', error)
    ElMessage.error('获取通知列表失败')
  } finally {
    loading.value = false
  }
}

function formatISODateTime(isoString) {
  const date = new Date(isoString)
  if (isNaN(date.getTime())) return '无效日期'
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  const h = String(date.getHours()).padStart(2, '0')
  const mi = String(date.getMinutes()).padStart(2, '0')
  const s = String(date.getSeconds()).padStart(2, '0')
  return `${y}-${m}-${d} ${h}:${mi}:${s}`
}

function handleResetQuery() {
  queryFormRef.value.resetFields()
  queryParams.page = 1
  handleQuery()
}

function handleSizeChange(val: number) {
  queryParams.page_size = val
  handleQuery()
}

function handleCurrentChange(val: number) {
  queryParams.page = val
  handleQuery()
}

const readLoading = ref(false)

async function markReadNotice(ids: number[]) {
  readLoading.value = true
  await noticeStore.markAsRead(ids);
  readLoading.value = false
  ids.forEach(id => selectedMap.delete(id)) // 清除已读项
  handleQuery()
}

function handleSelectionChange(selection: any[]) {
  /*const currentPageIds = pageData.value.map(item => item.id)
  // 清除当前页的选中项
  currentPageIds.forEach(id => selectedMap.delete(id))*/
  // 添加当前页新的选中项

  selection.forEach(item => {
    selectedMap.set(item.id, item)
  })
}

function openDetailDialog(row: any) {
  currentNotice.value = row
  detailDialog.visible = true
  if (!row.is_read) {
    markReadNotice([row.id])
  }
}

function closeDetailDialog() {
  detailDialog.visible = false
  currentNotice.value = {}
}

function handleJump(typ: string, url: string) {
  if (typ === 'internal') {
    router.push({ path: url })
  } else {
    window.open(url, '_blank')
  }
}

onMounted(() => {
  handleQuery()
  
  // 监听窗口大小变化
  window.addEventListener('resize', () => {
    if (isMobile.value) {
      isSearchCollapsed.value = true
    }
  })
})
</script>
<style scoped>
.notice-content {
  padding: 10px;
  background: #f5f5f5;
  border-radius: 4px;
}

/* 新增移动端适配样式 */
.search-container {
  margin-bottom: 15px;
  border: 1px solid var(--el-border-color-light);
  border-radius: 4px;
  background-color: var(--el-bg-color);
}

.search-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 15px;
  cursor: pointer;
  border-bottom: 1px solid var(--el-border-color-lighter);
  font-weight: bold;
}

.toggle-icon {
  transition: transform 0.3s;
}

.toggle-icon.is-active {
  transform: rotate(180deg);
}

.search-form {
  padding: 15px;
}

.button-group {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.mobile-buttons {
  flex-direction: column;
  width: 100%;
}

.mobile-buttons .el-button {
  width: 100%;
  margin-left: 0;
}

@media screen and (max-width: 768px) {
  .search-form :deep(.el-form-item) {
    margin-bottom: 10px;
    width: 100%;
  }
  
  .search-form :deep(.el-form-item__content) {
    width: 100%;
  }
  
  .search-form :deep(.el-input),
  .search-form :deep(.el-select) {
    width: 100% !important;
  }
  
  .pagination-container {
    overflow-x: auto;
  }
}
</style>
<style scoped>
.notice-drawer {
  --el-drawer-padding-primary: 20px;
  --el-drawer-bg-color: var(--el-bg-color);
  --el-drawer-title-font-size: 18px;
}

.drawer-content {
  padding: 0 20px;
  height: calc(100% - 60px);
  overflow-y: auto;
}

.description-label {
  width: 100px;
  font-weight: bold;
  color: var(--el-text-color-regular);
}

.description-value {
  color: var(--el-text-color-primary);
}

.notice-tag {
  font-size: 14px;
  padding: 0 10px;
  height: 28px;
  line-height: 28px;
}

.notice-content {
  padding: 15px;
  background: var(--el-bg-color-page);
  border-radius: 4px;
  line-height: 1.6;
  color: var(--el-text-color-primary);
  border: 1px solid var(--el-border-color-light);
}

.notice-content :deep(p) {
  margin: 0 0 10px 0;
}

.notice-content :deep(a) {
  color: var(--el-color-primary);
  text-decoration: none;
}

.notice-content :deep(a:hover) {
  text-decoration: underline;
}

.drawer-footer {
  margin-top: 20px;
  text-align: right;
  padding: 10px 0;
  border-top: 1px solid var(--el-border-color-light);
}

.close-button {
  width: 100px;
}

/* 暗夜模式特定样式 */
:root.dark .notice-content {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}

:root.dark .drawer-footer {
  border-top-color: var(--el-border-color);
}

/* 移动端弹窗适配 */
@media screen and (max-width: 768px) {
  .notice-drawer :deep(.el-drawer__header) {
    padding: 12px 15px;
    margin-bottom: 10px;
    font-size: 16px;
  }
  
  .drawer-content {
    padding: 0 10px;
  }
  
  .close-button {
    width: 100%;
  }
  
  .notice-drawer :deep(.el-descriptions__label) {
    padding: 8px;
  }
  
  .notice-drawer :deep(.el-descriptions__content) {
    padding: 8px;
  }
}
</style>
