const API_BASE = wx.getStorageSync('apiBase') || 'http://localhost:8080'

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

Page({
  data: {
    tableId: '',
    nickname: '',
    avatarUrl: ''
  },

  onLoad() {
    const userInfo = wx.getStorageSync('userInfo')
    if (userInfo) {
      this.setData({ nickname: userInfo.nickName, avatarUrl: userInfo.avatarUrl })
    }
  },

  onInputTableId(e) {
    this.setData({ tableId: e.detail.value })
  },

  getUserProfile() {
    return new Promise((resolve) => {
      wx.getUserProfile({
        desc: '用于游戏昵称和头像',
        success: (res) => {
          const userInfo = res.userInfo
          wx.setStorageSync('userInfo', userInfo)
          this.setData({
            nickname: userInfo.nickName,
            avatarUrl: userInfo.avatarUrl
          })
          resolve(userInfo)
        },
        fail: () => {
          resolve(null)
        }
      })
    })
  },

  async onJoin() {
    const { tableId } = this.data
    if (!tableId || tableId.length !== 4) {
      wx.showToast({ title: '请输入4位牌桌ID', icon: 'none' })
      return
    }

    let { nickname, avatarUrl } = this.data
    if (!nickname) {
      const info = await this.getUserProfile()
      if (info) {
        nickname = info.nickName
        avatarUrl = info.avatarUrl
      } else {
        nickname = '玩家' + Math.random().toString(36).slice(2, 6)
      }
      this.setData({ nickname, avatarUrl })
    }

    wx.request({
      url: `${API_BASE}/api/table/join`,
      method: 'POST',
      header: { 'content-type': 'application/json' },
      data: {
        table_id: tableId,
        open_id: this.getOpenId(),
        nickname: nickname,
        avatar_url: avatarUrl
      },
      success: (res) => {
        if (res.statusCode === 200) {
          wx.navigateTo({ url: `/pages/table/table?id=${tableId}` })
        } else {
          wx.showToast({ title: res.data.error || '加入失败', icon: 'none' })
        }
      },
      fail: () => {
        wx.showToast({ title: '连接失败', icon: 'none' })
      }
    })
  },

  async onCreate() {
    let { nickname, avatarUrl } = this.data
    if (!nickname) {
      const info = await this.getUserProfile()
      if (info) {
        nickname = info.nickName
        avatarUrl = info.avatarUrl
      } else {
        nickname = '玩家' + Math.random().toString(36).slice(2, 6)
      }
      this.setData({ nickname, avatarUrl })
    }

    wx.request({
      url: `${API_BASE}/api/table/create`,
      method: 'POST',
      header: { 'content-type': 'application/json' },
      data: {
        open_id: this.getOpenId(),
        nickname: nickname,
        avatar_url: avatarUrl
      },
      success: (res) => {
        if (res.statusCode === 200) {
          wx.navigateTo({ url: `/pages/table/table?id=${res.data.table_id}` })
        } else {
          wx.showToast({ title: res.data.error || '创建失败', icon: 'none' })
        }
      },
      fail: () => {
        wx.showToast({ title: '连接失败', icon: 'none' })
      }
    })
  },

  getOpenId() {
    let openId = wx.getStorageSync('openId')
    if (!openId) {
      openId = 'user_' + Math.random().toString(36).slice(2, 10)
      wx.setStorageSync('openId', openId)
    }
    return openId
  }
})