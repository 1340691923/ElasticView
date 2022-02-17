<template>
    <div>
      <div class="app-container" style="background: white;height: 95%;">
        <div class="filter-container">
          <el-tabs v-model="crudTab" >

            <el-tab-pane label="筛选" name="filter">

              <filter-where
                v-model="form.whereFilter"
                table-typ="2"
                :data-type-map="attrMap"
                :options="eventAttrOptions"
              />
            </el-tab-pane>
            <el-tab-pane label="排序" name="sort">

            </el-tab-pane>
          </el-tabs>
        </div>
        <el-card class="box-card">
          <div
            style="background: white;"
          >
            123
          </div>
        </el-card>
      </div>
    </div>
</template>

<script>
    export default {
        name: "crud",
      components: {
        'FilterWhere': () => import('@/components/AnalyseTools/FilterWhere/index'),
      },
        props:{
          indexName:{
            type:String,
            default:""
          }
        },
        data(){
          return {
            crudTab:'filter',
            eventAttrOptions: [ { "label": "事件", "options": [ { "value": "xwl_server_time", "label": "服务端入库时间" }, { "value": "xwl_city", "label": "用户所在城市" }, { "value": "xwl_distinct_id", "label": "访客ID" }, { "value": "xwl_ip", "label": "事件生成时用户所在IP" }, { "value": "xwl_lib", "label": "SDK语言" }, { "value": "xwl_lib_version", "label": "SDK版本" }, { "value": "xwl_network_type", "label": "用户网络类型" }, { "value": "xwl_os", "label": "用户操作系统" }, { "value": "xwl_os_version", "label": "用户操作系统版本" }, { "value": "xwl_province", "label": "用户所在省份" }, { "value": "xwl_reg_time", "label": "用户注册时间" }, { "value": "xwl_screen_height", "label": "用户屏幕高度" }, { "value": "xwl_screen_width", "label": "用户屏幕宽度" }, { "value": "xwl_update_time", "label": "用户信息最后修改时间" }, { "value": "xwl_device_model", "label": "用户机型" }, { "value": "xwl_manufacturer", "label": "设备提供商" }, { "value": "url_path", "label": "url_path" }, { "value": "xwl_kafka_partition", "label": "kafka分区" } ] } ],
            attrMap:{ "1": [ { "attribute_name": "xwl_distinct_id", "show_name": "访客ID", "data_type": 3, "attribute_type": 1, "data_type_format": "字符串类型", "attribute_source": 1, "status": 0 }, { "attribute_name": "xwl_reg_time", "show_name": "用户注册时间", "data_type": 4, "attribute_type": 1, "data_type_format": "时间类型", "attribute_source": 1, "status": 0 }, { "attribute_name": "xwl_update_time", "show_name": "用户信息最后修改时间", "data_type": 4, "attribute_type": 1, "data_type_format": "时间类型", "attribute_source": 1, "status": 0 }, { "attribute_name": "xwl_server_time", "show_name": "服务端入库时间", "data_type": 4, "attribute_type": 1, "data_type_format": "时间类型", "attribute_source": 1, "status": 0 }, { "attribute_name": "xwl_province", "show_name": "用户所在省份", "data_type": 3, "attribute_type": 1, "data_type_format": "字符串类型", "attribute_source": 1, "status": 0 }, { "attribute_name": "xwl_city", "show_name": "用户所在城市", "data_type": 3, "attribute_type": 1, "data_type_format": "字符串类型", "attribute_source": 1, "status": 0 }, { "attribute_name": "xwl_ip", "show_name": "事件生成时用户所在IP", "data_type": 3, "attribute_type": 1, "data_type_format": "字符串类型", "attribute_source": 1, "status": 0 } ], "2": [ { "attribute_name": "xwl_server_time", "show_name": "服务端入库时间", "data_type": 4, "attribute_type": 1, "data_type_format": "时间类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "xwl_city", "show_name": "用户所在城市", "data_type": 3, "attribute_type": 1, "data_type_format": "字符串类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "xwl_distinct_id", "show_name": "访客ID", "data_type": 3, "attribute_type": 1, "data_type_format": "字符串类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "xwl_ip", "show_name": "事件生成时用户所在IP", "data_type": 3, "attribute_type": 1, "data_type_format": "字符串类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "xwl_lib", "show_name": "SDK语言", "data_type": 3, "attribute_type": 1, "data_type_format": "字符串类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "xwl_lib_version", "show_name": "SDK版本", "data_type": 3, "attribute_type": 1, "data_type_format": "字符串类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "xwl_network_type", "show_name": "用户网络类型", "data_type": 3, "attribute_type": 1, "data_type_format": "字符串类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "xwl_os", "show_name": "用户操作系统", "data_type": 3, "attribute_type": 1, "data_type_format": "字符串类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "xwl_os_version", "show_name": "用户操作系统版本", "data_type": 3, "attribute_type": 1, "data_type_format": "字符串类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "xwl_province", "show_name": "用户所在省份", "data_type": 3, "attribute_type": 1, "data_type_format": "字符串类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "xwl_reg_time", "show_name": "用户注册时间", "data_type": 4, "attribute_type": 1, "data_type_format": "时间类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "xwl_screen_height", "show_name": "用户屏幕高度", "data_type": 1, "attribute_type": 1, "data_type_format": "数字类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "xwl_screen_width", "show_name": "用户屏幕宽度", "data_type": 1, "attribute_type": 1, "data_type_format": "数字类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "xwl_update_time", "show_name": "用户信息最后修改时间", "data_type": 4, "attribute_type": 1, "data_type_format": "时间类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "xwl_device_model", "show_name": "用户机型", "data_type": 3, "attribute_type": 1, "data_type_format": "字符串类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "xwl_manufacturer", "show_name": "设备提供商", "data_type": 3, "attribute_type": 1, "data_type_format": "字符串类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "url_path", "show_name": "", "data_type": 3, "attribute_type": 2, "data_type_format": "字符串类型", "attribute_source": 2, "status": 0 }, { "attribute_name": "xwl_kafka_partition", "show_name": "kafka分区", "data_type": 1, "attribute_type": 2, "data_type_format": "数字类型", "attribute_source": 2, "status": 0 } ] },
            form: {
              whereFilter: {
                filterType: 'COMPOUND',
                filts: [],
                relation: '且'
              },
            },

          }
        }
    }
</script>

<style scoped>

</style>
