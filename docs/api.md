Flicker-API 文档

## 基本格式

地址: 暂定

API 前缀: /api/v1，**文档中省略前缀**

### 数据交换格式

#### 成功

```json
{
  "success": true,
  "message": "",
  "data": {}
}
```

data 可以为任何类型，**_之后文档中只写出 data 的数据_**

#### 失败

```json
{
  "success": false,
  "message": "错误提示信息",
  "error": "错误显示信息",
  "data": null
}
```

失败会返回 200 以外的 http 状态码

**错误提示信息是供用户看的文字提示，出现错误时直接使用即可**

之后 API 文档中不会再写失败回复。

## 身份验证

使用 jwt 进行身份验证，token 位于 http 头的 Authorization 字段，以 Bearer 开头，例

Authorization: Bearer 一个巨长的 token

## 基本 /

### POST /verify 发送邮件验证码

- 无需 jwt

- 验证码请求间隔至少为 1 分钟
- 验证码缓存有效期暂定 15 分钟

#### 请求参数

- mail：字符串，必需，邮件地址

#### 请求示例

```json
{
  "mail": "12345678@qq.com"
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
  "mail": "123456789@qq.com",
  "username": "xxx",
  "password": "xxx",
  "code": "12345"
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
  "mail": "123456789@qq.com",
  "password": "xxx"
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
- avatar：字符串，可选，头像 url
  - 这个之后再加，需要考虑一下要不要上七牛

#### 请求示例

```json
{
  "username": "xxx",
  "password": "xxx",
  "code": "12345"
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

- user_id：路径参数，用户 id，缺省时获取当前登录用户信息

#### 响应参数

- username：字符串，可选，用户名
- avatar：字符串，可选，头像 url

#### 响应示例

```json
{
  "id": "6194e37bd786c6c07ea8a4fb",
  "username": "xxx",
  "avatar": "https://example.com/1.jpg",
  "favorite": [
    "6199cc3db4600da8e6102dac"
  ]
}
```

### POST /favorite

#### 请求参数

- cardset_id：字符串，必需，卡片集 ID
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

### GET /favorite 获取用户收藏的卡片集

#### 请求参数

无，用户id从jwt中读取

#### 响应参数

由以下字段组成的数组

- id：字符串，卡片集id
- owner_id：字符串，创建者id
- name：字符串，名称
- description：字符串，卡片集描述
- access：整数，访问权限

#### 响应示例

```json
[
  {
    "id": "6199cc3db4600da8e6102dac",
    "owner_id": "6194e37bd786c6c07ea8a4fb",
    "name": "test-cards",
    "description": "test test.",
    "access": 1
  },
  {
    "id": "6199cc3db4600da8e6102dac",
    "owner_id": "6194e37bd786c6c07ea8a4fb",
    "name": "test-cards",
    "description": "test test.",
    "access": 1
  }
]
```

### GET /created 获取用户创建的卡片集

#### 请求参数

无，用户id从jwt中读取

#### 响应参数

由以下字段组成的数组

- id：字符串，卡片集id
- owner_id：字符串，创建者id
- name：字符串，名称
- description：字符串，卡片集描述
- access：整数，访问权限

#### 响应示例

```json
[
  {
    "id": "6199cc3db4600da8e6102dac",
    "owner_id": "6194e37bd786c6c07ea8a4fb",
    "name": "test-cards",
    "description": "test test.",
    "access": 1
  },
  {
    "id": "6199cc3db4600da8e6102dac",
    "owner_id": "6194e37bd786c6c07ea8a4fb",
    "name": "test-cards",
    "description": "test test.",
    "access": 1
  }
]
```

## 卡片集 /cardset

- alias 题库

### POST / 新建空卡片集

#### 请求参数

- name：字符串，必需，名称
- description：字符串，可选，卡片集描述
- template：字符串，可选，模板
  - 模板一定程度上规定了卡片的样式与布局
  - 暂定为 HTML 格式，具体如何设计看前端
  - 当此字段为空时应使用默认的样式
- access：整数，可选，访问权限
  - 0-仅创建者可见，1-所有人可见
  - 默认权限为 0

#### 请求示例

```json
{
  "name": "test-cards",
  "description": "test test.",
  "access": 1
}
```

#### 响应参数

字符串，卡片集 id

#### 响应示例

```json
"id"
```

### PUT /:id 修改卡片集信息

#### 请求参数

- id：字符串，必需，卡片集 id
- name：字符串，可选，名称
  - 参数留空时表示不做修改，下同
- description：字符串，可选，卡片集描述
- template：字符串，可选，模板
  - 模板一定程度上规定了卡片的样式与布局
  - 暂定为 HTML 格式，具体如何设计看前端
  - 当此字段为空时应使用默认的样式
- access：整数，可选，访问权限

#### 请求示例

```json
{
  "name": "test-cards",
  "description": "test test.",
  "access": 1
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

- id：字符串，必需，卡片集 id

#### 响应参数

无额外参数

#### 响应示例

```json
"ok"
```

### GET /:id 获取卡片集

#### 请求参数

- id：字符串，必需，卡片集 id

#### 响应参数

- name：字符串，名称
- description：字符串，卡片集描述
- owner_id：字符串，创建者id
- owner_name：字符串，创建者用户名
- favorite_count：字符串，收藏数量
- visit_count：字符串，访问数量
- create_time：整数，创建时间戳
- access：整数，访问权限
- cards：字符串数组，卡片的 id 列表
- is_favorite: 布尔值，是否被当前登录用户喜欢

#### 响应示例

```json
{
  "id": "id",
  "name": "test-cards",
  "description": "test test.",
  "owner_id": "6194e37bd786c6c07ea8a4fb",
  "owner_name": "test",
  "favorite_count": 1,
  "visit_count": 7,
  "create_time": 1637469245,
  "access": 1,
  "cards": [],
  "is_favorite": true
}
```

### GET / 搜索卡片集

#### 请求参数

- keyword：字符串，必需，关键词
- skip：整数，可选，跳过的记录条数
  - 默认值为 0
- limit：整数，可选，返回的记录最大条数
  - 默认值为 10

#### 响应参数

由以下参数组成的数组

- name：字符串，名称
- description：字符串，卡片集描述
- favorite_count：字符串，收藏数量
- visit_count：字符串，访问数量
- create_time：整数，创建时间戳
- access：整数，访问权限
- cards：字符串数组，卡片的 id 列表

#### 响应示例

```json
[
  {
    "id": "id",
    "name": "test-cards",
    "description": "test test.",
    "favorite_count": 1,
    "visit_count": 7,
    "create_time": 1637469245,
    "access": 1
  }
]
```

### GET /random 获取随机卡片集

#### 请求参数

- count：整数，必需，随机卡片集数目

#### 响应参数

由以下参数组成的数组

- name：字符串，名称
- description：字符串，卡片集描述
- favorite_count：字符串，收藏数量
- visit_count：字符串，访问数量
- create_time：整数，创建时间戳
- access：整数，访问权限
- cards：字符串数组，卡片的 id 列表

#### 响应示例

```json
[
  {
    "id": "id",
    "name": "test-cards",
    "description": "test test.",
    "favorite_count": 1,
    "visit_count": 7,
    "create_time": 1637469245,
    "access": 1
  }
]
```

## 卡片 /cardset/:cardset_id/card

- cardset_id：路径参数，必需，卡片集 id

### POST / 新建卡片

#### 请求参数

- question：字符串，必需，题面
- answer：字符串，必需，答案
- image：字符串，可选，图片 url
- audio：字符串，可选，音频 url

#### 请求示例

```json
{
  "question": "question",
  "answer": "answer",
  "image": "https://example.com/1.jpg",
  "audio": "https://example.com/1.wav"
}
```

#### 响应参数

卡片 id

#### 响应示例

```json
"id"
```

### POST /many 批量新建卡片

#### 请求参数

- cards：对象数组，必需，由以下字段组成：
  - question：字符串，必需，题面
  - answer：字符串，必需，答案
  - image：字符串，可选，图片 url
  - audio：字符串，可选，音频 url

#### 请求示例

```json
{
  "cards": [
    {
      "question": "question",
      "answer": "answer",
      "image": "https://example.com/1.jpg",
      "audio": "https://example.com/1.wav"
    }
  ]
}
```

#### 响应参数

卡片 id 列表

#### 响应示例

```json
["id1","id2"]
```

### PUT/:id 修改卡片

#### 请求参数

- id：路径参数，必需，卡片 id

- question：字符串，必需，题面
- answer：字符串，必需，答案
- image：字符串，可选，图片 url
- audio：字符串，可选，音频 url

#### 请求示例

```json
{
  "question": "question",
  "answer": "answer",
  "image": "https://example.com/1.jpg",
  "audio": "https://example.com/1.wav"
}
```

#### 响应参数

卡片 id

#### 响应示例

```json
"id"
```

### GET /:id 获取卡片

#### 请求参数

- id：路径参数，必需，卡片 id

#### 响应参数

- id：字符串，卡片id
- question：字符串，题面
- answer：字符串，答案
- image：字符串，图片 url
- audio：字符串，音频 url

#### 响应示例

```json
{
  "id": "id",
  "question": "question",
  "answer": "answer",
  "image": "https://example.com/1.jpg",
  "audio": "https://example.com/1.wav"
}
```

### GET / 批量获取卡片

#### 请求参数

- ids：json字符串，必需，卡片 id 列表

#### 请求示例

```
GET /cardset/6193c1cfd9598aa1a050b041/card?ids=["6199cc46b4600da8e6102dad","619b3df9441373383bc6c6dc"]
```

#### 响应参数

由以下字段组成的数组

- id：字符串，卡片id
- question：字符串，题面
- answer：字符串，答案
- image：字符串，图片 url
- audio：字符串，音频 url

#### 响应示例

```json
[
  {
    "id": "id",
    "question": "question",
    "answer": "answer",
    "image": "https://example.com/1.jpg",
    "audio": "https://example.com/1.wav"
  }
]
```

### DELETE /:id 删除卡片

#### 请求参数

- id：路径参数，必需，卡片 id

#### 响应参数

无额外参数

#### 响应示例

```json
"ok"
```

### POST /:id/comment 发表评论

#### 请求参数

- id：路径参数，必需，卡片 id
- comment：字符串，必需，评论内容

#### 请求示例

```json
{
  "comment": "comment"
}
```

#### 响应参数

评论 id

#### 响应示例

```json
"id"
```

### GET /:id/comment 获取评论列表

#### 请求参数

- id：路径参数，必需，卡片 id

#### 响应参数

数组，每个元素包含**该条评论的 ID**, **发表评论的用户**, **评论内容**以及**评论最后一次更新时间**。
- liked：布尔值，表示当前用户是否点赞过评论
- likes：整数，评论的总点赞数目

注：评论最后一次更新时间在默认情况下为发表评论时间，以Unix时间戳表示，如`1638775217`表示`2021-12-06 15:20:17 +0800 CST`。

#### 响应示例

```json
[
  {
    "id": "id",
    "owner": {
      "id": "id",
      "username": "xxx",
      "avatar": "https://example.com/a.jpg"
    },
    "comment": "comment",
    "lastupdate": "1638775217",
    "liked": false,
    "likes": 0
  },
  {
    "id": "id",
    "owner": {
      "id": "id",
      "username": "xxx",
      "avatar": "https://example.com/a.jpg"
    },
    "comment": "comment",
    "lastupdate": "1638775217",
    "liked": false,
    "likes": 0
  }
]
```

### DELETE /:id/comment/:comment_id 删除评论

#### 请求参数

- id：路径参数，必需，卡片 id
- comment_id：路径参数，必需，评论 id

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

- url：字符串，文件 url
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

## 学习记录 /record

### POST /:cardset_id/:card_id 添加学习记录

#### 请求参数

- cardset_id：路径参数，卡片集id
- card_id：路径参数，卡片id
- status：整数，学习状态
  - 0：未掌握；1：已掌握

当后端接收到这个请求：

- 若这张卡是第一次学习，记录学习时间，并且学习次数为 1；
- 若这张卡已经学过，如果上次学习时间不在当天，则学习次数加 1；更新学习时间；

#### 响应参数

无额外参数

#### 响应示例

```json
"ok"
```

### GET /:cardset_id 获取卡片集学习记录

#### 请求参数

- cardset_id：路径参数，卡片集id

#### 响应参数

- records：由下列参数组成的对象数组：
  - card_id：字符串，卡片id
  - last_study：整数，unix时间戳，最后一次学习时间
  - study_times：整数，学习次数
  - status：整数，掌握状况
    - 0：未掌握；1：已掌握
- total：总卡片数

#### 响应示例

```json
{
  "total": 10,
  "records": [
    {
      "card_id": "xxx",
      "status": 0,
      "study_times": 1,
      "last_study": 0
    },
    {
      "card_id": "xxx",
      "status": 0,
      "study_times": 1,
      "last_study": 0
    }
  ]
}
```

### DELETE /:cardset_id 清空卡片集学习记录

#### 请求参数

- cardset_id：路径参数，卡片集id

#### 响应参数

无额外参数

#### 响应示例

```json
"ok"
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
