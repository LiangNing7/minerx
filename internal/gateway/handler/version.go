package handler

import (
	"context"

	"github.com/LiangNing7/goutils/pkg/version"

	v1 "github.com/LiangNing7/minerx/pkg/api/gateway/v1"
)

func (s *Handler) GetVersion(ctx context.Context, rq *v1.GetVersionRequest) (*v1.GetVersionResponse, error) {
	vinfo := version.Get()
	return &v1.GetVersionResponse{
		GitVersion:   vinfo.GitVersion,
		GitCommit:    vinfo.GitCommit,
		GitTreeState: vinfo.GitTreeState,
		BuildDate:    vinfo.BuildDate,
		GoVersion:    vinfo.GoVersion,
		Compiler:     vinfo.Compiler,
		Platform:     vinfo.Platform,
	}, nil
}
