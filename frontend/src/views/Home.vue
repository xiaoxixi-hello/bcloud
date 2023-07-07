<template>
  <div>
    <n-card :bordered="false" style="height: 100%;">
      <img src="../assets/images/marvel.gif" alt="" style="height: 300px">
      <br>
      <el-input v-model="shareLink" style="width: 90%" type="textarea" placeholder="分享链接"/>
      <br>
      <el-button type="primary" style="margin-top: 5px" :disabled="isButtonDisabled" @click="putShareLink">提交</el-button>
    </n-card>
  </div>
</template>

<script>
import {ElMessage} from "element-plus";
import {CreateShareLinkInfo} from "../../wailsjs/go/controller/App.js";
import {Drizzling} from "@element-plus/icons-vue";

export default {
  components: {Drizzling},
  data() {
    return {
      shareLink: null,
      isButtonDisabled: false,
    }
  },
  mounted() {},
  methods: {
    putShareLink(){
      if (this.shareLink === null){
        ElMessage({
          message: '当提交内容为空',
          type: 'warning',
        })
        return
      }
      this.isButtonDisabled = true; // 设置按钮为不可用状态
      setTimeout(() => {
        this.isButtonDisabled = false; // 10秒后恢复按钮可用状态
      }, 3000);

      this.getAuthorityCount()
    },
    getAuthorityCount(){
      CreateShareLinkInfo(this.shareLink).then(res =>{
        if (res === -1){
          ElMessage({
            message: '当前暂无提交次数',
            type: 'warning',
          })
        }
      })
    },
  }
}
</script>
