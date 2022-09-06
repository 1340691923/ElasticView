// http相关帮助函数
export function removeSpecialSymbols(data) {
  data = data.replace(/\%/g, '%25')
  data = data.replace(/\#/g, '%23')
  data = data.replace(/\&/g, '%26')
  data = data.replace(/\+/g, '%2B')
  data = data.replace(/\ /g, '%20')
  data = data.replace(/\//g, '%2F')
  data = data.replace(/\?/g, '%3F')
  console.log('httpdata', data)
  return data
}
