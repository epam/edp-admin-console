package k8s

import (
	"k8s.io/apimachinery/pkg/runtime"
)

func (in *BusinessApplication) DeepCopyInto(out *BusinessApplication) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

func (in *BusinessApplication) DeepCopy() *BusinessApplication {
	if in == nil {
		return nil
	}
	out := new(BusinessApplication)
	in.DeepCopyInto(out)
	return out
}

func (in *BusinessApplication) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *BusinessApplicationList) DeepCopyInto(out *BusinessApplicationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BusinessApplication, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

func (in *BusinessApplicationList) DeepCopy() *BusinessApplicationList {
	if in == nil {
		return nil
	}
	out := new(BusinessApplicationList)
	in.DeepCopyInto(out)
	return out
}

func (in *BusinessApplicationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *BusinessApplicationSpec) DeepCopyInto(out *BusinessApplicationSpec) {
	*out = *in
	if in.Repository != nil {
		in, out := &in.Repository, &out.Repository
		*out = new(Repository)
		**out = **in
	}
	if in.Route != nil {
		in, out := &in.Route, &out.Route
		*out = new(Route)
		**out = **in
	}
	if in.Database != nil {
		in, out := &in.Database, &out.Database
		*out = new(Database)
		**out = **in
	}
	return
}

func (in *BusinessApplicationSpec) DeepCopy() *BusinessApplicationSpec {
	if in == nil {
		return nil
	}
	out := new(BusinessApplicationSpec)
	in.DeepCopyInto(out)
	return out
}

func (in *BusinessApplicationStatus) DeepCopyInto(out *BusinessApplicationStatus) {
	*out = *in
	return
}

func (in *BusinessApplicationStatus) DeepCopy() *BusinessApplicationStatus {
	if in == nil {
		return nil
	}
	out := new(BusinessApplicationStatus)
	in.DeepCopyInto(out)
	return out
}

func (in *Database) DeepCopyInto(out *Database) {
	*out = *in
	return
}

func (in *Database) DeepCopy() *Database {
	if in == nil {
		return nil
	}
	out := new(Database)
	in.DeepCopyInto(out)
	return out
}

func (in *Repository) DeepCopyInto(out *Repository) {
	*out = *in
	return
}

func (in *Repository) DeepCopy() *Repository {
	if in == nil {
		return nil
	}
	out := new(Repository)
	in.DeepCopyInto(out)
	return out
}

func (in *Route) DeepCopyInto(out *Route) {
	*out = *in
	return
}

func (in *Route) DeepCopy() *Route {
	if in == nil {
		return nil
	}
	out := new(Route)
	in.DeepCopyInto(out)
	return out
}
