# minerx-blockchain-controller

minerx-blockchain-controller 用来实现私有链业务，私有链业务囊括向下承载 2 大技术板块：
- Kubernetes
- Web3

minerx-blockchain-controller 包括的控制器如下：
- **chain-controller**：管理区块链（Chain）资源的生命周期，负责创建、更新、删除链实例
- **chain-sync-controller**：将 Chain 资源状态同步到 MySQL 数据库
- **minerset-controller**：管理矿机集（MinerSet）资源，类似 Kubernetes 的 ReplicaSet，负责维护指定数量的矿机副本
- **minerset-sync-controller**：将 MinerSet 资源状态同步到 MySQL 数据库
- **miner-controller**：管理单个矿机（Miner）资源的生命周期，与底层 Kubernetes Provider 交互创建实际的 Pod
- **miner-sync-controller**：将 Miner 资源状态同步到 MySQL 数据库
- **resource-clean-controller**：清理已删除资源在数据库中的残留数据，保持数据一致性

minerx-blockchain-controller 的开发方式，是 controller-runtime 方式

将所有跟区块链相关的功能聚合在一个 controller 核心优点是：可以复用公共部分的源码

