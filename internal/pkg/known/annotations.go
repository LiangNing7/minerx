package known

const (
	// This exposes compute information based on the miner type.
	CPUAnnotation    = "apps.LiangNing7.io/vCPU"
	MemoryAnnotation = "apps.LiangNing7.io/memoryMb"
)

const (
	SkipVerifyAnnotation = "apps.LiangNing7.io/skip-verify"
)

var AllImmutableAnnotations = []string{
	CPUAnnotation,
	MemoryAnnotation,
}
