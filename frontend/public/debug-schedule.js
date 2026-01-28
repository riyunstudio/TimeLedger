// 診斷排課網格問題 - 將此代碼貼到瀏覽器控制台執行

(function diagnoseScheduleGrid() {
  console.log('=== 排課網格診斷 ===\n');

  // 查找 Vue 組件實例
  const scheduleGrids = document.querySelectorAll('.glass-card');
  console.log('找到的 ScheduleGrid 元素數量:', scheduleGrids.length);

  // 嘗試通過 Vue DevTools API 獲取組件實例
  let vm = null;
  if (window.__VUE_DEVTOOLS_GLOBAL_HOOK__) {
    const apps = window.__VUE_DEVTOOLS_GLOBAL_HOOK__?.apps;
    if (apps && apps.length > 0) {
      vm = apps[0]?.app?.config?.globalProperties;
      console.log('找到 Vue 實例');
    }
  }

  // 檢查頁面上的矩陣視圖元素
  const matrixView = document.querySelector('[class*="min-w-[800px]"]');
  console.log('矩陣視圖元素:', matrixView ? '存在' : '不存在');
  console.log('矩陣視圖可見性:', matrixView ? window.getComputedStyle(matrixView).display : 'N/A');

  // 檢查週曆視圖元素
  const calendarView = document.querySelector('[class*="min-w-[600px]"]');
  console.log('週曆視圖元素:', calendarView ? '存在' : '不存在');
  console.log('週曆視圖可見性:', calendarView ? window.getComputedStyle(calendarView).display : 'N/A');

  // 檢查資源面板
  const resourcePanel = document.querySelector('[class*="lg:w-80"]');
  console.log('資源面板元素:', resourcePanel ? '存在' : '不存在');

  // 檢查按鈕狀態
  const teacherTab = Array.from(document.querySelectorAll('button')).find(b => b.textContent.includes('老師列表'));
  console.log('老師列表按鈕:', teacherTab ? '找到' : '未找到');
  if (teacherTab) {
    console.log('按鈕被選中樣式:', teacherTab.classList.contains('ring-2') || teacherTab.classList.contains('bg-primary'));
  }

  // 檢查選擇狀態
  const selectedItems = document.querySelectorAll('[class*="ring-2"][class*="ring-primary"]');
  console.log('已選中項目數量:', selectedItems.length);

  // 檢查 localStorage 中的狀態
  console.log('\n=== localStorage 狀態 ===');
  const adminToken = localStorage.getItem('admin_token');
  console.log('管理員 token:', adminToken ? '存在' : '不存在');

  console.log('\n=== 建議操作 ===');
  console.log('1. 點擊"老師列表"標籤');
  console.log('2. 點擊一位老師名稱');
  console.log('3. 檢查頁面是否切換到矩陣視圖');
  console.log('4. 如果沒有，請重新執行此診斷腳本');

  return {
    matrixViewExists: !!matrixView,
    calendarViewExists: !!calendarView,
    resourcePanelExists: !!resourcePanel,
    selectedItemsCount: selectedItems.length
  };
})();
