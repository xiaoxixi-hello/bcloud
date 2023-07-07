<template>
  <el-table
      ref="multipleTableRef"
      :data="DownDataList"
      style="width: 100%;margin-top: 20px">
    <el-table-column label="文件名" :show-overflow-tooltip="true">
      <template #default="scope">
        <span class="truncate">{{ scope.row.Name }}</span>
      </template>
    </el-table-column>
    <el-table-column property="Size" label="文件大小" width="90px"/>
    <el-table-column property="Status" label="下载状态" width="150px"/>
    <el-table-column property="CreatedAt" label="下载时间" width="120px"/>
    <el-table-column label="操作" width="140px">
      <template #default="scope">
        <el-button size="small" type="info" v-loading.fullscreen.lock="fullscreenLoading" @click="retryClick(scope.row.ID)">重试</el-button>
        <el-button size="small" type="danger" v-loading.fullscreen.lock="fullscreenLoading" @click="deleteClick(scope.row.ID)">删除</el-button>
      </template>
    </el-table-column>
    <el-progress :percentage="50" />
  </el-table>
</template>
<script>
import {DeleteRetry, DownRetry, GetDownListDetail} from "../../wailsjs/go/controller/App.js";
import {Delete} from "@element-plus/icons-vue";
import {ElMessage} from "element-plus";

export default {
  computed: {
    Delete() {
      return Delete
    }
  },
  data(){
    return {
      DownDataList: null,
      fullscreenLoading: null,
    }
  },
  mounted() {
    this.getDownListDetail();
    setInterval(this.getDownListDetail, 1500);
  },
  methods: {
    getDownListDetail(){
      GetDownListDetail().then(res =>{
        this.DownDataList = res
      })
    },
    retryClick(id){
      ElMessage({
        message: '重试任务提交成功!',
        type: 'success',
      })
      DownRetry(id)
      this.getDownListDetail()
      this.fullscreenLoading = true
      setTimeout(() => {
        this.fullscreenLoading = false
      }, 1500)
    },
    deleteClick(item){
      ElMessage({
        message: '删除任务提交成功!',
        type: 'warning',
      })
      DeleteRetry(item)
      this.getDownListDetail()
      this.fullscreenLoading = true
      setTimeout(() => {
        this.fullscreenLoading = false
      }, 2500)
    },
  }
}
</script>
<style>
</style>