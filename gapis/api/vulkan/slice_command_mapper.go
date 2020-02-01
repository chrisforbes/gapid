// Copyright (C) 2020 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vulkan

import (
	"context"

	"github.com/google/gapid/core/log"
	"github.com/google/gapid/gapis/api"
	"github.com/google/gapid/gapis/api/transform"
)

type sliceCommandMapper struct{
	offset uint64
	submissionIds *[]uint64
}

func (t *sliceCommandMapper) Transform(ctx context.Context, id api.CmdID, cmd api.Cmd, out transform.Writer) error {
	ctx = log.Enter(ctx, "Command Vk handle mapper")
	s := out.State()

	switch cmd := cmd.(type) {
	case *VkCmdBeginRenderPass:
		cb := cmd.commandBuffer
		beginInfo := cmd.PRenderPassBegin().MustRead(ctx, cmd, s, nil)
		rp := beginInfo.RenderPass()
		fb := beginInfo.Framebuffer()
		log.I(ctx, "HANDLES %v %v %v %v", id, cb, fb, rp)
		return out.MutateAndWrite(ctx, id, cmd)
	case *VkQueueSubmit:
		if uint64(id) >= t.offset {
			*t.submissionIds = append(*t.submissionIds, uint64(id) - t.offset)
		}
		return out.MutateAndWrite(ctx, id, cmd)
	default:
		return out.MutateAndWrite(ctx, id, cmd)

	}

	return nil
}

func (t *sliceCommandMapper) PreLoop(ctx context.Context, out transform.Writer) {
	out.NotifyPreLoop(ctx)
}
func (t *sliceCommandMapper) PostLoop(ctx context.Context, out transform.Writer) {
	out.NotifyPostLoop(ctx)
}
func (t *sliceCommandMapper) Flush(ctx context.Context, out transform.Writer) error { return nil }
func (t *sliceCommandMapper) BuffersCommands() bool {
	return false
}

func newSliceCommandMapper(offset uint64, submissionIds *[]uint64) *sliceCommandMapper {
	return &sliceCommandMapper{offset: offset, submissionIds: submissionIds}
}