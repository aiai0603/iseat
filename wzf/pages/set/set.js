// pages/set/set.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    menuTapCurrent:0,
    erae:[[{
      id:"A区",
      max:50,
      now:12
    },{
      id:"B区",
      max:50,
      now:33
    },{
      id:"C区",
      max:50,
      now:50
    }],[{
      id:"A区",
      max:50,
      now:10
    },{
      id:"B区",
      max:50,
      now:20
    },{
      id:"C区",
      max:50,
      now:40
    }],[{
      id:"A区",
      max:50,
      now:2
    },{
      id:"B区",
      max:50,
      now:3
    },{
      id:"C区",
      max:50,
      now:0
    }]

    ]
  },

  menuTap:function(e){
    var current=e.currentTarget.dataset.current;//获取到绑定的数据
    //改变menuTapCurrent的值为当前选中的menu所绑定的数据
    this.setData({
    menuTapCurrent:current
    });
  
     
     
    },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {

  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady: function () {

  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow: function () {

  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide: function () {

  },

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload: function () {

  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh: function () {

  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom: function () {

  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage: function () {

  }

  
})