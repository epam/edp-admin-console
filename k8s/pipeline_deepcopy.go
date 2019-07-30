/*
 * Copyright 2019 EPAM Systems.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package k8s

import (
	"k8s.io/apimachinery/pkg/runtime"
)

func (in *CDPipeline) DeepCopyInto(out *CDPipeline) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

func (in *CDPipeline) DeepCopy() *CDPipeline {
	if in == nil {
		return nil
	}
	out := new(CDPipeline)
	in.DeepCopyInto(out)
	return out
}

func (in *CDPipeline) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *CDPipelineList) DeepCopyInto(out *CDPipelineList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CDPipeline, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

func (in *CDPipelineList) DeepCopy() *CDPipelineList {
	if in == nil {
		return nil
	}
	out := new(CDPipelineList)
	in.DeepCopyInto(out)
	return out
}

func (in *CDPipelineList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *CDPipelineSpec) DeepCopyInto(out *CDPipelineSpec) {
	*out = *in
	if in.CodebaseBranch != nil {
		in, out := &in.CodebaseBranch, &out.CodebaseBranch
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.InputDockerStreams != nil {
		in, out := &in.InputDockerStreams, &out.InputDockerStreams
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ThirdPartyServices != nil {
		in, out := &in.ThirdPartyServices, &out.ThirdPartyServices
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

func (in *CDPipelineSpec) DeepCopy() *CDPipelineSpec {
	if in == nil {
		return nil
	}
	out := new(CDPipelineSpec)
	in.DeepCopyInto(out)
	return out
}

func (in *CDPipelineStatus) DeepCopyInto(out *CDPipelineStatus) {
	*out = *in
	return
}

func (in *CDPipelineStatus) DeepCopy() *CDPipelineStatus {
	if in == nil {
		return nil
	}
	out := new(CDPipelineStatus)
	in.DeepCopyInto(out)
	return out
}
