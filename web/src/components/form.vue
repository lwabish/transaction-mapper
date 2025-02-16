<template>
  <div class="file-upload-form">
    <form @submit.prevent="handleSubmit">
      <div class="form-group">
        <label for="bank">银行：</label>
        <select id="bank" v-model="selectedBank" :disabled="isLoading">
          <option value="">请选择银行</option>
          <option v-for="bank in banks" :key="bank" :value="bank">
            {{ bank }}
          </option>
        </select>
      </div>

      <div class="form-group">
        <label for="app">记账app：</label>
        <select id="selectApp" v-model="selectedApp" :disabled="isLoading">
          <option value="">请选择记账app</option>
          <option v-for="app in apps" :key="app" :value="app">
            {{ app }}
          </option>
        </select>
      </div>

      <div class="form-group">
        <label for="account">导入账户：</label>
        <input type="text" id="account" v-model="account" :disabled="isLoading" placeholder="必填" required>
      </div>

      <div class="form-group">
        <label for="accountType">账户类型：</label>
        <input type="text" id="accountType" v-model="accountType" :disabled="isLoading" placeholder="可选">
      </div>

      <div class="form-group">
        <label for="file">选择文件：</label>
        <input type="file" id="file" ref="fileInput" @change="handleFileChange" required>
      </div>
      <div v-if="selectedFile" class="selected-file">
        已选择文件: {{ selectedFile.name }}
      </div>
      <button type="submit" :disabled="isUploading || !selectedBank || isLoading">
        {{ isUploading ? '上传中...' : '上传' }}
      </button>
    </form>
    <div v-if="status" :class="['status', status.type]">
      {{ status.message }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';

const fileInput = ref(null);
const selectedFile = ref(null);
const isUploading = ref(false);
const isLoading = ref(false);
const status = ref(null);

const banks = ref([]);
const selectedBank = ref('');

const apps = ref([]);
const selectedApp = ref('');

const account = ref('');

const accountType = ref('');

// 假设你的API基础URL是这个，请根据实际情况修改
let API_BASE_URL = '';

onMounted(async () => {
  const protocol = window.location.protocol;
  const host = window.location.host;
  API_BASE_URL = `${protocol}//${host}/api/v1`;
  if (API_BASE_URL.includes('localhost')) {
    API_BASE_URL = 'http://localhost:8080/api/v1';
  }
  await fetchRemoteData();
});

const fetchRemoteData = async () => {
  isLoading.value = true;
  status.value = null;
  try {
    const [banksResponse, appsResponse] = await Promise.all([
      axios.get(`${API_BASE_URL}/banks`),
      axios.get(`${API_BASE_URL}/apps`)
    ]);
    banks.value = banksResponse.data.data;
    apps.value = appsResponse.data.data;
    if (banks.value.length > 0) {
      selectedBank.value = banks.value[0];
    }
    if (apps.value.length > 0) {
      selectedApp.value = apps.value[0];
    }
  } catch (error) {
    console.error('获取远程数据失败:', error);
    status.value = { type: 'error', message: '获取远程数据失败，请刷新页面重试' };
  } finally {
    isLoading.value = false;
  }
};

const handleFileChange = (event) => {
  selectedFile.value = event.target.files[0];
};

const handleSubmit = async () => {
  if (!selectedFile.value || !selectedBank.value || !account.value) {
    status.value = { type: 'error', message: '请填写必填项' };
    return;
  }

  isUploading.value = true;
  status.value = null;

  const formData = new FormData();
  formData.append('input', selectedFile.value);
  formData.append('bank', selectedBank.value);
  formData.append('app', selectedApp.value);
  formData.append('account', account.value);
  formData.append('accountType', accountType.value);

  try {
    const response = await axios.post(`${API_BASE_URL}/transform`, formData, {
      responseType: 'blob', // 重要：设置响应类型为blob
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: (progressEvent) => {
        const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
        console.log('上传进度:', percentCompleted);
        // 这里可以添加上传进度显示的逻辑
      }
    });

    // 检查响应头，确定是文件下载还是错误消息
    const contentType = response.headers['content-type'];
    if (contentType && contentType.includes('application/json')) {
      // 如果是 JSON，可能是错误消息
      const reader = new FileReader();
      reader.onload = () => {
        const errorResponse = JSON.parse(reader.result);
        console.log(errorResponse);
        status.value = { type: 'error', message: errorResponse.error || '处理失败，请重试' };
      };
      reader.readAsText(response.data);
    } else {
      // 否则，作为文件下载处理
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      console.log(response.headers);
      const filename = response.headers['content-disposition']
        ? response.headers['content-disposition'].split('filename=')[1]
        : 'processed_file';
      link.setAttribute('download', filename);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      status.value = { type: 'success', message: '文件处理成功，正在下载' };
    }

    selectedFile.value = null;
    if (fileInput.value) {
      fileInput.value.value = '';
    }
  } catch (error) {
    console.error('处理错误:', error);
    status.value = { type: 'error', message: '文件处理失败，请重试' };
  } finally {
    isUploading.value = false;
  }
};
</script>

<style scoped>
.file-upload-form {
  max-width: 500px;
  margin: 0 auto;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
}

input[type="file"],
input[type="text"],
select {
  display: block;
  width: 98%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-sizing: border-box;
}

select:disabled {
  background-color: #f9f9f9;
  cursor: not-allowed;
}

button {
  padding: 10px 20px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.selected-file {
  margin-top: 10px;
  font-style: italic;
}

.status {
  margin-top: 15px;
  padding: 10px;
  border-radius: 4px;
}

.status.success {
  background-color: #d4edda;
  color: #155724;
}

.status.error {
  background-color: #f8d7da;
  color: #721c24;
}
</style>