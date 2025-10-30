package known

const (
	// This exposes compute information based on the miner type.
	CPUAnnotation    = "apps.onex.io/vCPU"
	MemoryAnnotation = "apps.onex.io/memoryMb"
)

const (
	SkipVerifyAnnotation = "apps.onex.io/skip-verify"
)

var AllImmutableAnnotations = []string{
	CPUAnnotation,
	MemoryAnnotation,
}
