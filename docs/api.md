Flicker-API文档

## 基本格式

地址: 暂定

API前缀: /api/v1，**文档中省略前缀**

### 数据交换格式

#### 成功

```json
{
    "success": true,
    "message": "",
    "data": {}
}
```

data 可以为任何类型，***之后文档中只写出data的数据***

#### 失败

```json
{
    "success": false,
    "message": "错误提示信息",
    "error": "错误显示信息",
    "data": null
}
```

失败会返回200以外的http状态码

**错误提示信息是供用户看的文字提示，出现错误时直接使用即可**

之后API文档中不会再写失败回复。

## 身份验证

使用jwt进行身份验证，token位于http头的Authorization字段，以Bearer开头，例

Authorization: Bearer 一个巨长的token

## 基本 /

### POST /verify 发送邮件验证码

- 无需jwt

- 验证码请求间隔至少为1分钟
- 验证码缓存有效期暂定15分钟

#### 请求参数

- mail：字符串，必需，邮件地址

#### 请求示例

```json
{
    "mail":"12345678@qq.com"
}
```

#### 响应参数

无额外参数

#### 响应示例

```json
"ok"
```

### POST /signup 新建用户

#### 请求参数

- mail：字符串，必需，邮件地址
- username：字符串，必需，用户名
- password：字符串，必需，密码
- code：字符串，必需，邮件验证码

#### 请求示例

```json
{
    "mail":"123456789@qq.com",
    "username":"xxx",
    "password":"xxx",
    "code":"12345"
}
```

#### 响应参数

无额外参数

#### 响应示例

```json
"ok"
```

### POST /login 登录

#### 请求参数

- mail：字符串，必需，邮件地址
- password：字符串，必需，密码

#### 请求示例

```json
{
    "mail":"123456789@qq.com",
    "password":"xxx"
}
```

#### 响应参数

字符串，JWT Token

#### 响应示例

```json
"jwt"
```

## 用户 /user

### PUT / 修改用户

- 暂不支持修改邮箱

#### 请求参数

- username：字符串，可选，用户名
- password：字符串，可选，密码
- code：字符串，当`password`不为空时必需，邮件验证码
- avatar：字符串，可选，头像url
    - 这个之后再加，需要考虑一下要不要上七牛

#### 请求示例

```json
{
    "username":"xxx",
    "password":"xxx",
    "code":"12345"
}
```

#### 响应参数

无额外参数

#### 响应示例

```json
"ok"
```



### GET /:user_id 获取用户

#### 请求参数

- user_id：路径参数，用户id，缺省时获取当前登录用户信息

#### 响应参数

- username：字符串，可选，用户名
- avatar：字符串，可选，头像url

#### 响应示例

```json
{
    "id": "6194e37bd786c6c07ea8a4fb",
    "username":"xxx",
    "avatar":"https://example.com/1.jpg",
    "favorite": [
        "6199cc3db4600da8e6102dac"
    ]
}
```



### POST /favorite

#### 请求参数

- cardset_id：字符串，必需，卡片集ID
- liked：字符串，可选，表示请求前用户是否收藏了该卡片集，请求将使此状态反转。默认值为`false`

#### 请求示例

```json
{
    "cardset_id": "xxx",
    "liked": false
}
```

#### 响应参数

无额外参数

#### 响应示例

```json
"ok"
```



## 卡片集 /cardset

- alias 题库

### POST / 新建空卡片集

#### 请求参数

- name：字符串，必需，名称
- description：字符串，可选，卡片集描述
- template：字符串，可选，模板
    - 模板一定程度上规定了卡片的样式与布局
    - 暂定为HTML格式，具体如何设计看前端
    - 当此字段为空时应使用默认的样式
- access：整数，可选，访问权限
    - 0-仅创建者可见，1-所有人可见
    - 默认权限为0

#### 请求示例

```json
{
    "name":"test-cards",
    "description":"test test.",
    "access":1
}
```

#### 响应参数

字符串，卡片集id

#### 响应示例

```json
"id"
```

### PUT /:id 修改卡片集信息

#### 请求参数

- id：字符串，必需，卡片集id
- name：字符串，可选，名称
    - 参数留空时表示不做修改，下同
- description：字符串，可选，卡片集描述
- template：字符串，可选，模板
    - 模板一定程度上规定了卡片的样式与布局
    - 暂定为HTML格式，具体如何设计看前端
    - 当此字段为空时应使用默认的样式
