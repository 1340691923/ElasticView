// 获取title
import defaultSettings from '@/settings'
import {i18nText} from "@/utils/lang";

const title = defaultSettings.title || 'ElasticView'


export default function getPageTitle(pageTitle) {

  if (pageTitle) {
    return `${i18nText(pageTitle)} - ${title}`
  }
  return `${title}`
}
