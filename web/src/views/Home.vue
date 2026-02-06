<template>
  <div class="home">
    <h1>教育链系统</h1>
    <p>欢迎使用 EduChain 教育区块链系统</p>
    <div v-if="userStore.token">
      <p>已登录</p>
      <button @click="handleLogout">退出登录</button>
    </div>
    <div v-else>
      <router-link to="/login">登录</router-link>
    </div>
  </div>
</template>

<script setup>
import { useUserStore } from '../stores/user'
import { useRouter } from 'vue-router'
import { logout } from '../api/user'

const userStore = useUserStore()
const router = useRouter()

async function handleLogout() {
  try {
    await logout()
  } catch (e) {
    console.error(e)
  }
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.home {
  text-align: center;
  padding: 50px;
}
</style>