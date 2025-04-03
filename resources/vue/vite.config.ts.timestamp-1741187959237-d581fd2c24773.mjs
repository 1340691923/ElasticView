// vite.config.ts
import vue from "file:///D:/eve/ev2/resources/vue/node_modules/.pnpm/@vitejs+plugin-vue@5.1.2_vite@5.3.5_@types+node@20.14.14_sass@1.77.8_terser@5.31.3__vue@3.4.35_typescript@5.5.4_/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import vueJsx from "file:///D:/eve/ev2/resources/vue/node_modules/.pnpm/@vitejs+plugin-vue-jsx@3.1.0_vite@5.3.5_@types+node@20.14.14_sass@1.77.8_terser@5.31.3__vue@3.4.35_typescript@5.5.4_/node_modules/@vitejs/plugin-vue-jsx/dist/index.mjs";
import { loadEnv, defineConfig } from "file:///D:/eve/ev2/resources/vue/node_modules/.pnpm/vite@5.3.5_@types+node@20.14.14_sass@1.77.8_terser@5.31.3/node_modules/vite/dist/node/index.js";
import { createHtmlPlugin } from "file:///D:/eve/ev2/resources/vue/node_modules/.pnpm/vite-plugin-html@3.2.2_vite@5.3.5_@types+node@20.14.14_sass@1.77.8_terser@5.31.3_/node_modules/vite-plugin-html/dist/index.mjs";
import AutoImport from "file:///D:/eve/ev2/resources/vue/node_modules/.pnpm/unplugin-auto-import@0.17.8_@vueuse+core@10.11.0_vue@3.4.35_typescript@5.5.4___rollup@4.20.0/node_modules/unplugin-auto-import/dist/vite.js";
import Components from "file:///D:/eve/ev2/resources/vue/node_modules/.pnpm/unplugin-vue-components@0.26.0_@babel+parser@7.25.3_rollup@4.20.0_vue@3.4.35_typescript@5.5.4_/node_modules/unplugin-vue-components/dist/vite.js";
import { ElementPlusResolver } from "file:///D:/eve/ev2/resources/vue/node_modules/.pnpm/unplugin-vue-components@0.26.0_@babel+parser@7.25.3_rollup@4.20.0_vue@3.4.35_typescript@5.5.4_/node_modules/unplugin-vue-components/dist/resolvers.js";
import Icons from "file:///D:/eve/ev2/resources/vue/node_modules/.pnpm/unplugin-icons@0.18.5_@vue+compiler-sfc@3.4.35/node_modules/unplugin-icons/dist/vite.js";
import IconsResolver from "file:///D:/eve/ev2/resources/vue/node_modules/.pnpm/unplugin-icons@0.18.5_@vue+compiler-sfc@3.4.35/node_modules/unplugin-icons/dist/resolver.js";
import { createSvgIconsPlugin } from "file:///D:/eve/ev2/resources/vue/node_modules/.pnpm/vite-plugin-svg-icons@2.0.1_vite@5.3.5_@types+node@20.14.14_sass@1.77.8_terser@5.31.3_/node_modules/vite-plugin-svg-icons/dist/index.mjs";
import mockDevServerPlugin from "file:///D:/eve/ev2/resources/vue/node_modules/.pnpm/vite-plugin-mock-dev-server@1.6.1_bufferutil@4.0.8_rollup@4.20.0_utf-8-validate@5.0.10_vite@5_m7qfccemwb525qavxcujshylnq/node_modules/vite-plugin-mock-dev-server/dist/index.js";
import { dynamicBase } from "file:///D:/eve/ev2/resources/vue/node_modules/.pnpm/vite-plugin-dynamic-base@1.1.0_vite@5.3.5_@types+node@20.14.14_sass@1.77.8_terser@5.31.3_/node_modules/vite-plugin-dynamic-base/dist/index.js";
import UnoCSS from "file:///D:/eve/ev2/resources/vue/node_modules/.pnpm/unocss@0.58.9_postcss@8.4.40_rollup@4.20.0_vite@5.3.5_@types+node@20.14.14_sass@1.77.8_terser@5.31.3_/node_modules/unocss/dist/vite.mjs";
import { resolve } from "path";

// package.json
var name = "ElasticView";
var version = "0.0.1";
var dependencies = {
  "@element-plus/icons-vue": "^2.3.1",
  "@vavt/cm-extension": "^1.5.0",
  "@vueuse/core": "^10.11.0",
  "@wangeditor/editor": "^5.1.23",
  "@wangeditor/editor-for-vue": "5.1.10",
  "animate.css": "^4.1.1",
  axios: "^1.7.2",
  centrifuge: "^5.3.2",
  color: "^4.2.3",
  dayjs: "^1.11.12",
  echarts: "^5.5.0",
  "element-plus": "^2.7.6",
  exceljs: "^4.4.0",
  "lodash-es": "^4.17.21",
  "md-editor-v3": "^4.20.1",
  "monaco-editor": "^0.50.0",
  net: "^1.0.2",
  nprogress: "^0.2.0",
  "path-browserify": "^1.0.1",
  "path-to-regexp": "^6.2.2",
  pinia: "^2.1.7",
  qiankun: "^2.10.16",
  "sockjs-client": "1.6.1",
  sortablejs: "^1.15.2",
  "sql-formatter": "^15.4.0",
  stompjs: "^2.3.3",
  "vite-plugin-dynamic-base": "^1.1.0",
  "vite-plugin-html": "^3.2.2",
  vue: "^3.4.31",
  "vue-i18n": "9.9.1",
  "vue-router": "^4.4.0",
  "vue3-json-viewer": "^2.2.2",
  "vue3-markdown-it": "^1.0.10"
};
var devDependencies = {
  "@commitlint/cli": "^18.6.1",
  "@commitlint/config-conventional": "^18.6.3",
  "@iconify-json/ep": "^1.1.15",
  "@types/color": "^3.0.6",
  "@types/lodash": "^4.17.6",
  "@types/node": "^20.14.10",
  "@types/nprogress": "^0.2.3",
  "@types/path-browserify": "^1.0.2",
  "@types/sockjs-client": "^1.5.4",
  "@types/sortablejs": "^1.15.8",
  "@types/stompjs": "^2.3.9",
  "@typescript-eslint/eslint-plugin": "^7.15.0",
  "@typescript-eslint/parser": "^7.15.0",
  "@vitejs/plugin-vue": "^5.0.5",
  "@vitejs/plugin-vue-jsx": "^3.1.0",
  autoprefixer: "^10.4.19",
  commitizen: "^4.3.0",
  "cz-git": "^1.9.3",
  eslint: "^8.57.0",
  "eslint-config-prettier": "^9.1.0",
  "eslint-plugin-import": "^2.29.1",
  "eslint-plugin-prettier": "^5.1.3",
  "eslint-plugin-vue": "^9.27.0",
  "fast-glob": "^3.3.2",
  husky: "^9.0.11",
  "lint-staged": "^15.2.7",
  postcss: "^8.4.39",
  "postcss-html": "^1.7.0",
  "postcss-scss": "^4.0.9",
  prettier: "^3.3.2",
  sass: "^1.77.6",
  stylelint: "^16.6.1",
  "stylelint-config-html": "^1.1.0",
  "stylelint-config-recess-order": "^4.6.0",
  "stylelint-config-recommended-scss": "^14.0.0",
  "stylelint-config-recommended-vue": "^1.5.0",
  "stylelint-config-standard": "^36.0.1",
  terser: "^5.31.1",
  typescript: "^5.5.3",
  unocss: "^0.58.9",
  "unplugin-auto-import": "^0.17.6",
  "unplugin-icons": "^0.18.5",
  "unplugin-vue-components": "^0.26.0",
  vite: "^5.3.3",
  "vite-plugin-mock-dev-server": "^1.5.1",
  "vite-plugin-monaco-editor": "^1.1.0",
  "vite-plugin-svg-icons": "^2.0.1",
  "vite-plugin-vue-devtools": "^7.3.5",
  "vue-tsc": "^2.0.26"
};
var engines = {
  node: ">=18.0.0"
};

