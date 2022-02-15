// guide.js
const steps = [
  {
    element: '.select-method',
    popover: {
      title: '选择请求方式',
      description: '系统内置了 POST，GET，PUT，DELETE，HEAD 五种请求方式'
    }
  },
  {
    element: '.select-path',
    popover: {
      title: '请输入请求地址',
      description: '系统增加了部分请求地址的提示 例如：_cat/indices,_cat/tasks 等请求地址 用户如需加上get参数请输入 Path?key=val&key=val的格式'
    }
  }, {
    element: '.go',
    popover: {
      title: '发起请求',
      description: '用户填好请求方式后，即可点击 ->GO 发起请求'
    }
  }, {
    element: '.sql-format',
    popover: {
      title: 'SQL转换器',
      description: '用户可点击此 来开启SQL编辑器'

    }
  }, {
    element: '.search-history',
    popover: {
      title: '历史记录',
      description: '当前用户所用过的GET请求的历史记录，显示离当前时间最近的20条记录'
    }
  }, {
    element: '.req-body',
    popover: {
      title: 'JSON 请求体',
      description: '当请求 需要JSON 请求体时用到，系统内置了部分提示 例如 term,bool，match等'
    }
  }, {
    element: '.res-body',
    popover: {
      title: 'JSON 响应体',
      description: '当请求发送成功后返回的JSON返回体'
    }
  }
]

export default steps
