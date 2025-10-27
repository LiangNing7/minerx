## 特性开关

- 用来集中声明、注册并读取开关状态

### 如何使用
- 启用/关闭功能（程序内）：
```go 
_ = feature.DefaultMutableFeatureGate.SetFromMap(map[string]bool{
	"MachinePool": true, // or false
})

```

- 在代码中判断是否开启：
```go 
if feature.DefaultFeatureGate.Enabled(feature.MachinePool) {
	// 执行开启后的逻辑
}
```
