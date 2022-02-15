export function filterData(data, searchContent) {
  // var input = this.searchContent && this.searchContent.toLowerCase();
  var input = searchContent.toLowerCase()
  var items = data
  var items1
  if (input) {
    items1 = items.filter(function(item) {
      return Object.keys(item).some(function(key1) {
        return String(item[key1])
          .toLowerCase()
          .match(input)
      })
    })
  } else {
    items1 = items
  }
  return items1
}
