# dolphin
[![Go Report Card](https://goreportcard.com/badge/github.com/ShiningRush/dolphin)](https://goreportcard.com/report/github.com/ShiningRush/dolphin)

dolphin 是一个轻量级的任务框架，主要用来解决微服务中那些复杂的跨服务统计查询与传统企业应用中的一些统计报表，你也可以把它当作一个普通的后台任务管理框架。

## Get Started

### 安装

```
go get -u github.com/ShiningRush/dolphin
```

### 使用

```
type testBatch struct {
	Name string
}

func (c *testBatch) GetName() string {
	return c.Name
}

func (c *testBatch) Begin(e *task.EtlTask) error {
	log.Println(c.Name)
	return nil
}

func (c *testBatch) Reset() error {
	return nil
}

func main() {
  oneShotTask := task.NewTask(&testBatch{Name: "oneShotTask"})
  planTask := task.NewTask(&testBatch{Name: "planTask"})
  planTask.PlanTime = "0-59/5 * * * * *" 
  if err := dolphin.Add(oneShotTask).Add(planTask).Build(); err != nil {
		panic("build faild!")
	}
  
  ....
}

```

上面的代码没有贴完全，但在执行完 Build 之后你应该能够看到 OneShotTask 的输出, 并且每隔 5 秒你就能看到 PlanTask的输出

## 工作方式
dolphin 主要由任务(task)组成，与普通的后台任务处理不同在于任务是由批处理(batch)与同步器(syncer)组成。

### 任务(Task)

dolphin 的核心角色，与正常后台任务的控制没有太大区别。
任务可以分成两类

- OneShot
- Plan

如字面意义上所表现的，OneShot任务只会执行一次，而Plan会定时执行。

同时你也可以指定任务在指定的时刻运行，只需要设置任务的 `PlanTime`属性

### 批处理(Batch)

批处理是任务执行所运行的动作，批处理必须实现以下接口

```
// Batch interface
type Batch interface {
	GetName() string
	Begin(t *EtlTask) error
	Reset() error
}
```

GetName 所取得的名字会作为任务的唯一ID，不允许使用相同的名字。
Begin 是任务运行时会执行的一些动作
Reset 可以执行一些消除任务副作用的动作，比如某个数据表每天会全量更新，那么就需要在Begin之前先Reset老的数据

> ### 提示
> 可以通过设置Task的ResetBeforeBegin属性来决定是否在Begin前Reset


### 同步器(syncer)

这是dolphin与其他任务调度框架不同的最大之处，为了解决构建数据仓库中实时同步的问题，我们需要一些同步器来同步实时数据到数据仓库。
同步器需要实现以下接口。

```
// Syncer must be singleton
type Syncer interface {
	Start()
	Stop()
}
```

注意上面的提示：同步器最好是单例的，意味着你在处理更新信号时，要保证使用同一个同步器。
同步器会在任务执行时会调用Start方法，在任务执行完毕后调用Stop方法。

dolphin内置了一个阻塞式的同步器，它提供了一个CheckStatus方法，这个方法会在任务执行时阻塞，一直到任务执行完毕才释放。
使用阻塞同步器只要在你的自定义同步器里组合`syncer.BlockSyncer`，同时在同步器执行时调用 CheckStatus 就行了。

```
// Tradesyncer will handle changes in trade table
type TestSyncer struct {
	syncer.BlockSyncer
}

```

## 仪表盘

dolphin提供了一个Dashboard查看各个任务的执行状态和一些基本的控制，如果要启动它，使用以下的代码：

```
import "github.com/ShiningRush/dolphin/dashserver"

func XXX(){
  hs := dashserver.Start("0.0.0.0:6060")
  ....
  if err := hs.Shutdown(nil); err != nil {
				log.Println("Get error when shutdown dashserver" + err.Error())
				return err
	}
}
```
