import { setupDirective } from "@/directive";
import { setupRouter } from "@/router";
import { setupStore } from "@/store";
import type { App } from "vue";
import { setupElIcons } from "./icons";
import { setupPermission } from "./permission";

export default {
  async install(app: App<Element>) {
    // 自定义指令(directive)
    setupDirective(app);
    // 路由(router)

    // 状态管理(store)
    await setupStore(app);

    // 国际化
    //setupI18n(app);
    // Element-plus图标
    setupElIcons(app);


    setupRouter(app);
    // 路由守卫
    setupPermission();
  },
};
