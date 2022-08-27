# RobKing_MIT6.824
# MIT 6.824 分布式系统 | Lab 1：MapReduce

## 项目总览

首先明确我们需要做的事情：就是统计电子书所有单词出现的次数。其次这整个框架我们需要用到的就是两个包，一个是`main`包，一个是`mr`包，其中`main`包下的`mrmaster.go`和`mrworker.go`是用来调用`mr`包下的`master.go`和`worker.go`。

这个框架给了我们一个非分布式的实现方法在`main/mrsequential.go`，而我们需要做的就是实现一个分布式系统用来统计。`master`相当于老板，`worker`相当于工人，每个工人会向老板索要任务，老板给他们分配任务，他们之间通过`RPC`进行通信，分配的任务有两种，一种是`Map`任务，即将电子书的所有单词分离出来，通过键值对保存，可能会有很多重复的，所以还有一种任务就是`Reduce`任务，即将重复的键值对进行合并，这样我们就可以得到所有单词出现的次数。

## 必备的背景知识

### `RPC`通信

简单介绍一下`RPC`通信，主要是服务端和客户端，服务端注册服务，客户端调用这个服务就相当于本地调用函数一样。

服务端注册的流程如下：

1. 通过`rpc.Register()`方法注册`RPC`服务，参数为结构体
2. 通过`rpc.HandleHTTP()`将`RPC`服务绑定到`HTTP`，没有参数
3. 通过`http.ListenAndServer()`监听一个端口，第一个参数是`ip+port`，第二个参数是`nil`
4. 实现RPC服务的功能，功能函数必须按照规范来，`func(结构体) 功能(Args, *Result){}`

客服端调用的流程如下：

1. 通过`rpc.DailHTTP()`建立连接，第一个参数是`TCP`协议，第一个参数是`ip+port`
2. 直接通过`Call()`方法调用服务器实现好的功能，第一个参数是功能(string)，第二个为传入的参数结构体，第三个为返回的结果结构体

### `go --plugin`(动态库)

了解`go`语言的`plugin`有利于更好地理解这之间的调用关系。首先会有一个`go`文件实现了一个方法，接下来我们使用如下指令生成`.so`文件

```golang
go build -race -buildmode=plugin *.go
```

