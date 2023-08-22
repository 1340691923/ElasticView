<template>
  <div class="top-nav">
    <div class="log">EV</div>

    <el-menu
      id="menuUl"
      :key="updateTopMenuActive"
      class="el-menu-demo"
      :active-text-color="variables.menuActiveText"
      :default-active="activeMenu"
      mode="horizontal"
      @select="handleSelect"
    >
      <div v-for="item in permission_routes" :key="item.path" class="nav-item">
        <template v-if="!item.hidden">
          <app-link :to="resolvePath(item)">
            <el-menu-item

              :index="item.path"
            >{{ item.meta ? item.meta.title : item.children[0].meta.title }}</el-menu-item>
          </app-link>
        </template>
      </div>
    </el-menu>

    <!--        <div class="swiperMenu" id="swiperMenu" :style="{'padding-right':moveBtnWidth+'px'}" >
          <el-menu
            :key="updateTopMenuActive"
            class="el-menu-demo"
            id="menuUl"
            :active-text-color="variables.menuActiveText"
            :default-active="activeMenu"
            mode="horizontal"
            @select="handleSelect"
          >
            <div v-for="item in permission_routes" :key="item.path" class="nav-item">
              <template v-if="!item.hidden">
                <app-link  :to="resolvePath(item)">
                  <el-menu-item

                    :index="item.path"
                  >{{ item.meta ? item.meta.title : item.children[0].meta.title }}</el-menu-item>
                </app-link>
              </template>
            </div>
          </el-menu>
          <div class="moveBtn" id="moveBtn" v-show="!(leftNum==0&&rightNum==0)" >
            <div class="move" @click="toLeftMove(leftNum)">
              <span>{{leftNum}}</span>
            </div>
            <div class="lineBox">
              <div></div>
            </div>
            <div class="move" @click="toRightMove(leftNum)">
              <span>{{rightNum}}</span>
            </div>
          </div>
        </div>-->

    <!--
    <a-menu :selected-keys="[activeMenu]" mode="horizontal" style="background-color: #475285;"
      @select="onSelect">
      <template v-for="(menuItem, index) in routesArr">
        <template v-if="menuItem.children && !menuItem.hidden && menuItem.children.length == 1">
          <a-menu-item :key="resolvePath(menuItem.path, menuItem.children[0].path)">
            <app-link :to="resolvePath(menuItem.path, menuItem.children[0].path)">
              <span>{{ menuItem.children[0].meta ? menuItem.children[0].meta.title : '' }}</span>
            </app-link>
          </a-menu-item>
        </template>
        <template v-else-if="menuItem.children && menuItem.children.length > 1 && !menuItem.hidden">
          <a-sub-menu :key="menuItem.path">
            <template #title>
              <span class="menu_title">
                <span>{{ menuItem.meta ? menuItem.meta.title : '' }}</span>
              </span>
            </template>
            <a-menu-item v-for="(child, childIndex) in menuItem.children" :key="resolvePath(menuItem.path, child.path)">
              <app-link :to="resolvePath(menuItem.path, child.path)">
                <span><i :class="child.meta.icon" style="margin-right: 10px;" />{{
                  child.meta ? child.meta.title : ''
                }}</span>
              </app-link>
            </a-menu-item>
          </a-sub-menu>
        </template>
      </template>
    </a-menu>

    -->

    <div class="right-menu">
      <select-link style="margin-top:4px">
        <el-dropdown slot="avatar" style="margin-left:30px" class="avatar-container" trigger="click">
          <div class="avatar-wrapper">
            <img :src="logo" class="user-avatar">
            <i class="el-icon-caret-bottom" />
          </div>
          <el-dropdown-menu slot="dropdown" class="user-dropdown">
            <router-link to="/">
              <el-dropdown-item>{{ $t('首页') }}</el-dropdown-item>
            </router-link>
            <el-dropdown-item divided>
              <span style="display: block;" @click="logout">{{ $t('注销') }}</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </select-link>

    </div>
  </div>
</template>

<script>
import logo from '@/assets/index.ico'

import { mapGetters } from 'vuex'
import AppLink from './Sidebar/Link'
import SelectLink from '@/components/SelectLink'

