<template>
    <div>
        <iframe src="view/JavascriptVisualRelease/about.html" width="100%" height="1000"></iframe>
    </div>
</template>
  
<script>
export default {
    methods: {
        goToAbout() {
            window.location.href = './about.html';
        }
    }
}
</script>
  
<style scoped>
/* 可以在这里添加组件的样式 */
</style>