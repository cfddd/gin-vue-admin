<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
        @keyup.enter="onSubmit">
        <el-form-item label="创建日期" prop="createdAt">
          <template #label>
            <span>
              创建日期
              <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
                <el-icon>
                  <QuestionFilled />
                </el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker v-model="searchInfo.startCreatedAt" type="date" placeholder="开始日期"
            :disabled-date="time => searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"></el-date-picker>
          —
          <el-date-picker v-model="searchInfo.endCreatedAt" type="date" placeholder="结束日期"
            :disabled-date="time => searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"></el-date-picker>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
        </el-form-item>

        <el-form-item label="查找用户">
          <template #label>
            <span>
              查找用户
              <el-tooltip content="搜索栏中填入的昵称必须完全匹配">
                <el-icon>
                  <QuestionFilled />
                </el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-input v-model="searchInfo.user_name" class="keyword" placeholder="请输入要查找的用户名" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" icon="search" @click="getTableData">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>

      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID"
        @selection-change="handleSelectionChange">
        <el-table-column type="index" align="center" :resizable="false" label="序号" width="100">
        </el-table-column>
        <el-table-column align="center" label="日期" width="180">
          <template #default="scope">{{ formatDate(scope.row.date).substring(0, 10) }}</template>
        </el-table-column>
        <el-table-column align="center" label="用户昵称" prop="user_name" width="180" />
        <el-table-column align="center" label="打卡详情" width="180">
          <template #default="scope">
            <el-button type="primary" link icon="info-filled" class="table-button"
              @click="showRecordFunc(scope.row)">查看</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination layout="total, sizes, prev, pager, next, jumper" :current-page="page" :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]" :total="total" @current-change="handleCurrentChange"
          @size-change="handleSizeChange" />
      </div>
    </div>
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :width="`61%`" :height="`100%`" destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80">
        <el-form-item label="代码:" prop="code">
          <el-scrollbar style="height: 400px; width: 100%;">
            <pre><code class="language-javascript">{{ formData.code }}</code></pre>
          </el-scrollbar>
        </el-form-item>
        <el-form-item label="题目链接:" prop="link">
          <el-link :href="formData.link" target="_blank">{{ formData.link }}</el-link>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="copyCode">复制代码</el-button>
          <el-button type="primary" @click="closeDialog">返 回</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import hljs from 'highlight.js';
import 'highlight.js/styles/default.css';

export default {
  name: 'TEST',
  mounted() {
    // 初始化highlight.js
    hljs.initHighlightingOnLoad();
  },
};
</script>

<script setup>
import {
  findDailyAlgorithmRecord,
  getDailyAlgorithmRecordList
} from '@/api/sysDailyAlgorithmRecord'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  code: '',
  link: '',
})

// 验证规则
const rule = reactive({
})

const searchRule = reactive({
  createdAt: [
    {
      validator: (rule, value, callback) => {
        if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
          callback(new Error('请填写结束日期'))
        } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
          callback(new Error('请填写开始日期'))
        } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
          callback(new Error('开始日期应当早于结束日期'))
        } else {
          callback()
        }
      }, trigger: 'change'
    }
  ],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async (valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async () => {
  const table = await getDailyAlgorithmRecordList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 更新行
const showRecordFunc = async (row) => {
  const res = await findDailyAlgorithmRecord({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.reDAR
    dialogFormVisible.value = true
  }
}

// 复制代码
const copyCode = () => {
  const codeToCopy = formData.value.code;

  // 使用Clipboard API复制代码
  navigator.clipboard.writeText(codeToCopy)
    .then(() => {
      ElMessage({
        type: 'success',
        message: '代码已成功复制到剪贴板'
      })
    })
    .catch((error) => {
      ElMessage({
        type: 'fail',
        message: error
      })
    });
}
// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    code: '',
    link: '',
  }
}


</script>

<style></style>
