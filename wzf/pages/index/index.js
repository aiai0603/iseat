const app = getApp();
const token = app.globalData.token;

Page({
  data: {

  },

  onLoad:function(){
   let data = wx.getStorageSync('token');
   if(data !== '')
   {
    wx.switchTab({
      url: '../notice/notice',
    })
   }
  },

 enter:function(){
  wx.setStorageSync('token', "1212121")
  wx.switchTab({
    url: '../notice/notice',
  })

 }
 
})
