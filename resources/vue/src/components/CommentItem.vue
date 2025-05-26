<template>
  <div class="comment-item">
    <div class="comment-header">
      <span class="comment-user">{{ comment.realname }}</span>
      <span class="comment-time">{{ formattedDate }}</span>
    </div>
    <div class="comment-content">{{ comment.content }}</div>
    <div class="comment-footer">
      <div class="comment-actions">
        <el-button
          type="text"
          size="small"
          class="like-button"
          :class="{ 'is-liked': comment.has_like }"
          @click.stop="handleLike"
        >
          <el-icon><Pointer /></el-icon>
          <span>{{ comment.like_count }}</span>
        </el-button>
        <el-button
          type="text"
          size="small"
          class="reply-button"
          @click="handleReply"
        >
          回复
        </el-button>
        <!-- 折叠/展开按钮 -->
        <el-button
          v-if="comment.children && comment.children.length > 0"
          type="text"
          size="small"
          class="toggle-button"
          @click="handleToggle"
        >
          <span v-if="isCollapsed">展开 ({{ comment.children.length }}条回复)</span>
          <span v-else>收起回复</span>
        </el-button>
      </div>
    </div>

    <!-- 回复输入框 -->
    <div v-if="replyVisible" class="reply-input-container">
      <div class="reply-to-tip">回复: {{ comment.realname }}</div>
      <el-input
        v-model="replyContent"
        type="textarea"
        :rows="2"
        placeholder="发表回复..."
        maxlength="200"
        show-word-limit
      />
      <div class="emoji-container">
        <el-button
          type="text"
          size="small"
          class="emoji-toggle"
          @click="showEmojiPicker = !showEmojiPicker"
        >
          😊 表情
        </el-button>
        <div v-if="showEmojiPicker" style="z-index:300000" class="emoji-picker">
          <div
            v-for="emoji in emojiList"
            :key="emoji"
            class="emoji-item"
            @click="insertEmoji(emoji)"
          >
            {{ emoji }}
          </div>
        </div>
      </div>
      <div class="reply-actions">
        <el-button size="small" @click="cancelReply">取消</el-button>
        <el-button type="primary" size="small" @click="submitReply">发送</el-button>
      </div>
    </div>

    <!-- 子评论 -->
    <div v-if="comment.children && comment.children.length > 0 && !isCollapsed && canShowChildren" class="comment-replies">
      <!-- 递归渲染子评论 -->
      <comment-item
        v-for="child in comment.children"
        :key="child.id"
        :comment="child"
        :plugin-id="pluginId"
        :level="level + 1"
        @refresh="$emit('refresh')"
      />
    </div>
    <import-ev-key v-model:visible="importEvkeyDialogVisible" ></import-ev-key>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Pointer } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { AddComment, LikeComment } from '@/api/plugins'
import { useAppStore } from "@/store"
import { DeviceEnum } from "@/enums/DeviceEnum"
import ImportEvKey from '@/components/ImportEvKey/index.vue'

const appStore = useAppStore()
const isMobile = computed(() => appStore.device === DeviceEnum.MOBILE)

const props = defineProps({
  comment: {
    type: Object,
    required: true
  },
  pluginId: {
    type: [Number, String],
    required: true
  },
  level: {
    type: Number,
    default: 1
  }
})

const importEvkeyDialogVisible = ref(false);

const openImportEvkeyDialogVisible = () => {
  importEvkeyDialogVisible.value = true;
};

const emit = defineEmits(['refresh'])

const replyVisible = ref(false)
const replyContent = ref('')
const isCollapsed = ref(true) // 默认折叠子评论
const showEmojiPicker = ref(false) // emoji选择器显示状态
const likeLoading = ref(false) // 点赞加载状态

// 最大嵌套层级 (移动端)
const MAX_MOBILE_LEVEL = 3

// 判断是否可以显示子评论（移动端限制层级）
const canShowChildren = computed(() => {
  if (!isMobile.value) return true
  return props.level < MAX_MOBILE_LEVEL
})

// 处理点赞
const handleLike = async () => {
  if (likeLoading.value) return

  console.log('点赞按钮被点击，评论ID:', props.comment.id)
  likeLoading.value = true

  // 先优化用户体验，立即改变UI状态
  const originalLikeState = props.comment.has_like
  const originalLikeCount = props.comment.like_count

  // 先在UI上反映变化
  if (props.comment.has_like) {
    props.comment.like_count = Math.max(0, props.comment.like_count - 1)
    props.comment.has_like = false
  } else {
    props.comment.like_count += 1
    props.comment.has_like = true
  }

  try {
    const res = await LikeComment({
      comment_id: props.comment.id
    })


    if (res.msg.indexOf('请前往')!==-1){
      openImportEvkeyDialogVisible()
      return
    }

    if (res.code !== 0) {
      ElMessage.error({
        type: 'error',
        offset: 60,
        message: res.msg || '操作失败'
      })

      // 恢复原始状态
      props.comment.has_like = originalLikeState
      props.comment.like_count = originalLikeCount
      return
    }

  } catch (error) {
    console.error('点赞失败', error)
    ElMessage.error({
      type: 'error',
      offset: 60,
      message: '点赞失败'
    })

    // 恢复原始状态
    props.comment.has_like = originalLikeState
    props.comment.like_count = originalLikeCount
  } finally {
    likeLoading.value = false
  }
}

// 处理折叠/展开按钮点击
const handleToggle = () => {
  if (isMobile.value && props.level >= MAX_MOBILE_LEVEL - 1 && props.comment.children?.length > 0) {
    // 移动端且评论层级太深，显示提示
    ElMessageBox.alert('请在PC端查看更多嵌套回复', '提示', {
      confirmButtonText: '确定',
      type: 'info'
    })
  } else {
    // 正常切换折叠状态
    isCollapsed.value = !isCollapsed.value
  }
}

