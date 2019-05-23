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

func (in *CodebaseBranch) DeepCopyInto(out *CodebaseBranch) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

func (in *CodebaseBranch) DeepCopy() *CodebaseBranch {
	if in == nil {
		return nil
	}
	out := new(CodebaseBranch)
	in.DeepCopyInto(out)
	return out
}

func (in *CodebaseBranch) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *CodebaseBranchList) DeepCopyInto(out *CodebaseBranchList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CodebaseBranch, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

func (in *CodebaseBranchList) DeepCopy() *CodebaseBranchList {
	if in == nil {
		return nil
	}
	out := new(CodebaseBranchList)
	in.DeepCopyInto(out)
	return out
}

func (in *CodebaseBranchList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *CodebaseBranchSpec) DeepCopyInto(out *CodebaseBranchSpec) {
	*out = *in
	return
}

func (in *CodebaseBranchSpec) DeepCopy() *CodebaseBranchSpec {
	if in == nil {
		return nil
	}
	out := new(CodebaseBranchSpec)
	in.DeepCopyInto(out)
	return out
}

func (in *CodebaseBranchStatus) DeepCopyInto(out *CodebaseBranchStatus) {
	*out = *in
	return
}

func (in *CodebaseBranchStatus) DeepCopy() *CodebaseBranchStatus {
	if in == nil {
		return nil
	}
	out := new(CodebaseBranchStatus)
	in.DeepCopyInto(out)
	return out
}