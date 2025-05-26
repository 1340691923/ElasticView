import { defineStore } from "pinia";
import { store } from "@/store";
import { GetList, MarkReadNotice } from '@/api/notice';
import { ElMessage } from 'element-plus';

export const useNoticeStore = defineStore("notice", () => {
  // 所有通知列表
  const noticeList = ref([]);
  
  // 获取通知列表
  async function fetchNoticeList(params) {
    try {
      const res = await GetList(params);
      if (res.code === 0) {
        if(res.data.list == null) res.data.list = [];
        noticeList.value = res.data.list;
      } else {
        ElMessage.error(res.msg || '获取通知列表失败');
      }
    } catch (error) {
      console.error('获取通知列表失败:', error);
      ElMessage.error('获取通知列表失败');
    }
  }

  // 添加新通知
  function addNotice(notice) {
    if (!noticeList.value.some(item => item.id == notice.id)) {
      noticeList.value.unshift(notice);
    }
  }

  // 标记通知为已读
  async function markAsRead(ids) {
    try {
      const res = await MarkReadNotice({ ids });
      if (res.code === 0) {
        // 从列表中移除已读通知
        noticeList.value = noticeList.value.filter(notice => !ids.includes(notice.id));
        return res;
      } else {
        ElMessage.error(res.msg || '标记已读失败');
        return res;
      }
    } catch (error) {
      console.error('标记已读失败:', error);
      ElMessage.error('标记已读失败');
      throw error;
    }
  }

  return {
    noticeList,
    fetchNoticeList,
    addNotice,
    markAsRead
  };
});

/**
 * 用于在组件外部使用 notice store
 */
export function useNoticeStoreHook() {
  return useNoticeStore(store);
} 