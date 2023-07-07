<template>
  <div>
    <n-card :bordered="false" style="height: 100%;">
      <n-tabs type="line" animated default-value="downSetting">
        <n-tab-pane name="downSetting" tab="下载设置">
          <div style="width: 90%">
            <div style="display: flex;">
              <div style="margin-top: 6px; margin-right: 10px;">默认下载路径为</div>
              <el-input v-model="downPath" :disabled="true" style="width: 80%">
                <template #suffix>
                  <el-icon>
                    <component :is="Edit" @click="changeDownPath"></component>
                  </el-icon>
                </template>
              </el-input>
            </div>
            <div style="margin-top: 10px;display: flex">
              <div style="margin-top: 6px; margin-right: 10px;">下载并行任务数</div>
              <el-select v-model="downProcess" placeholder="Select" size="default" style="margin-left: 2px; width: 60px" @change="changeDownProcess">
                <el-option
                    v-for="item in options" :value="item"/>
              </el-select>
            </div>
            <div style="margin-top: 10px;display: flex">
              <div style="margin-top: 6px; margin-right: 10px;">可用下载次数 {{ count }}</div>
            </div>
          </div>
        </n-tab-pane>
        <n-tab-pane name="instructions" tab="使用说明">
          <div style="margin-left: 10px;text-align: left">
            在主页输入框中附上自己想要下载内容的分享链接,提交待后台转存后, 可在文件列表中查看<br>
            每个设备默认只提供两次提交次数<br>
          </div>
        </n-tab-pane>
        <n-tab-pane name="feedback" tab="问题反馈">
          <img src="../assets/images/404.gif" alt="" style="margin-top: 10px"><br>
          如有故障联系 QQ 2025907338<br>
        </n-tab-pane>
      </n-tabs>
    </n-card>
  </div>
</template>

<script>

import {ChangeDownPath, GetAuthorityCount, GetConfig, UpdateConfigItem} from "../../wailsjs/go/controller/App.js";
import {Edit} from "@element-plus/icons-vue";
import {ElMessage} from "element-plus";

export default {
  computed: {
    Edit() {
      return Edit
    },
  },
  data(){
    return {
      downPath: null,
      downProcess: null,
      options:['1','2','3','4'],
      count: null
    }
  },
  mounted() {
    this.getDownPath()
    this.getCount()
  },
  methods: {
    getCount(){
      GetAuthorityCount().then(res =>{
        if (res < 0){
          this.count = 0
          return
        }
        this.count = res
      })
    },
    getDownPath(){
      GetConfig().then(res =>{
        this.downPath = res.DownPath
        this.downProcess = res.MaxDownProcess
      })
    },
    changeDownPath() {
      ChangeDownPath().then(res =>{
        this.downPath = res
      })
    },
    changeDownProcess(){
      console.log(this.downProcess)
      UpdateConfigItem({ name: "MaxDownProcess", value: this.downProcess}).then(res =>{
        ElMessage({
          message: '更改成功!',
          type: 'success',
        })
      })
    }
  },
}
</script>

<style scoped>
.card-radius-10 {
  border-radius: 10px;
}
</style>