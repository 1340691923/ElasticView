<template>
  <div  v-loading="connectLoading">
    <el-row  :gutter="40" >
      <el-col v-for="(v,k,index) in list" :key="index" :xs="18" :sm="18" :lg="12" class="card-panel-col">
        <el-card class="box-card" style="">
          <div slot="header" style="display: flex; align-items: center; justify-content: space-between;">
            <span><el-tag>节点名:{{v["name"]}}</el-tag></span>
            <div>

              <a-tooltip  v-if="v.master" placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>主节点</span>
                </template>
                <el-tag type="primary"><i class="el-icon-star-on"></i>主节点</el-tag>
              </a-tooltip>
              <a-tooltip  v-else-if="v['node.role'].includes('m')" placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>主节点候选</span>
                </template>
                <el-tag type="primary"><i class="el-icon-star-off"></i>主节点候选</el-tag>
              </a-tooltip>
              <a-tooltip v-if="v['node.role'].includes('d')"  placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>数据节点</span>
                </template>
                <el-tag type="success">
                  <i class="el-icon-bank-card"></i>数据节点
                </el-tag>
              </a-tooltip>
              <a-tooltip  v-if="v['node.role'].includes('i')" placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>预处理节点</span>
                </template>
                <el-tag type="warning">
                  <i class="el-icon-postcard"></i>预处理节点
                </el-tag>
              </a-tooltip>
              <a-tooltip  v-if="v['node.role'] === '-'" placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>
                    仅协调节点
                  </span>
                </template>
                <el-tag type="warning">
                  仅协调节点
                </el-tag>
              </a-tooltip>
            </div>
          </div>
          <div style="display: flex; align-items: center; justify-content: space-between;">
            <div style="width: 48%;border-bottom: 1px solid white">
              <div style="display: flex; align-items: center; justify-content: space-between;min-height: 50px;">
                <div>
                  {{ $t('IP地址') }}:
                </div>
                <div>
                  <el-tag>{{v.ip}}</el-tag>
                </div>
              </div>
              <div style="display: flex; align-items: center; justify-content: space-between;min-height: 50px;">
                <div>
                  {{ $t('主节点') }}:
                </div>
                <div>
                  <template v-if="v.master"><el-tag>yes</el-tag></template>
                  <template v-else-if="v.masterEligible"><el-tag>eligible</el-tag></template>
                  <template v-else><el-tag>no</el-tag></template>
                </div>
              </div>
              <div style="display: flex; align-items: center; justify-content: space-between;min-height: 50px;">
                <div>
                  {{ $t('节点角色') }}:
                </div>
                <div>
                  <el-tag>{{v["node.role"]}}</el-tag>
                </div>
              </div>
              <div style="display: flex; align-items: center; justify-content: space-between;min-height: 50px;">
                <div>
                  {{ $t('负载') }}:
                </div>
                <div>
                 <el-tag>1m:{{ v.load_1m }} / 5m:{{ v.load_5m }} / 15m:{{ v.load_15m }}</el-tag>
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
                    {{ $t('内存') }}:&nbsp;&nbsp;<el-tag>{{ v["ram.current"] }}/{{ v["ram.max"] }}</el-tag>
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
                    {{ $t('堆内存') }}:&nbsp;&nbsp;<el-tag>{{ v["heap.current"] }}/ {{ v["heap.max"] }}</el-tag>
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
                    {{ $t('磁盘') }}:&nbsp;&nbsp;<el-tag>{{ v["disk.used"] }}/ {{ v["disk.total"] }}</el-tag>
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
      value: 90,
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
        opt.progress.color = "#4ad47f"
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
