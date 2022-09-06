// guide.js
const steps = [
  {
    element: '#index-health-status',
    popover: {
      title: '选择索引健康状态',
      description: '提供有 green，yellow，red 三种状态筛选'
    }
  },
  {
    element: '#index-keyword',
    popover: {
      title: '在这里可进行关键词筛选',
      description: '可输入索引名，索引id等关键词 进行筛选'
    }
  }, {
    element: '#index-search',
    popover: {
      title: '刷新索引列表',
      description: '点击这个按钮用来刷新索引列表'
    }
  }, {
    element: '#new-index',
    popover: {
      title: '新建索引',
      description: '用户可点击此来打开新建索引界面'

    }
  }, {
    element: '#readOnlyAllowDelete',
    popover: {
      title: '将所有索引状态变为可写状态',
      description: '由于一些不可抗的因素，例如磁盘空间满了，ES会采取保护机制，将所有索引变为只读状态，需要用户手动将所有索引状态变为可写状态'
    }
  }, {
    element: '#flushIndex',
    popover: {
      title: 'flush刷新操作',
      description: 'flush刷新操作：确保原来只保留在事务日志（transaction log）中的数据，得以真正的保存到Lucene索引中'
    }
  }, {
    element: '#patch-operate',
    popover: {
      title: '批量操作',
      description: '这个区域的按钮需要配合表数据前的多选框进行批量操作（例如批量关闭索引）'
    }
  }, {
    element: '#patchCloseIndex',
    popover: {
      title: '批量关闭索引',
      description: '索引处于open状态，就会占用内存+磁盘，如果将索引close，就只会占用磁盘'
    }
  }, {
    element: '#patchOpenIndex',
    popover: {
      title: '批量开启索引',
      description: '索引关闭后， 对集群的相关开销基本降低为 0，但是无法被读取和搜索，当需要的时候， 可以重新打开，索引恢复正常'
    }
  }, {
    element: '#patchForcemergeIndex',
    popover: {
      title: '批量强制合并索引',
      description: 'forcemerge操作,手动释放磁盘空间，用于删除文档后，进行合并索引操作，清理磁盘'
    }
  }, {
    element: '#patchRefreshIndex',
    popover: {
      title: '批量刷新索引',
      description: '提高ES的实时性，使添加文档尽可能快的被搜索到，同时又避免频繁fsync带来性能开销'
    }
  }, {
    element: '#patchFlushIndex',
    popover: {
      title: '批量flush索引',
      description: '将os cache的索引文件(segment file)持久化到磁盘'
    }
  }, {
    element: '#patchCacheClear',
    popover: {
      title: '批量清除缓存',
      description: '批量清除缓存'
    }
  }, {
    element: '#patchDeleteIndex',
    popover: {
      title: '批量删除索引',
      description: '不同于删除文档操作，删除索引后磁盘空间将会立即释放'
    }
  }
]

export default steps
