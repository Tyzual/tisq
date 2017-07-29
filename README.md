# TISQ

# API说明
## 统一要求
1. HTTP Method  
	POST
2. 参数要求  
	先将key和value进行URL编码（js调用encodeURIComponent()函数）  
	然后再将 key和value以```key1=value1&key2=value2```的形式作为body发送  
	比如需要传递的参数为
	```
		"content:"测试content++"
		"user":"测试用户%%"
	```
	则向body发送的字符串为
	```
	content=%E6%B5%8B%E8%AF%95content%2B%2B&user=%E6%B5%8B%E8%AF%95%E7%94%A8%E6%88%B7%25%25
	```
### 添加评论接口
#### URL
	addComment
#### 参数以及含义

| 参数 | 含义 |
| --- | --- |
| domain | 博客的域名 |
| email | 评论者的email |
| displayname | 评论者显示的名字(昵称)，可不传 |
| site | 评论者的主页，可不传 |
| articlekey | 评论文章的特征码(可以使用文章的URL地址) |
| content | 评论内容 |
| replyid | 如果评论是回复某条评论，则在这里填写回复评论的评论ID，如果不是回复，不要传递这个参数 |
| lastcommentid | 客户端最新一条评论的评论id，若不传这个值，服务器会返回articlekey下的所有评论，否则返回lastcommentid以后的评论。若找不到lastcommentid所对应的评论。则返回错误。 |

#### 返回值
若成功，HTTP状态码200，Body为JSON格式内容。示例：
```
{
	/* 评论信息，内容为一个数组，数组已经以评论的先后顺序排好序，最后的评论在第一个元素。 */
	Comment: [
		{
			/* 评论的用户，UserID为用户信息中的key */
			UserID: "94267782ed9f6739daf63c9ab2964767", 
			/* 评论内容 */
			Content: "测试评论1", 
			/* 评论ID，在查询或者添加评论时可作为参数使用 */
			CommentID: 25, 
			/* 评论的时间，格式为UNIX时间戳 */
			CreateTime: 1501340545,
			/* 若有这个字段，这个字段表示这条评论是回复CommentID为13的评论（可以没有） */
			ReplyCommentID: 13
		}, 
		{
			UserID: "bc9e6c3187434cbc0826698490abe8c0", 
			Content: "测试评论2", 
			CommentID: 26, 
			CreateTime: 1501343024
		}
	],
	/* 用户信息 */
	User: {
		/* 用户ID，对应评论信息中的UserID字段 */
		bc9e6c3187434cbc0826698490abe8c0: {
			/* 用户的Email */
			Email: "e.tyzual@gmail.com", 
			/* 用户的昵称 (可以没有) */
			DisplayName: "tyzual",
			/* 用户的网站 (可以没有) */
			Site: "tyzual.com"
		},
		94267782ed9f6739daf63c9ab2964767: {
			Email: "echizen@foxmail.com", 
			DisplayName: "echizen"
		} 
	} 
}
```

若失败，HTTP状态码不为200，Body为JSON格式内容。
JSON个字段意义如下
```
{
	/* 错误码 */
	ErrorNum: 1000，
	/* 错误描述 */
	ErrorMsg: "未知错误"
}
```

### 获取评论列表接口
#### URL
	commentList
#### 参数以及含义

| 参数 | 含义 |
| --- | --- |
|domain | 博客的域名 |
| articlekey | 评论文章的特征码(可以使用文章的URL地址) |
| lastcommentid | 客户端最新一条评论的评论id，若不传这个值，服务器会返回articlekey下的所有评论，否则返回lastcommentid以后的评论。若找不到lastcommentid所对应的评论。则返回错误。 |

#### 返回值
与 添加评论接口 返回值一样

## FAQ
### Q: TISQ 是啥  
> A: 目前，TISQ 的定位是一个简单的评论管理后台。支持多站点管理。

### Q: TISQ 只有后台吗？不开发前台吗？
> A: 是的，TISQ 只有后台，并且通过RESTAPI提供服务。

### Q: 为啥没有前台
> A: 因为我不会前台🌚

### Q: 可是我看这个后台也很挫的样子
> A: 是很挫。挫到连我自己都不用🌚。如果哪天搜狐畅言停止服务的话说不定我会更加有动力来更新这个后台。

### Q: 以后会不会加上XX功能
> A: 不知道。因为现在有搜狐畅言可以用，所以我并不急着完善这个后台系统。
