package known

const (
	// This exposes compute information based on the miner type.
	CPUAnnotation    = "apps.liangning7.cn/vCPU"
	MemoryAnnotation = "apps.liangning7.cn/memoryMb"
)

const (
	SkipVerifyAnnotation = "apps.liangning7.cn/skip-verify"
)

var AllImmutableAnnotations = []string{
	CPUAnnotation,
	MemoryAnnotation,
}
