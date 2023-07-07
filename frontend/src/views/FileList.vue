<template>
  <div class="container">
    <div class="top">
      <div style="float: left; flex: 1;">
        <div style="position: relative;margin-top: 15px; float: left; margin-left: 15px; margin-bottom: 3px">
         <span v-for="(item, index) in topfileList" :key="index" @click="handleClick(item)">
           <span :class="{ 'other-letter': index !== topfileList.length - 1 }">{{ item }}</span>
           <span v-if="index !== topfileList.length - 1" style="margin-right: 3px">></span>
        </span>
        </div>
      </div>
      <div style="width: 100px;flex:1; margin-top: 10px;">
        <el-button @click="downFile">下载内容</el-button>
      </div>
    </div>
    <div class="bottom">
      <el-table
          v-loading="loading"
          ref="multipleTableRef"
          :data="tableData"
          style="width: 100%"
          @row-click="handleRowClick"
          @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="30px" />
        <el-table-column label="文件名">
          <template #default="scope">
            <div v-if="scope.row.FileType === '文件夹'">
              <el-icon><Folder /></el-icon>
              <span style="margin-left: 10px">{{ scope.row.FileName }}</span>
            </div>
            <div v-else>
              <el-icon><Document /></el-icon>
              <span style="margin-left: 10px">{{ scope.row.FileName }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column property="FileType" label="文件类型" width="100px"/>
        <el-table-column property="FileSize" label="文件大小" width="100px"/>
        <el-table-column property="FileTime" label="修改时间" width="180px" show-overflow-tooltip />
      </el-table>
    </div>
  </div>
</template>

<script>

import {
  DownFilePre,
  GetConfig,
  GetFileList,
  GetFileOtherList,
  GetTopList
} from "../../wailsjs/go/controller/App.js";
import {Document, Folder} from "@element-plus/icons-vue";
import {ElMessage} from "element-plus";

export default {
  components: {Document, Folder},
  data(){
    return {
      topfileList: [], // 文件夹列表,调用方法获得
      currentPath: null, // 当前文件夹位置
      tableData: [],  // 表格展示的数据
      selectedRows: [], // 存储选中的行
      downData: [], // 定义需要下载的数据
      loading:true,
    }
  },
  mounted() {
    this.getConfigMap()
  },
  methods: {
    getConfigMap(){
      GetConfig().then((res) => {
        this.currentPath = res.PanID
        this.getTopList()
        this.getFileList()
      })
    },
    getTopList(){
      GetTopList(this.currentPath).then(res =>{
        this.topfileList = res
      })
    },
    getFileList(){
      this.loading = true
      GetFileList(this.currentPath).then(res =>{
        this.tableData = res
        this.loading = false
      })
    },
    handleClick(item) {
      GetFileOtherList(this.currentPath,item).then(res =>{
        this.currentPath = res
        this.getTopList()
        this.getFileList()
      })
    },
    handleSelectionChange(data){
      this.downData = data
    },
    handleRowClick(row){
      if (row.FileType === "文件夹"){
        // const loading = ElLoading.service({
        //   lock: true,
        //   text: 'Loading',
        //   background: 'rgba(0, 0, 0, 0.7)',
        // })
        this.currentPath = row.FilePath
        this.getTopList()
        this.getFileList()

      }
    },
    downFile(){
      this.loading = true
      DownFilePre(this.downData).then(res =>{
        this.loading = false
        if (res === -1){
          ElMessage.error('单次提交下载文件个数最大为500,您已超限')
        }
        if (res === -2){
          ElMessage.error('每分钟单次请求超限,请一分钟后重试')
        }
      })
    },
  }
}
</script>



<style>
.container {
  display: flex;
  flex-direction: column;
  height: 100%; /* 设置容器高度填满整个视口 */

}
.other-letter{
  color: #06a7FF;
  margin-right: 3px;
}

.top {
  height: 100px;
  position: fixed;
  top: 0;
  width: 100%;
  display: flex;
}

.bottom {
  margin-top: 45px;
  width: 100%
}
</style>