// vite.config.ts
var __vite_injected_original_dirname = "D:\\eve\\ev2\\resources\\vue";
var __APP_INFO__ = {
  pkg: { name, version, engines, dependencies, devDependencies },
  buildTimestamp: Date.now()
};
var pathSrc = resolve(__vite_injected_original_dirname, "src");
var vite_config_default = defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd());
  return {
    resolve: {
      alias: {
        "@": pathSrc
      }
    },
    base: process.env.NODE_ENV === "production" ? "/__dynamic_base__/" : "/",
    css: {
      // CSS 预处理器
      preprocessorOptions: {
        // 定义全局 SCSS 变量
        scss: {
          javascriptEnabled: true,
          additionalData: `
            @use "@/styles/variables.scss" as *;
          `
        }
      }
    },
    server: {
      // 允许IP访问
      host: "0.0.0.0",
      // 应用端口 (默认:3000)
      port: Number(env.VITE_APP_PORT),
      // 运行是否自动打开浏览器
      open: true,
      proxy: {
        /** 代理前缀为 /dev-api 的请求  */
        [env.VITE_APP_BASE_API]: {
          changeOrigin: true,
          // 接口地址
          target: env.VITE_APP_API_URL,
          rewrite: (path) => path.replace(new RegExp("^" + env.VITE_APP_BASE_API), "")
        }
      }
    },
    plugins: [
      vue(),
      dynamicBase({
        publicPath: "window.appSubUrl",
        transformIndexHtml: true
      }),
      createHtmlPlugin({
        inject: {
          data: {
            isProd: env.VITE_ISPROD,
            VITE_APPURL: env.VITE_APPURL,
            VITE_APPSubUrl: env.VITE_APPSubUrl,
            VITE_APPVersion: env.VITE_APPVersion,
            VITE_Lang: env.VITE_Lang,
            title: "ElasticView"
          }
        }
      }),
      // jsx、tsx语法支持
      vueJsx(),
      // MOCK 服务
      env.VITE_MOCK_DEV_SERVER === "true" ? mockDevServerPlugin() : null,
      UnoCSS({
        hmrTopLevelAwait: false
      }),
      // 自动导入参考： https://github.com/sxzz/element-plus-best-practices/blob/main/vite.config.ts
      AutoImport({
        // 自动导入 Vue 相关函数，如：ref, reactive, toRef 等
        imports: ["vue", "@vueuse/core", "pinia", "vue-router", "vue-i18n"],
        resolvers: [
          // 自动导入 Element Plus 相关函数，如：ElMessage, ElMessageBox... (带样式)
          ElementPlusResolver(),
          // 自动导入图标组件
          IconsResolver({})
        ],
        eslintrc: {
          // 是否自动生成 eslint 规则，建议生成之后设置 false
          enabled: false,
          // 指定自动导入函数 eslint 规则的文件
          filepath: "./.eslintrc-auto-import.json",
          globalsPropValue: true
        },
        // 是否在 vue 模板中自动导入
        vueTemplate: true,
        // 指定自动导入函数TS类型声明文件路径 (false:关闭自动生成)
        dts: false
        // dts: "src/types/auto-imports.d.ts",
      }),
      Components({
        resolvers: [
          // 自动导入 Element Plus 组件
          ElementPlusResolver(),
          // 自动注册图标组件
          IconsResolver({
            // element-plus图标库，其他图标库 https://icon-sets.iconify.design/
            enabledCollections: ["ep"]
          })
        ],
        // 指定自定义组件位置(默认:src/components)
        dirs: ["src/components", "src/**/components"],
        // 指定自动导入组件TS类型声明文件路径 (false:关闭自动生成)
        dts: false
        // dts: "src/types/components.d.ts",
      }),
      Icons({
        // 自动安装图标库
        autoInstall: true
      }),
      createSvgIconsPlugin({
        // 指定需要缓存的图标文件夹
        iconDirs: [resolve(pathSrc, "assets/icons")],
        // 指定symbolId格式
        symbolId: "icon-[dir]-[name]"
      })
      /* VueDevTools({
        openInEditorHost: `http://localhost:${env.VITE_APP_PORT}`,
      }), */
    ],
    // 预加载项目必需的组件
    optimizeDeps: {
      include: [
        "vue",
        "vue-router",
        "pinia",
        "axios",
        "@vueuse/core",
        "sortablejs",
        "path-to-regexp",
        "echarts",
        "@wangeditor/editor",
        "@wangeditor/editor-for-vue",
        "vue-i18n",
        "path-browserify",
        `monaco-editor/esm/vs/language/json/json.worker`,
        `monaco-editor/esm/vs/language/css/css.worker`,
        `monaco-editor/esm/vs/language/html/html.worker`,
        `monaco-editor/esm/vs/language/typescript/ts.worker`,
        `monaco-editor/esm/vs/editor/editor.worker`,
        "element-plus/es/components/form/style/css",
        "element-plus/es/components/form-item/style/css",
        "element-plus/es/components/button/style/css",
        "element-plus/es/components/input/style/css",
        "element-plus/es/components/input-number/style/css",
        "element-plus/es/components/switch/style/css",
        "element-plus/es/components/upload/style/css",
        "element-plus/es/components/menu/style/css",
        "element-plus/es/components/col/style/css",
        "element-plus/es/components/icon/style/css",
        "element-plus/es/components/row/style/css",
        "element-plus/es/components/tag/style/css",
        "element-plus/es/components/dialog/style/css",
        "element-plus/es/components/loading/style/css",
        "element-plus/es/components/radio/style/css",
        "element-plus/es/components/radio-group/style/css",
        "element-plus/es/components/popover/style/css",
        "element-plus/es/components/scrollbar/style/css",
        "element-plus/es/components/tooltip/style/css",
        "element-plus/es/components/dropdown/style/css",
        "element-plus/es/components/dropdown-menu/style/css",
        "element-plus/es/components/dropdown-item/style/css",
        "element-plus/es/components/sub-menu/style/css",
        "element-plus/es/components/menu-item/style/css",
        "element-plus/es/components/divider/style/css",
        "element-plus/es/components/card/style/css",
        "element-plus/es/components/link/style/css",
        "element-plus/es/components/breadcrumb/style/css",
        "element-plus/es/components/breadcrumb-item/style/css",
        "element-plus/es/components/table/style/css",
        "element-plus/es/components/tree-select/style/css",
        "element-plus/es/components/table-column/style/css",
        "element-plus/es/components/select/style/css",
        "element-plus/es/components/option/style/css",
        "element-plus/es/components/pagination/style/css",
        "element-plus/es/components/tree/style/css",
        "element-plus/es/components/alert/style/css",
        "element-plus/es/components/radio-button/style/css",
        "element-plus/es/components/checkbox-group/style/css",
        "element-plus/es/components/checkbox/style/css",
        "element-plus/es/components/tabs/style/css",
        "element-plus/es/components/tab-pane/style/css",
        "element-plus/es/components/rate/style/css",
        "element-plus/es/components/date-picker/style/css",
        "element-plus/es/components/notification/style/css",
        "element-plus/es/components/image/style/css",
        "element-plus/es/components/statistic/style/css",
        "element-plus/es/components/watermark/style/css",
        "element-plus/es/components/config-provider/style/css",
        "element-plus/es/components/text/style/css",
        "element-plus/es/components/drawer/style/css",
        "element-plus/es/components/color-picker/style/css",
        "element-plus/es/components/backtop/style/css",
        "element-plus/es/components/message-box/style/css",
        "element-plus/es/components/skeleton/style/css",
        "element-plus/es/components/skeleton/style/css",
        "element-plus/es/components/skeleton-item/style/css",
        "element-plus/es/components/badge/style/css"
      ]
    },
    // 构建配置
    build: {
      outDir: "../views/dist",
      assetsDir: "static",
      chunkSizeWarningLimit: 2e3,
      // 消除打包大小超过500kb警告
      minify: "terser",
      // Vite 2.6.x 以上需要配置 minify: "terser", terserOptions 才能生效
      terserOptions: {
        compress: {
          keep_infinity: true
          // 防止 Infinity 被压缩成 1/0，这可能会导致 Chrome 上的性能问题
          //drop_console: true, // 生产环境去除 console
          //drop_debugger: true, // 生产环境去除 debugger
        },
        format: {
          comments: false
          // 删除注释
        }
      },
      rollupOptions: {
        output: {
          // manualChunks: {
          //   "vue-i18n": ["vue-i18n"],
          // },
          // 用于从入口点创建的块的打包输出格式[name]表示文件名,[hash]表示该文件内容hash值
          entryFileNames: "js/[name].[hash].js",
          // 用于命名代码拆分时创建的共享块的输出命名
          chunkFileNames: "js/[name].[hash].js",
          // 用于输出静态资源的命名，[ext]表示文件扩展名
          assetFileNames: (assetInfo) => {
            const info = assetInfo.name.split(".");
            let extType = info[info.length - 1];
            if (/\.(mp4|webm|ogg|mp3|wav|flac|aac)(\?.*)?$/i.test(assetInfo.name)) {
              extType = "media";
            } else if (/\.(png|jpe?g|gif|svg)(\?.*)?$/.test(assetInfo.name)) {
              extType = "img";
            } else if (/\.(woff2?|eot|ttf|otf)(\?.*)?$/i.test(assetInfo.name)) {
              extType = "fonts";
            }
            return `${extType}/[name].[hash].[ext]`;
          }
        }
      }
    },
    define: {
      __APP_INFO__: JSON.stringify(__APP_INFO__)
    }
  };
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiLCAicGFja2FnZS5qc29uIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyJjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZGlybmFtZSA9IFwiRDpcXFxcZXZlXFxcXGV2MlxcXFxyZXNvdXJjZXNcXFxcdnVlXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCJEOlxcXFxldmVcXFxcZXYyXFxcXHJlc291cmNlc1xcXFx2dWVcXFxcdml0ZS5jb25maWcudHNcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfaW1wb3J0X21ldGFfdXJsID0gXCJmaWxlOi8vL0Q6L2V2ZS9ldjIvcmVzb3VyY2VzL3Z1ZS92aXRlLmNvbmZpZy50c1wiO2ltcG9ydCB2dWUgZnJvbSBcIkB2aXRlanMvcGx1Z2luLXZ1ZVwiO1xuaW1wb3J0IHZ1ZUpzeCBmcm9tIFwiQHZpdGVqcy9wbHVnaW4tdnVlLWpzeFwiO1xuaW1wb3J0IHsgVXNlckNvbmZpZywgQ29uZmlnRW52LCBsb2FkRW52LCBkZWZpbmVDb25maWcgfSBmcm9tIFwidml0ZVwiO1xuaW1wb3J0IHsgY3JlYXRlSHRtbFBsdWdpbiB9IGZyb20gXCJ2aXRlLXBsdWdpbi1odG1sXCJcblxuaW1wb3J0IEF1dG9JbXBvcnQgZnJvbSBcInVucGx1Z2luLWF1dG8taW1wb3J0L3ZpdGVcIjtcbmltcG9ydCBDb21wb25lbnRzIGZyb20gXCJ1bnBsdWdpbi12dWUtY29tcG9uZW50cy92aXRlXCI7XG5pbXBvcnQgeyBFbGVtZW50UGx1c1Jlc29sdmVyIH0gZnJvbSBcInVucGx1Z2luLXZ1ZS1jb21wb25lbnRzL3Jlc29sdmVyc1wiO1xuaW1wb3J0IEljb25zIGZyb20gXCJ1bnBsdWdpbi1pY29ucy92aXRlXCI7XG5pbXBvcnQgSWNvbnNSZXNvbHZlciBmcm9tIFwidW5wbHVnaW4taWNvbnMvcmVzb2x2ZXJcIjtcblxuaW1wb3J0IHsgY3JlYXRlU3ZnSWNvbnNQbHVnaW4gfSBmcm9tIFwidml0ZS1wbHVnaW4tc3ZnLWljb25zXCI7XG5pbXBvcnQgbW9ja0RldlNlcnZlclBsdWdpbiBmcm9tIFwidml0ZS1wbHVnaW4tbW9jay1kZXYtc2VydmVyXCI7XG5pbXBvcnQgeyBkeW5hbWljQmFzZSB9IGZyb20gJ3ZpdGUtcGx1Z2luLWR5bmFtaWMtYmFzZSdcblxuaW1wb3J0IFVub0NTUyBmcm9tIFwidW5vY3NzL3ZpdGVcIjtcbmltcG9ydCB7IHJlc29sdmUgfSBmcm9tIFwicGF0aFwiO1xuaW1wb3J0IHtcbiAgbmFtZSxcbiAgdmVyc2lvbixcbiAgZW5naW5lcyxcbiAgZGVwZW5kZW5jaWVzLFxuICBkZXZEZXBlbmRlbmNpZXMsXG59IGZyb20gXCIuL3BhY2thZ2UuanNvblwiO1xuXG4vLyBodHRwczovL2RldnRvb2xzLW5leHQudnVlanMub3JnL1xuaW1wb3J0IFZ1ZURldlRvb2xzIGZyb20gXCJ2aXRlLXBsdWdpbi12dWUtZGV2dG9vbHNcIjtcblxuLyoqIFx1NUU3M1x1NTNGMFx1NzY4NFx1NTQwRFx1NzlGMFx1MzAwMVx1NzI0OFx1NjcyQ1x1MzAwMVx1OEZEMFx1ODg0Q1x1NjI0MFx1OTcwMFx1NzY4NGBub2RlYFx1NzI0OFx1NjcyQ1x1MzAwMVx1NEY5RFx1OEQ1Nlx1MzAwMVx1Njc4NFx1NUVGQVx1NjVGNlx1OTVGNFx1NzY4NFx1N0M3Qlx1NTc4Qlx1NjNEMFx1NzkzQSAqL1xuY29uc3QgX19BUFBfSU5GT19fID0ge1xuICBwa2c6IHsgbmFtZSwgdmVyc2lvbiwgZW5naW5lcywgZGVwZW5kZW5jaWVzLCBkZXZEZXBlbmRlbmNpZXMgfSxcbiAgYnVpbGRUaW1lc3RhbXA6IERhdGUubm93KCksXG59O1xuXG5jb25zdCBwYXRoU3JjID0gcmVzb2x2ZShfX2Rpcm5hbWUsIFwic3JjXCIpO1xuLy8gIGh0dHBzOi8vY24udml0ZWpzLmRldi9jb25maWdcbmV4cG9ydCBkZWZhdWx0IGRlZmluZUNvbmZpZygoeyBtb2RlIH06IENvbmZpZ0Vudik6IFVzZXJDb25maWcgPT4ge1xuICBjb25zdCBlbnYgPSBsb2FkRW52KG1vZGUsIHByb2Nlc3MuY3dkKCkpO1xuICByZXR1cm4ge1xuICAgIHJlc29sdmU6IHtcbiAgICAgIGFsaWFzOiB7XG4gICAgICAgIFwiQFwiOiBwYXRoU3JjLFxuICAgICAgfSxcbiAgICB9LFxuICAgIGJhc2U6IHByb2Nlc3MuZW52Lk5PREVfRU5WID09PSBcInByb2R1Y3Rpb25cIiA/IFwiL19fZHluYW1pY19iYXNlX18vXCIgOiBcIi9cIixcblxuICAgIGNzczoge1xuICAgICAgLy8gQ1NTIFx1OTg4NFx1NTkwNFx1NzQwNlx1NTY2OFxuICAgICAgcHJlcHJvY2Vzc29yT3B0aW9uczoge1xuICAgICAgICAvLyBcdTVCOUFcdTRFNDlcdTUxNjhcdTVDNDAgU0NTUyBcdTUzRDhcdTkxQ0ZcbiAgICAgICAgc2Nzczoge1xuICAgICAgICAgIGphdmFzY3JpcHRFbmFibGVkOiB0cnVlLFxuICAgICAgICAgIGFkZGl0aW9uYWxEYXRhOiBgXG4gICAgICAgICAgICBAdXNlIFwiQC9zdHlsZXMvdmFyaWFibGVzLnNjc3NcIiBhcyAqO1xuICAgICAgICAgIGAsXG4gICAgICAgIH0sXG4gICAgICB9LFxuICAgIH0sXG4gICAgc2VydmVyOiB7XG4gICAgICAvLyBcdTUxNDFcdThCQjhJUFx1OEJCRlx1OTVFRVxuICAgICAgaG9zdDogXCIwLjAuMC4wXCIsXG4gICAgICAvLyBcdTVFOTRcdTc1MjhcdTdBRUZcdTUzRTMgKFx1OUVEOFx1OEJBNDozMDAwKVxuICAgICAgcG9ydDogTnVtYmVyKGVudi5WSVRFX0FQUF9QT1JUKSxcbiAgICAgIC8vIFx1OEZEMFx1ODg0Q1x1NjYyRlx1NTQyNlx1ODFFQVx1NTJBOFx1NjI1M1x1NUYwMFx1NkQ0Rlx1ODlDOFx1NTY2OFxuICAgICAgb3BlbjogdHJ1ZSxcbiAgICAgIHByb3h5OiB7XG4gICAgICAgIC8qKiBcdTRFRTNcdTc0MDZcdTUyNERcdTdGMDBcdTRFM0EgL2Rldi1hcGkgXHU3Njg0XHU4QkY3XHU2QzQyICAqL1xuICAgICAgICBbZW52LlZJVEVfQVBQX0JBU0VfQVBJXToge1xuICAgICAgICAgIGNoYW5nZU9yaWdpbjogdHJ1ZSxcbiAgICAgICAgICAvLyBcdTYzQTVcdTUzRTNcdTU3MzBcdTU3NDBcbiAgICAgICAgICB0YXJnZXQ6IGVudi5WSVRFX0FQUF9BUElfVVJMLFxuICAgICAgICAgIHJld3JpdGU6IChwYXRoKSA9PlxuICAgICAgICAgICAgcGF0aC5yZXBsYWNlKG5ldyBSZWdFeHAoXCJeXCIgKyBlbnYuVklURV9BUFBfQkFTRV9BUEkpLCBcIlwiKSxcbiAgICAgICAgfSxcbiAgICAgIH0sXG4gICAgfSxcbiAgICBwbHVnaW5zOiBbXG4gICAgICB2dWUoKSxcbiAgICAgIGR5bmFtaWNCYXNlKHtcbiAgICAgICAgcHVibGljUGF0aDogJ3dpbmRvdy5hcHBTdWJVcmwnLFxuICAgICAgICB0cmFuc2Zvcm1JbmRleEh0bWw6ICB0cnVlXG4gICAgICB9KSxcbiAgICAgIGNyZWF0ZUh0bWxQbHVnaW4oe1xuICAgICAgICBpbmplY3Q6IHtcbiAgICAgICAgICBkYXRhOiB7XG4gICAgICAgICAgICBpc1Byb2Q6IGVudi5WSVRFX0lTUFJPRCxcbiAgICAgICAgICAgIFZJVEVfQVBQVVJMOiBlbnYuVklURV9BUFBVUkwsXG4gICAgICAgICAgICBWSVRFX0FQUFN1YlVybDogZW52LlZJVEVfQVBQU3ViVXJsLFxuICAgICAgICAgICAgVklURV9BUFBWZXJzaW9uOiBlbnYuVklURV9BUFBWZXJzaW9uLFxuICAgICAgICAgICAgVklURV9MYW5nOiBlbnYuVklURV9MYW5nLFxuICAgICAgICAgICAgdGl0bGU6J0VsYXN0aWNWaWV3J1xuICAgICAgICAgIH0sXG4gICAgICAgIH0sXG4gICAgICB9KSxcbiAgICAgIC8vIGpzeFx1MzAwMXRzeFx1OEJFRFx1NkNENVx1NjUyRlx1NjMwMVxuICAgICAgdnVlSnN4KCksXG4gICAgICAvLyBNT0NLIFx1NjcwRFx1NTJBMVxuICAgICAgZW52LlZJVEVfTU9DS19ERVZfU0VSVkVSID09PSBcInRydWVcIiA/IG1vY2tEZXZTZXJ2ZXJQbHVnaW4oKSA6IG51bGwsXG4gICAgICBVbm9DU1Moe1xuICAgICAgICBobXJUb3BMZXZlbEF3YWl0OiBmYWxzZSxcbiAgICAgIH0pLFxuICAgICAgLy8gXHU4MUVBXHU1MkE4XHU1QkZDXHU1MTY1XHU1M0MyXHU4MDAzXHVGRjFBIGh0dHBzOi8vZ2l0aHViLmNvbS9zeHp6L2VsZW1lbnQtcGx1cy1iZXN0LXByYWN0aWNlcy9ibG9iL21haW4vdml0ZS5jb25maWcudHNcbiAgICAgIEF1dG9JbXBvcnQoe1xuICAgICAgICAvLyBcdTgxRUFcdTUyQThcdTVCRkNcdTUxNjUgVnVlIFx1NzZGOFx1NTE3M1x1NTFGRFx1NjU3MFx1RkYwQ1x1NTk4Mlx1RkYxQXJlZiwgcmVhY3RpdmUsIHRvUmVmIFx1N0I0OVxuICAgICAgICBpbXBvcnRzOiBbXCJ2dWVcIiwgXCJAdnVldXNlL2NvcmVcIiwgXCJwaW5pYVwiLCBcInZ1ZS1yb3V0ZXJcIiwgXCJ2dWUtaTE4blwiXSxcbiAgICAgICAgcmVzb2x2ZXJzOiBbXG4gICAgICAgICAgLy8gXHU4MUVBXHU1MkE4XHU1QkZDXHU1MTY1IEVsZW1lbnQgUGx1cyBcdTc2RjhcdTUxNzNcdTUxRkRcdTY1NzBcdUZGMENcdTU5ODJcdUZGMUFFbE1lc3NhZ2UsIEVsTWVzc2FnZUJveC4uLiAoXHU1RTI2XHU2ODM3XHU1RjBGKVxuICAgICAgICAgIEVsZW1lbnRQbHVzUmVzb2x2ZXIoKSxcbiAgICAgICAgICAvLyBcdTgxRUFcdTUyQThcdTVCRkNcdTUxNjVcdTU2RkVcdTY4MDdcdTdFQzRcdTRFRjZcbiAgICAgICAgICBJY29uc1Jlc29sdmVyKHt9KSxcbiAgICAgICAgXSxcbiAgICAgICAgZXNsaW50cmM6IHtcbiAgICAgICAgICAvLyBcdTY2MkZcdTU0MjZcdTgxRUFcdTUyQThcdTc1MUZcdTYyMTAgZXNsaW50IFx1ODlDNFx1NTIxOVx1RkYwQ1x1NUVGQVx1OEJBRVx1NzUxRlx1NjIxMFx1NEU0Qlx1NTQwRVx1OEJCRVx1N0Y2RSBmYWxzZVxuICAgICAgICAgIGVuYWJsZWQ6IGZhbHNlLFxuICAgICAgICAgIC8vIFx1NjMwN1x1NUI5QVx1ODFFQVx1NTJBOFx1NUJGQ1x1NTE2NVx1NTFGRFx1NjU3MCBlc2xpbnQgXHU4OUM0XHU1MjE5XHU3Njg0XHU2NTg3XHU0RUY2XG4gICAgICAgICAgZmlsZXBhdGg6IFwiLi8uZXNsaW50cmMtYXV0by1pbXBvcnQuanNvblwiLFxuICAgICAgICAgIGdsb2JhbHNQcm9wVmFsdWU6IHRydWUsXG4gICAgICAgIH0sXG4gICAgICAgIC8vIFx1NjYyRlx1NTQyNlx1NTcyOCB2dWUgXHU2QTIxXHU2NzdGXHU0RTJEXHU4MUVBXHU1MkE4XHU1QkZDXHU1MTY1XG4gICAgICAgIHZ1ZVRlbXBsYXRlOiB0cnVlLFxuICAgICAgICAvLyBcdTYzMDdcdTVCOUFcdTgxRUFcdTUyQThcdTVCRkNcdTUxNjVcdTUxRkRcdTY1NzBUU1x1N0M3Qlx1NTc4Qlx1NThGMFx1NjYwRVx1NjU4N1x1NEVGNlx1OERFRlx1NUY4NCAoZmFsc2U6XHU1MTczXHU5NUVEXHU4MUVBXHU1MkE4XHU3NTFGXHU2MjEwKVxuICAgICAgICBkdHM6IGZhbHNlLFxuICAgICAgICAvLyBkdHM6IFwic3JjL3R5cGVzL2F1dG8taW1wb3J0cy5kLnRzXCIsXG4gICAgICB9KSxcbiAgICAgIENvbXBvbmVudHMoe1xuICAgICAgICByZXNvbHZlcnM6IFtcbiAgICAgICAgICAvLyBcdTgxRUFcdTUyQThcdTVCRkNcdTUxNjUgRWxlbWVudCBQbHVzIFx1N0VDNFx1NEVGNlxuICAgICAgICAgIEVsZW1lbnRQbHVzUmVzb2x2ZXIoKSxcbiAgICAgICAgICAvLyBcdTgxRUFcdTUyQThcdTZDRThcdTUxOENcdTU2RkVcdTY4MDdcdTdFQzRcdTRFRjZcbiAgICAgICAgICBJY29uc1Jlc29sdmVyKHtcbiAgICAgICAgICAgIC8vIGVsZW1lbnQtcGx1c1x1NTZGRVx1NjgwN1x1NUU5M1x1RkYwQ1x1NTE3Nlx1NEVENlx1NTZGRVx1NjgwN1x1NUU5MyBodHRwczovL2ljb24tc2V0cy5pY29uaWZ5LmRlc2lnbi9cbiAgICAgICAgICAgIGVuYWJsZWRDb2xsZWN0aW9uczogW1wiZXBcIl0sXG4gICAgICAgICAgfSksXG4gICAgICAgIF0sXG4gICAgICAgIC8vIFx1NjMwN1x1NUI5QVx1ODFFQVx1NUI5QVx1NEU0OVx1N0VDNFx1NEVGNlx1NEY0RFx1N0Y2RShcdTlFRDhcdThCQTQ6c3JjL2NvbXBvbmVudHMpXG4gICAgICAgIGRpcnM6IFtcInNyYy9jb21wb25lbnRzXCIsIFwic3JjLyoqL2NvbXBvbmVudHNcIl0sXG4gICAgICAgIC8vIFx1NjMwN1x1NUI5QVx1ODFFQVx1NTJBOFx1NUJGQ1x1NTE2NVx1N0VDNFx1NEVGNlRTXHU3QzdCXHU1NzhCXHU1OEYwXHU2NjBFXHU2NTg3XHU0RUY2XHU4REVGXHU1Rjg0IChmYWxzZTpcdTUxNzNcdTk1RURcdTgxRUFcdTUyQThcdTc1MUZcdTYyMTApXG4gICAgICAgIGR0czogZmFsc2UsXG4gICAgICAgIC8vIGR0czogXCJzcmMvdHlwZXMvY29tcG9uZW50cy5kLnRzXCIsXG4gICAgICB9KSxcbiAgICAgIEljb25zKHtcbiAgICAgICAgLy8gXHU4MUVBXHU1MkE4XHU1Qjg5XHU4OEM1XHU1NkZFXHU2ODA3XHU1RTkzXG4gICAgICAgIGF1dG9JbnN0YWxsOiB0cnVlLFxuICAgICAgfSksXG4gICAgICBjcmVhdGVTdmdJY29uc1BsdWdpbih7XG4gICAgICAgIC8vIFx1NjMwN1x1NUI5QVx1OTcwMFx1ODk4MVx1N0YxM1x1NUI1OFx1NzY4NFx1NTZGRVx1NjgwN1x1NjU4N1x1NEVGNlx1NTkzOVxuICAgICAgICBpY29uRGlyczogW3Jlc29sdmUocGF0aFNyYywgXCJhc3NldHMvaWNvbnNcIildLFxuICAgICAgICAvLyBcdTYzMDdcdTVCOUFzeW1ib2xJZFx1NjgzQ1x1NUYwRlxuICAgICAgICBzeW1ib2xJZDogXCJpY29uLVtkaXJdLVtuYW1lXVwiLFxuICAgICAgfSksXG4gICAgICAvKiBWdWVEZXZUb29scyh7XG4gICAgICAgIG9wZW5JbkVkaXRvckhvc3Q6IGBodHRwOi8vbG9jYWxob3N0OiR7ZW52LlZJVEVfQVBQX1BPUlR9YCxcbiAgICAgIH0pLCAqL1xuICAgIF0sXG4gICAgLy8gXHU5ODg0XHU1MkEwXHU4RjdEXHU5ODc5XHU3NkVFXHU1RkM1XHU5NzAwXHU3Njg0XHU3RUM0XHU0RUY2XG4gICAgb3B0aW1pemVEZXBzOiB7XG4gICAgICBpbmNsdWRlOiBbXG4gICAgICAgIFwidnVlXCIsXG4gICAgICAgIFwidnVlLXJvdXRlclwiLFxuICAgICAgICBcInBpbmlhXCIsXG4gICAgICAgIFwiYXhpb3NcIixcbiAgICAgICAgXCJAdnVldXNlL2NvcmVcIixcbiAgICAgICAgXCJzb3J0YWJsZWpzXCIsXG4gICAgICAgIFwicGF0aC10by1yZWdleHBcIixcbiAgICAgICAgXCJlY2hhcnRzXCIsXG4gICAgICAgIFwiQHdhbmdlZGl0b3IvZWRpdG9yXCIsXG4gICAgICAgIFwiQHdhbmdlZGl0b3IvZWRpdG9yLWZvci12dWVcIixcbiAgICAgICAgXCJ2dWUtaTE4blwiLFxuICAgICAgICBcInBhdGgtYnJvd3NlcmlmeVwiLFxuXG4gICAgICAgIGBtb25hY28tZWRpdG9yL2VzbS92cy9sYW5ndWFnZS9qc29uL2pzb24ud29ya2VyYCxcbiAgICAgICAgYG1vbmFjby1lZGl0b3IvZXNtL3ZzL2xhbmd1YWdlL2Nzcy9jc3Mud29ya2VyYCxcbiAgICAgICAgYG1vbmFjby1lZGl0b3IvZXNtL3ZzL2xhbmd1YWdlL2h0bWwvaHRtbC53b3JrZXJgLFxuICAgICAgICBgbW9uYWNvLWVkaXRvci9lc20vdnMvbGFuZ3VhZ2UvdHlwZXNjcmlwdC90cy53b3JrZXJgLFxuICAgICAgICBgbW9uYWNvLWVkaXRvci9lc20vdnMvZWRpdG9yL2VkaXRvci53b3JrZXJgLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2Zvcm0vc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvZm9ybS1pdGVtL3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2J1dHRvbi9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9pbnB1dC9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9pbnB1dC1udW1iZXIvc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvc3dpdGNoL3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3VwbG9hZC9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9tZW51L3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2NvbC9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9pY29uL3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3Jvdy9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy90YWcvc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvZGlhbG9nL3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2xvYWRpbmcvc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvcmFkaW8vc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvcmFkaW8tZ3JvdXAvc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvcG9wb3Zlci9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9zY3JvbGxiYXIvc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvdG9vbHRpcC9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9kcm9wZG93bi9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9kcm9wZG93bi1tZW51L3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2Ryb3Bkb3duLWl0ZW0vc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvc3ViLW1lbnUvc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvbWVudS1pdGVtL3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2RpdmlkZXIvc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvY2FyZC9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9saW5rL3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2JyZWFkY3J1bWIvc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvYnJlYWRjcnVtYi1pdGVtL3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3RhYmxlL3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3RyZWUtc2VsZWN0L3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3RhYmxlLWNvbHVtbi9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9zZWxlY3Qvc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvb3B0aW9uL3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3BhZ2luYXRpb24vc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvdHJlZS9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9hbGVydC9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9yYWRpby1idXR0b24vc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvY2hlY2tib3gtZ3JvdXAvc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvY2hlY2tib3gvc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvdGFicy9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy90YWItcGFuZS9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9yYXRlL3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2RhdGUtcGlja2VyL3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL25vdGlmaWNhdGlvbi9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9pbWFnZS9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9zdGF0aXN0aWMvc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvd2F0ZXJtYXJrL3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2NvbmZpZy1wcm92aWRlci9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy90ZXh0L3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2RyYXdlci9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9jb2xvci1waWNrZXIvc3R5bGUvY3NzXCIsXG4gICAgICAgIFwiZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvYmFja3RvcC9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9tZXNzYWdlLWJveC9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9za2VsZXRvbi9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9za2VsZXRvbi9zdHlsZS9jc3NcIixcbiAgICAgICAgXCJlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9za2VsZXRvbi1pdGVtL3N0eWxlL2Nzc1wiLFxuICAgICAgICBcImVsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2JhZGdlL3N0eWxlL2Nzc1wiLFxuICAgICAgXSxcbiAgICB9LFxuICAgIC8vIFx1Njc4NFx1NUVGQVx1OTE0RFx1N0Y2RVxuICAgIGJ1aWxkOiB7XG4gICAgICBvdXREaXI6Jy4uL3ZpZXdzL2Rpc3QnLFxuICAgICAgYXNzZXRzRGlyOiAnc3RhdGljJyxcbiAgICAgIGNodW5rU2l6ZVdhcm5pbmdMaW1pdDogMjAwMCwgLy8gXHU2RDg4XHU5NjY0XHU2MjUzXHU1MzA1XHU1OTI3XHU1QzBGXHU4RDg1XHU4RkM3NTAwa2JcdThCNjZcdTU0NEFcbiAgICAgIG1pbmlmeTogXCJ0ZXJzZXJcIiwgLy8gVml0ZSAyLjYueCBcdTRFRTVcdTRFMEFcdTk3MDBcdTg5ODFcdTkxNERcdTdGNkUgbWluaWZ5OiBcInRlcnNlclwiLCB0ZXJzZXJPcHRpb25zIFx1NjI0RFx1ODBGRFx1NzUxRlx1NjU0OFxuICAgICAgdGVyc2VyT3B0aW9uczoge1xuICAgICAgICBjb21wcmVzczoge1xuICAgICAgICAgIGtlZXBfaW5maW5pdHk6IHRydWUsIC8vIFx1OTYzMlx1NkI2MiBJbmZpbml0eSBcdTg4QUJcdTUzOEJcdTdGMjlcdTYyMTAgMS8wXHVGRjBDXHU4RkQ5XHU1M0VGXHU4MEZEXHU0RjFBXHU1QkZDXHU4MUY0IENocm9tZSBcdTRFMEFcdTc2ODRcdTYwMjdcdTgwRkRcdTk1RUVcdTk4OThcbiAgICAgICAgICAvL2Ryb3BfY29uc29sZTogdHJ1ZSwgLy8gXHU3NTFGXHU0RUE3XHU3M0FGXHU1ODgzXHU1M0JCXHU5NjY0IGNvbnNvbGVcbiAgICAgICAgICAvL2Ryb3BfZGVidWdnZXI6IHRydWUsIC8vIFx1NzUxRlx1NEVBN1x1NzNBRlx1NTg4M1x1NTNCQlx1OTY2NCBkZWJ1Z2dlclxuICAgICAgICB9LFxuICAgICAgICBmb3JtYXQ6IHtcbiAgICAgICAgICBjb21tZW50czogZmFsc2UsIC8vIFx1NTIyMFx1OTY2NFx1NkNFOFx1OTFDQVxuICAgICAgICB9LFxuICAgICAgfSxcbiAgICAgIHJvbGx1cE9wdGlvbnM6IHtcbiAgICAgICAgb3V0cHV0OiB7XG4gICAgICAgICAgLy8gbWFudWFsQ2h1bmtzOiB7XG4gICAgICAgICAgLy8gICBcInZ1ZS1pMThuXCI6IFtcInZ1ZS1pMThuXCJdLFxuICAgICAgICAgIC8vIH0sXG4gICAgICAgICAgLy8gXHU3NTI4XHU0RThFXHU0RUNFXHU1MTY1XHU1M0UzXHU3MEI5XHU1MjFCXHU1RUZBXHU3Njg0XHU1NzU3XHU3Njg0XHU2MjUzXHU1MzA1XHU4RjkzXHU1MUZBXHU2ODNDXHU1RjBGW25hbWVdXHU4ODY4XHU3OTNBXHU2NTg3XHU0RUY2XHU1NDBELFtoYXNoXVx1ODg2OFx1NzkzQVx1OEJFNVx1NjU4N1x1NEVGNlx1NTE4NVx1NUJCOWhhc2hcdTUwM0NcbiAgICAgICAgICBlbnRyeUZpbGVOYW1lczogXCJqcy9bbmFtZV0uW2hhc2hdLmpzXCIsXG4gICAgICAgICAgLy8gXHU3NTI4XHU0RThFXHU1NDdEXHU1NDBEXHU0RUUzXHU3ODAxXHU2MkM2XHU1MjA2XHU2NUY2XHU1MjFCXHU1RUZBXHU3Njg0XHU1MTcxXHU0RUFCXHU1NzU3XHU3Njg0XHU4RjkzXHU1MUZBXHU1NDdEXHU1NDBEXG4gICAgICAgICAgY2h1bmtGaWxlTmFtZXM6IFwianMvW25hbWVdLltoYXNoXS5qc1wiLFxuICAgICAgICAgIC8vIFx1NzUyOFx1NEU4RVx1OEY5M1x1NTFGQVx1OTc1OVx1NjAwMVx1OEQ0NFx1NkU5MFx1NzY4NFx1NTQ3RFx1NTQwRFx1RkYwQ1tleHRdXHU4ODY4XHU3OTNBXHU2NTg3XHU0RUY2XHU2MjY5XHU1QzU1XHU1NDBEXG4gICAgICAgICAgYXNzZXRGaWxlTmFtZXM6IChhc3NldEluZm86IGFueSkgPT4ge1xuICAgICAgICAgICAgY29uc3QgaW5mbyA9IGFzc2V0SW5mby5uYW1lLnNwbGl0KFwiLlwiKTtcbiAgICAgICAgICAgIGxldCBleHRUeXBlID0gaW5mb1tpbmZvLmxlbmd0aCAtIDFdO1xuICAgICAgICAgICAgLy8gY29uc29sZS5sb2coJ1x1NjU4N1x1NEVGNlx1NEZFMVx1NjA2RicsIGFzc2V0SW5mby5uYW1lKVxuICAgICAgICAgICAgaWYgKFxuICAgICAgICAgICAgICAvXFwuKG1wNHx3ZWJtfG9nZ3xtcDN8d2F2fGZsYWN8YWFjKShcXD8uKik/JC9pLnRlc3QoYXNzZXRJbmZvLm5hbWUpXG4gICAgICAgICAgICApIHtcbiAgICAgICAgICAgICAgZXh0VHlwZSA9IFwibWVkaWFcIjtcbiAgICAgICAgICAgIH0gZWxzZSBpZiAoL1xcLihwbmd8anBlP2d8Z2lmfHN2ZykoXFw/LiopPyQvLnRlc3QoYXNzZXRJbmZvLm5hbWUpKSB7XG4gICAgICAgICAgICAgIGV4dFR5cGUgPSBcImltZ1wiO1xuICAgICAgICAgICAgfSBlbHNlIGlmICgvXFwuKHdvZmYyP3xlb3R8dHRmfG90ZikoXFw/LiopPyQvaS50ZXN0KGFzc2V0SW5mby5uYW1lKSkge1xuICAgICAgICAgICAgICBleHRUeXBlID0gXCJmb250c1wiO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgcmV0dXJuIGAke2V4dFR5cGV9L1tuYW1lXS5baGFzaF0uW2V4dF1gO1xuICAgICAgICAgIH0sXG4gICAgICAgIH0sXG4gICAgICB9LFxuICAgIH0sXG4gICAgZGVmaW5lOiB7XG4gICAgICBfX0FQUF9JTkZPX186IEpTT04uc3RyaW5naWZ5KF9fQVBQX0lORk9fXyksXG4gICAgfSxcbiAgfTtcbn0pO1xuIiwgIntcbiAgXCJuYW1lXCI6IFwiRWxhc3RpY1ZpZXdcIixcbiAgXCJ2ZXJzaW9uXCI6IFwiMC4wLjFcIixcbiAgXCJwcml2YXRlXCI6IHRydWUsXG4gIFwidHlwZVwiOiBcIm1vZHVsZVwiLFxuICBcInNjcmlwdHNcIjoge1xuICAgIFwiZGV2XCI6IFwidml0ZVwiLFxuICAgIFwiYnVpbGRcIjogXCJ2aXRlIGJ1aWxkIC0tbW9kZSBwcm9kdWN0aW9uIC0tZW1wdHlPdXREaXJcIixcbiAgICBcInByZXZpZXdcIjogXCJ2aXRlIHByZXZpZXdcIixcbiAgICBcImJ1aWxkLW9ubHlcIjogXCJ2aXRlIGJ1aWxkXCIsXG4gICAgXCJ0eXBlLWNoZWNrXCI6IFwidnVlLXRzYyAtLW5vRW1pdFwiLFxuICAgIFwibGludDplc2xpbnRcIjogXCJlc2xpbnQgIC0tZml4IC0tZXh0IC50cywuanMsLnZ1ZSAuL3NyYyBcIixcbiAgICBcImxpbnQ6cHJldHRpZXJcIjogXCJwcmV0dGllciAtLXdyaXRlIFxcXCIqKi8qLntqcyxjanMsdHMsanNvbix0c3gsY3NzLGxlc3Msc2Nzcyx2dWUsaHRtbCxtZH1cXFwiXCIsXG4gICAgXCJsaW50OnN0eWxlbGludFwiOiBcInN0eWxlbGludCAgXFxcIioqLyoue2NzcyxzY3NzLHZ1ZX1cXFwiIC0tZml4XCIsXG4gICAgXCJsaW50OmxpbnQtc3RhZ2VkXCI6IFwibGludC1zdGFnZWRcIixcbiAgICBcInByZWluc3RhbGxcIjogXCJucHggb25seS1hbGxvdyBwbnBtXCIsXG4gICAgXCJwcmVwYXJlXCI6IFwiaHVza3lcIixcbiAgICBcImNvbW1pdFwiOiBcImdpdC1jelwiXG4gIH0sXG4gIFwiY29uZmlnXCI6IHtcbiAgICBcImNvbW1pdGl6ZW5cIjoge1xuICAgICAgXCJwYXRoXCI6IFwibm9kZV9tb2R1bGVzL2N6LWdpdFwiXG4gICAgfVxuICB9LFxuICBcImxpbnQtc3RhZ2VkXCI6IHtcbiAgICBcIioue2pzLHRzfVwiOiBbXG4gICAgICBcImVzbGludCAtLWZpeFwiLFxuICAgICAgXCJwcmV0dGllciAtLXdyaXRlXCJcbiAgICBdLFxuICAgIFwiKi57Y2pzLGpzb259XCI6IFtcbiAgICAgIFwicHJldHRpZXIgLS13cml0ZVwiXG4gICAgXSxcbiAgICBcIioue3Z1ZSxodG1sfVwiOiBbXG4gICAgICBcImVzbGludCAtLWZpeFwiLFxuICAgICAgXCJwcmV0dGllciAtLXdyaXRlXCIsXG4gICAgICBcInN0eWxlbGludCAtLWZpeFwiXG4gICAgXSxcbiAgICBcIioue3Njc3MsY3NzfVwiOiBbXG4gICAgICBcInN0eWxlbGludCAtLWZpeFwiLFxuICAgICAgXCJwcmV0dGllciAtLXdyaXRlXCJcbiAgICBdLFxuICAgIFwiKi5tZFwiOiBbXG4gICAgICBcInByZXR0aWVyIC0td3JpdGVcIlxuICAgIF1cbiAgfSxcbiAgXCJkZXBlbmRlbmNpZXNcIjoge1xuICAgIFwiQGVsZW1lbnQtcGx1cy9pY29ucy12dWVcIjogXCJeMi4zLjFcIixcbiAgICBcIkB2YXZ0L2NtLWV4dGVuc2lvblwiOiBcIl4xLjUuMFwiLFxuICAgIFwiQHZ1ZXVzZS9jb3JlXCI6IFwiXjEwLjExLjBcIixcbiAgICBcIkB3YW5nZWRpdG9yL2VkaXRvclwiOiBcIl41LjEuMjNcIixcbiAgICBcIkB3YW5nZWRpdG9yL2VkaXRvci1mb3ItdnVlXCI6IFwiNS4xLjEwXCIsXG4gICAgXCJhbmltYXRlLmNzc1wiOiBcIl40LjEuMVwiLFxuICAgIFwiYXhpb3NcIjogXCJeMS43LjJcIixcbiAgICBcImNlbnRyaWZ1Z2VcIjogXCJeNS4zLjJcIixcbiAgICBcImNvbG9yXCI6IFwiXjQuMi4zXCIsXG4gICAgXCJkYXlqc1wiOiBcIl4xLjExLjEyXCIsXG4gICAgXCJlY2hhcnRzXCI6IFwiXjUuNS4wXCIsXG4gICAgXCJlbGVtZW50LXBsdXNcIjogXCJeMi43LjZcIixcbiAgICBcImV4Y2VsanNcIjogXCJeNC40LjBcIixcbiAgICBcImxvZGFzaC1lc1wiOiBcIl40LjE3LjIxXCIsXG4gICAgXCJtZC1lZGl0b3ItdjNcIjogXCJeNC4yMC4xXCIsXG4gICAgXCJtb25hY28tZWRpdG9yXCI6IFwiXjAuNTAuMFwiLFxuICAgIFwibmV0XCI6IFwiXjEuMC4yXCIsXG4gICAgXCJucHJvZ3Jlc3NcIjogXCJeMC4yLjBcIixcbiAgICBcInBhdGgtYnJvd3NlcmlmeVwiOiBcIl4xLjAuMVwiLFxuICAgIFwicGF0aC10by1yZWdleHBcIjogXCJeNi4yLjJcIixcbiAgICBcInBpbmlhXCI6IFwiXjIuMS43XCIsXG4gICAgXCJxaWFua3VuXCI6IFwiXjIuMTAuMTZcIixcbiAgICBcInNvY2tqcy1jbGllbnRcIjogXCIxLjYuMVwiLFxuICAgIFwic29ydGFibGVqc1wiOiBcIl4xLjE1LjJcIixcbiAgICBcInNxbC1mb3JtYXR0ZXJcIjogXCJeMTUuNC4wXCIsXG4gICAgXCJzdG9tcGpzXCI6IFwiXjIuMy4zXCIsXG4gICAgXCJ2aXRlLXBsdWdpbi1keW5hbWljLWJhc2VcIjogXCJeMS4xLjBcIixcbiAgICBcInZpdGUtcGx1Z2luLWh0bWxcIjogXCJeMy4yLjJcIixcbiAgICBcInZ1ZVwiOiBcIl4zLjQuMzFcIixcbiAgICBcInZ1ZS1pMThuXCI6IFwiOS45LjFcIixcbiAgICBcInZ1ZS1yb3V0ZXJcIjogXCJeNC40LjBcIixcbiAgICBcInZ1ZTMtanNvbi12aWV3ZXJcIjogXCJeMi4yLjJcIixcbiAgICBcInZ1ZTMtbWFya2Rvd24taXRcIjogXCJeMS4wLjEwXCJcbiAgfSxcbiAgXCJkZXZEZXBlbmRlbmNpZXNcIjoge1xuICAgIFwiQGNvbW1pdGxpbnQvY2xpXCI6IFwiXjE4LjYuMVwiLFxuICAgIFwiQGNvbW1pdGxpbnQvY29uZmlnLWNvbnZlbnRpb25hbFwiOiBcIl4xOC42LjNcIixcbiAgICBcIkBpY29uaWZ5LWpzb24vZXBcIjogXCJeMS4xLjE1XCIsXG4gICAgXCJAdHlwZXMvY29sb3JcIjogXCJeMy4wLjZcIixcbiAgICBcIkB0eXBlcy9sb2Rhc2hcIjogXCJeNC4xNy42XCIsXG4gICAgXCJAdHlwZXMvbm9kZVwiOiBcIl4yMC4xNC4xMFwiLFxuICAgIFwiQHR5cGVzL25wcm9ncmVzc1wiOiBcIl4wLjIuM1wiLFxuICAgIFwiQHR5cGVzL3BhdGgtYnJvd3NlcmlmeVwiOiBcIl4xLjAuMlwiLFxuICAgIFwiQHR5cGVzL3NvY2tqcy1jbGllbnRcIjogXCJeMS41LjRcIixcbiAgICBcIkB0eXBlcy9zb3J0YWJsZWpzXCI6IFwiXjEuMTUuOFwiLFxuICAgIFwiQHR5cGVzL3N0b21wanNcIjogXCJeMi4zLjlcIixcbiAgICBcIkB0eXBlc2NyaXB0LWVzbGludC9lc2xpbnQtcGx1Z2luXCI6IFwiXjcuMTUuMFwiLFxuICAgIFwiQHR5cGVzY3JpcHQtZXNsaW50L3BhcnNlclwiOiBcIl43LjE1LjBcIixcbiAgICBcIkB2aXRlanMvcGx1Z2luLXZ1ZVwiOiBcIl41LjAuNVwiLFxuICAgIFwiQHZpdGVqcy9wbHVnaW4tdnVlLWpzeFwiOiBcIl4zLjEuMFwiLFxuICAgIFwiYXV0b3ByZWZpeGVyXCI6IFwiXjEwLjQuMTlcIixcbiAgICBcImNvbW1pdGl6ZW5cIjogXCJeNC4zLjBcIixcbiAgICBcImN6LWdpdFwiOiBcIl4xLjkuM1wiLFxuICAgIFwiZXNsaW50XCI6IFwiXjguNTcuMFwiLFxuICAgIFwiZXNsaW50LWNvbmZpZy1wcmV0dGllclwiOiBcIl45LjEuMFwiLFxuICAgIFwiZXNsaW50LXBsdWdpbi1pbXBvcnRcIjogXCJeMi4yOS4xXCIsXG4gICAgXCJlc2xpbnQtcGx1Z2luLXByZXR0aWVyXCI6IFwiXjUuMS4zXCIsXG4gICAgXCJlc2xpbnQtcGx1Z2luLXZ1ZVwiOiBcIl45LjI3LjBcIixcbiAgICBcImZhc3QtZ2xvYlwiOiBcIl4zLjMuMlwiLFxuICAgIFwiaHVza3lcIjogXCJeOS4wLjExXCIsXG4gICAgXCJsaW50LXN0YWdlZFwiOiBcIl4xNS4yLjdcIixcbiAgICBcInBvc3Rjc3NcIjogXCJeOC40LjM5XCIsXG4gICAgXCJwb3N0Y3NzLWh0bWxcIjogXCJeMS43LjBcIixcbiAgICBcInBvc3Rjc3Mtc2Nzc1wiOiBcIl40LjAuOVwiLFxuICAgIFwicHJldHRpZXJcIjogXCJeMy4zLjJcIixcbiAgICBcInNhc3NcIjogXCJeMS43Ny42XCIsXG4gICAgXCJzdHlsZWxpbnRcIjogXCJeMTYuNi4xXCIsXG4gICAgXCJzdHlsZWxpbnQtY29uZmlnLWh0bWxcIjogXCJeMS4xLjBcIixcbiAgICBcInN0eWxlbGludC1jb25maWctcmVjZXNzLW9yZGVyXCI6IFwiXjQuNi4wXCIsXG4gICAgXCJzdHlsZWxpbnQtY29uZmlnLXJlY29tbWVuZGVkLXNjc3NcIjogXCJeMTQuMC4wXCIsXG4gICAgXCJzdHlsZWxpbnQtY29uZmlnLXJlY29tbWVuZGVkLXZ1ZVwiOiBcIl4xLjUuMFwiLFxuICAgIFwic3R5bGVsaW50LWNvbmZpZy1zdGFuZGFyZFwiOiBcIl4zNi4wLjFcIixcbiAgICBcInRlcnNlclwiOiBcIl41LjMxLjFcIixcbiAgICBcInR5cGVzY3JpcHRcIjogXCJeNS41LjNcIixcbiAgICBcInVub2Nzc1wiOiBcIl4wLjU4LjlcIixcbiAgICBcInVucGx1Z2luLWF1dG8taW1wb3J0XCI6IFwiXjAuMTcuNlwiLFxuICAgIFwidW5wbHVnaW4taWNvbnNcIjogXCJeMC4xOC41XCIsXG4gICAgXCJ1bnBsdWdpbi12dWUtY29tcG9uZW50c1wiOiBcIl4wLjI2LjBcIixcbiAgICBcInZpdGVcIjogXCJeNS4zLjNcIixcbiAgICBcInZpdGUtcGx1Z2luLW1vY2stZGV2LXNlcnZlclwiOiBcIl4xLjUuMVwiLFxuICAgIFwidml0ZS1wbHVnaW4tbW9uYWNvLWVkaXRvclwiOiBcIl4xLjEuMFwiLFxuICAgIFwidml0ZS1wbHVnaW4tc3ZnLWljb25zXCI6IFwiXjIuMC4xXCIsXG4gICAgXCJ2aXRlLXBsdWdpbi12dWUtZGV2dG9vbHNcIjogXCJeNy4zLjVcIixcbiAgICBcInZ1ZS10c2NcIjogXCJeMi4wLjI2XCJcbiAgfSxcbiAgXCJyZXBvc2l0b3J5XCI6IFwiaHR0cHM6Ly9naXRlZS5jb20veW91bGFpb3JnL3Z1ZTMtZWxlbWVudC1hZG1pbi5naXRcIixcbiAgXCJhdXRob3JcIjogXCJcdTY3MDlcdTY3NjVcdTVGMDBcdTZFOTBcdTdFQzRcdTdFQzdcIixcbiAgXCJsaWNlbnNlXCI6IFwiTUlUXCIsXG4gIFwiZW5naW5lc1wiOiB7XG4gICAgXCJub2RlXCI6IFwiPj0xOC4wLjBcIlxuICB9LFxuICBcInBhY2thZ2VNYW5hZ2VyXCI6IFwicG5wbUA5LjEuMytzaGE1MTIuN2MyZWEwODllMWE2YWYzMDY0MDljNGZjOGM0ZjA4OTdiZGFjMzJiNzcyMDE2MTk2YzQ2OWQ5NDI4ZjFmZTJkNWEyMWRhZjhhZDY1MTI3NjI2NTRhYzY0NWI1ZDkxMzZiYjIxMGVjOWEwMGFmYThkYmM0Njc3ODQzYmEzNjJlY2RcIixcbiAgXCJfX25wbWluc3RhbGxfZG9uZVwiOiBmYWxzZVxufVxuIl0sCiAgIm1hcHBpbmdzIjogIjtBQUFvUSxPQUFPLFNBQVM7QUFDcFIsT0FBTyxZQUFZO0FBQ25CLFNBQWdDLFNBQVMsb0JBQW9CO0FBQzdELFNBQVMsd0JBQXdCO0FBRWpDLE9BQU8sZ0JBQWdCO0FBQ3ZCLE9BQU8sZ0JBQWdCO0FBQ3ZCLFNBQVMsMkJBQTJCO0FBQ3BDLE9BQU8sV0FBVztBQUNsQixPQUFPLG1CQUFtQjtBQUUxQixT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Q0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLE1BQ0Y7QUFBQSxJQUNGO0FBQUE7QUFBQSxJQUVBLE9BQU87QUFBQSxNQUNMLFFBQU87QUFBQSxNQUNQLFdBQVc7QUFBQSxNQUNYLHVCQUF1QjtBQUFBO0FBQUEsTUFDdkIsUUFBUTtBQUFBO0FBQUEsTUFDUixlQUFlO0FBQUEsUUFDYixVQUFVO0FBQUEsVUFDUixlQUFlO0FBQUE7QUFBQTtBQUFBO0FBQUEsUUFHakI7QUFBQSxRQUNBLFFBQVE7QUFBQSxVQUNOLFVBQVU7QUFBQTtBQUFBLFFBQ1o7QUFBQSxNQUNGO0FBQUEsTUFDQSxlQUFlO0FBQUEsUUFDYixRQUFRO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxVQUtOLGdCQUFnQjtBQUFBO0FBQUEsVUFFaEIsZ0JBQWdCO0FBQUE7QUFBQSxVQUVoQixnQkFBZ0IsQ0FBQyxjQUFtQjtBQUNsQyxrQkFBTSxPQUFPLFVBQVUsS0FBSyxNQUFNLEdBQUc7QUFDckMsZ0JBQUksVUFBVSxLQUFLLEtBQUssU0FBUyxDQUFDO0FBRWxDLGdCQUNFLDZDQUE2QyxLQUFLLFVBQVUsSUFBSSxHQUNoRTtBQUNBLHdCQUFVO0FBQUEsWUFDWixXQUFXLGdDQUFnQyxLQUFLLFVBQVUsSUFBSSxHQUFHO0FBQy9ELHdCQUFVO0FBQUEsWUFDWixXQUFXLGtDQUFrQyxLQUFLLFVBQVUsSUFBSSxHQUFHO0FBQ2pFLHdCQUFVO0FBQUEsWUFDWjtBQUNBLG1CQUFPLEdBQUcsT0FBTztBQUFBLFVBQ25CO0FBQUEsUUFDRjtBQUFBLE1BQ0Y7QUFBQSxJQUNGO0FBQUEsSUFDQSxRQUFRO0FBQUEsTUFDTixjQUFjLEtBQUssVUFBVSxZQUFZO0FBQUEsSUFDM0M7QUFBQSxFQUNGO0FBQ0YsQ0FBQzsiLAogICJuYW1lcyI6IFtdCn0K
