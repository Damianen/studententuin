package apps

import (
	"context"
	"errors"
	"testing"

	"servermanager/internal/domain"
	"servermanager/internal/mocks"

	"go.uber.org/mock/gomock"
)

const testAppID = "8b9f5e9e-9a3a-4b7e-9a59-1d2f3a4b5c6d"

func testLimits() Limits {
	return Limits{
		DefaultMemoryBytes: 256 * 1024 * 1024,
		MaxMemoryBytes:     1024 * 1024 * 1024,
		DefaultNanoCPUs:    500_000_000,
		MaxNanoCPUs:        2_000_000_000,
		DefaultPidsLimit:   256,
		DefaultRuntime:     domain.RuntimeRunc,
	}
}

func validRunInput() RunInput {
	return RunInput{
		AppID:        testAppID,
		Image:        "busybox:1.37",
		StartCommand: []string{"sleep", "3600"},
	}
}

// expectCreateStart wires a successful create+start and captures the spec the
// usecase built, so cases can assert defaulting and cap handling.
func expectCreateStart(rt *mocks.MockContainerRuntime, captured *domain.ContainerSpec) {
	rt.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, spec domain.ContainerSpec) (string, error) {
			*captured = spec
			return "cid-1", nil
		})
	rt.EXPECT().Start(gomock.Any(), "cid-1").Return(nil)
}

func TestRunApp_Execute_SpecBuilding(t *testing.T) {
	tests := []struct {
		name      string
		input     func() RunInput
		checkSpec func(t *testing.T, spec domain.ContainerSpec)
	}{
		{
			name:  "defaults applied when limits empty",
			input: validRunInput,
			checkSpec: func(t *testing.T, spec domain.ContainerSpec) {
				l := testLimits()
				if spec.MemoryBytes != l.DefaultMemoryBytes {
					t.Errorf("MemoryBytes = %d, want default %d", spec.MemoryBytes, l.DefaultMemoryBytes)
				}
				if spec.NanoCPUs != l.DefaultNanoCPUs {
					t.Errorf("NanoCPUs = %d, want default %d", spec.NanoCPUs, l.DefaultNanoCPUs)
				}
				if spec.PidsLimit != l.DefaultPidsLimit {
					t.Errorf("PidsLimit = %d, want default %d", spec.PidsLimit, l.DefaultPidsLimit)
				}
				if spec.Runtime != domain.RuntimeRunc {
					t.Errorf("Runtime = %q, want runc", spec.Runtime)
				}
				if spec.AppID != testAppID {
					t.Errorf("AppID = %q, want %q", spec.AppID, testAppID)
				}
			},
		},
		{
			name: "explicit limits parsed",
			input: func() RunInput {
				in := validRunInput()
				in.MemoryLimit = "512m"
				in.CpuLimit = "1.5"
				in.Runtime = domain.RuntimeKata
				return in
			},
			checkSpec: func(t *testing.T, spec domain.ContainerSpec) {
				if spec.MemoryBytes != 512*1024*1024 {
					t.Errorf("MemoryBytes = %d, want 512m in bytes", spec.MemoryBytes)
				}
				if spec.NanoCPUs != 1_500_000_000 {
					t.Errorf("NanoCPUs = %d, want 1.5 CPUs", spec.NanoCPUs)
				}
				if spec.Runtime != domain.RuntimeKata {
					t.Errorf("Runtime = %q, want kata", spec.Runtime)
				}
			},
		},
		{
			name: "pids limit is always server-controlled",
			input: func() RunInput {
				in := validRunInput()
				in.Env = map[string]string{"NODE_ENV": "production"}
				in.Port = 3000
				return in
			},
			checkSpec: func(t *testing.T, spec domain.ContainerSpec) {
				if spec.PidsLimit != testLimits().DefaultPidsLimit {
					t.Errorf("PidsLimit = %d, want server default", spec.PidsLimit)
				}
				if spec.Port != 3000 || spec.Env["NODE_ENV"] != "production" {
					t.Errorf("spec = %+v, want port/env passed through", spec)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			rt := mocks.NewMockContainerRuntime(ctrl)

			var spec domain.ContainerSpec
			expectCreateStart(rt, &spec)

			svc := NewService(Dependencies{Runtime: rt, Limits: testLimits()})
			cid, err := svc.Run.Execute(context.Background(), tt.input())
			if err != nil {
				t.Fatalf("Execute returned error: %v", err)
			}
			if cid != "cid-1" {
				t.Errorf("Execute = %q, want cid-1", cid)
			}
			tt.checkSpec(t, spec)
		})
	}
}

func TestRunApp_Validate(t *testing.T) {
	// Validate must never touch the runtime — gomock fails on any call.
	ctrl := gomock.NewController(t)
	svc := NewService(Dependencies{Runtime: mocks.NewMockContainerRuntime(ctrl), Limits: testLimits()})

	if err := svc.Run.Validate(validRunInput()); err != nil {
		t.Errorf("Validate(valid) = %v, want nil", err)
	}

	bad := validRunInput()
	bad.Port = 99999
	if err := svc.Run.Validate(bad); !errors.Is(err, domain.ErrInvalid) {
		t.Errorf("Validate(bad port) = %v, want ErrInvalid", err)
	}
}

func TestRunApp_Execute_Rejections(t *testing.T) {
	tests := []struct {
		name  string
		input func() RunInput
	}{
		{"memory above server max", func() RunInput { in := validRunInput(); in.MemoryLimit = "2g"; return in }},
		{"memory below docker floor", func() RunInput { in := validRunInput(); in.MemoryLimit = "1m"; return in }},
		{"cpu above server max", func() RunInput { in := validRunInput(); in.CpuLimit = "4"; return in }},
		{"unparseable memory", func() RunInput { in := validRunInput(); in.MemoryLimit = "lots"; return in }},
		{"unparseable cpu", func() RunInput { in := validRunInput(); in.CpuLimit = "fast"; return in }},
		{"runtime not on allowlist", func() RunInput { in := validRunInput(); in.Runtime = "gvisor"; return in }},
		{"bind mount volume", func() RunInput { in := validRunInput(); in.Volumes = []string{"/host/etc:/etc"}; return in }},
		{"malformed volume", func() RunInput { in := validRunInput(); in.Volumes = []string{"data"}; return in }},
		{"bad image reference", func() RunInput { in := validRunInput(); in.Image = "not a valid image!!"; return in }},
		{"empty image", func() RunInput { in := validRunInput(); in.Image = ""; return in }},
		{"bad env key", func() RunInput { in := validRunInput(); in.Env = map[string]string{"1BAD-KEY": "v"}; return in }},
		{"port out of range", func() RunInput { in := validRunInput(); in.Port = 70000; return in }},
		{"negative port", func() RunInput { in := validRunInput(); in.Port = -1; return in }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			rt := mocks.NewMockContainerRuntime(ctrl) // no expectations: runtime must not be touched

			svc := NewService(Dependencies{Runtime: rt, Limits: testLimits()})
			if _, err := svc.Run.Execute(context.Background(), tt.input()); !errors.Is(err, domain.ErrInvalid) {
				t.Errorf("Execute error = %v, want ErrInvalid", err)
			}
		})
	}
}

