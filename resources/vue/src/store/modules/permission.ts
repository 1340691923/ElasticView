import { RouteRecordRaw } from "vue-router";
import { constantRoutes } from "@/router";
import { store } from "@/store";
import MenuAPI, { RouteVO } from "@/api/menu";
import {asyncRoutes} from "@/utils/asyncRoutes";
import path from "path-browserify";
import {NoAuthRoute, SetRootUrl} from "@/api/user";

const modules = import.meta.glob("../../views/**/**.vue");
const Layout = () => import("@/layout/index.vue");

export const usePermissionStore = defineStore("permission", () => {
  /**
   * 应用中所有的路由列表，包括静态路由和动态路由
   */
  const routes = ref<RouteRecordRaw[]>([]);
  /**
   * 混合模式左侧菜单列表
   */
  const mixLeftMenus = ref<RouteRecordRaw[]>([]);

  /**
   * 生成动态路由
   */
  function generateRoutes() {

    return new Promise<RouteRecordRaw[]>((resolve, reject) => {
      MenuAPI.getRoutes({baseRoutes:asyncRoutes})
        .then(async (res) => {

          if (res.code !== 0) {
            ElMessage.error(res.msg)
            return
          }

          let evDownloadUrl = res.data.evDownloadUrl
          let evLatestVersion = res.data.evLatestVersion
          let evUpdateAvailable = res.data.evUpdateAvailable

          if (evUpdateAvailable) {
            ElNotification({
              title: `您的Ev落后于官网最新版本:${evLatestVersion}`,
              dangerouslyUseHTMLString: true,
              duration: 0, // 不自动关闭
              showClose: true,
              message: ' <a target="_blank" href="'+evDownloadUrl+'">点我前往下载页面</a>',
            })
          }

          let data = JSON.parse(res.data.list)

          const dynamicRoutes = transformRoutes(data);
          routes.value = constantRoutes.concat(dynamicRoutes);
          resolve({
            dynamicRoutes:dynamicRoutes,
            qiankunMicroApps:res.data.qiankunMicroApps
          });
        })
        .catch((error) => {
          reject(error);
        });
    });
  }

  function generateNoAuthRoutes(fullpath) {
    return new Promise<RouteRecordRaw[]>((resolve, reject) => {
      NoAuthRoute({})
        .then((res) => {

          if (res.code !== 0) {
            ElMessage.error(res.msg)
            return
          }

          let data = res.data.list

          const dynamicRoutes = transformRoutes(data);
          resolve({
            dynamicRoutes:dynamicRoutes,
            qiankunMicroApps:res.data.qiankunMicroApps,
            isNoAuthRoute:res.data.isNoAuthRoute
          });
        })
        .catch((error) => {
          reject(error);
        });
    });
  }

  /**
   * 混合模式菜单下根据顶部菜单路径设置左侧菜单
   *
   * @param topMenuPath - 顶部菜单路径
   */
  const setMixLeftMenus = (topMenuPath: string) => {
    const matchedItem = routes.value.find((item) => item.path === topMenuPath);
    if (matchedItem && matchedItem.children) {
      mixLeftMenus.value = matchedItem.children;
    }
  };

  return {
    routes,
    generateRoutes,
    generateNoAuthRoutes,
    mixLeftMenus,
    setMixLeftMenus,
  };
});

/**
 * 转换路由数据为组件
 */
export const transformRoutes = (routes: RouteVO[]) => {
  const asyncRoutes: RouteRecordRaw[] = [];
  if (routes == null){
    return asyncRoutes;
  }
  routes.forEach((route) => {
    const tmpRoute = { ...route } as RouteRecordRaw;
    // 顶级目录，替换为 Layout 组件
    if (tmpRoute.component?.toString().toLowerCase() == "layout") {
      tmpRoute.component = Layout;
    } else {
      let component = null
      // 其他菜单，根据组件路径动态加载组件
      if (tmpRoute.component && tmpRoute.component!=''){
        component = modules[`../../${tmpRoute.component}.vue`];
      }

      if (component) {
        tmpRoute.component = component;
      }
    }

    if (tmpRoute.children) {
      tmpRoute.children = transformRoutes(route.children);
    }

    asyncRoutes.push(tmpRoute);
  });

  return asyncRoutes;
};



/**
 * 用于在组件外部（如在Pinia Store 中）使用 Pinia 提供的 store 实例。
 * 官方文档解释了如何在组件外部使用 Pinia Store：
 * https://pinia.vuejs.org/core-concepts/outside-component-usage.html#using-a-store-outside-of-a-component
 */
export function usePermissionStoreHook() {
  return usePermissionStore(store);
}
