<template>
    <div>
        <el-form>
            <el-form-item label="cwd">
                <el-input v-model="sets.cwd" disabled></el-input>
            </el-form-item>
            <el-form-item label="confPath">
                <el-input v-model="sets.confPath" disabled></el-input>
            </el-form-item>
            <el-form-item label="logDir">
                <div style="display: flex;">
                    <el-input v-model="sets.LOG_DIR" style="width: 500px;margin-right: 20px;"></el-input>
                    <el-button type="primary" @click="setFolder('LOG_DIR')">导入</el-button>
                </div>
            </el-form-item>

            <el-form-item label="操作">
                    <el-button type="primary" @click="submit">提交</el-button>
                </el-form-item>
        </el-form>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import {GetAppSet, SetAppSet,GetFolderPath} from '../../../wailsjs/go/apis/App'
import { ElNotification } from 'element-plus';
let sets = ref<Record<string,string>>({})
onMounted(()=>{
    GetAppSet().then(res=>{
        sets.value = res
        console.log('res',res);
    })
})

function submit(){
    SetAppSet(sets.value).then(res=>{
        ElNotification.success({
            title: '成功',
            message: '设置成功',
            duration: 1000,
        })
    })
}

function setFolder(key:string){
    GetFolderPath().then(res=>{
        sets.value[key] = res
    })
}
</script>

<style scoped>

</style>