<template>
  <div class="dashboard-container">
    <github-corner class="github-corner" />

    <el-card shadow="never">
      <el-row justify="space-between">
        <el-col :span="18" :xs="24">
          <div class="flex h-full items-center">
<!--            <img
              class="w-20 h-20 mr-5 rounded-full"
              :src="userStore.user.avatar + '?imageView2/1/w/80/h/80'"
            />-->
            <div>
              <p>{{ greetings }}</p>
              <p class="text-sm text-gray">
                但凡成功者，并非都出类拔萃，而是他们相信勤能补拙，只要忍受了，挺住了，成功早晚露出笑容
              </p>
            </div>
          </div>
        </el-col>

<!--        <el-col :span="6" :xs="24">
          <div class="flex h-full items-center justify-around">
            <el-statistic
              v-for="item in statisticData"
              :key="item.key"
              :value="item.value"
            >
              <template #title>
                <div class="flex items-center">
                  <svg-icon :icon-class="item.iconClass" size="20px" />
                  <span class="text-[16px] ml-1">{{ item.title }}</span>
                </div>
              </template>
              <template v-if="item.suffix" #suffix>/100</template>
            </el-statistic>
          </div>
        </el-col>-->
      </el-row>
      <el-row :gutter="10" class="mt-5">
        <el-col :xs="24" :sm="24" :lg="12" :span="8">
          <el-card>
            <template #header>
              <div class="flex-x-between">
                <div class="flex-y-center">
                  最新动态<el-icon class="ml-1"><Notification /></el-icon>
                </div>
              </div>
            </template>

            <el-scrollbar >
              <div
                v-for="(item, index) in articleData.list.Data"
                :key="index"
                class="flex-y-center py-3"
              >
                <el-tag :type="item.typ" size="small">
                  {{ item.tag_name }}
                </el-tag>
                <el-text
                  truncated
                  class="!mx-2 flex-1 !text-xs !text-[var(--el-text-color-secondary)]"
                >
                  {{ item.title }}
                </el-text>
                <a :href="item.link" target="_blank">
                  <el-icon class="text-sm"><View /></el-icon>
                </a>
              </div>
            </el-scrollbar>
          </el-card>
        </el-col>

      </el-row>

      <es-dashbord v-if="isEs"></es-dashbord>

    </el-card>

  </div>
</template>

<script setup lang="ts">
import {GetWxArticleList} from "@/api/plugins";

defineOptions({
  name: "Dashboard",
  inheritAttrs: false,
});

import { useUserStore } from "@/store/modules/user";
import { NoticeTypeEnum, getNoticeLabel } from "@/enums/NoticeTypeEnum";
import {IndexsCountAction,CatAction} from "@/api/es";
import {GetEsConnect,GetEsConnectVer} from "@/utils/es_link";
import EsDashbord from "@/views/dashboard/components/EsDashbord.vue";



const isEs = computed(()=>{
  return `${GetEsConnectVer()}`.indexOf('elasticsearch') !==-1
})


const userStore = useUserStore();

const date: Date = new Date();
const greetings = computed(() => {
  const hours = date.getHours();
  if (hours >= 6 && hours < 8) {
    return "晨起披衣出草堂，轩窗已自喜微凉🌅！";
  } else if (hours >= 8 && hours < 12) {
    return "上午好，" + userStore.user.username + "！";
  } else if (hours >= 12 && hours < 18) {
    return "下午好，" + userStore.user.username + "！";
  } else if (hours >= 18 && hours < 24) {
    return "晚上好，" + userStore.user.username + "！";
  } else {
    return "偷偷向银河要了一把碎星，只等你闭上眼睛撒入你的梦中，晚安🌛！";
  }
});

const notices = ref([
  {
    level: 2,
    type: NoticeTypeEnum.SYSTEM_UPGRADE,
    title: "v2.12.0 新增系统日志，访问趋势统计功能。",
  },
  {
    level: 0,
    type: NoticeTypeEnum.COMPANY_NEWS,
    title: "公司将在 7 月 1 日举办年中总结大会，请各部门做好准备。",
  },
  {
    level: 3,
    type: NoticeTypeEnum.HOLIDAY_NOTICE,
    title: "端午节假期从 6 月 12 日至 6 月 14 日放假，共 3 天。",
  },

  {
    level: 2,
    type: NoticeTypeEnum.SECURITY_ALERT,
    title: "最近发现一些钓鱼邮件，请大家提高警惕，不要点击陌生链接。",
  },
  {
    level: 2,
    type: NoticeTypeEnum.SYSTEM_MAINTENANCE,
    title: "系统将于本周六凌晨 2 点进行维护，预计维护时间为 2 小时。",
  },
  {
    level: 0,
    type: NoticeTypeEnum.OTHER,
    title: "公司新规章制度发布，请大家及时查阅。",
  },
  {
    level: 3,
    type: NoticeTypeEnum.HOLIDAY_NOTICE,
    title: "中秋节假期从 9 月 22 日至 9 月 24 日放假，共 3 天。",
  },
  {
    level: 1,
    type: NoticeTypeEnum.COMPANY_NEWS,
    title: "公司将在 10 月 15 日举办新产品发布会，敬请期待。",
  },
  {
    level: 2,
    type: NoticeTypeEnum.SECURITY_ALERT,
    title:
      "请注意，近期有恶意软件通过即时通讯工具传播，请勿下载不明来源的文件。",
  },
  {
    level: 2,
    type: NoticeTypeEnum.SYSTEM_MAINTENANCE,
    title: "系统将于下周日凌晨 3 点进行升级，预计维护时间为 1 小时。",
  },
  {
    level: 3,
    type: NoticeTypeEnum.OTHER,
    title: "公司年度体检通知已发布，请各位员工按时参加。",
  },
]);

const articleData = reactive({
  list:[],
})



const GetArticleData = async ()=>{
  //await SearchBigMode({content:"mysql如何拼接字符串"})

  let res =  await GetWxArticleList({})
  if(res.code != 0){
    ElMessage.error(res.msg);
    return
  }

  articleData.list = res.data

  return
}


onMounted(()=>{
  GetArticleData()
})

</script>

<style lang="scss" scoped>
.dashboard-container {
  position: relative;
  padding: 24px;

  .github-corner {
    position: absolute;
    top: 0;
    right: 0;
    z-index: 1;
    border: 0;
  }
}
</style>
