const API_BASE = 'http://localhost:8080'

Page({
  data: {
    tableId: '',
    status: 'waiting',
    players: [],
    is_owner: false,
    showPayModal: false,
    showSettleModal: false,
    payOptions: [100, 500, 1000, 2000, 3000, 5000, 10000],
    payIndex: -1,
    targetOpenId: '',
    results: []
  },

  onLoad(options) {
    const tableId = options.id
    this.setData({ tableId })
    this.fetchTable()
    this.connectWebSocket()
  },

  onUnload() {
    if (this.ws) {
      this.ws.close()
    }
  },

  async fetchTable() {
    try {
      const res = await wx.request({
        url: `${API_BASE}/api/table/${this.data.tableId}`
      })
      if (res.statusCode === 200) {
        const { status, players } = res.data
        const openId = this.getOpenId()
        const is_owner = players.some(p => p.is_owner && p.open_id === openId)
        this.setData({ status, players, is_owner })
      }
    } catch (err) {
      console.error('fetch table error:', err)
    }
  },

  connectWebSocket() {
    const openId = this.getOpenId()
    this.ws = wx.connectSocket({
      url: `${API_BASE.replace('http', 'ws')}/ws/table/${this.data.tableId}/${openId}`
    })

    this.ws.onMessage((res) => {
      const data = JSON.parse(res.data)
      if (data.type === 'state_update') {
        this.setData({
          status: data.table.status,
          players: data.table.players
        })
      }
    })

    this.ws.onClose(() => {
      console.log('ws closed')
    })
  },

  getOpenId() {
    return wx.getStorageSync('openId') || ''
  },

  onPlayerTap(e) {
    const openId = e.currentTarget.dataset.openid
    if (openId !== this.getOpenId()) {
      this.setData({
        showPayModal: true,
        targetOpenId: openId
      })
    }
  },

  onPay() {
    const openId = this.getOpenId()
    const players = this.data.players
    if (players.length > 1) {
      const other = players.find(p => p.open_id !== openId)
      if (other) {
        this.setData({
          showPayModal: true,
          targetOpenId: other.open_id
        })
      }
    }
  },

  onPayChange(e) {
    this.setData({ payIndex: e.detail.value })
  },

  async onConfirmPay() {
    const { payIndex, payOptions, tableId, targetOpenId } = this.data
    if (payIndex === -1) {
      wx.showToast({ title: '请选择金额', icon: 'none' })
      return
    }

    const amount = payOptions[payIndex]
    try {
      const res = await wx.request({
        url: `${API_BASE}/api/table/pay`,
        method: 'POST',
        data: {
          table_id: tableId,
          from_open_id: this.getOpenId(),
          to_open_id: targetOpenId,
          amount: amount
        }
      })

      if (res.statusCode === 200) {
        wx.showToast({ title: '支付成功', icon: 'success' })
      }
    } catch (err) {
      wx.showToast({ title: '支付失败', icon: 'none' })
    }

    this.setData({ showPayModal: false, payIndex: -1 })
  },

  onCancelPay() {
    this.setData({ showPayModal: false, payIndex: -1 })
  },

  async onSettle() {
    try {
      const res = await wx.request({
        url: `${API_BASE}/api/table/settle`,
        method: 'POST',
        data: { table_id: this.data.tableId }
      })

      if (res.statusCode === 200) {
        this.setData({
          showSettleModal: true,
          results: res.data.results
        })
      }
    } catch (err) {
      wx.showToast({ title: '结算失败', icon: 'none' })
    }
  },

  async onDismiss() {
    try {
      await wx.request({
        url: `${API_BASE}/api/table/dismiss`,
        method: 'POST',
        data: { table_id: this.data.tableId }
      })
    } catch (err) {}

    wx.navigateBack()
  },

  onNextRound() {
    this.setData({ showSettleModal: false })
    this.fetchTable()
  },

  async onLeave() {
    try {
      await wx.request({
        url: `${API_BASE}/api/table/leave`,
        method: 'POST',
        data: {
          table_id: this.data.tableId,
          open_id: this.getOpenId()
        }
      })
    } catch (err) {}

    wx.navigateBack()
  }
})