const API_BASE = 'http://localhost:8080'

Page({
  data: {
    tableId: '',
    nickname: ''
  },

  onLoad() {
    const userInfo = wx.getStorageSync('userInfo')
    if (userInfo) {
      this.setData({ nickname: userInfo.nickname })
    }
  },

  onInputTableId(e) {
    this.setData({ tableId: e.detail.value })
  },

  async onJoin() {
    const { tableId, nickname } = this.data
    if (!tableId || tableId.length !== 4) {
      wx.showToast({ title: '请输入4位牌桌ID', icon: 'none' })
      return
    }

    if (!nickname) {
      const nick = wx.getStorageSync('userInfo')?.nickname || '玩家' + Math.random().toString(36).slice(2, 6)
      this.setData({ nickname: nick })
    }

    try {
      const res = await wx.request({
        url: `${API_BASE}/api/table/join`,
        method: 'POST',
        data: {
          table_id: tableId,
          open_id: this.getOpenId(),
          nickname: this.data.nickname
        }
      })

      if (res.statusCode === 200) {
        wx.navigateTo({ url: `/pages/table/table?id=${tableId}` })
      } else {
        wx.showToast({ title: res.data.error || '加入失败', icon: 'none' })
      }
    } catch (err) {
      wx.showToast({ title: '连接失败', icon: 'none' })
    }
  },

  async onCreate() {
    const nickname = this.data.nickname || '玩家' + Math.random().toString(36).slice(2, 6)
    this.setData({ nickname })

    try {
      const res = await wx.request({
        url: `${API_BASE}/api/table/create`,
        method: 'POST',
        data: {
          open_id: this.getOpenId(),
          nickname: nickname
        }
      })

      if (res.statusCode === 200) {
        const tableId = res.data.table_id
        wx.navigateTo({ url: `/pages/table/table?id=${tableId}` })
      } else {
        wx.showToast({ title: res.data.error || '创建失败', icon: 'none' })
      }
    } catch (err) {
      wx.showToast({ title: '连接失败', icon: 'none' })
    }
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