App({
  onLaunch() {
    const userInfo = wx.getStorageSync('userInfo')
    if (!userInfo) {
      wx.login({
        success: res => {
          console.log('login code:', res.code)
        }
      })
    }
  }
})