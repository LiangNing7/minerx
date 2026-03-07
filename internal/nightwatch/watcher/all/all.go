package all

//nolint: golint
import (
	_ "github.com/LiangNing7/minerx/internal/nightwatch/watcher/cronjob/cronjob"
	_ "github.com/LiangNing7/minerx/internal/nightwatch/watcher/cronjob/statesync"
	_ "github.com/LiangNing7/minerx/internal/nightwatch/watcher/job/llmtrain"
	_ "github.com/LiangNing7/minerx/internal/nightwatch/watcher/secretsclean"
	_ "github.com/LiangNing7/minerx/internal/nightwatch/watcher/user"
)
