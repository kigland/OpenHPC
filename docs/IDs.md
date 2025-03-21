# Short Code, Service Tag & Container ID

KHS 总共使用三种 ID 来标识一个容器：

1. 容器 ID: Container ID (cid)
2. 服务标签: Service Tag (svcTag)
3. 短代码: Short Code (sc)

## 容器 ID/Container ID/CID

每一个 VNode 对应着一比一地对应着底层容器，这个底层容器的 ID 就是容器 ID。目前 KHS 使用 Docker 的容器 ID 来标识一个容器。

容器 ID 通常拥有12位和64位两种格式，12位被称呼为 `SCID`。两者在没有哈希碰撞时是等价的。请参阅 Docker 文档。

## 服务标签/Service Tag/SvcTag

Service Tag 类似于 CID，也是用于确定唯一容器的标识符。其格式类似于

```
[Identifier]-[Owner]-[Project]-[Random ID]
```

- `Identifier` 是 KHS 的标识符，默认为 `KHS`。
- `Owner` 是用户名。
- `Project` 是项目名。
- `Random ID` 是随机字符串。

其中 `Project` 是可选的，当用户没有指定项目时，`Project` 为空，内容为：

```
[Identifier]-[Owner]-[Random ID]
```

在实现中，服务标签被标注为容器名（Container Name）。

## 短代码 ShortCode/SC

短代码是对服务标签的简化：

```
[Random ID]@[Owner]/[Project]
```

如果 `Project` 为空，则可以简化为：

```
[Random ID]@[Owner]
```

> 当部署多个 Scheduler 时，需要使用不同 Identifier 来区分不同管理节点。SC不具备 Identifier，因此无法区分不同管理节点。因此在此情况下，需要使用 SvcTag。