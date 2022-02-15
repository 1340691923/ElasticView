/** 重置message，防止重复点击重复弹出message弹框 */

let messageInstance = null
const mainMessage = function DoneMessage(options) {
  // 如果弹窗已存在先关闭
  if (messageInstance) {
    messageInstance.close()
  }
  messageInstance = ELEMENT.Message(options)
}
const arr = ['success', 'warning', 'info', 'error']
arr.forEach(function(type) {
  mainMessage[type] = function(options) {
    if (typeof options === 'string') {
      options = {
        message: options
      }
    }
    options.type = type
    return mainMessage(options)
  }
})
export const message = mainMessage
