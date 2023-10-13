<template>
  <div>
    <!-- 表单项 -->
    <el-form ref="registerForm" :model="registerFormData" :rule="rules" :validate-on-rule-change="false"
      @keyup.enter="submitForm">
      <el-form-item prop="username" class="mb-6">
        <el-input v-model="registerFormData.username" size="large" placeholder="请输入用户名" suffix-icon="user" />
      </el-form-item>
      <el-form-item prop="password" class="mb-6">
        <el-input v-model="registerFormData.password" show-password size="large" type="password" placeholder="请输入密码" />
      </el-form-item>
      <el-form-item v-if="registerFormData.openCaptcha" prop="captcha" class="mb-6">
        <div class="flex w-full justify-between">
          <el-input v-model="registerFormData.captcha" placeholder="请输入验证码" size="large" class="flex-1 mr-5" />
          <div class="w-1/3 h-11 bg-[#c3d4f2] rounded">
            <img v-if="picPath" class="w-full h-full" :src="picPath" alt="请输入验证码" @click="registerVerify()">
          </div>
        </div>
      </el-form-item>
      <el-form-item class="mb-6">
        <el-button class="shadow shadow-blue-600 h-11 w-full" type="primary" size="large"
          @click="submitForm">注册</el-button>
      </el-form-item>
      <!-- 其他表单项 -->
    </el-form>


  </div>
</template>

<script setup>
import { captcha } from '@/api/user'

import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'

const router = useRouter()
// 验证函数
const checkUsername = (rule, value, callback) => {
  if (value.length < 5) {
    return callback(new Error('请输入正确的用户名'))
  } else {
    callback()
  }
}
const checkPassword = (rule, value, callback) => {
  if (value.length < 6) {
    return callback(new Error('请输入正确的密码'))
  } else {
    callback()
  }
}

// 获取验证码
const registerVerify = () => {
  captcha({}).then(async (ele) => {
    rules.captcha.push({
      max: ele.data.captchaLength,
      min: ele.data.captchaLength,
      message: `请输入${ele.data.captchaLength}位验证码`,
      trigger: 'blur',
    })
    picPath.value = ele.data.picPath
    registerFormData.captchaId = ele.data.captchaId
    registerFormData.openCaptcha = ele.data.openCaptcha
  })
}
registerVerify()



const registerForm = ref(null)
const picPath = ref('')
const registerFormData = reactive({
  username: '',
  password: '',
  captcha: '',
  captchaId: '',
})

const rules = reactive({
  username: [{ validator: checkUsername, trigger: 'blur' }],
  password: [{ validator: checkPassword, trigger: 'blur' }],
  captcha: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    {
      message: '验证码格式不正确',
      trigger: 'blur',
    },
  ],
})

const userStore = useUserStore()

const register = async () => {
  console.log(registerFormData.value)
  return await userStore.Register(registerFormData)
}
const submitForm = () => {
  registerForm.value.validate(async (v) => {
    if (v) {
      const flag = await register()
      if (!flag) {
        registerVerify()
      }
      console.log(flag)
    } else {
      ElMessage({
        type: 'error',
        message: '请正确填写注册信息',
        showClose: true,
      })
      registerVerify()
      return false
    }
  })
}


</script>