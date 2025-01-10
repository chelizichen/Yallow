<template>
    <div>
        <el-form style="margin-bottom: 20px;">
            <el-button type="primary" @click="loadTarRelease">导入服务配置</el-button>
            <el-button type="success" @click="loadConf">加载配置</el-button>
        </el-form>

        <el-table :data="confList" border>
            <el-table-column prop="filePath" label="文件地址"></el-table-column>
            <el-table-column prop="serverType" label="服务类型"></el-table-column>
            <el-table-column prop="APP" label="应用名称"></el-table-column>
            <el-table-column prop="SERVER" label="服务名称"></el-table-column>
            <el-table-column prop="PROJECT_VERSION" label="项目版本"></el-table-column>
            <el-table-column prop="RANDOM_HASH" label="进程ID"></el-table-column>
            <el-table-column label="操作">
                <template #default="scope">
                    <el-button type="text" @click="edit(scope.row, false)">查看</el-button>
                    <el-button type="text" @click="edit(scope.row, true)">编辑</el-button>
                    <el-button type="text" @click="publish(scope.row, true)">编译发布</el-button>
                    <el-button type="text" @click="publish(scope.row, false)">发布</el-button>
                    <el-button type="text" @click="getLog(scope.row)">查看日志</el-button>
                    <el-button type="text" @click="openProject(scope.row)">打开项目</el-button>
                </template>
            </el-table-column>
        </el-table>

        <el-dialog title="确认服务配置" v-model="mergeVisible">
            <el-form :model="form" label-width="100px">
                <el-form-item label="服务地址" prop="name">
                    <el-input v-model="form.filePath"></el-input>
                </el-form-item>
                <el-form-item label="服务类型" prop="url">
                    <el-input v-model="form.serverType"></el-input>
                </el-form-item>
                <el-form-item label="应用名称" prop="APP">
                    <el-input v-model="form.APP" disabled></el-input>
                </el-form-item>
                <el-form-item label="服务器名称" prop="SERVER">
                    <el-input v-model="form.SERVER" disabled></el-input>
                </el-form-item>
                <el-form-item label="备注" prop="COMMENT">
                    <el-input v-model="form.COMMENT"></el-input>
                </el-form-item>
                <el-form-item label="包路径" prop="PACKAGE_PATH">
                    <el-input v-model="form.PACKAGE_PATH"></el-input>
                </el-form-item>
                <el-form-item label="项目版本" prop="PROJECT_VERSION">
                    <el-input v-model="form.PROJECT_VERSION"></el-input>
                </el-form-item>
                <el-form-item label="TAF TICKET" prop="TAF_SERVER_TICKET">
                    <el-input v-model="form.TAF_SERVER_TICKET"></el-input>
                </el-form-item>
                <el-form-item label="TAF地址" prop="TAF_SERVER_PATH">
                    <el-input v-model="form.TAF_SERVER_PATH"></el-input>
                </el-form-item>
                <el-form-item label="编译命令" prop="TAF_SERVER_PATH">
                    <el-input v-model="form.BUILD_CMD"></el-input>
                </el-form-item>
                <el-form-item label="操作" v-show="isEdit">
                    <el-button type="primary" @click="handleCofirmMerge">确认</el-button>
                    <el-button @click="mergeVisible = false">取消</el-button>
                </el-form-item>
                <el-form-item label="发布" v-show="isPublish">
                    <el-button type="warning" @click="handleReleaseBeforeBuild">发布</el-button>
                    <el-button @click="mergeVisible = false">取消</el-button>
                </el-form-item>
            </el-form>
        </el-dialog>

        <el-dialog title="查看日志" v-model="logvisible" style="width: 700px;">
            <div style="background-color: black; padding:10px;">
                <div v-for="item in logdata" style="color:aliceblue">
                    {{ item }}
                </div>
            </div>
        </el-dialog>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import { OpenTarsReleaseFile, LoadConf, MergeConf, RunRelease, RunReleaseBeforeBuild, CheckBuildLog, OpenProject } from '../../../wailsjs/go/apis/App'
import { ElNotification } from 'element-plus';
import _ from 'lodash';

const mergeVisible = ref(false)
watch(mergeVisible, (v) => {
    if (v === false) {
        isEdit.value = false
        isPublish.value = false
    }
})
const originFormObj = {
    filePath: '',
    serverType: '',
    APP: '',
    SERVER: '',
    COMMENT: '',
    PACKAGE_PATH: '',
    PROJECT_VERSION: '',
    TAF_SERVER_TICKET: '',
    TAF_SERVER_PATH: '',
    BUILD_CMD: '',
}
type formType = typeof originFormObj

const form = ref(_.cloneDeep(originFormObj))
function loadTarRelease() {
    OpenTarsReleaseFile().then(res => {
        if (res == "") {
            return
        }
        mergeVisible.value = true
        isEdit.value = true
        const conf = JSON.parse(res)
        form.value = conf
        console.log('form.value', form.value);
    })
}

onMounted(() => {
    loadConf()
})
const confList = ref<formType[]>([])
function loadConf() {
    LoadConf().then(res => {
        console.log('res', res);
        if (res == '') {
            confList.value = []
        }
        else {
            confList.value = JSON.parse(res).map((v: formType) => {
                return _.assign({}, originFormObj, v)
            })
        }
    })
}

function handleCofirmMerge() {
    let item = confList.value.find(v => v.filePath == form.value.filePath)
    if (!item) {
        confList.value.push(form.value)
    }
    else {
        console.log('form.value', form.value);
        _.assign(item, form.value)
    }
    console.log('confList', confList.value);
    MergeConf(JSON.stringify(confList.value)).then(res => {
        ElNotification.success({
            title: '成功',
            message: '合并成功',
        })
        mergeVisible.value = false
    })

}
const isEdit = ref(false)
function edit(row: formType, t: boolean) {
    mergeVisible.value = true;
    isEdit.value = t
    form.value = _.cloneDeep(row)
}

const isPublish = ref(false)
function publish(row: formType, t: boolean) {
    if (t) {
        isPublish.value = t
        mergeVisible.value = true;
        form.value = _.cloneDeep(row)
    } else {
        RunRelease(row.filePath)
    }
}

function handleReleaseBeforeBuild() {
    if (!form.value.filePath || !form.value.BUILD_CMD) {
        ElNotification.error({
            title: '失败',
            message: '请输入文件路径和编译命令',
        })
        mergeVisible.value = false;
        return
    }
    RunReleaseBeforeBuild(form.value.filePath, form.value.BUILD_CMD).then(res => {
    }).catch(err => {
        ElNotification.error({
            title: '失败',
            message: err,
        })
    })
    mergeVisible.value = false;
    ElNotification.warning({
        title: '运行中',
        message: '运行中，请耐心等待',
    })

}

const logvisible = ref(false)
const logdata = ref<string[]>([])
let timer:any = null
watch(logvisible,function(val){
    if(!val){
        clearTimeout(timer)
    }
})
function getLog(row: formType) {
    logvisible.value = true
    timer = setInterval(()=>{
        CheckBuildLog(row.filePath).then(res => {
            console.log('res',res);
            logdata.value = res
        })
    },3000)
}

function openProject(row:formType){
    OpenProject(row.filePath).then(res=>{
        if(res !== ""){
            ElNotification.error({
                title: '失败',
                message: '打开失败' + res,
            })
        }
    })
}
</script>

<style scoped></style>