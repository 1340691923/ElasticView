<template>
  <div  v-loading="connectLoading">
    <el-row  :gutter="40" >
      <el-col v-for="(v,k,index) in list" :key="index" :xs="18" :sm="18" :lg="12" class="card-panel-col">
        <el-card class="box-card" style="">
          <div slot="header" style="display: flex; align-items: center; justify-content: space-between;">
            <span>节点名:{{v["name"]}}</span>
            <div>

              <a-tooltip  v-if="v.master" placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>主节点</span>
                </template>
                <el-button type="warning" icon="el-icon-star-on" circle/>
              </a-tooltip>
              <a-tooltip  v-else-if="v['node.role'].includes('m')" placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>主节点候选</span>
                </template>
                <el-button type="warning" icon=" el-icon-star-off" circle/>
              </a-tooltip>
              <a-tooltip v-if="v['node.role'].includes('d')"  placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>数据节点</span>
                </template>
                <el-button icon="el-icon-bank-card" circle />
              </a-tooltip>
              <a-tooltip  v-if="v['node.role'].includes('i')" placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>预处理节点</span>
                </template>
                <el-button type="success" icon="el-icon-postcard" circle />
              </a-tooltip>
              <a-tooltip  v-if="v['node.role'] === '-'" placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>仅协调节点</span>
                </template>
                <el-button type="success" icon="el-icon-s-flag" circle />
              </a-tooltip>
            </div>
          </div>
          <div style="display: flex; align-items: center; justify-content: space-between;">
            <div style="width: 48%;border-bottom: 1px solid white">

              <div style="display: flex; align-items: center; justify-content: space-between;min-height: 50px;">
                <div>
                  {{ $t('IP地址') }}
                </div>
                <div>
                  {{v.ip}}
                </div>
              </div>
              <div style="display: flex; align-items: center; justify-content: space-between;min-height: 50px;">
                <div>
                  {{ $t('主节点') }}
                </div>
                <div>
                  <template v-if="v.master">yes</template>
                  <template v-else-if="v.masterEligible">eligible</template>
                  <template v-else>no</template>
                </div>
              </div>
              <div style="display: flex; align-items: center; justify-content: space-between;min-height: 50px;">
                <div>
                  {{ $t('节点角色') }}
                </div>
                <div>
                  {{v["node.role"]}}
                </div>
              </div>
              <div style="display: flex; align-items: center; justify-content: space-between;min-height: 50px;">
                <div>
                  {{ $t('负载') }}
                </div>
                <div>
                 1m:{{ v.load_1m }} / 5m:{{ v.load_5m }} / 15m:{{ v.load_15m }}
                </div>
              </div>
            </div>
            <div style="width: 48%;border-bottom: 1px solid white">

              <div
                style="min-height: 50px;display: flex;align-items: center;text-align: center;width: 100%">

                <div style="display: flex; align-items: center; justify-content: space-between;min-height: 50px;">
                  <div>
                    {{ $t('CPU') }}:
                  </div>
                  <div style="margin-left: 10px">
                    <progress-bar
                      style="margin-right: 1px;background: white"
                      :options="getOpt(v.cpu)"
                      :value="v.cpu"
                    />
                  </div>
                </div>

              </div>
              <div
                style="min-height: 50px;display: flex;align-items: center;text-align: center;width: 100%">

                <div style="display: flex; align-items: center; justify-content: space-between;min-height: 50px;">
                  <div>
                    {{ $t('内存') }}:{{ v["ram.current"] }}/ {{ v["ram.max"] }}
                  </div>
                  <div style="margin-left: 10px">
                    <progress-bar
                      style="margin-right: 1px;background: white"
                      :options="getOpt(v['ram.percent'])"
                      :value="v['ram.percent']"
                    />
                  </div>
                </div>

              </div>
              <div
                style="min-height: 50px;display: flex;align-items: center;text-align: center;width: 100%">

                <div style="display: flex; align-items: center; justify-content: space-between;min-height: 50px;">
                  <div>
                    {{ $t('堆内存') }}:{{ v["heap.current"] }}/ {{ v["heap.max"] }}
                  </div>
                  <div style="margin-left: 10px">
                    <progress-bar
                      style="margin-right: 1px;background: white"
                      :options="getOpt(v['heap.percent'])"
                      :value="v['heap.percent']"
                    />
                  </div>
                </div>

              </div>
              <div
                style="min-height: 50px;display: flex;align-items: center;text-align: center;width: 100%">

                <div style="display: flex; align-items: center; justify-content: space-between;min-height: 50px;">
                  <div>
                    {{ $t('磁盘') }}:{{ v["disk.used"] }}/ {{ v["disk.total"] }}
                  </div>
                  <div style="margin-left: 10px">
                    <progress-bar
                      style="background: white"
                      :options="getOpt(v['disk.used_percent'])"
                      :value="v['disk.used_percent']"
                    />

                  </div>
                </div>
              </div>
            </div>
          </div>

        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import {CatAction} from '@/api/es'

export default {
  name: "Node",
  data() {
    return {
      list: [],
      connectLoading: false,
      options: {
        text: {
          color: '#FFFFFF',
          shadowEnable: true,
          shadowColor: '#000000',
          fontSize: 12,
          fontFamily: 'Helvetica',
          dynamicPosition: false,
          hideText: false
        },
        progress: {
          color: '#ff9800',
          backgroundColor: '#614215'
        },
        layout: {
          height: 15,
          width: 140,
          verticalTextAlign: 61,
          horizontalTextAlign: 43,
          zeroOffset: 0,
          strokeWidth: 30,
          progressPadding: 0,
          type: 'line'
        }
      },
      value: 90
    }
  },
  beforeMount() {
    this.searchData()
  },
  methods: {
    getOpt(v){
      let numV = Number(v)
      const opt = JSON.parse(JSON.stringify(this.options))
      if(numV<60){
        opt.progress.color = "green"
        opt.progress.backgroundColor = "grey"
      }else if (numV>=60 && numV <=80){
        opt.progress.color = "#ff9800"
        opt.progress.backgroundColor = "grey"
      }else{
        opt.progress.color = "red"
        opt.progress.backgroundColor = "grey"
      }
      return opt
    },
    async searchData() {
      this.connectLoading = true
      const form = {
        cat: "Node",
        es_connect: this.$store.state.baseData.EsConnectID
      }
      const res = await CatAction(form)
      if (res.code != 0) {
        this.$message({
          type: 'error',
          message: res.msg
        })
        return
      }
      this.list = res.data

      this.connectLoading = false
    }
  }
}
</script>

<style scoped>
.card-panel-col{
  margin-bottom: 30px;
}
</style>
