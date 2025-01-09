<template>
    <el-container>
        <el-container>
            <el-aside width="150px">
                <menus></menus>
            </el-aside>
            <el-container>
                <el-main>
                    <router-view></router-view>
                </el-main>
            </el-container>
        </el-container>
    </el-container>
</template>

<script setup lang="ts">
import menus from './menus.vue';
</script>

<style scoped></style>