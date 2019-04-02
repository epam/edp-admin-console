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

func (in *ApplicationBranch) DeepCopyInto(out *ApplicationBranch) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

func (in *ApplicationBranch) DeepCopy() *ApplicationBranch {
	if in == nil {
		return nil
	}
	out := new(ApplicationBranch)
	in.DeepCopyInto(out)
	return out
}

func (in *ApplicationBranch) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *ApplicationBranchList) DeepCopyInto(out *ApplicationBranchList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ApplicationBranch, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

func (in *ApplicationBranchList) DeepCopy() *ApplicationBranchList {
	if in == nil {
		return nil
	}
	out := new(ApplicationBranchList)
	in.DeepCopyInto(out)
	return out
}

func (in *ApplicationBranchList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *ApplicationBranchSpec) DeepCopyInto(out *ApplicationBranchSpec) {
	*out = *in
	return
}

func (in *ApplicationBranchSpec) DeepCopy() *ApplicationBranchSpec {
	if in == nil {
		return nil
	}
	out := new(ApplicationBranchSpec)
	in.DeepCopyInto(out)
	return out
}

func (in *ApplicationBranchStatus) DeepCopyInto(out *ApplicationBranchStatus) {
	*out = *in
	return
}

func (in *ApplicationBranchStatus) DeepCopy() *ApplicationBranchStatus {
	if in == nil {
		return nil
	}
	out := new(ApplicationBranchStatus)
	in.DeepCopyInto(out)
	return out
}