// 格式化日期显示
const formattedDate = computed(() => {
  if (!props.comment.created_at) return ''

  if (isMobile.value) {
    // 移动端只显示月-日
    const date = new Date(props.comment.created_at)
    const month = (date.getMonth() + 1).toString().padStart(2, '0')
    const day = date.getDate().toString().padStart(2, '0')
    return `${month}-${day}`
  } else {
    // 桌面端显示完整日期
    return props.comment.created_at
  }
})

const handleReply = () => {
  replyVisible.value = true
}

const cancelReply = () => {
  replyVisible.value = false
  replyContent.value = ''
}

const submitReply = async () => {
  if (!replyContent.value.trim()) {
    ElMessage.warning({
      type: 'warning',
      offset: 60,
      message: '回复内容不能为空'
    })
    return
  }

  try {
    const res = await AddComment({
      plugin_id: props.pluginId,
      content: replyContent.value,
      parent_id: props.comment.id
    })


    if (res.msg.indexOf('请前往')!==-1){
      openImportEvkeyDialogVisible()
      return
    }

    if (res.code !== 0) {
      ElMessage.error({
        type: 'error',
        offset: 60,
        message: res.msg
      })
      return
    }

    ElMessage.success({
      type: 'success',
      offset: 60,
      message: '回复成功'
    })

    // 回复成功后自动展开子评论
    isCollapsed.value = false

    cancelReply()
    emit('refresh') // 通知父组件刷新评论列表
  } catch (error) {
    console.error('回复失败', error)
    ElMessage.error({
      type: 'error',
      offset: 60,
      message: '回复失败'
    })
  }
}

// emoji表情列表
const emojiList = [
  '😀', '😃', '😄', '😁', '😆', '😅', '😂', '🤣', '😊', '😇',
  '🙂', '🙃', '😉', '😌', '😍', '🥰', '😘', '😗', '😙', '😚',
  '😋', '😛', '😝', '😜', '🤪', '🤨', '🧐', '🤓', '😎', '🤩',
  '🥳', '😏', '😒', '😞', '😔', '😟', '😕', '🙁', '☹️', '😣',
  '😖', '😫', '😩', '🥺', '😢', '😭', '😤', '😠', '😡', '🤬',
  '🤯', '😳', '🥵', '🥶', '😱', '😨', '😰', '😥', '😓', '🤗',
  '🤔', '🤭', '🤫', '🤥', '😶', '😐', '😑', '😬', '🙄', '😯',
  '😦', '😧', '😮', '😲', '🥱', '😴', '🤤', '😪', '😵', '🤐',
  '🥴', '🤢', '🤮', '🤧', '😷', '🤒', '🤕', '🤑', '🤠', '👍',
  '👎', '👏', '🙏', '🤝', '💪', '❤️', '💔', '💯', '✨', '🔥'
]

// 向回复框插入emoji
const insertEmoji = (emoji) => {
  replyContent.value += emoji
  showEmojiPicker.value = false // 选择后关闭emoji选择器
}


</script>

<style lang="scss" scoped>
.comment-item {
  @apply p-4;
  @apply bg-white/50 dark:bg-gray-800/50;
  @apply backdrop-blur-sm;
  @apply border border-gray-200 dark:border-gray-700;
  @apply rounded-lg;
  @apply transition-all duration-300;
  @apply hover:shadow-md;
  @apply mb-4;
}

.comment-header {
  @apply flex items-center justify-between;
  @apply mb-2;

  .comment-user {
    @apply font-medium;
    @apply text-gray-900 dark:text-gray-100;
  }

  .comment-time {
    @apply text-sm;
    @apply text-gray-500 dark:text-gray-400;
  }
}

.comment-content {
  @apply py-2;
  @apply text-gray-700 dark:text-gray-300;
}

.comment-footer {
  @apply flex justify-between items-center;
  @apply mt-2;

  .comment-actions {
    @apply flex items-center gap-2;
  }

  .like-button, .reply-button, .toggle-button {
    @apply flex items-center gap-1;
    @apply text-gray-500 dark:text-gray-400;
    @apply hover:text-blue-500 dark:hover:text-blue-400;
    @apply transition-colors duration-300;

    &.is-liked {
      @apply text-red-500 dark:text-red-400;
      @apply font-medium;
    }
  }

  .toggle-button {
    @apply text-blue-500 dark:text-blue-400;
  }
}

.reply-input-container {
  @apply mt-3 mb-3 pl-4 border-l-2 border-gray-200 dark:border-gray-700;

  .reply-to-tip {
    @apply text-sm font-medium text-blue-500 mb-2;
  }

  .reply-actions {
    @apply flex justify-end mt-2 gap-2;
  }
}

.emoji-container {
  @apply relative mt-2;

  .emoji-toggle {
    @apply text-blue-500;
  }

      .emoji-picker {
      @apply bg-white dark:bg-gray-800;
      @apply rounded-lg shadow-lg;
      @apply border border-gray-200 dark:border-gray-700;
      @apply flex flex-wrap;
      @apply overflow-y-auto;
      @apply absolute bottom-full left-0;
      @apply p-2;
      @apply mb-2;
      width: 280px;
      max-height: 180px;

      .emoji-item {
        @apply p-2 cursor-pointer text-xl;
        @apply hover:bg-gray-100 dark:hover:bg-gray-700;
        @apply rounded transition-colors;
        @apply flex items-center justify-center;
        width: 40px;
        height: 40px;
      }
  }
}

.comment-replies {
  @apply mt-4 ml-2;
  @apply space-y-4;
  @apply border-l-2 border-gray-200 dark:border-gray-700;
  @apply pl-2;
}
</style>
