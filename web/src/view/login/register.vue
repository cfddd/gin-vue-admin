<template>
  <div>
    <!-- 表单项 -->
    <el-form ref="registerForm" :model="registerFormData" :rules="rules" :validate-on-rule-change="false"
      @keyup.enter="submitForm">
      <el-form-item prop="username" class="mb-6">
        <el-input v-model="registerFormData.username" size="large" placeholder="请输入用户名" suffix-icon="user" />
      </el-form-item>
      <el-form-item prop="password" class="mb-6">
        <el-input v-model="registerFormData.password" show-passsword size="large" type="password" placeholder="请输入密码" />
      </el-form-item>
      <el-form-item prop="repassword" class="mb-6">
        <el-input v-model="registerFormData.repassword" show-password size="large" type="password" placeholder="请再次输入密码" />
      </el-form-item>
      <el-form-item prop="qq" class="mb-6">
        <el-input v-model="registerFormData.qq" size="large" placeholder="请输入QQ号" />
      </el-form-item>
      <el-form-item prop="email" class="mb-6">
        <el-input v-model="registerFormData.email" size="large" placeholder="请输入邮箱" />
      </el-form-item>
      <el-form-item prop="phone" class="mb-6">
        <el-input v-model="registerFormData.phone" size="large" placeholder="请输入手机号" />
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

// 在这里导入其他需要的模块或组件

const router = useRouter()
const registerForm = ref(null)
const picPath = ref('')
const userStore = useUserStore()

// 验证函数
const validators = {
  checkUsername: (rule, value, callback) => {
    value.length < 5 ? callback(new Error('请输入正确的用户名')) : callback()
  },
  checkPassword: (rule, value, callback) => {
    value.length < 6 ? callback(new Error('请输入正确的密码')) : callback()
  },
  checkRePassword: (rule, value, callback) => {
    registerFormData.password !== value ? callback(new Error('密码不一致')) : callback()
  }
}

// 获取验证码
const registerVerify = async () => {
  const ele = await captcha({})
  rules.captcha.push({
    max: ele.data.captchaLength,
    min: ele.data.captchaLength,
    message: `请输入${ele.data.captchaLength}位验证码`,
    trigger: 'blur',
  })
  picPath.value = ele.data.picPath
  registerFormData.captchaId = ele.data.captchaId
  registerFormData.openCaptcha = ele.data.openCaptcha
}

registerVerify()

const registerFormData = reactive({
  username: '',
  password: '',
  repassword: '',
  qq: '',
  email: '',
  phone: '',
  captcha: '',
  captchaId: '',
})

const rules = reactive({
  username: [{ required: true,validator: validators.checkUsername, trigger: 'blur' }],
  password: [{ required: true,validator: validators.checkPassword, trigger: 'blur' }],
  repassword: [{ required: true,validator: validators.checkRePassword, trigger: 'blur' }],
  qq:[{required: true, message: '请输入QQ号', trigger: 'blur'}],
  email:[{required: true, message: '请输入邮箱', trigger: 'blur'}],
  phone:[{required: true, message: '请输入验证码', trigger: 'blur'}],
  captcha: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { message: '验证码格式不正确', trigger: 'blur' },
  ],
})

// 注册函数
const register = async () => {
  return await userStore.Register(registerFormData)
}

// 提交表单函数
const submitForm = () => {
  registerForm.value.validate(async (valid) => {
    if (valid) {
      const flag = await register()
      if (!flag) {
        registerVerify()
      }
    } else {
      ElMessage.error({ message:'请正确填写注册信息', showClose:true })
      registerVerify()
      return false
    }
  })
}
</script>