import {
  NavigationGuardNext,
  RouteLocationNormalized,
  RouteRecordRaw,
} from "vue-router";

import NProgress from "@/utils/nprogress";
import router from "@/router";
import { usePermissionStore, useUserStore } from "@/store";
import {RegisterMicroApps} from "@/utils/plugin";

import {getToken, removeToken} from "@/utils/auth";
import path from "path-browserify";
import {disconnectCentrifuge} from "@/utils/centrifuge";

const generateRoutes = (routes, findPath,basePath = '/',isNoAuthRoute)=> {
  let res = []
  for (let route of routes) {

    const onlyOneShowingChild = onlyOneShowingChildFn(route.children, route)

    if (route.children && route.children.length > 0  && onlyOneShowingChild && !route.alwaysShow) {
      route = onlyOneShowingChild
    }

    const data = {
      path: path.resolve(basePath, route.path),
      title: route.meta && route.meta.title
    }

    console.log("findPath",findPath,data.path)
    if ((route.meta.hidden && !route.meta.service) && findPath == data.path) {
      isNoAuthRoute.value = true
    }

    if (route.children && route.children.length >0 ) {
      return generateRoutes(route.children, findPath,data.path,false)
    }
    res.push(route)
  }

  return res

}

const onlyOneShowingChildFn = (children = [], parent) => {
  let onlyOneChild = null
  const showingChildren = children

  if (showingChildren.length === 1) {
    onlyOneChild = showingChildren[0]
    onlyOneChild.path =  path.resolve(parent.path, onlyOneChild.path)
    return onlyOneChild
  }

  if (showingChildren.length === 0) {
    onlyOneChild = { ...parent, path: '', noShowingChildren: true }
    return onlyOneChild
  }

  return false
}

let isInitNoLogin = false


const getNoAuthRoute = async () => {
  try {
    const permissionStore = usePermissionStore();
    const res = await permissionStore.generateNoAuthRoutes();

    let dynamicRoutes = res.dynamicRoutes
    let qiankunMicroApps = res.qiankunMicroApps

    dynamicRoutes.forEach((route: RouteRecordRaw) =>
      router.addRoute(route)
    );

    if(qiankunMicroApps !=null &&qiankunMicroApps.length >0){
      RegisterMicroApps(qiankunMicroApps)
    }

  } catch (error) {
    console.error(error)
  }
};

export async function setupPermission() {
  // 白名单路由
  const whiteList = ["/login"]

  router.beforeEach(async (to, from, next) => {
    NProgress.start();

    const hasToken = getToken();

    if (hasToken) {
      if (to.path === "/login") {
        // 如果已登录，跳转到首页
        next({ path: "/" });
        NProgress.done();
      }   else {
        const userStore = useUserStore();
        const hasRoles =
          userStore.user.roles && userStore.user.roles.length > 0;

        if (hasRoles) {
          // 如果未匹配到任何路由，跳转到404页面
          if (to.matched.length === 0) {
            next(from.name ? { name: from.name } : "/404");
          } else {
            // 如果路由参数中有 title，覆盖路由元信息中的 title
            window.document.title = to.meta.title
            next();
          }
        } else {
          const permissionStore = usePermissionStore();
          try {
            await userStore.getUserInfo();
            const res = await permissionStore.generateRoutes();
            let dynamicRoutes = res.dynamicRoutes
            let qiankunMicroApps = res.qiankunMicroApps

            dynamicRoutes.forEach((route: RouteRecordRaw) =>
              router.addRoute(route)
            );

            if(qiankunMicroApps !=null &&qiankunMicroApps.length >0){
              RegisterMicroApps(qiankunMicroApps)
            }

            next({ ...to, replace: true });
          } catch (error) {
            console.log("err",error)
            // 移除 token 并重定向到登录页，携带当前页面路由作为跳转参数
            removeToken()
            disconnectCentrifuge()
            redirectToLogin(to, next);
            NProgress.done();
          }
        }
      }
    } else {
      if (!isInitNoLogin){
        try {
          await getNoAuthRoute()
          next({ ...to, replace: true });
        } catch (error) {
          console.log("err",error)
          // 移除 token 并重定向到登录页，携带当前页面路由作为跳转参数
        }
        isInitNoLogin = true
      }else{
        if (whiteList.includes(to.path) || (to.meta.hidden && !to.meta.service)) {
          next(); // 在白名单，直接进入
        } else {
          // 不在白名单，重定向到登录页
          redirectToLogin(to, next);
          NProgress.done();
        }
      }
    }
  });

  router.afterEach(() => {
    NProgress.done();
  });
}

/** 重定向到登录页 */
function redirectToLogin(
  to: RouteLocationNormalized,
  next: NavigationGuardNext
) {
  const params = new URLSearchParams(to.query as Record<string, string>);
  const queryString = params.toString();
  const redirect = queryString ? `${to.path}?${queryString}` : to.path;
  next(`/login?redirect=${encodeURIComponent(redirect)}`);
}

/** 判断是否有权限 */
export function hasAuth(
  value: string | string[],
  type: "button" | "role" = "button"
) {
  const { roles, perms } = useUserStore().user;

  // 超级管理员 拥有所有权限
  if (type === "button" && roles.includes("ROOT")) {
    return true;
  }

  const auths = type === "button" ? perms : roles;
  return typeof value === "string"
    ? auths.includes(value)
    : value.some((perm) => auths.includes(perm));
}
