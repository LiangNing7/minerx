# minerx-job-controller

minerx-job-controller 用来实现私有链业务，私有链业务囊括向下承载 2 大技术板块：
- Kubernetes
- Web3

minerx-job-controller 包括的控制器如下：
- cronjob-controller: 负责管理定时任务(CronJob)的控制器，支持按照 Cron 表达式定期创建和执行 Job

minerx-job-controller 的开发方式，是 controller-runtime 方式

将所有跟区块链相关的功能聚合在一个 controller 核心优点是：可以复用公共部分的源码

