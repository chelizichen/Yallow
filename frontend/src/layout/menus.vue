<template>
    <div>
        <el-menu
        style="height: 100vh;"
        active-text-color="#ffd04b"
        background-color="black"
        default-active="1"
        text-color="#fff"
        @select="handleSelect"
      >
        <el-menu-item index="0" disabled>
          <span>TAF发布平台</span>
        </el-menu-item>
        <el-menu-item index="1">
          <span >导入文件</span>
        </el-menu-item>
        <el-menu-item index="2">
          <span >JSON判断</span>
        </el-menu-item>
        <el-menu-item index="3">
          <span >设置</span>
        </el-menu-item>
      </el-menu>
    </div>
</template>

<script setup lang="ts">
import router, { routes } from '../router';

function handleSelect(index:string){
    const idx = Number(index) - 1
    const route = routes[idx]
    router.push(route.path)
}
</script>

<style scoped>

</style>