import { constantRoutes } from '@/router'
import variables from '@/styles/variables.scss'
import { isExternal } from '@/utils/validate'

export default {
  name: 'Topbar',
  components: {
    AppLink,
    SelectLink
  },
  data() {
    return {
      logo: logo,
      updateTopMenuActive: 0,
      show: true,
      copy_Menus: [],
      resizeTick: false,
      firstMenuChange: false,
      templateMoreAlias: '_templateMoreMenu', // 临时更多目录的别名
      activeMenuIndex: sessionStorage.activeMenuIndex || '/home', // 默认选择首页
      allWidth: 0, // 菜单ul总宽度
      leftNum: 0, // 左边菜单隐藏个数
      rightNum: 0, // 右边菜单隐藏个数
      boxLength: 0, // 可视窗口宽度
      moveBtnWidth: 0, // 按钮盒子宽度
      firstMenuWidth: 200, // 第一个菜单的宽度（首页）
      eveyMenuWidth: 114 // 每一个菜单的宽度
    }
  },
  watch: {
    boxLength(val) {
      if (this.allWidth > val) {
        // 右边隐藏个数 = 总共菜单的个数 - 左边隐藏的个数 - 右边隐藏宽度/114
        this.rightNum = this.permission_routes.length - Math.ceil((this.boxLength - this.firstMenuWidth) / this.eveyMenuWidth) - this.leftNum
      } else {
        this.rightNum = 0
      }
    }
  },
  created() {
    /* this.initMenus()
    // 监听页面窗口变化
    window.addEventListener('resize',function(){
      this.boxLength = document.getElementById('swiperMenu').offsetWidth;
    })*/
  },
  computed: {
    ...mapGetters([
      'permission_routes',
      'sidebar'
    ]),
    routes() {
      return this.$store.state.permission_routes
    },
    activeMenu() {
      const route = this.$route
      const { meta, path } = route
      // if set path, the sidebar will highlight the path you set
      if (meta.activeMenu) {
        return meta.activeMenu
      }
      // 如果是首页，首页高亮
      if (path === '/dashboard') {
        return '/'
      }
      // 如果不是首页，高亮一级菜单
      const activeMenu = '/' + path.split('/')[1]
      return activeMenu
    },
    variables() {
      return variables
    },
    sidebar() {
      return this.$store.state.app.sidebar
    }
  },
  mounted() {
    this.initCurrentRoutes()
  },
  methods: {
    onSelect() {
      console.log(123)
    },
    async logout() {
      // this.$websocket.close();
      await this.$store.dispatch('user/logout')
      this.$router.push(`/login?redirect=${this.$route.fullPath}`)
    },
    // 通过当前路径找到二级菜单对应项，存到store，用来渲染左侧菜单
    initCurrentRoutes() {
      const { path } = this.$route
      let route = this.permission_routes.find(
        item => item.path === '/' + path.split('/')[1]
      )
      // 如果找不到这个路由，说明是首页
      if (!route) {
        route = this.permission_routes.find(item => item.path === '/')
      }
      this.$store.commit('permission/SET_CURRENT_ROUTES', route)
      this.setSidebarHide(route)
    },
    // 判断该路由是否只有一个子项或者没有子项，如果是，则在一级菜单添加跳转路由
    isOnlyOneChild(item) {
      if (item.children && item.children.length === 1) {
        return true
      }
      return false
    },
    resolvePath(item) {
      // 如果是个完成的url直接返回
      if (isExternal(item.path)) {
        return item.path
      }
      // 如果是首页，就返回重定向路由
      if (item.path === '/') {
        const path = item.redirect
        return path
      }

      // 如果有子项，默认跳转第一个子项路由
      let path = ''
      /**
       * item 路由子项
       * parent 路由父项
       */
      const getDefaultPath = (item, parent) => {
        // 如果path是个外部链接（不建议），直接返回链接，存在个问题：如果是外部链接点击跳转后当前页内容还是上一个路由内容
        if (isExternal(item.path)) {
          path = item.path
          return
        }
        // 第一次需要父项路由拼接，所以只是第一个传parent
        if (parent) {
          path += (parent.path + '/' + item.path)
        } else {
          path += ('/' + item.path)
        }
        // 如果还有子项，继续递归
        if (item.children) {
          getDefaultPath(item.children[0])
        }
      }

      if (item.children) {
        getDefaultPath(item.children[0], item)

        return path
      }

      return item.path
    },
    handleSelect(key, keyPath) {
      // 把选中路由的子路由保存store
      const route = this.permission_routes.find(item => item.path === key)
      this.$store.commit('permission/SET_CURRENT_ROUTES', route)
      this.setSidebarHide(route)
    },
    // 设置侧边栏的显示和隐藏
    setSidebarHide(route) {
      if (!route.children || route.children.length === 1) {
        this.$store.dispatch('app/toggleSideBarHide', true)
      } else {
        this.$store.dispatch('app/toggleSideBarHide', false)
      }
    },
    // 菜单右侧移动按钮
    toRightMove(num) {
      const menuscrollLeft = document.getElementById('menuUl').scrollLeft
      if (num > 0) {
        if (menuscrollLeft == 0) {
          // document.getElementById('menuUl').scrollLeft = this.firstMenuWidth
          document.getElementById('menuUl').scrollLeft = this.eveyMenuWidth
        } else {
          document.getElementById('menuUl').scrollLeft += this.eveyMenuWidth
        }
        this.rightNum--
        this.leftNum++
      }
    },
    // 菜单左侧移动按钮
    toLeftMove(num) {
      console.log(num)
      const menuscrollLeft = document.getElementById('menuUl').scrollLeft
      if (num > 0) {
        if (menuscrollLeft > this.firstMenuWidth) {
          document.getElementById('menuUl').scrollLeft -= this.eveyMenuWidth
        } else {
          document.getElementById('menuUl').scrollLeft = 0
        }
        this.rightNum++
        this.leftNum--
      }
    },

    initMenus(newVal, oldVal) {
      // ---------------------此处省略其他代码-----------------------------
      this.$nextTick(() => {
        // 按钮盒子宽度
        this.moveBtnWidth = document.getElementById('moveBtn').offsetWidth
        // 可视窗口宽度 = 大盒子宽度 - 按钮盒子宽度
        this.boxLength = document.getElementById('swiperMenu').offsetWidth - this.moveBtnWidth
        // 菜单总长度
        this.allWidth = (this.permission_routes.length - 1) * this.eveyMenuWidth + this.firstMenuWidth

        if (this.allWidth > this.boxLength) {
          this.rightNum = this.permission_routes.length - Math.ceil((this.boxLength - this.firstMenuWidth) / this.eveyMenuWidth)
        } else {
          this.rightNum = 0
        }
      })
    }
  }
}
</script>