func TestRunApp_Execute_RuntimeFailures(t *testing.T) {
	t.Run("create conflict propagates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		rt := mocks.NewMockContainerRuntime(ctrl)
		rt.EXPECT().Create(gomock.Any(), gomock.Any()).Return("", domain.ErrConflict)

		svc := NewService(Dependencies{Runtime: rt, Limits: testLimits()})
		if _, err := svc.Run.Execute(context.Background(), validRunInput()); !errors.Is(err, domain.ErrConflict) {
			t.Errorf("Execute error = %v, want ErrConflict", err)
		}
	})

	t.Run("failed start removes the container but keeps the network", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		rt := mocks.NewMockContainerRuntime(ctrl)
		startErr := errors.New("boom")
		rt.EXPECT().Create(gomock.Any(), gomock.Any()).Return("cid-1", nil)
		rt.EXPECT().Start(gomock.Any(), "cid-1").Return(startErr)
		rt.EXPECT().RemoveContainer(gomock.Any(), "cid-1").Return(nil)

		svc := NewService(Dependencies{Runtime: rt, Limits: testLimits()})
		if _, err := svc.Run.Execute(context.Background(), validRunInput()); !errors.Is(err, startErr) {
			t.Errorf("Execute error = %v, want the start error", err)
		}
	})
}

func TestStartStopApp_Execute(t *testing.T) {
	wantName := domain.AppContainerName(testAppID)

	t.Run("start uses derived container name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		rt := mocks.NewMockContainerRuntime(ctrl)
		rt.EXPECT().Start(gomock.Any(), wantName).Return(nil)

		svc := NewService(Dependencies{Runtime: rt, Limits: testLimits()})
		if err := svc.Start.Execute(context.Background(), testAppID); err != nil {
			t.Errorf("Execute returned error: %v", err)
		}
	})

	t.Run("stop passes not-found through", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		rt := mocks.NewMockContainerRuntime(ctrl)
		rt.EXPECT().Stop(gomock.Any(), wantName).Return(domain.ErrNotFound)

		svc := NewService(Dependencies{Runtime: rt, Limits: testLimits()})
		if err := svc.Stop.Execute(context.Background(), testAppID); !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("Execute error = %v, want ErrNotFound", err)
		}
	})
}

