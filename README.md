# mycha
一个学习项目
在观看golang多并发编程这本书中,书中的一个爬虫软件让我十分感兴趣.也算是学习的一种,这个项目是临摹的一个项目.尚在学习中,所以许多的测试文件都没有专门的写.尚未完成
总之
项目的思路大概是这样
程序中分涉及的要数为 四种角色  一个总的调度器
三种角色 分别为Downloader  Analyer Pipeline  error  分别对应的功能为  获取到请求对象 响应对象 处理分析响应对象 处理错误对象
由调度器 scheduler 为每个对象进行注册 登记 和启动服务


功能就很简单 下载器 downloader 负责请求一个url 并将请求的对象放入 downloader的缓存池中 这里的操作 scheduler 会开启一个独立的goroutine来进行处理并且支持多个 downloader对象运行 每个对象在独立的goroutine中处理

同样analyer则是将请求的对象转化对应文本结构和解析 并保存analyer的缓存池中 也是独立goroutine来进行操作 支持多个analyer运行

pipeline 则把文本转话成对应的永久数据存储 如文件 如数组

当程序启动时,会有许多的goroutine同时进行 这其中涉及到的并发抢夺情况很多 所以程序的内容是很有学习意义的.

下次更新会将整个程序临摹完 

