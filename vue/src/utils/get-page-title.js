// 获取title
import defaultSettings from '@/settings'

const title = defaultSettings.title || 'ElasticView'

export default function getPageTitle(pageTitle) {
  if (pageTitle) {
    return `${pageTitle} - ${title}`
  }
  return `${title}`
}