func TestRemoveApp_Execute(t *testing.T) {
	wantName := domain.AppContainerName(testAppID)

	t.Run("not found is idempotent success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		rt := mocks.NewMockContainerRuntime(ctrl)
		rt.EXPECT().Remove(gomock.Any(), wantName).Return(domain.ErrNotFound)

		svc := NewService(Dependencies{Runtime: rt, Limits: testLimits()})
		if err := svc.Remove.Execute(context.Background(), testAppID); err != nil {
			t.Errorf("Execute returned error: %v, want nil", err)
		}
	})

	t.Run("other errors propagate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		rt := mocks.NewMockContainerRuntime(ctrl)
		rmErr := errors.New("daemon down")
		rt.EXPECT().Remove(gomock.Any(), wantName).Return(rmErr)

		svc := NewService(Dependencies{Runtime: rt, Limits: testLimits()})
		if err := svc.Remove.Execute(context.Background(), testAppID); !errors.Is(err, rmErr) {
			t.Errorf("Execute error = %v, want %v", err, rmErr)
		}
	})
}

func TestGetAppStatus_Execute(t *testing.T) {
	wantName := domain.AppContainerName(testAppID)

	t.Run("state passes through", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		rt := mocks.NewMockContainerRuntime(ctrl)
		want := &domain.ContainerState{Exists: true, Running: true, Status: "running"}
		rt.EXPECT().Inspect(gomock.Any(), wantName).Return(want, nil)

		svc := NewService(Dependencies{Runtime: rt, Limits: testLimits()})
		got, err := svc.Status.Execute(context.Background(), testAppID)
		if err != nil {
			t.Fatalf("Execute returned error: %v", err)
		}
		if got != want {
			t.Errorf("Execute = %+v, want %+v", got, want)
		}
	})

	t.Run("not found means exists false", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		rt := mocks.NewMockContainerRuntime(ctrl)
		rt.EXPECT().Inspect(gomock.Any(), wantName).Return(nil, domain.ErrNotFound)

		svc := NewService(Dependencies{Runtime: rt, Limits: testLimits()})
		got, err := svc.Status.Execute(context.Background(), testAppID)
		if err != nil {
			t.Fatalf("Execute returned error: %v", err)
		}
		if got == nil || got.Exists {
			t.Errorf("Execute = %+v, want Exists false", got)
		}
	})
}

func TestGetAppLogs_Execute(t *testing.T) {
	wantName := domain.AppContainerName(testAppID)
	opts := domain.LogOptions{Tail: 50}

	t.Run("lines pass through", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		rt := mocks.NewMockContainerRuntime(ctrl)
		want := []domain.LogLine{{Stream: "stdout", Message: "hello"}}
		rt.EXPECT().Logs(gomock.Any(), wantName, opts).Return(want, nil)

		svc := NewService(Dependencies{Runtime: rt, Limits: testLimits()})
		got, err := svc.Logs.Execute(context.Background(), testAppID, opts)
		if err != nil {
			t.Fatalf("Execute returned error: %v", err)
		}
		if len(got) != 1 || got[0] != want[0] {
			t.Errorf("Execute = %+v, want %+v", got, want)
		}
	})

	t.Run("not found propagates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		rt := mocks.NewMockContainerRuntime(ctrl)
		rt.EXPECT().Logs(gomock.Any(), wantName, opts).Return(nil, domain.ErrNotFound)

		svc := NewService(Dependencies{Runtime: rt, Limits: testLimits()})
		if _, err := svc.Logs.Execute(context.Background(), testAppID, opts); !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("err = %v, want ErrNotFound", err)
		}
	})
}

func TestFollowAppLogs_Execute(t *testing.T) {
	wantName := domain.AppContainerName(testAppID)
	opts := domain.LogOptions{Tail: 50}

	t.Run("stream passes through", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		rt := mocks.NewMockContainerRuntime(ctrl)
		ch := make(chan domain.LogLine, 1)
		ch <- domain.LogLine{Stream: "stdout", Message: "live"}
		close(ch)
		rt.EXPECT().FollowLogs(gomock.Any(), wantName, opts).Return((<-chan domain.LogLine)(ch), nil)

		svc := NewService(Dependencies{Runtime: rt, Limits: testLimits()})
		got, err := svc.Follow.Execute(context.Background(), testAppID, opts)
		if err != nil {
			t.Fatalf("Execute returned error: %v", err)
		}
		if line, open := <-got; !open || line.Message != "live" {
			t.Errorf("first line = %+v open = %v", line, open)
		}
		if _, open := <-got; open {
			t.Error("channel still open, want closed")
		}
	})

	t.Run("not found propagates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		rt := mocks.NewMockContainerRuntime(ctrl)
		rt.EXPECT().FollowLogs(gomock.Any(), wantName, opts).Return(nil, domain.ErrNotFound)

		svc := NewService(Dependencies{Runtime: rt, Limits: testLimits()})
		if _, err := svc.Follow.Execute(context.Background(), testAppID, opts); !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("err = %v, want ErrNotFound", err)
		}
	})
}
