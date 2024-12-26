<template>
    <div
      ref="floatDrag"
      id="floatDrag"
      class="float-position"
      :style="{ top: data.top + 'px', right: data.right + 'px !important' }"
      @touchmove.prevent
      @mousemove.prevent
      @mousedown="mouseDown"
      @mouseup="mouseUp"
    >
      <div style="color: white" class="content">
        AI
      </div>
    </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, reactive, nextTick, ref, onBeforeMount } from 'vue';
interface dataOptionItem {
  clientWidth: null | number;
  clientHeight: null | number;
  left: number;
  top: number;
  right: number;
  timer: null | number;
  currentTop: 0 | number;
  mousedownX: 0 | number;
  mousedownY: 0 | number;
}
export default defineComponent({
  name: 'DragBall',
  components: {

  },
  props: {
    distanceRight: {
      type: Number,
      default: 36,
    },
    distanceBottom: {
      type: Number,
      default: 100,
    },
    isScrollHidden: {
      type: Boolean,
      default: false,
    },
    isCanDraggable: {
      type: Boolean,
      default: true,
    },
  },
  emits: ['handlepaly'],
  setup(props, { emit }) {
    const data = reactive<dataOptionItem>({
      clientWidth: null,
      clientHeight: null,
      left: 0,
      top: 0,
      right: 0,
      timer: null,
      currentTop: 0,
      mousedownX: 0,
      mousedownY: 0,
    });
    // const floatDrag = ref();
    const floatDragDom = ref();
    data.clientWidth = document.documentElement.clientWidth;
    data.clientHeight = document.documentElement.clientHeight;
    onMounted(() => {
      props.isCanDraggable &&
      nextTick(() => {
        // 获取元素位置属性
        floatDragDom.value = document.getElementById('floatDrag');
        // 设置初始位置
        data.right = props.distanceRight;
        data.top =
          Number(data.clientHeight) - floatDragDom.value?.offsetHeight - props.distanceBottom;
        initDraggable();
      });
      // this.isScrollHidden && window.addEventListener('scroll', this.handleScroll);
      window.addEventListener('resize', handleResize);
    });
    onBeforeMount(() => {
      // window.removeEventListener('scroll', this.handleScroll);
      window.removeEventListener('resize', handleResize);
    });
    /**
            * 初始化draggable
            */
    const initDraggable = () => {
      floatDragDom.value.addEventListener('touchstart', toucheStart);
      floatDragDom.value.addEventListener('touchmove', (e) => touchMove(e));
      floatDragDom.value.addEventListener('touchend', touchEnd);
    };
    const handleResize = () => {
      checkDraggablePosition();
    };
    /**
            * 判断元素显示位置
            * 在窗口改变和move end时调用
            */
    const checkDraggablePosition = () => {
      data.clientWidth = document.documentElement.clientWidth;
      data.clientHeight = document.documentElement.clientHeight;
      if (data.right + floatDragDom.value.offsetWidth / 2 >= data.clientWidth / 2) {
        // 判断位置是往左往右滑动
        data.right = data.clientWidth - floatDragDom.value.offsetWidth;
      } else {
        data.right = 0;
      }
      if (data.top < 0) {
        // 判断是否超出屏幕上沿
        data.top = 0;
      }
      if (data.top + floatDragDom.value.offsetHeight >= data.clientHeight) {
        // 判断是否超出屏幕下沿
        data.top = data.clientHeight - floatDragDom.value.offsetHeight;
      }
    };
    const canClick = ref(false);
    const mouseDown = (e) => {
      const event = e || window.event;
      data.mousedownX = event.screenX;
      data.mousedownY = event.screenY;
      let floatDragWidth = floatDragDom.value.offsetWidth / 2;
      let floatDragHeight = floatDragDom.value.offsetHeight / 2;
      if (event.preventDefault) {
        event.preventDefault();
      }
      canClick.value = false;
      floatDragDom.value.style.transition = 'none';
      document.onmousemove = function (e) {
        var event = e || window.event;

        // 这里的判断是为了保证按钮只能在浏览器内拖动，不会超出
        data.right = Number(data.clientWidth) - event.clientX - floatDragWidth;
        data.top = event.clientY - floatDragHeight;
        if (data.right < 0) data.right = 0;
        if (data.top < 0) data.top = 0;
        // 鼠标移出可视区域后给按钮还原
        if (
          event.clientY < 0 ||
          event.clientY > Number(data.clientHeight) ||
          event.clientX > Number(data.clientWidth) ||
          event.clientX < 0
        ) {
          data.right = 0;
          data.top =
            Number(data.clientHeight) - floatDragDom.value?.offsetHeight - props.distanceBottom;
          document.onmousemove = null;
          floatDragDom.value.style.transition = 'all 0.3s';
          return;
        }
        if (data.right >= document.documentElement.clientWidth - floatDragWidth * 2) {
          data.right = document.documentElement.clientWidth - floatDragWidth * 2;
        }
        if (data.top >= Number(data.clientHeight) - floatDragHeight * 2) {
          data.top = Number(data.clientHeight) - floatDragHeight * 2;
        }
      };
    };
    const mouseUp = (e) => {
      const event = e || window.event;
      //判断只是单纯的点击，没有拖拽
      if (data.mousedownY == event.screenY && data.mousedownX == event.screenX) {
        emit('handlepaly');
      }
      document.onmousemove = null;
      checkDraggablePosition();
      floatDragDom.value.style.transition = 'all 0.3s';
    };
    const toucheStart = () => {
      canClick.value = false;
      floatDragDom.value.style.transition = 'none';
    };
    const touchMove = (e) => {
      canClick.value = true;
      if (e.targetTouches.length === 1) {
        // 单指拖动
        let touch = e.targetTouches[0];
        data.right =
          Number(data.clientWidth) - touch.clientX - floatDragDom.value.offsetWidth / 2;
        data.top = touch.clientY - floatDragDom.value.offsetHeight / 2;
      }
    };
    const touchEnd = () => {
      if (!canClick.value) return; // 解决点击事件和touch事件冲突的问题
      floatDragDom.value.style.transition = 'all 0.3s';
      checkDraggablePosition();
    };
    return {
      data,
      handleResize,
      checkDraggablePosition,
      mouseUp,
      mouseDown,
    };
  },
});
</script>



<style>
  html,
  body {
      overflow: hidden;
    }
</style>
<style scoped lang="scss">

.float-position:hover{
  cursor: pointer;
}

.float-position {
  position: fixed;
  z-index: 99999999 !important;
  right: 0;
  top: 70%;
  width: 70px;
  height: 70px;
  display: flex;
  align-items: center;
  justify-content: center;
  user-select: none;

.content {
  width: 50px;
  height: 50px;
  background:rgb(64 201 198);
  border-radius: 50%;
  position: relative;
  padding: 0.8em;
  display: flex;
  align-content: center;
  justify-content: center;
}

.close {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  background: white;
  position: absolute;
  right: -10px;
  top: -12px;
  cursor: pointer;
}
}

.cart {
  border-radius: 50%;
  width: 5em;
  height: 5em;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-notice {
  display: inline-block;
  transition: all 0.3s;

span {
  vertical-align: initial;
}

.notice-badge {
  color: inherit;

.header-notice-icon {
  font-size: 16px;
  padding: 4px;
}
}
}

.drag-ball .drag-content {
  overflow-wrap: break-word;
  font-size: 14px;
  color: #fff;
  letter-spacing: 2px;
}
</style>