- access：整数，可选，访问权限

#### 请求示例

```json
{
    "name":"test-cards",
    "description":"test test.",
    "access":1
}
```

#### 响应参数

无额外参数

#### 响应示例

```json
"ok"
```

### DELETE /:id 删除卡片集

#### 请求参数

- id：字符串，必需，卡片集id

#### 响应参数

无额外参数

#### 响应示例

```json
"ok"
```

### GET /:id 获取卡片集

#### 请求参数

- id：字符串，必需，卡片集id

#### 响应参数

- name：字符串，名称
- description：字符串，卡片集描述
- access：整数，访问权限
- cards：字符串数组，卡片的id列表

#### 响应示例

```json
{
    "id":"id",
    "name":"test-cards",
    "description":"test test.",
    "access":1,
    "cards":[]
}
```

## 卡片 /cardset/:cardset_id/card

- cardset_id：路径参数，必需，卡片集id

### POST / 新建卡片

#### 请求参数

- question：字符串，必需，题面
- answer：字符串，必需，答案
- image：字符串，可选，图片url
- audio：字符串，可选，音频url

#### 请求示例

```json
{
    "question":"question",
    "answer":"answer",
    "image":"https://example.com/1.jpg",
    "audio":"https://example.com/1.wav"
}
```

#### 响应参数

卡片id

#### 响应示例

```json
"id"
```

### PUT/:id 修改卡片

#### 请求参数

- id：路径参数，必需，卡片id

- question：字符串，必需，题面
- answer：字符串，必需，答案
- image：字符串，可选，图片url
- audio：字符串，可选，音频url

#### 请求示例

```json
{
    "question":"question",
    "answer":"answer",
    "image":"https://example.com/1.jpg",
    "audio":"https://example.com/1.wav"
}
```

#### 响应参数

卡片id

#### 响应示例

```json
"id"
```

### GET /:id 获取卡片

#### 请求参数

- id：路径参数，必需，卡片id

#### 响应参数

- question：字符串，题面
- answer：字符串，答案
- image：字符串，图片url
- audio：字符串，音频url

#### 响应示例

```json
{
    "question":"question",
    "answer":"answer",
    "image":"https://example.com/1.jpg",
    "audio":"https://example.com/1.wav"
}
```

### DELETE /:id 删除卡片

#### 请求参数

- id：路径参数，必需，卡片id

#### 响应参数

无额外参数

#### 响应示例

```json
"ok"
```

### POST /:id/comment 发表评论

#### 请求参数

- id：路径参数，必需，卡片id
- comment：字符串，必需，评论内容

#### 请求示例

```json
{
    "comment":"comment"
}
```

#### 响应参数

评论id

#### 响应示例

```json
"ok"
```

### GET /:id/comment 获取评论列表

#### 请求参数

- id：路径参数，必需，卡片id

#### 响应参数

评论数组

#### 响应示例

```json
[
    {
    	"comment":"comment"
	},
    {
    	"comment":"comment"
	}
]
```

### DELETE /:id/comment/:comment_id 删除评论

#### 请求参数

- id：路径参数，必需，卡片id
- comment_id：路径参数，必需，评论id

#### 响应参数

无额外参数

#### 响应示例

```json
"ok"
```

## 服务 /service

### GET /upload_token 获取七牛上传凭证

#### 请求参数

- type：字符串，必需，文件类型
  - 值为`avatar`,`image`,`audio`之一

#### 响应参数

- url：字符串，文件url
- resource_key：字符串，资源键
- token：字符串，七牛上传凭证

#### 响应示例

```json
{
    "url": "https://flicker-static.hust.online/avatar/4aac8aa5-7965-4833-be01-b8d9ab7a6f56",
    "resource_key": "avatar/4aac8aa5-7965-4833-be01-b8d9ab7a6f56",
    "token": "_vgfaQkb3E3MjAE0k9aDcOmezkbXBcFX4bqA2WSS:bV9mvzzIvMGxGqxPqqAm_0Mx6kE=:eyJzY29wZSI6ImZsaWNrZXI6YXZhdGFyLzRhYWM4YWE1LTc5NjUtNDgzMy1iZTAxLWI4ZDlhYjdhNmY1NiIsImRlYWRsaW5lIjoxNjM3NDMzNzUyfQ=="
}
```

### tmp

#### 请求参数

#### 请求示例

```json

```

#### 响应参数

#### 响应示例

```json

```

