<template>
  <nav class="navbar">
    <div class="navbar-brand">
      <router-link to="/">EduChain</router-link>
    </div>
    <div class="navbar-menu">
      <template v-if="userStore.token">
        <span class="navbar-item">欢迎, {{ username }}</span>
        <button class="navbar-item logout-btn" @click="handleLogout">退出</button>
      </template>
      <template v-else>
        <router-link class="navbar-item" to="/login">登录</router-link>
      </template>
    </div>
  </nav>
</template>

<script setup>
import { computed } from 'vue'
import { useUserStore } from '../stores/user'
import { useRouter } from 'vue-router'
import { logout } from '../api/user'

const userStore = useUserStore()
const router = useRouter()

const username = computed(() => {
  return userStore.userInfo?.username || '用户'
})

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
.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background-color: #42b983;
  color: white;
}

.navbar-brand a {
  font-size: 1.5rem;
  font-weight: bold;
  color: white;
  text-decoration: none;
}

.navbar-menu {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.navbar-item {
  color: white;
  text-decoration: none;
}

.logout-btn {
  background: transparent;
  border: 1px solid white;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
}

.logout-btn:hover {
  background: white;
  color: #42b983;
}
</style>