import { useClipboard,usePermission } from "@vueuse/core";
import {ElMessage} from "element-plus";

const clipboardSuccess = (text: any) => {

  ElMessage({
    type: 'success',
    message: `拷贝成功`
  })
}

const clipboardError = (text: any) => {

  ElMessage({
    type: 'error',
    message: `拷贝${text}失败`
  })
}

export default function handleClipboard(text: string) {
  const { isSupported, copy } = useClipboard()
  if (!isSupported) usePermission('clipboard-write')

  copy(text)
    .then(() => {
      clipboardSuccess(text)
    })
    .catch(() => {
      clipboardError(text)
    })
}