接下来主函数(执行 `go run main.go *.so）`可以通过`*.so`文件调用其中的函数或者仅仅获取函数。这其中有两个重要的`API`

- 一个是 `pdll, err := plugin.Open("*.so")`用来打开动态库
- 另一个是 `funcName, err := pdll.LookUp("funcName")`用来获取方法，返回的是`interface{}`类型
- 通过这个方法我们可以得到 `funcName.(函数参数类型)（函数参数）`来调用函数

## 非分布式的实现方法思路

这个是官方提供的在`main/mrsequential.go`，具体是这样实现的。

首先是读取电子书的内容，接着调用`Map`函数执行`Map`任务，`Map`函数是通过将读取到的文件内容转化为单词数组的形式，然后遍历这个数组，将每个单词的`value`记为1，返回键值对。这样我们就获得了键值对，紧接着我们对这些键值对进行排序，使相同的放在一块，遍历这些键值对，将相同的单词执行`Reduce`任务，`Reduce`函数就是传入字符串数组，返回他的长度，这个长度就是我们要的每个单词出现的次数，这个数组主要是将相同的单词的`value`加入到里面，其实就是很多个`“1”`，同时将执行完`Reduce`任务的键值对输入到`mr-out-0`文件即可。

实现起来也是相当的简单，我们可以参考里面读写文件的操作来实现我们的分布式系统。接下来就是分布式的具体实现。

## 基本配置

首先我们要清楚需要文件之间的调用关系以及我们需要关注哪些文件。`main/mrmaster.go`主要是调用`mr/master.go`的`MakeMaster`方法，执行需要的参数是电子书，可以是多本；`main/mrworker.go`主要是调用`mr/worker.go`的`MakeWorker`方法，执行的参数是`wc.so`，这个文件是通过执行`go build -buildmode=plugin .. /mrapps/wc.go`执行生成的，并且每次执行之前都需要运行这条指令更新`wc.so`文件。我们的目标就是编写`mr/master.go`和`mr/worker.go`及`rpc.go`这几个文件来实现分布式的`MapReduce`。

然后进行环境的配置。

我们可以将生成wc.so文件的执行写入到脚本文件build-wc.sh文件中

```shell
PATH=$PATH:/home/robking/AGolang/go1.18/go/bin
rm -f mr-*
rm -f wc.so
go build -buildmode=plugin .. /mrapps/wc.go
```

第一行代码是go语言的执行器

![](https://raw.githubusercontent.com/RobKing9/Blog_Pic/master/Git/20220827222134.png)

`main/mrmaster.go`配置如下

![](https://raw.githubusercontent.com/RobKing9/Blog_Pic/master/Git/20220827221851.png)

`main/mrmaster`的配置如下

![](https://raw.githubusercontent.com/RobKing9/Blog_Pic/master/Git/20220827222048.png)

## 分布式实现

首先我们先进行系统的设计，我们一开始必须初始化Master和Worker，Worker很简单，需要一个id唯一标识，另外就是mapf和reducef处理Map和Reduce任务的两项能力；Master掌管所有的事情，需要files标记总文件数，NReduce标记总共的Reduce任务数，另外每一项任务都有他的状态，准备中，放入管道了，正在运行，运行完成，运行出错，还需要它的开始时间用来判断一个任务的执行时间，还需要的字段就是工号，将任务和Worker联系起来，那么Master就需要管理所有任务的状态，需要一个TaskStats字段，还有就是发送任务的管道，发放的工号，锁，是否完成了所有的任务。

我们按照以下思路实现

1. `Master`分配`Map`任务，`Worker`拿到之后打印出需要处理的文件名
2. `Worker`通过传过来的文件名做`Map`任务，并输出结果
3. `Master`感知`Map`任务转`Reduce`的时机
4. `Master`分配`Reduce`任务，`Worker`做`Reduce`任务
5. 向`Master`汇报完成，`Master`宣布完成，结束进程

### `Master`分配`Map`任务，`Worker`拿到之后打印出需要处理的文件名

首先从整体上来看，`worker`需要向`master`申请一个工号，这样就可以将工号和处理的任务相对应起来，`Master`的处理逻辑是将工号每次进行+1发放给`Worker`，之后`Worker`不断地根据自己的`id`向`Master`申请任务，而`Master`开始运行的时候就会初始化`Map`任务，然后通过一个单独的协程来根据任务的状态给`Worker`分配任务，并且这个过程是持续的。一开始都处于准备中，我们根据任务`id`进行任务初始化以及放入到任务管道中，并标记状态变为在管道中，接下来`Master`会取出管道中的任务分配给`Worker`，并标记当前时间，将任务id和`workerid`对应起来，同时标记任务状态为执行中，`Worker`通过`rpc`通信获得`Master`的任务开始执行，打印出需要处理的文件名，这一步我们就成功完成了！

### `Worker`通过传过来的文件名做`Map`任务，并输出结果

接下来就是做`Map`任务。根据文件名打开文件，读取文件的内容，调用`mapf`生成很多的键值对`kvs`，然后我们将所有的键值对分成`NReduce`个用来之后执行`Reduce`任务，具体操作是对键值取哈希值，通过这个哈希值对`NReduce`取余，其中相同的键哈希值肯定是一样的，这样所有的键值就分别到`NReduce`个中了。之后我们将每一个`reduce`结果保存到中间文件中，文件名是`mr-TaskId-reduceId`，通过将内容编码写入到文件中即可。

在任务执行过程中我们需要处理一件事情就是是否超时（10s内完成），我们可以`Master`的调度函数，如果任务的状态是正在运行，我们就通过`time.Now().Sub(m.taskStats[taskid].startTime) > MaxTaskRunTime`这样的判断，如果超时了，那么重新标记任务的状态为在管道中，之后将会进行重新的分配。

### `Master`感知`Map`任务转`Reduce`的时机

`Worker`完成了`Map`任务，那么就需要告知`Master`，如果所有的`Map`任务都完成了，那么`Master`将初始化`Reduce`任务然后进行分配。`Worker`汇报任务的时候需要告知`Master`是否完成，完成的任务`id`，你的`workerId`，以及任务类型，之后`Master`会通过`Worker`传过来的信息进行判断，如果没有完成出错了，那么将继续重新分配，继续`Master`的调度函数。`Master`需要一个`allFinish`的全局判断，只有所有的任务都是完成的状态，才能进行初始化`Reduce`任务。

### `Master`分配`Reduce`任务，`Worker`做`Reduce`任务

所有`Map`任务完成了，初始化`Reduce`任务，标记好类型，以及所有的任务状态。继续任务的调度，`Worker`接收到`Reduce`任务开始执行。根据所有的中间文件对文件内容进行解码，取出键值对放入到一个哈希表中，哈希表的`Key`是键值对的`Key`，值是`[]string`类型，为键值对的`Value`数组，然后通过`reducef`函数，计算出每个`Key`出现的次数。将这些内容保存到最终的文件即可

### 向`Master`汇报完成，`Master`宣布完成，结束进程

`Worker`向`Master`汇报`Reduce`任务完成，`Master`还是根据`allFinish`判断，如果当前任务类型是`Reduce`，那么标记`done`字段为`true`，结束进程，所以任务都已经完成。

## 分布式实现总结

`Master`需要启动`server`方法来监听端口实现`rpc`通信，而`Worker`通过`call`方法即可调用`Master`实现好的函数。

加锁机制需要在以下几个方面

- 申请workId的时候
- 申请任务的时候
- 分配任务的时候
- 接受Worker汇报任务的时候

加锁都是防止多个Worker竞争同一资源

多协程机制，分配任务的时候需要启动一个协程，不断的去分配