<style lang="scss" >

.el-menu-demo {
  overflow: hidden;
  white-space: nowrap;
  scroll-behavior: smooth;
                              .el-submenu {
                                display: inline-block;
                                float: none;
                                width: 114px;
                              }
}

.swiperMenu {
  float: left;
                  position: relative;
                  overflow: hidden;
                  padding-right: 74px !important;
.moveBtn {
  position: absolute;
  right: 0;
  top: 0;
  height: 40px;
  z-index: 1;
.move {
  width: 36px;
  height: 100%;
  background: #304156;
  float: left;
  display: grid;
  grid-template-rows: 70% 30%;
  justify-items: center;
&:hover {
.moveImg {
  width: 15px;
  height: 24px;
  margin-top: 18px;
  opacity: 1;
  cursor: pointer;
}
.moveImgnone {
  opacity: 0.2;
  cursor: auto;
}
}
.moveImg {
  width: 15px;
  height: 24px;
  margin-top: 18px;
  opacity: 0.5;
}
.moveImgnone {
  opacity: 0.2;
}
span {
  display: inline-block;
  width: 14px;
  height: 14px;
  background: #304156;
  border-radius: 1px;
  font-size: 12px;
  font-family: Helvetica;
  color: #FFFFFF;
  text-align: center;
}
}
.lineBox {
  float: left;
  background: #304156;
  height: 40px;
div {
  background: #fff;
  opacity: 0.2;
  width: 1px;
  height: 48px;
  margin-top: 14px;
}
}
}
}
</style>
