
export function translateRouteTitle(title: any) {
  const { te,t } = useI18n()
  // 判断是否存在国际化配置，如果没有原生返回
  const hasKey = te(title);
  if (hasKey) {
    const translatedTitle = t(title);
    return translatedTitle;
  }
  return title;
}

