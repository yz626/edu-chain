import request from './request'

export function login(data) {
  return request.post('/api/v1/user/login', data)
}

export function getUserInfo() {
  return request.get('/api/v1/user/info')
}

export function logout() {
  return request.post('/api/v1/user/logout')
}