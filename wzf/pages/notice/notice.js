Page({
  
  /**
  * 页面的初始数据
  */
  data: {
  ntnews: [
    {
      content:"近日，学校为保护学生，预防“新冠疫情”，采取了每天上午定点为学生宿舍喷洒消毒剂的措施。这一措施实行后，每天上午来图书馆的学生数量有较大增加，造成公共阅览区和部分阅览室人员密集。为了更好地保护广大同学的身体健康，经学校预防“非典”领导小组批准，每天中午将在全馆的公共阅览区域增加喷洒消毒剂一次。为此，图书馆需要在中午12：00—14：00间闭馆喷洒消毒药剂，敬请各位读者理解并支持我们的工作。新、老馆（包括老馆普通阅览室）开馆时间从4月28日起暂时调整如下：",
      time:"2020/01/01 15:12:22",
    },{
      content:"12111111111111111111111111111111111111111111111111111111111111111111111111111111",
      time:"2020/01/01 15:12:22",
    }
  ],
  /*轮播图 配置*/
  imgUrls: [
  '../../img/q.jpg',
  '../../img/t.jpg'
  ],
  indicatorDots: true, // 是否显示面板指示点
  autoplay: true, // 是否自动切换
  interval: 5000, // 自动切换时间间隔
  duration: 500, // 滑动动画时长
  circular: true, // 是否采用衔接滑动
  /*自定义轮播图 配置*/
  slider: [
  { id: '0', linkUrl: 'pages/index/index', picUrl: '../../img/q.jpg' },
  { id: '1', linkUrl: 'pages/index/index', picUrl: '../../img/t.jpg' }
  ],
  swiperCurrent: 0
  },
   
  /**
  * 生命周期函数--监听页面加载
  */
  onLoad: function (options) {
   
  },
   
  //轮播图的切换事件 
  swiperChange: function (e) {
  //只要把切换后当前的index传给<swiper>组件的current属性即可 
  this.setData({
  swiperCurrent: e.detail.current
  })
  },
  //点击指示点切换 
  chuangEvent: function (e) {
  this.setData({
  swiperCurrent: e.currentTarget.id
  })
  }
 })