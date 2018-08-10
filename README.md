# dolphin
[![Go Report Card](https://goreportcard.com/badge/github.com/ShiningRush/dolphin)](https://goreportcard.com/report/github.com/ShiningRush/dolphin)

dolphin 是一个轻量级的任务框架，主要用来微服务中那些复杂的跨服务统计查询与传统企业应用中的一些统计报表，你也可以把它当作一个普通的后台任务管理框架。

## 工作方式
dolphin 主要由任务(task)组成，与普通的后台任务处理不同在于任务是由批处理(batch)与同步者(syncer)组成。
